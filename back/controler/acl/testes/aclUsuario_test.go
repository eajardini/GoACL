package acltestes

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	modelacl "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/gin-gonic/gin"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
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
var (
	bd       bancoDeDados.BDCon
	msgErros mensagensErros.LIBErroMSGRetorno
	// msgErro     string
	//erroRetorno error
)

// var ACLUserTest modelacl.ACLUsuario

//Configura Gin com as rotas para teste. A função abaixo é uma
// cópia do arquivo routeACL.go

func ConfigRouter() *gin.Engine {
	r := gin.Default()
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.POST("/NovoUsuario", aclcontroler.NovoUsuario)
		acl.POST("/DesativaUsuario", aclcontroler.DesativaUsuario)
		acl.POST("/RemoveFisicamenteUsuario", aclcontroler.RemoveFisicamenteUsuario)
		acl.POST("/AtivaUsuario", aclcontroler.AtivaUsuario)
		acl.POST("/AlteraUsuario", aclcontroler.AlteraUsuario)
		acl.GET("/BuscaUsuarioPorLogin", aclcontroler.BuscaUsuarioPorLogin)
		acl.GET("/BuscaTodosUsuario", aclcontroler.BuscaTodosUsuario)
		acl.GET("/BuscaTodosUsuariosAtivos", aclcontroler.BuscaTodosUsuariosAtivos)

	}
	return r
}

//**Cria novo usuário
func TestModelCriaNovoUsuario(t *testing.T) {
	now.TimeFormats = append(now.TimeFormats, "02/01/2006")
	assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	//Garantindo que o sistema vai barrar se um usuário já existir
	dataAtual := time.Now()
	novoLoginA := "a" + dataAtual.Format("02/01/200615:04:05")
	ACLUserTest := modelacl.ACLUsuario{Login: novoLoginA, Password: "teste", Datacriacao: dataAtual,
		Datavalidade: dataAtual, Userbloqueado: 1, Userativo: 1}

	//Aqui ele tem que criar
	err := modelacl.CriaNovoUsuario(ACLUserTest, bd)
	assert.Equal(err, nil, "[MAUCNU001 | modelAclUsuario.go|CriaNovousuario N.01] Erro ao criar no usuário")

	//Aqui ele não pode criar
	err = modelacl.CriaNovoUsuario(ACLUserTest, bd)
	assert.NotEqual(err, nil, "[modelAclUsuario.go|CriaNovousuario N.03] Erro de insert: Usuário já existe)")

	//Garantino que se tiver usuário novo, o sistema não vai barrar
	novoLogin := "n" + dataAtual.Format("02/01/200615:04:05")
	fmt.Println("[aclUsuario_test|TestModelCriaNovoUsuario01] Valor do NOVO Login:", novoLogin)

	ACLUserTest = modelacl.ACLUsuario{Login: novoLogin, Password: "teste", Datacriacao: dataAtual,
		Datavalidade: dataAtual, Userbloqueado: 1, Userativo: 1}
	err = modelacl.CriaNovoUsuario(ACLUserTest, bd)
	assert.Equal(err, nil, err)

}

// VerificaSeLoginJaExisteNoBD : zz
func TestVerificaSeLoginJaExisteNoBD(t *testing.T) {
	assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	dataAtual := time.Now()
	novoLoginA := "aa" + dataAtual.Format("02/01/200615:04:05")
	ACLUserTest := modelacl.ACLUsuario{Login: novoLoginA, Password: "teste", Datacriacao: dataAtual,
		Datavalidade: dataAtual, Userbloqueado: 1, Userativo: 1}
	err := modelacl.CriaNovoUsuario(ACLUserTest, bd)
	assert.Equal(err, nil, "[MAUCNU001 | modelAclUsuario.go|CriaNovousuario N.01] Erro ao criar no usuário")

	//Verificando com um usuário que existe no banco de dados
	achou := modelacl.VerificaSeLoginJaExisteNoBD(novoLoginA, bd)
	assert.GreaterOrEqual(achou, 1, "[aclUsuario_test|TestVerificaSeLoginJaExisteNoBD 01 ]Login encontrado no BD")

	// achou = modelacl.VerificaSeLoginJaExisteNoBD("joaoxx", bd)
	// assert.Equal(achou, 0, "[aclUsuario_test|TestVerificaSeLoginJaExisteNoBD 02 ]Login encontrado no BD")

}

