package acl

//Este arquivo não está sendo usado, pois foi substituido pelo
// login contido dentro da Lib do gin-jwt

import (
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

//ACLLogin : contém o username (credencial) e senha (chave) e
type ACLLogin struct {
	Credencial string `json:"Credencial" validate:"required"`
	Chave      string `json:"Chave" validate:"required"`
}

//Login : Realiza o login no sistema.Login
func Login(c *gin.Context) {
	var (
		// menuLocal        []modelACL.ItemsNivel1
		// menuLocalRetorno []modelACL.ItemsNivel1
		loginLocal       ACLLogin
		erroRetornoLocal mensagensErros.LIBErroMSGRetorno
		msgErroLocal     mensagensErros.LIBErroMSGSGBD
		moduloLocal      string
		validateLocal    *validator.Validate
		// usuario          string
		// senha string
	)

	erroRetornoLocal.Erro = c.ShouldBindJSON(&loginLocal)
	if erroRetornoLocal.Erro != nil {
		moduloLocal = "[aclLogin.go|Login|ERRO01] "
		erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(3, moduloLocal)
		c.JSON(200, erroRetornoLocal.Mensagem)
		return
	}

	validateLocal = validator.New()
	erroRetornoLocal.Erro = validateLocal.Struct(loginLocal)
	if erroRetornoLocal.Erro != nil {
		moduloLocal = "[aclLogin.go|Login|ERRO02] "
		erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(141, moduloLocal)
		c.JSON(200, erroRetornoLocal.Mensagem)
		return
	}

	if len(loginLocal.Chave) < 8 {
		moduloLocal = "[aclLogin.go|Login|ERRO03] "
		erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(142, moduloLocal)
		c.JSON(200, erroRetornoLocal.Mensagem)
		return
	}
	bd.AbreConexao()
	defer bd.FechaConexao()

	moduloLocal = ""
	erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(143, moduloLocal)

	c.JSON(200, erroRetornoLocal.Mensagem)
}
