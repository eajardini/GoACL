package modelacl

import (
	"fmt"
	"log"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	"github.com/jinzhu/gorm"
)

//ACLGrupoAcessaMenu : zz
type ACLGrupoAcessaMenu struct {
	GrupoAcessaMenuID uint64 `gorm:"primary_key; type: bigserial ;column:grupoacessamenuid"`
	GrupoID           uint64 `gorm:"type:bigint; column:grupoid;  unique_index:idx_grupo_acessa_menu" validate:"required"`
	MenuID            uint64 `gorm:"type:bigint; column:menuid; unique_index:idx_grupo_acessa_menu" validate:"required"`
}

//ACLGrupoAcessaMenuJSON : zz
type ACLGrupoAcessaMenuJSON struct {
	CodigoGrupo string   `json:"CodigoGrupo"`
	CodigoMenu  []string `json:"CodigoMenu"`
}

//ExecutaCreateStructGrupoAcessaMenu : Função desenvolvida para criar (inserir) os dados da estrutura no BD
//Foi criada por que quando se tenta reutilizar uma variável que é Model o Gorm não consegue
// Criar uma nova Chave Primária. Assim, aqui só recebe os dados, atribui a variável Model
// e executa o Create
func ExecutaCreateStructGrupoAcessaMenu(idMenuPar uint64, idGrupoPar uint64, txPar *gorm.DB) ACLGrupoAcessaMenu {
	var (
		ACLGrupoAcessaMenuLocal ACLGrupoAcessaMenu
	)

	ACLGrupoAcessaMenuLocal.MenuID = idMenuPar
	ACLGrupoAcessaMenuLocal.GrupoID = idGrupoPar

	qtdadeRegistrosCriados := txPar.Create(&ACLGrupoAcessaMenuLocal)

	if qtdadeRegistrosCriados.RowsAffected == 0 {
		msgErro = "[modelAclMenu.go|InsereNoBDMenuDoJSON| ERRO01] Erro ao inserir menu nível 1" + string(ACLGrupoAcessaMenuLocal.GrupoID)
		log.Println(msgErro)
		// erroRetorno = errors.New(msgErro)
	}
	return ACLGrupoAcessaMenuLocal
}

// ConcedeAcessoAoMenu : Concede a um grupo o acesso aos menu do sistema.
func ConcedeAcessoAoMenu(codigosMenuPar []string, codigoGrupoPar string, bdPar bancoDeDados.BDCon) mensagensErros.LIBErroMSGRetorno {
	var (
		erroRetornoLocal            mensagensErros.LIBErroMSGRetorno
		msgErroLocal                mensagensErros.LIBErroMSGSGBD
		moduloLocal                 string
		qtdadeRegistrosAchadosLocal int
		ACLGrupoAcessaMenuLocal     ACLGrupoAcessaMenu

		// numerosDeRegistrosCriados *gorm.DB
	)
	fmt.Println("[modelAclgrupoAcessaMenu] Valor de codigosMenuPar", codigosMenuPar)
	fmt.Println("[modelAclgrupoAcessaMenu] Valor de codigoGrupoPar", codigoGrupoPar)

	erroRetornoLocal.Erro = bdPar.BD.Table("acl_grupo").Select("grupoid").Where("codigogrupo = ?", codigoGrupoPar).Row().Scan(&ACLGrupoAcessaMenuLocal.GrupoID)

	if erroRetornoLocal.Erro != nil {
		moduloLocal = "[modelAclGrupoAcessaMenu|ConcedeAcessoAoMenu|ERRO01] "
		log.Println(moduloLocal + erroRetornoLocal.Erro.Error())
		erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(61, moduloLocal)
		return erroRetornoLocal
	}
	tx := bdPar.BD.Begin()
	//Remove todos os direitos antes de atribui-los novamente
	tx.Where("grupoid = ?", ACLGrupoAcessaMenuLocal.GrupoID).Delete(ACLGrupoAcessaMenu{})
	tx.Commit()

	tx = bdPar.BD.Begin()
	for _, codigoMenu := range codigosMenuPar {
		erroRetornoLocal.Erro = bdPar.BD.Table("acl_menu").Select("menuid").Where("codigo = ?", codigoMenu).Row().Scan(&ACLGrupoAcessaMenuLocal.MenuID)
		if erroRetornoLocal.Erro != nil {
			moduloLocal = "[modelAclGrupoAcessaMenu|ConcedeAcessoAoMenu|ERRO02] "
			log.Println(moduloLocal + erroRetornoLocal.Erro.Error())
			erroRetornoLocal = msgErroLocal.BuscaMensagemPeloCodigo(42, moduloLocal)
		} else {
			bdPar.BD.Table("acl_grupo_acessa_menu").Where("menuid = ? and grupoid = ?", ACLGrupoAcessaMenuLocal.MenuID, ACLGrupoAcessaMenuLocal.GrupoID).Count(&qtdadeRegistrosAchadosLocal)
			if qtdadeRegistrosAchadosLocal == 0 {
				ExecutaCreateStructGrupoAcessaMenu(ACLGrupoAcessaMenuLocal.MenuID, ACLGrupoAcessaMenuLocal.GrupoID, tx)
			}
		}

	}

	tx.Commit()
	return erroRetornoLocal
}