// GinFazRequisicao : zz
func GinFazRequisicao(t *testing.T, login string, password string, datacriacao string,
	datavalidade string, userbloqueado int, userativo int, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s", "password":"%s","datacriacao": "%s", "datavalidade": "%s", "userbloqueado": %d, "userativo": %d
	}`, login, password, datacriacao, datavalidade, userbloqueado, userativo)

	// fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("POST", "/acl/NovoUsuario", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)
}

func TestGinNovoUsuario(t *testing.T) {

	now.TimeFormats = append(now.TimeFormats, "02/01/2006")
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	// bd.AbreConexao()
	// defer bd.FechaConexao()

	dataAtual := time.Now()
	novoLogin := "ab" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLogin, "321", "31/12/2020", "01/01/0001", 0, 1, "Usuário Criado com Sucesso")

	dataAtual = time.Now()
	novoLogin = "ba" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLogin, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")
}

//**Remove Usuário
func TestRemoveUsuario(t *testing.T) {
	assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	dataAtual := time.Now()
	novoLoginA := "c" + dataAtual.Format("02/01/200615:04:05")
	ACLUserTest := modelacl.ACLUsuario{Login: novoLoginA, Password: "teste", Datacriacao: dataAtual,
		Datavalidade: dataAtual, Userbloqueado: 1, Userativo: 1}
	err := modelacl.CriaNovoUsuario(ACLUserTest, bd)
	assert.Equal(err, nil, "[MAUCNU001 | modelAclUsuario.go|CriaNovousuario N.01] Erro ao criar no usuário")

	// dataAtual = time.Now()
	// novoLoginB := "b" + dataAtual.Format("02/01/200615:04:05")
	// GinFazRequisicao(t, novoLoginB, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")

	err = modelacl.RemoveUsuarioPorLogin(novoLoginA, bd)
	assert.Equal(err, nil, err)

}

// GinFazRequisicaoParaRemocao : zz
func GinFazRequisicaoParaRemocao(t *testing.T, login string, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s"
	}`, login)

	fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("POST", "/acl/DesativaUsuario", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)

}

func TestGinRemoveUsuario(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	//Cria dois usuários para depois removê-los
	dataAtual := time.Now()
	novoLoginA := "d" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginA, "321", "31/12/2020", "01/01/0001", 0, 1, "Usuário Criado com Sucesso")
	dataAtual = time.Now()
	novoLoginB := "e" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginB, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")

	GinFazRequisicaoParaRemocao(t, novoLoginA, "[AUSINFRUR002 | aclUsuario|RemoveUsuario 02]  Usuário Removido com Sucesso")
	GinFazRequisicaoParaRemocao(t, novoLoginB, "[AUSINFRUR002 | aclUsuario|RemoveUsuario 02]  Usuário Removido com Sucesso")
	GinFazRequisicaoParaRemocao(t, "xxxx", "[MAUINFRPL001 | modelAclUsuario.go|RemoveUsuarioPorLogin 01] Usuário não localizado")

}

