package acl

import (
	"log"

	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/gin-gonic/gin"
)

// ConcedeAcessoAoMenu : inclui usuario em um grupo
func ConcedeAcessoAoMenu(c *gin.Context) {

	var (
		modulo string
		// libMSG                      mensagensErros.LIBErroMSGSGBD
		MsgErrosRetornoLocal        mensagensErros.LIBErroMSGRetorno
		ACLGrupoAcessaMenuJSONLocal modelACL.ACLGrupoAcessaMenuJSON
	)
	MsgErrosRetornoLocal.Erro = c.ShouldBindJSON(&ACLGrupoAcessaMenuJSONLocal)

	if MsgErrosRetornoLocal.Erro != nil {
		modulo = "[aclGrupoAcessaMenu|ConcedeAcessoAoMenu|ERRO01] "
		MsgErrosRetornoLocal = libMSG.BuscaMensagemPeloCodigo(3, modulo)
		// MsgErrosRetornoLocal.Mensagem = MsgErrosRetornoLocal.Erro.Error()
		c.JSON(200, MsgErrosRetornoLocal)
		return
	}
	bd.AbreConexao()
	defer bd.FechaConexao()

	MsgErrosRetornoLocal = modelACL.ConcedeAcessoAoMenu(ACLGrupoAcessaMenuJSONLocal.CodigoMenu, ACLGrupoAcessaMenuJSONLocal.CodigoGrupo, bd)

	if MsgErrosRetornoLocal.Erro != nil {
		log.Println("Erro que veio:", MsgErrosRetornoLocal.Erro.Error())
		modulo = "[aclGrupoAcessaMenu|ConcedeAcessoAoMenu|ERRO02] "
		MsgErrosRetornoLocal.Mensagem = modulo + MsgErrosRetornoLocal.Mensagem
		c.JSON(200, MsgErrosRetornoLocal)
		return
	}

	MsgErrosRetornoLocal = libMSG.BuscaMensagemPeloCodigo(120, "")
	// MsgErrosRetornoLocal.Mensagem = MsgErrosRetornoLocal.Erro.Error()
	c.JSON(200, MsgErrosRetornoLocal)
	/*
		validate := validator.New()
		erroLocal = validate.Struct(ACLUsuarioGrupoJSONLocal)

		if erroLocal != nil {
			Resp.Mensagem = "[ASUERRIUG002 | SetaUsuarioEmGrupo.go|InsereUsuarioEmGrupo002] Valor do Parâmetro ACLUsuarioGrupoPar não validado" + erroLocal.Error()
			// log.Println(Resp.Mensagem)
			c.JSON(200, Resp)
			return
		}

		bd.AbreConexao()
		defer bd.FechaConexao()

		erroLocal = modelACL.InsereUsuarioEmGrupo(ACLUsuarioGrupoJSONLocal.Login, ACLUsuarioGrupoJSONLocal.CodigoGrupo, bd)
		if erroLocal != nil {
			Resp.Mensagem = erroLocal.Error()
			c.JSON(200, Resp)
			return
		}

		Resp.Mensagem = "[ASUINFIUG001 | SetaUsuarioEmGrupo.go|InsereUsuarioEmGrupo001] Usuário inserido no Grupo com sucesso"
		// Resp.Mensagem = "121"
		c.JSON(200, Resp)
	*/
}
