package acltestes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	aclControlerModel "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
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

var (
	bdMenu bancoDeDados.BDCon
)

// 	msgErro     string
// 	erroRetorno error
// )

// var ACLUserTest modelacl.ACLUsuario

func ConfigRouterMenu() *gin.Engine {
	r := gin.Default()
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.GET("/MontaMenu", aclcontroler.MontaMenu)

		// acl.GET("/BuscaTodosUsuario", aclcontroler.BuscaTodosUsuario)
		// acl.GET("/BuscaTodosUsuariosAtivos", aclcontroler.BuscaTodosUsuariosAtivos)
	}
	return r
}

// GinFazRequisicaoMontaMenu : zz
func GinFazRequisicaoMontaMenu(t *testing.T, ComparacaoRetorno string) {
	var (
		ItemsNivel1Locais []modelACL.ItemsNivel1
	)

	r := ConfigRouterMenu()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/acl/MontaMenu", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &ItemsNivel1Locais)
	log.Println("[aclMenu_test.go|GinFazRequisicaoMontaMenu N.01|INFO] Valor retornado:", ItemsNivel1Locais)
	assert.NotEqual(t, ComparacaoRetorno, ItemsNivel1Locais)
	assert.Greater(t, len(ItemsNivel1Locais), 0)
}

func TestGinFazRequisicaoMontaMenu(t *testing.T) {

	bdMenu.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bdMenu.IniciaConexao()
	bdMenu.AbreConexao()
	defer bdMenu.FechaConexao()

	GinFazRequisicaoMontaMenu(t, "Nenhum menu localizado")
}

func TestGinMenuFromArquivoJSON(t *testing.T) {
	var (
		ACLMenuFromJSONLocal aclControlerModel.ACLMenuFromJSON
		erroRetorno          error
	)

	bdMenu.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bdMenu.IniciaConexao()
	bdMenu.AbreConexao()
	defer bdMenu.FechaConexao()

	caminho := "../../../config/menu.json"
	ACLMenuFromJSONLocal, erroRetorno = aclControlerModel.CarregaMenuDoJSON(caminho)

	for i := 0; i < len(ACLMenuFromJSONLocal); i++ {
		fmt.Println("Código Menu: ", ACLMenuFromJSONLocal[i].Items[0].Items[0].Label)
		fmt.Println("Label Menu: ", ACLMenuFromJSONLocal[i].Label)
	}

	assert.Equal(t, erroRetorno, nil)

	erroRetorno = aclControlerModel.InsereNoBDMenuDoJSON(ACLMenuFromJSONLocal, bdMenu)

	assert.Equal(t, erroRetorno, nil)
}