// GinFazRequisicaoParaRemocaoFisica : zz
func GinFazRequisicaoParaRemocaoFisica(t *testing.T, login string, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s"
	}`, login)

	fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("POST", "/acl/RemoveFisicamenteUsuario", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)

}

// TestGinRemoveFisicamenteUsuario :
func TestGinRemoveFisicamenteUsuario(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	//Cria dois usuários para depois removê-los
	dataAtual := time.Now()
	novoLoginA := "f" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginA, "321", "31/12/2020", "01/01/0001", 0, 1, "Usuário Criado com Sucesso")
	dataAtual = time.Now()
	novoLoginB := "g" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginB, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")

	//Remove os usuários
	GinFazRequisicaoParaRemocaoFisica(t, novoLoginA, "[AUSINFRFU002 | aclUsuario|RemoveFisicamenteUsuario 02]  Usuário Removido Fisicamente com Sucesso")
	GinFazRequisicaoParaRemocaoFisica(t, novoLoginB, "[AUSINFRFU002 | aclUsuario|RemoveFisicamenteUsuario 02]  Usuário Removido Fisicamente com Sucesso")
	GinFazRequisicaoParaRemocaoFisica(t, "xxxx", "[MAUINFRFL001 | modelAclUsuario.go|RemoveFisicamenteUsuarioPorLogin 01] Usuário não localizado")

}

//Ativa Usuário
// GinFazRequisicaoParaAtivacaoUsuario : zz
func GinFazRequisicaoParaAtivacaoUsuario(t *testing.T, login string, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s"
	}`, login)

	req, _ := http.NewRequest("POST", "/acl/AtivaUsuario", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)

}

func TestGinAtivaUsuario(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	//Cria dois usuários para depois removê-los Logicamente e depois ativalos
	dataAtual := time.Now()
	novoLoginA := "h" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginA, "321", "31/12/2020", "01/01/0001", 0, 1, "Usuário Criado com Sucesso")
	dataAtual = time.Now()
	novoLoginB := "i" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginB, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")

	GinFazRequisicaoParaRemocao(t, novoLoginA, "[AUSINFRUR002 | aclUsuario|RemoveUsuario 02]  Usuário Removido com Sucesso")
	GinFazRequisicaoParaRemocao(t, novoLoginB, "[AUSINFRUR002 | aclUsuario|RemoveUsuario 02]  Usuário Removido com Sucesso")
	GinFazRequisicaoParaRemocao(t, "xxxx", "[MAUINFRPL001 | modelAclUsuario.go|RemoveUsuarioPorLogin 01] Usuário não localizado")

	//Ativa os usuários
	GinFazRequisicaoParaAtivacaoUsuario(t, novoLoginA, "[AUSINFATU002 | aclUsuario|AtivaUsuario 02]  Usuário Ativado com Sucesso")
	GinFazRequisicaoParaAtivacaoUsuario(t, novoLoginB, "[AUSINFATU002 | aclUsuario|AtivaUsuario 02]  Usuário Ativado com Sucesso")
	GinFazRequisicaoParaAtivacaoUsuario(t, "xxxx", "[MAUINFAPL001 | modelAclUsuario.go|AtivaUsuarioPorLogin 01] Usuário não localizado")

}

