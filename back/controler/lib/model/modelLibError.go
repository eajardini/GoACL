package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/jinzhu/gorm"
)

// LIBErroMSGSGBD : Estrutura de gerenciar os erros do SGBD
type LIBErroMSGSGBD struct {
	ErroID           uint64 `gorm:"type: bigint; primary_key; type: bigserial;column:erroid" json:"ErroID"`
	CodigoErro       int64  `gorm:"type:bigint; unique; not null; column:codigoerro" json:"CodigoErro"`
	MensagemErroPort string `gorm:"type:varchar(100); column:mensagemerroport" json:"MensagemErroPort"`
	MensagemErroIng  string `gorm:"type:varchar(100); column:mensagemerroing" json:"MensagemErroIng"`
	MensagemErroEsp  string `gorm:"type:varchar(100); column:mensagemerroesp" json:"MensagemErroEsp`
}

// LIBErroMSGSGBDJSON : Estrutura de gerenciar os erros do SGBD
type LIBErroMSGSGBDJSON struct {
	CodigoErro       int64  `json:"CodigoErro"`
	MensagemErroPort string `json:"MensagemErroPort"`
	MensagemErroIng  string `json:"MensagemErroIng"`
	MensagemErroEsp  string `json:"MensagemErroEsp"`
}

// LIBErroMSGRetorno : Estrutura para ser usada nos retornos de erros do sistema
type LIBErroMSGRetorno struct {
	CodigoErro int64
	Erro       error
}

var (
	LIBErroMSGSGBDMapGlobal map[int64]LIBErroMSGSGBD
)

