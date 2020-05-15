package modelacl

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/jinzhu/gorm"
)

//ACLMenuFromJSON : zz
type ACLMenuFromJSON []struct {
	Label     string `json:"label"`
	Codigo    string `json:"codigo"`
	Nivel     int    `json:"nivel"`
	Codigopai string `json:"codigopai"`
	Icon      string `json:"icon"`
	Items     []struct {
		Label     string `json:"label"`
		Codigo    string `json:"codigo"`
		Nivel     int    `json:"nivel"`
		Codigopai string `json:"codigopai"`
		Icon      string `json:"icon"`
		Items     []struct {
			Label     string `json:"label"`
			Codigo    string `json:"codigo"`
			Nivel     int    `json:"nivel"`
			Codigopai string `json:"codigopai"`
			Icon      string `json:"icon"`
			To        string `json:"to"`
		} `json:"items"`
	} `json:"items"`
}

// ACLMenu : Estrutura para armazenar cada item do menu no banco de dados.
type ACLMenu struct {
	MenuID    uint64         `gorm:"type: bigint; primary_key; type: bigserial;column:menuid"`
	Label     string         `gorm:"type:varchar(40); unique; not null; column:label"`
	Codigo    string         `gorm:"type:varchar(10); unique; not null; column:codigo"`
	Nivel     int            `gorm:"type:smallint; not null; column:nivel"`
	Codigopai sql.NullString `gorm:"type:varchar(10); column:codigopai"`
	Icon      sql.NullString `gorm:"type:varchar(20); column:icon"`
	To        sql.NullString `gorm:"type:varchar(20); column:to"`
}

//ACLMenuJSON : zz
type ACLMenuJSON struct {
	CodigoGrupo      string `json:"CodigoGrupo" validate:"required"`
	DescricaoGrupo   string `json:"DescricaoGrupo" validate:"required"`
	DataCriacaoGrupo string `json:"DataCriacaoGrupo" `
	SoftDelete       int    `json:"SoftDelete" validate:"number,gte=0,lte=1"`
	DataTeste        string `json:"DataTeste"`
}

//ItemsNivel1 : Eu aproveitei rotina que já havia feito anteriormente
type ItemsNivel1 struct {
	Label string        `json:"label"`
	Icone string        `json:"icone"`
	Items []ItemsNivel2 `json:"items"`
}

//ItemsNivel2 : zz
type ItemsNivel2 struct {
	Label string        `json:"label"`
	Icone string        `json:"icone"`
	Items []ItemsNivel3 `json:"items"`
}

//ItemsNivel3 : zz
type ItemsNivel3 struct {
	Label string `json:"label"`
	Icone string `json:"icone"`
	To    string `json:"to"`
}

//CarregaMenuDoJSON : carrega o menu de JSON que sta no arquivo ./config/menu.json
func CarregaMenuDoJSON(caminho string) (ACLMenuFromJSON, error) {
	var (
		ACLMenuFromJSONLocal ACLMenuFromJSON
		erroRetorno          error
		msgErro              string
	)

	menuJSON, erroLocal := ioutil.ReadFile(caminho)

	if erroLocal != nil {
		msgErro = "[MAMERRCMJ001 | modelAclMenu.go|CarregaMenuDoJSON N.01] Erro de carregar o menu de arquvio JSON"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return nil, erroRetorno
	}

	json.Unmarshal([]byte(menuJSON), &ACLMenuFromJSONLocal)
	return ACLMenuFromJSONLocal, erroLocal
}

//ExecutaCreateStruct : Função desenvolvida para criar (inserir) os dados da estrutura no BD
//Foi criada por que quando se tenta reutilizar uma variável que é Model o Gorm não consegue
// Criar uma nova Chave Primária. Assim, aqui só recebe os dados, atribui a variável Model
// e executa o Create
func ExecutaCreateStruct(ACLMenuLocalPar ACLMenu, txPar *gorm.DB) error {
	var (
		ACLMenuLocal ACLMenu
		erroRetorno  error
		msgErro      string
	)

	ACLMenuLocal = ACLMenuLocalPar

	qtdadeRegistrosCriados := txPar.Create(&ACLMenuLocal)

	if qtdadeRegistrosCriados.RowsAffected == 0 {
		msgErro = "[MAGERRIBM001 | modelAclMenu.go|InsereNoBDMenuDoJSON N.01] Erro ao inserir menu nível 1" + ACLMenuLocal.Label
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
	}
	return erroRetorno
}

