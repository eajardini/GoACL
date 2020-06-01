package libtestes

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	libJWT "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib"
	modelLib "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//********
//** Para limpar o cache dos testes:
//** a) go test -count=1 aclUsuario_test.go
//** b) go clean -testcache
//** Para o VSCode por reconhecer a limpeza do cache, edite o /etc/profile e faça:
//** 1) coloque a linha:
//** 1.1) export GOFLAGS="-count=1"
//** 2) Salve o arquivo
//** 3) Encerre a sessão
//** 4) Login novamente

var ()

//Configura Gin com as rotas para teste. A função abaixo é uma
// cópia do arquivo routeACL.go

//ConfigRouterJWT : zz
func ConfigRouterJWT() *gin.Engine {
	r := gin.Default()
	//***********
	//**** JWT
	//***********
	authMiddleware := libJWT.ConfiguraJWT()

	lib := r.Group("/lib")
	{
		lib.POST("/Login", authMiddleware.LoginHandler)
	}
	return r
}

// GinFazRequisicaoLoginJWT : zz
func GinFazRequisicaoLoginJWT(t *testing.T, nomeUsuarioPar string, senhaPar string, ComparacaoRetornoPar string) {
	var (
		// erroRetorno        mensagensErros.LIBErroMSGRetorno
		loginJWT             libJWT.LoginJWT
		RetornoTokenJWTLocal libJWT.RetornoTokenJWT
	)

	loginJWT.Username = nomeUsuarioPar
	loginJWT.Password = senhaPar

	dadosLogin, _ := json.Marshal(loginJWT)

	// log.Println("[libJWT_test.go|GinFazRequisicaoLoginJWT|INFO01] Valor dadosLogin:", string(dadosLogin))

	r := ConfigRouterJWT() //Configura as rotas deste arquivo de testes
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/lib/Login", strings.NewReader(string(dadosLogin)))
	r.ServeHTTP(w, req)
	codRetorno, _ := strconv.Atoi(ComparacaoRetornoPar)
	assert.Equal(t, codRetorno, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &RetornoTokenJWTLocal)
	log.Println("[libJWT_test.go|GinFazRequisicaoLoginJWT|INFO02] Valor RetornoTokenJWTLocal:", RetornoTokenJWTLocal)
	assert.Equal(t, codRetorno, RetornoTokenJWTLocal.Code)
}

//TestLibLogin :
func TestLibLogin(t *testing.T) {
	var (
		// erroRetorno modelLib.LIBErroMSGRetorno
		bdJWTLocal     bancoDeDados.BDCon
		libMSGJWTLocal modelLib.LIBErroMSGSGBD
	)

	bdJWTLocal.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bdJWTLocal.IniciaConexao()
	bdJWTLocal.AbreConexao()
	libMSGJWTLocal.CarregaTodosAsMensagensDeErro(bdJWTLocal)
	bdJWTLocal.FechaConexao()

	GinFazRequisicaoLoginJWT(t, "internal", "Intern@l", "200")
	GinFazRequisicaoLoginJWT(t, "admin", "admi", "401")

}
