package acl

import (
	"fmt"
	"strings"
	"testing"
	"time"

	modelacl "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
)

var (
	bd          bancoDeDados.BDCon
	msgErro     string
	erroRetorno error
)

// var ACLUserTest modelacl.ACLUsuario

//Configura Gin com as rotas para teste. A função abaixo é uma
// cópia do arquivo routeACL.go

func ConfigRouter() *gin.Engine {
	r := gin.Default()
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.POST("/NovoUsuario", NovoUsuario)
	}
	return r
}

func TestModelCriaNovoUsuario(t *testing.T) {
	now.TimeFormats = append(now.TimeFormats, "02/01/2006")
	assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	dataAtual := time.Now()

	ACLUserTest := modelacl.ACLUsuario{Login: "joao", Password: "teste", Datacriacao: dataAtual,
		Datavalidade: dataAtual, Userbloqueado: false, Userativo: true}

	err := modelacl.CriaNovoUsuario(ACLUserTest)
	assert.NotEqual(err, nil, "[aclUsuario_test|TestCriaNovoUsuario N.01] Erro ao inserir um novo usuário")

}

func TestGinNovoUsuario(t *testing.T) {
	var Resposta resposta

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s", "password":"%s"
	}`, "joao", "123456")

	fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("POST", "/acl/NovoUsuario", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resposta)
	assert.Equal(t, "Usuário Criado com Sucesso", Resposta.Mensagem)
}