//InsereNoBDMenuDoJSON : após o arquivo ./config/menu.json ter sido carregado, as opções do menu são inseridas no BD
func InsereNoBDMenuDoJSON(ACLMenuFromJSONPar ACLMenuFromJSON, bdPar bancoDeDados.BDCon) error {
	var (
		ACLMenuLocal           ACLMenu
		erroRetorno            error
		qtdadeRegistrosAchados int
	)
	tx := bdPar.BD.Begin()

	for _, regN1 := range ACLMenuFromJSONPar {
		//log.Println("Código Nivel 1:", regN1.Codigo, " CódigoPai Nível 1:", regN1.Codigopai, " Label Nível 1:", regN1.Label)
		bdPar.BD.Table("acl_menu").Where("codigo = ?", regN1.Codigo).Count(&qtdadeRegistrosAchados)
		if qtdadeRegistrosAchados == 0 {
			ACLMenuLocal.Label = regN1.Label
			ACLMenuLocal.Codigo = regN1.Codigo
			ACLMenuLocal.Nivel = regN1.Nivel
			ACLMenuLocal.Codigopai.String = ""
			ACLMenuLocal.Codigopai.Valid = true
			ACLMenuLocal.Icon.String = regN1.Icon
			ACLMenuLocal.Icon.Valid = true
			ACLMenuLocal.To.String = ""
			ACLMenuLocal.To.Valid = true
			erroRetorno = ExecutaCreateStruct(ACLMenuLocal, tx)
			if erroRetorno != nil {
				tx.Rollback()
				return erroRetorno
			}
		}
		for _, regN2 := range regN1.Items {
			qtdadeRegistrosAchados = 0
			// log.Println("Código Nivel 2:", regN2.Codigo, " CódigoPai Nível 2:", regN2.Codigopai, " Label Nível 2:", regN2.Label)
			bdPar.BD.Table("acl_menu").Where("codigo = ?", regN2.Codigo).Count(&qtdadeRegistrosAchados)
			if qtdadeRegistrosAchados == 0 {
				ACLMenuLocal.Label = regN2.Label
				ACLMenuLocal.Codigo = regN2.Codigo
				ACLMenuLocal.Nivel = regN2.Nivel
				ACLMenuLocal.Codigopai.String = regN2.Codigopai
				ACLMenuLocal.Codigopai.Valid = true
				ACLMenuLocal.Icon.String = regN2.Icon
				ACLMenuLocal.Icon.Valid = true
				ACLMenuLocal.To.String = ""
				ACLMenuLocal.To.Valid = true
				erroRetorno = ExecutaCreateStruct(ACLMenuLocal, tx)
				if erroRetorno != nil {
					tx.Rollback()
					return erroRetorno
				}
			}
			for _, regN3 := range regN2.Items {
				qtdadeRegistrosAchados = 0
				// log.Println("Código Nivel 3:", regN3.Codigo, " CódigoPai Nível 2:", regN3.Codigopai, " Label Nível 3:", regN2.Label)
				bdPar.BD.Table("acl_menu").Where("codigo = ?", regN3.Codigo).Count(&qtdadeRegistrosAchados)
				if qtdadeRegistrosAchados == 0 {
					if qtdadeRegistrosAchados == 0 {
						ACLMenuLocal.Label = regN3.Label
						ACLMenuLocal.Codigo = regN3.Codigo
						ACLMenuLocal.Nivel = regN3.Nivel
						ACLMenuLocal.Codigopai.String = regN3.Codigopai
						ACLMenuLocal.Codigopai.Valid = true
						ACLMenuLocal.Icon.String = regN3.Icon
						ACLMenuLocal.Icon.Valid = true
						ACLMenuLocal.To.String = regN3.To
						ACLMenuLocal.To.Valid = true
						erroRetorno = ExecutaCreateStruct(ACLMenuLocal, tx)
						if erroRetorno != nil {
							tx.Rollback()
							return erroRetorno
						}
					}
					log.Println(" Código Nivel 1:", regN1.Codigo, " CódigoPai Nível 2:", regN2.Codigo, " Código Nivel 3:", regN3.Codigo, " to Nível 3:", regN3.To)
				}
			}
		}
	}
	tx.Commit()
	// msgErro = "[MAMERRIMJ001 | modelAclMenu.go|InsereNoBDMenuDoJSON N.01] Erro de inserir o menu no Banco de Dados"
	// erroRetorno = errors.New(msgErro)
	return erroRetorno
}

