package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	acl "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	modelACL "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
	modelLIB "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/jinzhu/now"
)

var (
	bd                      bancoDeDados.BDCon
	sysConfig               bancoDeDados.SysConfig
	caminhoDoArquivoToml    string
	erroLocal               error
	ACLMenuFromJSONLocal    modelACL.ACLMenuFromJSON
	ACLUsuariolocal         modelACL.ACLUsuario
	senha                   string
	dataNow                 *now.Now
	LIBErroMSGRetornoGlobal modelLIB.LIBErroMSGRetorno
	libMSG                  modelLIB.LIBErroMSGSGBD
	modulo                  string
	userACLJSONGlobal       modelACL.ACLUsuarioJSON
)

//ConfiguraMensagensDeErro : zz
func ConfiguraMensagensDeErro(bdPar bancoDeDados.BDCon) modelLIB.LIBErroMSGRetorno {
	var (
		LIBErroMSGRetornoLocal modelLIB.LIBErroMSGRetorno
		modulo                 string
	)

	LIBErroMSGRetornoLocal = libMSG.InsereErroNoSGBD(0, "0 - Erro indefinido",
		"Mensagem em Inglês", "Mensagem em Espanhol", bdPar)

	if LIBErroMSGRetornoLocal.CodigoErro != 92 {
		if LIBErroMSGRetornoLocal.Erro != nil {
			modulo = "[migrate.go|ConfiguraMensagensDeErro|ERRO 001] "
			log.Println(modulo, LIBErroMSGRetornoLocal.Erro)
		}
	}

	return LIBErroMSGRetornoLocal
}

//CarregaMensagensErroDoArquivoJSON : zz
func CarregaMensagensErroDoArquivoJSON(caminhoArquivoJSONPar string, bdPar bancoDeDados.BDCon) modelLIB.LIBErroMSGRetorno {
	var (
		LIBErroMSGRetornoLocal modelLIB.LIBErroMSGRetorno
		//modulo                 string
	)
	LIBErroMSGRetornoLocal = libMSG.CarregaMensagemDoJSON(caminhoArquivoJSONPar, bdPar)
	return LIBErroMSGRetornoLocal
}

func main() {
	caminhoDoArquivoToml = "../config/ConfigBancoDados.toml"
	bd.ConfiguraStringDeConexao(caminhoDoArquivoToml)
	bd.IniciaConexao()

	bd.AbreConexao()
	defer bd.FechaConexao()

	//Realizando a migração das tabelas dos models
	bd.BD.SingularTable(true)
	bd.BD.AutoMigrate(&modelACL.ACLUsuario{}, &modelACL.ACLGrupo{}, &modelACL.ACLUsuarioGrupo{})
	bd.BD.AutoMigrate(&modelACL.ACLMenu{}, &modelLIB.LIBErroMSGSGBD{}, &modelACL.ACLGrupoAcessaMenu{})
	// Adicionando Foreign Key
	bd.BD.Model(&modelACL.ACLUsuarioGrupo{}).AddForeignKey("usuarioid", "acl_usuario", "RESTRICT", "RESTRICT")
	bd.BD.Model(&modelACL.ACLUsuarioGrupo{}).AddForeignKey("grupoid", "acl_grupo", "RESTRICT", "RESTRICT")
	bd.BD.Model(&modelACL.ACLGrupoAcessaMenu{}).AddForeignKey("grupoid", "acl_grupo", "RESTRICT", "RESTRICT")
	bd.BD.Model(&modelACL.ACLGrupoAcessaMenu{}).AddForeignKey("menuid", "acl_menu", "RESTRICT", "RESTRICT")

	log.Println("[migrate.go|main|INFO 001] Tabelas criado com sucesso!")

	//*****Insere as mensagens de erro do sistema *****//
	log.Println("[migrate.go|main|INFO 002] Carrengando Mensagens para o Sistema!")

	LIBErroMSGRetornoGlobal = CarregaMensagensErroDoArquivoJSON(bd.RetornaCaminhoMSGErros(), bd)
	if LIBErroMSGRetornoGlobal.Erro != nil {
		modulo = "[migrate.go|main|ERRO 009] "
		log.Println(modulo, LIBErroMSGRetornoGlobal.Erro.Error())
	}

	log.Println("[migrate.go|main|INFO 003] Mensagens carregadas com sucesso!")

	LIBErroMSGRetornoGlobal = libMSG.CarregaTodosAsMensagensDeErro(bd)
	//*****************************//

	// Cria o super usuário internal de senha Intern@l

	h := md5.New()
	h.Write([]byte("Intern@l"))
	senha = hex.EncodeToString(h.Sum(nil))

	dataAtual := time.Now()
	now.TimeFormats = append(now.TimeFormats, "02/01/2006")
	datavalidade, _ := now.Parse("01/01/0001")

	userACLJSONGlobal.Login = "internal"
	userACLJSONGlobal.Password = senha
	userACLJSONGlobal.Datacriacao = dataAtual.Format("02/01/2006")
	userACLJSONGlobal.Datavalidade = datavalidade.String()
	userACLJSONGlobal.Userbloqueado = 0
	userACLJSONGlobal.Userativo = 1

	LIBErroMSGRetornoGlobal = acl.NovoUsuarioACL(userACLJSONGlobal, bd)
	if LIBErroMSGRetornoGlobal.Erro != nil {
		log.Println("[migrate.go|main|INFO 006]Erro ao cadastrar o usuário internal", LIBErroMSGRetornoGlobal.Erro.Error())
		//	panic(0)
	}
	log.Println("[migrate.go|main|INFO 006] Usuário criado com sucesso!")

	// Insere no BD os items de menu
	_, erroLocal = toml.DecodeFile(caminhoDoArquivoToml, &sysConfig)
	if erroLocal != nil {
		log.Println("[migrate.go|main|ERRO 002] Erro ao abrir arquiv TOML:", erroLocal.Error())
		panic(0)
	}
	ACLMenuFromJSONLocal, LIBErroMSGRetornoGlobal.Erro = modelACL.CarregaMenuDoJSON(sysConfig.Principal.CaminhoMenuJSON)
	if erroLocal != nil {
		log.Println("[migrate.go|main|ERRO 003] Erro carregar arquivo do menu:", LIBErroMSGRetornoGlobal.Erro.Error())
		panic(0)
	}

	LIBErroMSGRetornoGlobal.Erro = modelACL.InsereNoBDMenuDoJSON(ACLMenuFromJSONLocal, bd)
	if LIBErroMSGRetornoGlobal.Erro != nil {
		log.Println("[migrate.go|main|ERRO 005] Erro inserir items do menu:", LIBErroMSGRetornoGlobal.Erro.Error())
		panic(0)
	}

	log.Println("[migrate.go|main|INFO 004] Itens do menu inseridos com sucesso!")
}
