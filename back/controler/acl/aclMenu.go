package acl

import (
	"log"

	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"
)

//*** Listar os dados dos grupos

//MontaMenu : Retorna o Menu do usuário
func MontaMenu(c *gin.Context) {
	var (
		ItemsNivel1Locais []modelACL.ItemsNivel1
		erroRetorno       error
		Resp              Resposta

		// itensDaTabelaMenu []modelmenu.ItensDaTabelaMenu
		// menuLocal         []modelmenu.ItemsNivel1
		// menuItemNivel1    modelmenu.ItemsNivel1
		// menuItemNivel2    modelmenu.ItemsNivel2
		// menuItemNivel3    modelmenu.ItemsNivel3
		// posN1, posN2      int
	)

	bd.AbreConexao()
	defer bd.FechaConexao()

	ItemsNivel1Locais, erroRetorno = modelACL.MontaMenu(bd)

	if erroRetorno != nil {
		Resp.Mensagem = erro.Error()
		ItemsNivel1Locais[0].Label = "Nenhum menu localizado"
	}

	log.Println("[aclMenu.go|MontaMenu|INFO 001] valor menuLocal:", ItemsNivel1Locais)

	c.JSON(200, ItemsNivel1Locais)
	// c.JSON(200, gin.H{
	// 	"resposta": ItemsNivel1Locais,
	// })
}
