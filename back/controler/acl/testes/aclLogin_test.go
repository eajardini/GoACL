package acltestes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"

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

// 	msgErro     string
// 	erroRetorno error
// )

var (
	bdLogin bancoDeDados.BDCon
)

func ConfigRouterLogin() *gin.Engine {
	r := gin.Default()
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.POST("/Login", aclcontroler.Login)

		// acl.GET("/BuscaTodosUsuario", aclcontroler.BuscaTodosUsuario)
		// acl.GET("/BuscaTodosUsuariosAtivos", aclcontroler.BuscaTodosUsuariosAtivos)
	}
	return r
}

// GinFazRequisicaoMontaMenu : zz
func GinFazRequisicaoLogin(t *testing.T, nomeUsuarioPar string, senhaPar string, ComparacaoRetorno string) {
	var (
		MsgErrosLocal mensagensErros.LIBErroMSGRetorno
		login         aclcontroler.ACLLogin
	)

	login.Credencial = nomeUsuarioPar
	login.Chave = senhaPar

	dadosLogin, _ := json.Marshal(login)

	r := ConfigRouterLogin()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/acl/Login", strings.NewReader(string(dadosLogin)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &MsgErrosLocal.Mensagem)
	assert.Equal(t, ComparacaoRetorno, MsgErrosLocal.Mensagem)
}

func TestGinLogin(t *testing.T) {
	var ()

	bdLogin.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bdLogin.IniciaConexao()
	bdLogin.AbreConexao()
	libMSG.CarregaTodosAsMensagensDeErro(bdLogin)
	bdLogin.FechaConexao()

	GinFazRequisicao(t, "login", "login123", "31/12/2020", "", 0, 1, "[MAUERRCNU003 | modelAclUsuario.go|CriaNovousuario N.03] Erro de insert: Usuário já existe: login")
	GinFazRequisicaoLogin(t, "login", "login123", "ok")
	GinFazRequisicaoLogin(t, "", "login123", "[aclLogin.go|Login|ERRO02] 141 - Nome do usuário ou senha não podem estar em branco.")
	GinFazRequisicaoLogin(t, "login", "", "[aclLogin.go|Login|ERRO02] 141 - Nome do usuário ou senha não podem estar em branco.")
	GinFazRequisicaoLogin(t, "login", "login", "[aclLogin.go|Login|ERRO03] 142 - A senha deve conter no mínimo 8 caracteres.")
}