// GinFazRequisicaoParaAlterarUsuario : zz
func GinFazRequisicaoParaAlterarUsuario(t *testing.T, login string, password string, datacriacao string,
	datavalidade string, userbloqueado int, userativo int, ComparacaoRetorno string) {
	var Resp aclcontroler.Resposta

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s", "password":"%s","datacriacao": "%s", "datavalidade": "%s", "userbloqueado": %d, "userativo": %d
	}`, login, password, datacriacao, datavalidade, userbloqueado, userativo)

	// fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("POST", "/acl/AlteraUsuario", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &Resp)
	assert.Equal(t, ComparacaoRetorno, Resp.Mensagem)
}

func TestGinAlteraUsuario(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	//Cria dois usuários para depois removê-los Logicamente e depois ativalos
	dataAtual := time.Now()
	novoLoginA := "a1" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginA, "321", "31/12/2020", "01/01/0001", 0, 1, "Usuário Criado com Sucesso")
	dataAtual = time.Now()
	novoLoginB := "a2" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginB, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")

	// //Ativa os usuários
	GinFazRequisicaoParaAlterarUsuario(t, novoLoginA, "123", "01/05/2015", "10/09/2019", 1, 0, "[AUSINFAUS002 | aclUsuario|AlteraUsuario N.02]Usuário alterado com sucesso")
	GinFazRequisicaoParaAlterarUsuario(t, novoLoginB, "456", "01/05/2015", "10/09/2019", 1, 0, "[AUSINFAUS002 | aclUsuario|AlteraUsuario N.02]Usuário alterado com sucesso")
	GinFazRequisicaoParaAlterarUsuario(t, "xxxx", "123", "01/05/2015", "10/09/2019", 1, 0, "[MAUINFAUS003 |  modelAclUsuario.go|AlteraUsuario N.03] Usuário não existe: "+"xxxx")

}

// GinFazRequisicaoParaBuscarTodosOsUsuario : zz
func GinFazRequisicaoParaBuscarTodosOsUsuario(t *testing.T, ComparacaoRetorno string) {
	var UserACLRetorno []modelacl.ACLUsuario

	r := ConfigRouter()
	w := httptest.NewRecorder()

	// fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("GET", "/acl/BuscaTodosUsuario", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &UserACLRetorno)
	result, _ := strconv.Atoi(ComparacaoRetorno)
	assert.Greater(t, len(UserACLRetorno), result)
}

func TestGinBuscaTodosOsUsuario(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	GinFazRequisicaoParaBuscarTodosOsUsuario(t, "0")
}

// GinFazRequisicaoParaBuscarTodosOsUsuarioAtivos : zz
func GinFazRequisicaoParaBuscarTodosOsUsuariosAtivos(t *testing.T, ComparacaoRetorno string) {
	var UserACLRetorno []modelacl.ACLUsuario

	r := ConfigRouter()
	w := httptest.NewRecorder()

	// fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("GET", "/acl/BuscaTodosUsuariosAtivos", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &UserACLRetorno)
	result, _ := strconv.Atoi(ComparacaoRetorno)
	assert.Greater(t, len(UserACLRetorno), result)
}

func TestGinBuscaTodosOsUsuariosAtivos(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	GinFazRequisicaoParaBuscarTodosOsUsuariosAtivos(t, "0")
}

// GinFazRequisicaoParaBuscarUsuariosPorLogin : zz
func GinFazRequisicaoParaBuscarUsuariosPorLogin(t *testing.T, loginPAR string) {
	var UserACLRetorno modelacl.ACLUsuario

	r := ConfigRouter()
	w := httptest.NewRecorder()
	dadosUsuario := fmt.Sprintf(`
	{
		"login": "%s"}`, loginPAR)
	// fmt.Println("Valor dos dadosUsuario:", dadosUsuario)
	req, _ := http.NewRequest("GET", "/acl/BuscaUsuarioPorLogin", strings.NewReader(dadosUsuario))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal([]byte(w.Body.String()), &UserACLRetorno)
	// result, _ := strconv.Atoi(ComparacaoRetorno)
	// assert.Greater(t, UserACLRetorno.Login , result)
	assert.Equal(t, loginPAR, UserACLRetorno.Login)
}

func TestGinBuscaUsuarioPorLogin(t *testing.T) {
	// assert := assert.New(t)
	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	//Cria dois usuários para depois removê-los Logicamente e depois ativalos
	dataAtual := time.Now()
	novoLoginA := "pl" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginA, "321", "31/12/2020", "01/01/0001", 0, 1, "Usuário Criado com Sucesso")
	dataAtual = time.Now()
	novoLoginB := "pm" + dataAtual.Format("02/01/200615:04:05")
	GinFazRequisicao(t, novoLoginB, "321", "31/12/2020", "", 0, 1, "Usuário Criado com Sucesso")

	GinFazRequisicaoParaBuscarUsuariosPorLogin(t, novoLoginA)
	GinFazRequisicaoParaBuscarUsuariosPorLogin(t, novoLoginB)
	GinFazRequisicaoParaBuscarUsuariosPorLogin(t, "Usuário não encontrado")

}
