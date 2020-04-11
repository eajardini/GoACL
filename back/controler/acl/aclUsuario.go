package acl

import (
	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
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

// msgErro     string
// erroRetorno error
)

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

//NovoUsuario : responsável por receber os dados via requisição, criar objeto modelacl.ACLUsuario e chamar a
// função CriaNovoUsuario() para salvar no BD
func NovoUsuario(c *gin.Context) {

	var (
		userACLJSON modelACL.ACLUsuarioJSON
		userACL     modelACL.ACLUsuario
	)
	erro = c.ShouldBindJSON(&userACLJSON)
	if erro != nil {
		Resp.Mensagem = "[AUSERRNUS001 | aclUsuario|NovoUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	userACL = atribuiDadosUsuario(userACLJSON)
	erro = modelACL.CriaNovoUsuario(userACL, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "Usuário Criado com Sucesso"
	}
	c.JSON(200, Resp)
}

//RemoveUsuario : responsável por receber os dados via requisição e atribuir valor 0 no campo UsuarioAtivo
// e assim apagar logicamente o susuário
func RemoveUsuario(c *gin.Context) {

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

// AtivaUsuario : responsável por receber os dados via requisição e atribuir valor 0 no campo UsuarioAtivo
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