//MontaMenu : Retorna o menu já na ordem para ser exibido
func MontaMenu(bdPar bancoDeDados.BDCon) ([]ItemsNivel1, error) {
	var (
		// ACLMenuFromJSONLocais     ACLMenuFromJSON
		// itensDaTabelaMenu         []modelmenu.ItensDaTabelaMenu
		// ItemsNivel1Local          ItemsNivel1
		// ItemsNivel1Locais         []ItemsNivel1
		menuLocal                 []ItemsNivel1
		menuItemNivel1            ItemsNivel1
		menuItemNivel2            ItemsNivel2
		menuItemNivel3            ItemsNivel3
		ACLMenuLocais             []ACLMenu
		sql                       string
		erroRetorno               error
		msgErro                   string
		qtdadeRegistrosRetornados *gorm.DB
		posN1, posN2              int
	)
	sql = `
				WITH RECURSIVE submenus AS (
					SELECT  menuid, label, codigo,  nivel, codigopai, icon, acl_menu.to,  CAST(label As varchar(1000)) As Itens_Menu
					FROM    acl_menu  
					where codigopai = ''
					UNION
					SELECT  m.menuid, m.label, m.codigo, m.nivel, m.codigopai, m.icon, m.to, CAST(s.Itens_Menu || '->' || m.label As varchar(1000)) As Itens_Menu
					FROM acl_menu m
					INNER JOIN submenus s ON s.codigo = m.codigopai
			) 
			SELECT  menuid, label, codigo, nivel, codigopai, icon , submenus.to 
				FROM submenus
			ORDER BY Itens_Menu;
	`

	qtdadeRegistrosRetornados = bdPar.BD.Raw(sql).Scan(&ACLMenuLocais) //.Count(&qtdadeRegistrosRetornados)

	if qtdadeRegistrosRetornados.RowsAffected == 0 {
		msgErro = "[MAMERRMTM001 | modelAclMenu.go|MontaMenu N.01] Não foi encontrado nenhum item de menu"
		erroRetorno = errors.New(msgErro)
		return nil, erroRetorno
	}

	posN1 = -1
	posN2 = -1

	for _, itensMenu := range ACLMenuLocais {
		if itensMenu.Nivel == 1 {
			posN1++

			menuItemNivel1 = ItemsNivel1{
				Label: itensMenu.Label,
			}

			menuLocal = append(menuLocal, menuItemNivel1)
			posN2 = -1

		} else if itensMenu.Nivel == 2 {
			posN2++
			menuItemNivel2 = ItemsNivel2{
				Label: itensMenu.Label,
			}
			menuLocal[posN1].Items = append(menuLocal[posN1].Items, menuItemNivel2)

		} else if itensMenu.Nivel == 3 {

			menuItemNivel3 = ItemsNivel3{
				Label: itensMenu.Label,
				To:    itensMenu.To.String,
			}
			menuLocal[posN1].Items[posN2].Items = append(menuLocal[posN1].Items[posN2].Items, menuItemNivel3)

		}
	}

	log.Println("[modelAclMenu.go|MontaMenu N.04|INFO] valor menuLocal:", menuLocal)
	return menuLocal, nil

}
