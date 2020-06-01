package acltestes

import (
	"fmt"
	"strings"
	"testing"
	"time"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
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
	bdAcessaMenu bancoDeDados.BDCon
	//	msgErros     mensagensErros.LIBErroMSGRetorno
	libMSG mensagensErros.LIBErroMSGSGBD
)

// var ACLUserTest modelacl.ACLUsuario

//Configura Gin com as rotas para teste. A função abaixo é uma
// cópia do arquivo routeACL.go

//ConfigRouterAcessaMenu :
func ConfigRouterAcessaMenu() *gin.Engine {
	r := gin.Default()
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.GET("/ConcedeAcessoAoMenu", aclcontroler.ConcedeAcessoAoMenu)

		// acl.GET("/BuscaTodosUsuario", aclcontroler.BuscaTodosUsuario)
		// acl.GET("/BuscaTodosUsuariosAtivos", aclcontroler.BuscaTodosUsuariosAtivos)
	}
	return r
}

// GinFazRequisicaoParaAcessoAoMenu : zz
func GinFazRequisicaoParaAcessoAoMenu(t *testing.T, menusPar modelACL.ACLGrupoAcessaMenuJSON, ComparacaoRetorno string) {
	var (
		// Resp          aclcontroler.Resposta
		MsgErrosLocal mensagensErros.LIBErroMSGRetorno
	)

	r := ConfigRouterAcessaMenu()
	w := httptest.NewRecorder()

	dadosGrupo, _ := json.Marshal(menusPar)

	// dadosGrupo := fmt.Sprintf(`
	// {
	// 	"CodigoGrupo": "%s", "CodigoMenu":"%s"
	// }`, menus.CodigoGrupo, menus.CodigoMenu)

	fmt.Println("Valor dos dadosUsuario:", string(dadosGrupo))
	req, _ := http.NewRequest("GET", "/acl/ConcedeAcessoAoMenu", strings.NewReader(string(dadosGrupo)))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &MsgErrosLocal)
	fmt.Println("retorno", MsgErrosLocal)
	assert.Equal(t, ComparacaoRetorno, MsgErrosLocal.Mensagem)
}

//TestGinFazRequisicaoParaAcessoAoMenu :
func TestGinFazRequisicaoParaAcessoAoMenu(t *testing.T) {
	var (
		menuLocal modelACL.ACLGrupoAcessaMenuJSON
	)
	bdAcessaMenu.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bdAcessaMenu.IniciaConexao()

	bdAcessaMenu.AbreConexao()
	libMSG.CarregaTodosAsMensagensDeErro(bdAcessaMenu)
	bdAcessaMenu.FechaConexao()

	dataAtual := time.Now()
	novoGrupoA := "ap" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoA, novoGrupoA, "", 0, "g", "Grupo Criado com Sucesso")

	menuLocal.CodigoGrupo = novoGrupoA
	menuLocal.CodigoMenu = []string{"admi111", "admi112", "ctbl123", "ctbl131"}
	GinFazRequisicaoParaAcessoAoMenu(t, menuLocal, "120 - Direitos de Acesso concedidos com sucesso!")

	dataAtual = time.Now()
	novoGrupoA = "ao" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoA, novoGrupoA, "", 0, "g", "Grupo Criado com Sucesso")

	menuLocal.CodigoGrupo = novoGrupoA
	menuLocal.CodigoMenu = []string{"admi111", "admi112", "ctbl131"}
	GinFazRequisicaoParaAcessoAoMenu(t, menuLocal, "120 - Direitos de Acesso concedidos com sucesso!")
}
