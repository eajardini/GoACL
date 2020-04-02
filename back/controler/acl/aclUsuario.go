package acl

import (
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

type resposta struct {
	Mensagem string
}

var (
	Resposta resposta
	erro     error

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
		Resposta.Mensagem = "[aclUsuario|NovoUsuario N.01]Houve erro ao fazer Bind do JSON"
		c.JSON(200, Resposta)
		return
	}

	userACL = atribuiDadosUsuario(userACLJSON)
	erro = modelACL.CriaNovoUsuario(userACL)

	if erro != nil {
		Resposta.Mensagem = "Houve erro ao criar o usuário"
	} else {
		Resposta.Mensagem = "Usuário Criado com Sucesso"
	}
	c.JSON(200, Resposta)

}
