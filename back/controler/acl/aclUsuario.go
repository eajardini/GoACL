package acl

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"

	// mensagensErrosLib "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

//Resposta : Armazena as mensagens de resposta das funções
type Resposta struct {
	Mensagem string
}

var (
	Resp Resposta
	erro error
	bd   bancoDeDados.BDCon
	// libMSG mensagensErros.LIBErroMSGSGBD

// msgErro     string
// erroRetorno error
)

//FazAutenticacao : recebe o login e senha e confronta com os dados do BD.
// se coincidir, retorna true; senao retorn false.
func FazAutenticacao(userNamePar string, passwordPar string) bool {
	var (
		bdLocal          bancoDeDados.BDCon
		AutenticadoLocal bool
		senhaLocal       string
	)

	// bdLocal = bd
	bdLocal.AbreConexao()
	defer bdLocal.FechaConexao()

	h := md5.New()
	h.Write([]byte(passwordPar))
	senhaLocal = hex.EncodeToString(h.Sum(nil))

	AutenticadoLocal = modelACL.FazAutenticacao(userNamePar, senhaLocal, bdLocal)

	return AutenticadoLocal

}

//criaGrupoParaOUsuario : todo usuário deve ter um grupo associado a ele. É pelo
// grupo que é feito a atribuicao dos direitos dos usuários
func criaGrupoParaOUsuario(userACLJSONPar modelACL.ACLUsuarioJSON) modelACL.ACLGrupo {
	var (
		ACLGrupo modelACL.ACLGrupo
		// msgErroLocal         mensagensErros.LIBErroMSGSGBD
		// erroRetornoLocal   mensagensErros.LIBErroMSGRetorno
	)

	now.TimeFormats = append(now.TimeFormats, "02/01/2006")

	ACLGrupo.CodigoGrupo = userACLJSONPar.Login
	ACLGrupo.DescricaoGrupo = userACLJSONPar.Login
	ACLGrupo.DataCriacaoGrupo, _ = modelACL.StringParaData(userACLJSONPar.Datacriacao)
	ACLGrupo.SoftDelete = 0
	ACLGrupo.TipoOrigemGrupo = "u"

	return ACLGrupo
}

func atribuiDadosUsuario(userACLJSON modelACL.ACLUsuarioJSON) modelACL.ACLUsuario {
	var (
		userACL modelACL.ACLUsuario
	)

	now.TimeFormats = append(now.TimeFormats, "02/01/2006")

	userACL.Login = userACLJSON.Login
	userACL.Password = userACLJSON.Login
	userACL.Login = userACLJSON.Login
	userACL.Password = userACLJSON.Password
	userACL.Datacriacao, _ = now.Parse(userACLJSON.Datacriacao)
	userACL.Datavalidade, _ = now.Parse(userACLJSON.Datavalidade)
	userACL.Userbloqueado = userACLJSON.Userbloqueado
	userACL.Userativo = userACLJSON.Userativo

	return userACL

}

//NovoUsuarioACL : recebe os dados da função NovoUsuario e realiza as atividades
// de criação do usuário, grupo do usuario e insere o usuario em seu grupo.
func NovoUsuarioACL(userACLJSONPar modelACL.ACLUsuarioJSON, bdPar bancoDeDados.BDCon) mensagensErros.LIBErroMSGRetorno {
	var (
		ACLGrupo modelACL.ACLGrupo
		// userACLJSON      modelACL.ACLUsuarioJSON
		userACL          modelACL.ACLUsuario
		msgErroLocal     mensagensErros.LIBErroMSGSGBD
		erroRetornoLocal mensagensErros.LIBErroMSGRetorno
		moduloLocal      string
	)

	// bdPar.AbreConexao()
	// defer bdPar.FechaConexao()

	userACL = atribuiDadosUsuario(userACLJSONPar)
	log.Println("[aclUsuario.go|valor userACL]" + userACL.Datacriacao.String())
	erro = modelACL.CriaNovoUsuario(userACL, bdPar)

	if erro != nil {
		erroRetornoLocal.Mensagem = erro.Error()
		// c.JSON(200, erroRetornoLocal.Mensagem)
		return erroRetornoLocal
	}

	//Cria o grupo do usuário
	ACLGrupo = criaGrupoParaOUsuario(userACLJSONPar)
	erroRetornoLocal.Erro = modelACL.CriaNovoGrupo(ACLGrupo, bdPar)
	if erroRetornoLocal.Erro != nil {
		moduloLocal = "[aclUsuario.go|NovoUsuario|ERRO02] "
		log.Println(moduloLocal + erroRetornoLocal.Erro.Error())
		erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(22, moduloLocal)
		//c.JSON(200, erroRetornoLocal.Mensagem)
		return erroRetornoLocal
	}

	//Insere o usuário em seu grupo
	erroRetornoLocal.Erro = modelACL.InsereUsuarioEmGrupo(userACL.Login, ACLGrupo.CodigoGrupo, bdPar)
	if erroRetornoLocal.Erro != nil {
		moduloLocal = "[aclUsuario.go|NovoUsuario|ERRO03] "
		log.Println(moduloLocal + erroRetornoLocal.Erro.Error())
		erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(24, moduloLocal)
		// c.JSON(200, erroRetornoLocal.Mensagem)
		return erroRetornoLocal
	}

	moduloLocal = ""
	erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(23, moduloLocal)
	log.Println("[aclusuario.go|NovoUsuario] Valor da mensagem:" + erroRetornoLocal.Mensagem)
	return erroRetornoLocal

}

