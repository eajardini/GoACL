package libtestes

import (
	"testing"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	modelLib "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
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
	bd     bancoDeDados.BDCon
	libMSG modelLib.LIBErroMSGSGBD

// 	msgErro     string
// 	erroRetorno error
)

// var ACLUserTest modelacl.ACLUsuario

//Configura Gin com as rotas para teste. A função abaixo é uma
// cópia do arquivo routeACL.go

// func ConfigRouterGrupo() *gin.Engine {
// 	r := gin.Default()
// 	lib := r.Group("/acl")
// {
//	lib.POST("/InsereErroNoSGBD", modelLIB.InsereErroNoSGBD)
// lib.POST("/ReverteSoftDeleteGrupo", aclcontroler.ReverteSoftDeleteGrupo)
// lib.DELETE("/SoftDeleteGrupo", aclcontroler.SoftDeleteGrupo)
// lib.GET("/ListaTodosOsGrupos", aclcontroler.ListaTodosOsGrupos)
// 	}
// 	return r
// }

func TestModelInsereErroNoSGBD(t *testing.T) {
	var (
		erroRetorno modelLib.LIBErroMSGRetorno
	)

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	erroRetorno = libMSG.InsereErroNoSGBD(0, "0 - Erro indefinido", "msg ingles", "msg espanhol", bd)

	if erroRetorno.CodigoErro != 92 {
		assert.Equal(t, nil, erroRetorno.Erro)
	}

}

func TestCarregaTodosAsMensagensDeErro(t *testing.T) {
	var (
		erroRetorno modelLib.LIBErroMSGRetorno
	)

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	erroRetorno = libMSG.CarregaTodosAsMensagensDeErro(bd)
	// fmt.Println("[libError_test.go|TestCarregaTodosAsMensagensDeErro 001] Valor LIBErroMSGSGBDMapGlobal:", modelLib.libErroMSGSGBDMapGlobal[0].MensagemErroPort)

	assert.Equal(t, nil, erroRetorno.Erro)
}

func TestBuscaMensagemPeloCodigo(t *testing.T) {
	var (
		erroRetorno modelLib.LIBErroMSGRetorno
	)

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	erroRetorno = libMSG.CarregaTodosAsMensagensDeErro(bd)
	// fmt.Println("[libError_test.go|TestCarregaTodosAsMensagensDeErro 001] Valor LIBErroMSGSGBDMapGlobal:", modelLib.libErroMSGSGBDMapGlobal[0].MensagemErroPort)
	modulo := "[libError_test.go|TestBuscaMensagemPeloCodigo|ERRO 001] "
	erroRetorno = libMSG.BuscaMensagemPeloCodigo(1, modulo)
	assert.Equal(t, modulo+"1 - Erro ao ler arquivo.", erroRetorno.Erro.Error())

	erroRetorno = libMSG.BuscaMensagemPeloCodigo(41, modulo)
	assert.Equal(t, modulo+"41 - Item do Menu não cadastrado.", erroRetorno.Erro.Error())

	erroRetorno = libMSG.BuscaMensagemPeloCodigo(92, modulo)
	assert.Equal(t, modulo+"92 - Código de Mensagem de Erro já cadastrado.", erroRetorno.Erro.Error())

	erroRetorno = libMSG.BuscaMensagemPeloCodigo(1192, modulo)
	assert.Equal(t, modulo+"0 - Erro indefinido.", erroRetorno.Erro.Error())
}