//InsereErroNoSGBD : criar os erros do sistema
func InsereErroNoSGBD(codigoErroPar int64, MensagemErroPortPar string, MensagemErroIngPar string,
	MensagemErroEspPar string, bdPar bancoDeDados.BDCon) LIBErroMSGRetorno {
	var (
		LIBErroMSGSGBDLocal    LIBErroMSGSGBD
		erroRetorno            LIBErroMSGRetorno
		qtdadeRegistrosAchados int
		qtdadeRegistroCriados  *gorm.DB
		modulo                 string
	)

	bdPar.BD.Table("lib_erro_msgsgbd").Where("codigoerro = ?", codigoErroPar).Count(&qtdadeRegistrosAchados)
	if qtdadeRegistrosAchados > 0 {
		modulo = "[modelLibError.go|InsereErro|ERRO 001] "
		LIBErroMSGSGBDLocal.CodigoErro = 92
		LIBErroMSGSGBDLocal.MensagemErroPort = modulo + "92 - Mensagem de Erro já cadastrado"
		erroRetorno.CodigoErro = LIBErroMSGSGBDLocal.CodigoErro
		erroRetorno.Erro = errors.New(LIBErroMSGSGBDLocal.MensagemErroPort)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()

	LIBErroMSGSGBDLocal.CodigoErro = codigoErroPar
	LIBErroMSGSGBDLocal.MensagemErroPort = MensagemErroPortPar
	LIBErroMSGSGBDLocal.MensagemErroIng = MensagemErroIngPar
	LIBErroMSGSGBDLocal.MensagemErroEsp = MensagemErroEspPar

	qtdadeRegistroCriados = tx.Create(&LIBErroMSGSGBDLocal)

	if qtdadeRegistroCriados.RowsAffected == 0 {
		modulo = "[modelLibError.go|InsereErro|ERRO 002] "
		LIBErroMSGSGBDLocal.CodigoErro = 93
		LIBErroMSGSGBDLocal.MensagemErroPort = modulo + "93 - Erro ao inserir uma nova mensagem"
		erroRetorno.CodigoErro = LIBErroMSGSGBDLocal.CodigoErro
		erroRetorno.Erro = errors.New(LIBErroMSGSGBDLocal.MensagemErroPort)
		return erroRetorno
	}

	tx.Commit()

	return erroRetorno
}

//ExecutaCreateStruct : Função desenvolvida para criar (inserir) os dados da estrutura no BD
//Foi criada por que quando se tenta reutilizar uma variável que é Model o Gorm não consegue
// Criar uma nova Chave Primária. Assim, aqui só recebe os dados, atribui a variável Model
// e executa o Create
func ExecutaCreateStruct(LIBErroMSGSGBDPar LIBErroMSGSGBD, txPar *gorm.DB) LIBErroMSGRetorno {
	var (
		LIBErroMSGRetornoLocal LIBErroMSGRetorno
		LIBErroMSGSGBDLocal    LIBErroMSGSGBD
		modulo                 string
	)

	LIBErroMSGSGBDLocal = LIBErroMSGSGBDPar

	qtdadeRegistrosCriados := txPar.Create(&LIBErroMSGSGBDLocal)

	if qtdadeRegistrosCriados.RowsAffected == 0 {
		modulo = "[modelLibError.go|ExecutaCreateStruct|ERRO 001] "
		LIBErroMSGRetornoLocal.CodigoErro = 94
		LIBErroMSGRetornoLocal.Erro = errors.New(modulo + "94 - Erro ao inserir uma nova mensagem.")
		return LIBErroMSGRetornoLocal
	}
	return LIBErroMSGRetornoLocal
}

// CarregaMensagemDoJSON :
func CarregaMensagemDoJSON(caminhoArquivoJSONPar string, bdPar bancoDeDados.BDCon) LIBErroMSGRetorno {
	var (
		LIBErroMSGRetornoLocal  LIBErroMSGRetorno
		LIBErroMSGSGBDLocais    []LIBErroMSGSGBD
		mensagensErrosDoArquivo []byte
		modulo                  string
		qtdadeRegistrosAchados  int
		// qtdadeRegistroCriados   *gorm.DB
	)

	mensagensErrosDoArquivo, LIBErroMSGRetornoLocal.Erro = ioutil.ReadFile(caminhoArquivoJSONPar)
	if LIBErroMSGRetornoLocal.Erro != nil {
		modulo = "[modelLibError.go|CarregaMensagemDoJSON|ERRO 001] "
		LIBErroMSGRetornoLocal.CodigoErro = 1
		LIBErroMSGRetornoLocal.Erro = errors.New(modulo + "1 - Erro ao ler arquivo" + ": " + caminhoArquivoJSONPar)
		return LIBErroMSGRetornoLocal
	}

	LIBErroMSGRetornoLocal.Erro = json.Unmarshal([]byte(mensagensErrosDoArquivo), &LIBErroMSGSGBDLocais)

	if LIBErroMSGRetornoLocal.Erro != nil {
		modulo = "[modelLibError.go|CarregaMensagemDoJSON|ERRO 002] "
		LIBErroMSGRetornoLocal.CodigoErro = 2
		LIBErroMSGRetornoLocal.Erro = errors.New(modulo + "2 - Erro de realizar Unmarshal em estrutura.")
		return LIBErroMSGRetornoLocal
	}

	tx := bdPar.BD.Begin()
	for _, reg := range LIBErroMSGSGBDLocais {
		bdPar.BD.Table("lib_erro_msgsgbd").Where("codigoerro = ?", reg.CodigoErro).Count(&qtdadeRegistrosAchados)
		if qtdadeRegistrosAchados == 0 {
			LIBErroMSGRetornoLocal = ExecutaCreateStruct(reg, tx)
		}
	}
	tx.Commit()
	return LIBErroMSGRetornoLocal
}

//CarregaTodosAsMensagensDeErro :
func CarregaTodosAsMensagensDeErro(bdPar bancoDeDados.BDCon) LIBErroMSGRetorno {
	var (
		LIBErroMSGSGBDLocal    LIBErroMSGSGBD
		LIBErroMSGSGBDLocais   []LIBErroMSGSGBD
		LIBErroMSGSGBDMapLocal map[int64]LIBErroMSGSGBD
		erroRetorno            LIBErroMSGRetorno
		qtdadeRegistrosAchados *gorm.DB
		modulo                 string
	)

	qtdadeRegistrosAchados = bdPar.BD.Find(&LIBErroMSGSGBDLocais)
	if qtdadeRegistrosAchados.RowsAffected == 0 {
		modulo = "[modelLibError.go|CarregaTodosAsMensagensDeErro|ERRO 001] "
		LIBErroMSGSGBDLocal.CodigoErro = 91
		LIBErroMSGSGBDLocal.MensagemErroPort = modulo + "91 - Nenhum registro de Mensagem de Erro não encontrado."
		erroRetorno.CodigoErro = LIBErroMSGSGBDLocal.CodigoErro
		erroRetorno.Erro = errors.New(LIBErroMSGSGBDLocal.MensagemErroPort)
	}

	LIBErroMSGSGBDMapLocal = make(map[int64]LIBErroMSGSGBD)
	for _, reg := range LIBErroMSGSGBDLocais {
		LIBErroMSGSGBDMapLocal[reg.CodigoErro] = reg
	}
	LIBErroMSGSGBDMapGlobal = LIBErroMSGSGBDMapLocal
	return erroRetorno
}

//*** Lista de erros definidos ***
/*
0 - Erro indefinido

--Usuarios
21 - Usuário não cadastrado


-- Menu
41 - Item do Menu não cadastrado

-- Erro
91 - Nenhum registro de Mensagem de Erro não encontrado.
92 - Mensagem de Erro já cadastrado
93 - Erro ao inserir uma nova mensagem
*/