//NovoUsuario : responsável por receber os dados via requisição e repassar a função NovoUsuarioACL para que o
// usuario, grupo e a insercão do usuaro em seu grupo seja realizada.
func NovoUsuario(c *gin.Context) {

	var (
		userACLJSON      modelACL.ACLUsuarioJSON
		erroRetornoLocal mensagensErros.LIBErroMSGRetorno
	)

	erro = c.ShouldBindJSON(&userACLJSON)
	if erro != nil {
		Resp.Mensagem = "[AUSERRNUS001 | aclUsuario|NovoUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	erroRetornoLocal = NovoUsuarioACL(userACLJSON, bd)
	c.JSON(200, erroRetornoLocal.Mensagem)
}

//DesativaUsuario : responsável por receber os dados via requisição e atribuir valor 0 no campo UsuarioAtivo
// e assim desativar (apagar logicamente, fazer soft delete) o susuário
func DesativaUsuario(c *gin.Context) {

	var (
		ACLUsuarioLoginJSON modelACL.ACLUsuarioLoginJSON
	)
	erro = c.ShouldBindJSON(&ACLUsuarioLoginJSON)
	if erro != nil {
		Resp.Mensagem = "[AUSINFRUR001 | aclUsuario|RemoveUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	erro = modelACL.RemoveUsuarioPorLogin(ACLUsuarioLoginJSON.Login, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "[AUSINFRUR002 | aclUsuario|RemoveUsuario 02]  Usuário Removido com Sucesso"
	}
	c.JSON(200, Resp)
}

// RemoveFisicamenteUsuario : zz
func RemoveFisicamenteUsuario(c *gin.Context) {

	var (
		ACLUsuarioLoginJSON modelACL.ACLUsuarioLoginJSON
	)
	erro = c.ShouldBindJSON(&ACLUsuarioLoginJSON)
	if erro != nil {
		Resp.Mensagem = "[AUSINFRFU001 | aclUsuario|RemoveFisicamenteUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	erro = modelACL.RemoveFisicamenteUsuarioPorLogin(ACLUsuarioLoginJSON.Login, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "[AUSINFRFU002 | aclUsuario|RemoveFisicamenteUsuario 02]  Usuário Removido Fisicamente com Sucesso"
	}
	c.JSON(200, Resp)
}

// Ativa Usuário

// AtivaUsuario : responsável por receber os dados via requisição e atribuir valor 1 no campo UsuarioAtivo
// e assim apagar logicamente o susuário
func AtivaUsuario(c *gin.Context) {

	var (
		ACLUsuarioLoginJSON modelACL.ACLUsuarioLoginJSON
	)
	erro = c.ShouldBindJSON(&ACLUsuarioLoginJSON)
	if erro != nil {
		Resp.Mensagem = "[AUSINFATU001 | aclUsuario |AtivaUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	erro = modelACL.AtivaUsuarioPorLogin(ACLUsuarioLoginJSON.Login, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "[AUSINFATU002 | aclUsuario|AtivaUsuario 02]  Usuário Ativado com Sucesso"
	}
	c.JSON(200, Resp)
}

// AlteraUsuario : responsável por receber os dados via requisição e alterar os dados do usuário
func AlteraUsuario(c *gin.Context) {

	var (
		userACLJSON modelACL.ACLUsuarioJSON
		userACL     modelACL.ACLUsuario
	)
	erro = c.ShouldBindJSON(&userACLJSON)
	if erro != nil {
		Resp.Mensagem = "[AUSERRAUS001 | aclUsuario|AlteraUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	userACL = atribuiDadosUsuario(userACLJSON)
	erro = modelACL.AlteraUsuario(userACL, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "[AUSINFAUS002 | aclUsuario|AlteraUsuario N.02]Usuário alterado com sucesso"
	}
	c.JSON(200, Resp)
}

// BuscaTodosUsuario : responsável por buscar todos os dados dos usuários
func BuscaTodosUsuario(c *gin.Context) {
	var (
		UserACLRetorno []modelACL.ACLUsuario
		// errolocal      error
	)

	bd.AbreConexao()
	defer bd.FechaConexao()

	UserACLRetorno, _ = modelACL.BuscaTodosUsuario(bd)

	c.JSON(200, UserACLRetorno)
}

// BuscaTodosUsuariosAtivos : responsável por buscar todos os dados dos usuários
func BuscaTodosUsuariosAtivos(c *gin.Context) {
	var (
		UserACLRetorno []modelACL.ACLUsuario
		// errolocal      error
	)
	bd.AbreConexao()
	defer bd.FechaConexao()

	UserACLRetorno, _ = modelACL.BuscaTodosUsuariosAtivos(bd)

	c.JSON(200, UserACLRetorno)
}

// BuscaUsuarioPorLogin : responsável por buscar todos os dados dos usuários
func BuscaUsuarioPorLogin(c *gin.Context) {
	var (
		UserACLRetorno      modelACL.ACLUsuario
		ACLUsuarioLoginJSON modelACL.ACLUsuarioLoginJSON
		errolocal           error
	)
	bd.AbreConexao()
	defer bd.FechaConexao()

	errolocal = c.ShouldBindJSON(&ACLUsuarioLoginJSON)
	if errolocal != nil {
		UserACLRetorno.Login = "[AUSERRBPL001 | aclUsuario |BuscaUsuarioPorLogin N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, UserACLRetorno)
		return
	}

	UserACLRetorno, errolocal = modelACL.BuscaUsuarioPorLogin(ACLUsuarioLoginJSON.Login, bd)

	if errolocal != nil {
		UserACLRetorno.Login = "Usuário não encontrado"
		c.JSON(200, UserACLRetorno)
		return
	}
	c.JSON(200, UserACLRetorno)
}
