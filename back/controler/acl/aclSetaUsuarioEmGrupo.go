package acl

import (
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// InsereUsuarioEmGrupo : inclui usuario em um grupo
func InsereUsuarioEmGrupo(c *gin.Context) {

	var (
		erroLocal                error
		Resp                     Resposta
		ACLUsuarioGrupoJSONLocal modelACL.ACLUsuarioGrupoJSON
	)

	erroLocal = c.ShouldBindJSON(&ACLUsuarioGrupoJSONLocal)
	if erroLocal != nil {
		Resp.Mensagem = "[AUSINFRFU001 | aclUsuario|RemoveFisicamenteUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

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
}

//*** Listar os dados dos grupos

//ListaTodosOsGruposX : chama a rotina do model que realiza o softdele
func ListaTodosOsGruposX(c *gin.Context) {
	var (
		ACLGrupoLocal []modelACL.ACLGrupoJSON
	)

	bd.AbreConexao()
	defer bd.FechaConexao()

	ACLGrupoLocal, erro = modelACL.ListaTodosOsGrupos(bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
		ACLGrupoLocal[0].DescricaoGrupo = "Nenhum grupo localizado"
	}
	c.JSON(200, ACLGrupoLocal)
}
