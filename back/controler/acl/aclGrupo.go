package acl

import (
	"log"

	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

func atribuiDadosGrupo(grupoACLJSON modelACL.ACLGrupoJSON) (modelACL.ACLGrupo, error) {
	var (
		grupoACL modelACL.ACLGrupo
		erro     error
	)

	now.TimeFormats = append(now.TimeFormats, "02/01/2006")

	grupoACL.CodigoGrupo = grupoACLJSON.CodigoGrupo
	grupoACL.DescricaoGrupo = grupoACLJSON.DescricaoGrupo
	// grupoACL.DataCriacaoGrupo.Time, erro = now.Parse(grupoACLJSON.DataTeste)
	grupoACL.DataCriacaoGrupo, erro = modelACL.StringParaData(grupoACLJSON.DataCriacaoGrupo) //grupoACLJSON.DataCriacaoGrupo)
	// erro = modelACL.StringParaData(&grupoACL.DataCriacaoGrupo, grupoACLJSON.DataCriacaoGrupo) //grupoACLJSON.DataCriacaoGrupo)
	// grupoACL.DataCriacaoGrupo.Valid = true
	if erro != nil {
		log.Println("[AGPERRADG001 | aclGrupo|atribuiDadosGrupo 001] Houve erro ao fazer Bind do JSON" + erro.Error())
		return grupoACL, erro
	}
	grupoACL.SoftDelete = grupoACLJSON.SoftDelete
	grupoACL.TipoOrigemGrupo = grupoACLJSON.TipoOrigemGrupo
	return grupoACL, erro
}

//NovoGrupo : responsável por receber os dados via requisição, criar objeto modelacl.ACLGrupo e chamar a
// função CriaNovoGrupo() para salvar no BD
func NovoGrupo(c *gin.Context) {

	var (
		grupoACLJSON modelACL.ACLGrupoJSON
		grupoACL     modelACL.ACLGrupo
		resp         Resposta
	)
	erro = c.ShouldBindJSON(&grupoACLJSON)
	if erro != nil {
		resp.Mensagem = "[AGPERRNGP002 | aclGrupo|NovoGrupo 002] Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	grupoACL, erro = atribuiDadosGrupo(grupoACLJSON)
	if erro != nil {
		resp.Mensagem = "[AGPERRNGP001 | aclGrupo|NovoGrupo 001] Houve erro ao fazer Bind do JSON" + erro.Error()
		c.JSON(200, Resp)
		return
	}
	erro = modelACL.CriaNovoGrupo(grupoACL, bd)

	if erro != nil {
		resp.Mensagem = "[AGPERRNGP002 | aclGrupo.go|NovoGrupo 002] Houve erro ao fazer Bind do JSON" + erro.Error()
		c.JSON(200, Resp)
		return

	}
	Resp.Mensagem = "Grupo Criado com Sucesso"
	c.JSON(200, Resp)
}

//*** SoftDelete

//SoftDeleteGrupo : chama a rotina do model que realiza o softdele
func SoftDeleteGrupo(c *gin.Context) {
	var (
		ACLCodigoGrupoJSONLocal modelACL.ACLCodigoGrupoJSON
	)
	erro = c.ShouldBindJSON(&ACLCodigoGrupoJSONLocal)
	if erro != nil {
		Resp.Mensagem = "[AGPERRSDL001 | aclGrupo.go|SoftDelete N.01]Houve erro ao fazer Bind do JSON no parâmetro de busca"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	erro = modelACL.SoftdeleteGrupo(ACLCodigoGrupoJSONLocal.CodigoGrupo, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "[AGPINFSDL002 | aclGrupo.go|SoftDelete 02]  Grupo Removido (SoftDelete) com Sucesso"
	}
	c.JSON(200, Resp)

}

//*** Reverte SoftDelete

//ReverteSoftDeleteGrupo : chama a rotina do model que realiza o softdele
func ReverteSoftDeleteGrupo(c *gin.Context) {
	var (
		ACLCodigoGrupoJSONLocal modelACL.ACLCodigoGrupoJSON
	)
	erro = c.ShouldBindJSON(&ACLCodigoGrupoJSONLocal)
	if erro != nil {
		Resp.Mensagem = "[AGPERRRSD001 | aclGrupo.go|ReverteSoftDeleteGrupo N.01]Houve erro ao fazer Bind do JSON no parâmetro de busca"
		c.JSON(200, Resp)
		return
	}

	bd.AbreConexao()
	defer bd.FechaConexao()

	erro = modelACL.ReverteSoftdeleteGrupo(ACLCodigoGrupoJSONLocal.CodigoGrupo, bd)

	if erro != nil {
		Resp.Mensagem = erro.Error()
	} else {
		Resp.Mensagem = "[AGPINFRSD002 | aclGrupo.go|ReverteSoftDeleteGrupo 02]  Grupo Ativado (reverteu SoftDelete) com Sucesso"
	}
	c.JSON(200, Resp)

}

//*** Listar os dados dos grupos

//ListaTodosOsGrupos : chama a rotina do model que realiza o softdele
func ListaTodosOsGrupos(c *gin.Context) {
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
