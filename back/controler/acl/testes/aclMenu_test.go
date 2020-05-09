package acltestes

import (
	"fmt"
	"testing"

	aclControlerModel "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"

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
// bd bancoDeDados.BDCon

// 	msgErro     string
// 	erroRetorno error
// )

// var ACLUserTest modelacl.ACLUsuario

func TestGinMenuFromArquivoJSON(t *testing.T) {
	var (
		ACLMenuFromJSONLocal aclControlerModel.ACLMenuFromJSON
		erroRetorno          error
	)

	bd.ConfiguraStringDeConexao("../../../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
	bd.AbreConexao()
	defer bd.FechaConexao()

	caminho := "../../../config/menu.json"
	ACLMenuFromJSONLocal, erroRetorno = aclControlerModel.CarregaMenuDoJSON(caminho)

	for i := 0; i < len(ACLMenuFromJSONLocal); i++ {
		fmt.Println("Código Menu: ", ACLMenuFromJSONLocal[i].Items[0].Items[0].Label)
		fmt.Println("Label Menu: ", ACLMenuFromJSONLocal[i].Label)
	}

	assert.Equal(t, erroRetorno, nil)

	erroRetorno = aclControlerModel.InsereNoBDMenuDoJSON(ACLMenuFromJSONLocal, bd)

	assert.Equal(t, erroRetorno, nil)
}
