package acltestes

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	modelacl "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	"github.com/gin-gonic/gin"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/jinzhu/now"
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

// var (
// 	bd          bancoDeDados.BDCon
// 	msgErro     string
// 	erroRetorno error
// )

// var ACLUserTest modelacl.ACLUsuario

//Configura Gin com as rotas para teste. A função abaixo é uma
// cópia do arquivo routeACL.go

func ConfigRouterGrupo() *gin.Engine {
	r := gin.Default()
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.POST("/NovoGrupo", aclcontroler.NovoGrupo)
		acl.POST("/ReverteSoftDeleteGrupo", aclcontroler.ReverteSoftDeleteGrupo)
		acl.DELETE("/SoftDeleteGrupo", aclcontroler.SoftDeleteGrupo)
		acl.GET("/ListaTodosOsGrupos", aclcontroler.ListaTodosOsGrupos)
		// acl.GET("/BuscaTodosUsuario", aclcontroler.BuscaTodosUsuario)
		// acl.GET("/BuscaTodosUsuariosAtivos", aclcontroler.BuscaTodosUsuariosAtivos)

	}
	return r
}

// GinFazRequisicao : zz
func GinFazRequisicaoNovoGrupo(t *testing.T, CodigoGrupoPar string, DescricaoGrupoPar string,
	DataCriacaoGrupoPar string, SoftDeletePar int, TipoOrigemGrupoPar string, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouterGrupo()
	w := httptest.NewRecorder()
	dadosGrupo := fmt.Sprintf(`
	{
		"CodigoGrupo": "%s", "DescricaoGrupo":"%s","DataCriacaoGrupo": "%s",
		"SoftDelete": %d, "TipoOrigemGrupo":"%s"
	}`, CodigoGrupoPar, DescricaoGrupoPar, DataCriacaoGrupoPar, SoftDeletePar, TipoOrigemGrupoPar)

	fmt.Println("Valor dos dadosUsuario:", dadosGrupo)
	req, _ := http.NewRequest("POST", "/acl/NovoGrupo", strings.NewReader(dadosGrupo))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)
}

func TestGinNovoGrupo(t *testing.T) {

	now.TimeFormats = append(now.TimeFormats, "02/01/2006")
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	// bd.AbreConexao()
	// defer bd.FechaConexao()

	dataAtual := time.Now()
	novoGrupoA := "aa" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoA, novoGrupoA, "", 0, "g", "Grupo Criado com Sucesso")

	dataAtual = time.Now()
	novoGrupoB := "ab" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoB, novoGrupoB, "09/12/2023", 0, "g", "Grupo Criado com Sucesso")
}

//** SoftdeleteGrupo

// GinFazRequisicao : zz
func GinSoftDeleteGrupo(t *testing.T, CodigoGrupoPar string, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouterGrupo()
	w := httptest.NewRecorder()
	dadosGrupo := fmt.Sprintf(`
	{
		"CodigoGrupo": "%s"}`, CodigoGrupoPar)

	fmt.Println("Valor dos dadosGrupo:", dadosGrupo)
	req, _ := http.NewRequest("DELETE", "/acl/SoftDeleteGrupo", strings.NewReader(dadosGrupo))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)
}

func TestGinSoftDeleteGrupo(t *testing.T) {

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	dataAtual := time.Now()
	novoGrupoA := "ac" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoA, novoGrupoA, "31/12/2020", 0, "g", "Grupo Criado com Sucesso")

	dataAtual = time.Now()
	novoGrupoB := "ad" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoB, novoGrupoB, "31/12/2020", 0, "g", "Grupo Criado com Sucesso")

	GinSoftDeleteGrupo(t, novoGrupoA, "[AGPINFSDL002 | aclGrupo.go|SoftDelete 02]  Grupo Removido (SoftDelete) com Sucesso")
	GinSoftDeleteGrupo(t, novoGrupoB, "[AGPINFSDL002 | aclGrupo.go|SoftDelete 02]  Grupo Removido (SoftDelete) com Sucesso")
	GinSoftDeleteGrupo(t, "xxxx", "[MAGERRSDL001 | modelAclGrupos.go|Softdelete 01] Grupo não localizado")
}

//** Reverte SoftdeleteGrupo

// GinFazRequisicao : zz
func GinReverteSoftDeleteGrupo(t *testing.T, CodigoGrupoPar string, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouterGrupo()
	w := httptest.NewRecorder()
	dadosGrupo := fmt.Sprintf(`
	{
		"CodigoGrupo": "%s"}`, CodigoGrupoPar)

	fmt.Println("Valor dos dadosGrupo:", dadosGrupo)
	req, _ := http.NewRequest("POST", "/acl/ReverteSoftDeleteGrupo", strings.NewReader(dadosGrupo))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)
}

func TestGinReveteSoftDeleteGrupo(t *testing.T) {

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	dataAtual := time.Now()
	novoGrupoA := "ae" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoA, novoGrupoA, "31/12/2020", 0, "g", "Grupo Criado com Sucesso")

	dataAtual = time.Now()
	novoGrupoB := "af" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicaoNovoGrupo(t, novoGrupoB, novoGrupoB, "31/12/2020", 0, "g", "Grupo Criado com Sucesso")

	GinSoftDeleteGrupo(t, novoGrupoA, "[AGPINFSDL002 | aclGrupo.go|SoftDelete 02]  Grupo Removido (SoftDelete) com Sucesso")
	GinSoftDeleteGrupo(t, novoGrupoB, "[AGPINFSDL002 | aclGrupo.go|SoftDelete 02]  Grupo Removido (SoftDelete) com Sucesso")

	GinReverteSoftDeleteGrupo(t, novoGrupoA, "[AGPINFRSD002 | aclGrupo.go|ReverteSoftDeleteGrupo 02]  Grupo Ativado (reverteu SoftDelete) com Sucesso")
	GinReverteSoftDeleteGrupo(t, novoGrupoB, "[AGPINFRSD002 | aclGrupo.go|ReverteSoftDeleteGrupo 02]  Grupo Ativado (reverteu SoftDelete) com Sucesso")
	GinReverteSoftDeleteGrupo(t, "xxxx", "[MAGERRRSD001 | modelAclGrupos.go|ReverteSoftdeleteGrupo 01] Grupo não localizado")
}

//** Lista os dados dos gruposd

// GinFazRequisicao : zz
func GinListaTodosOsGrupos(t *testing.T, ComparacaoRetorno string) {
	var (
		ACLGrupoLocal []modelacl.ACLGrupo
	)

	r := ConfigRouterGrupo()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/acl/ListaTodosOsGrupos", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &ACLGrupoLocal)
	result, _ := strconv.Atoi(ComparacaoRetorno)
	assert.Greater(t, len(ACLGrupoLocal), result)
}

func TestGinListaTodosOsGrupo(t *testing.T) {

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	GinListaTodosOsGrupos(t, "0")
}
