package modelacl

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/now"
)

//ACLGrupo : zz
type ACLGrupo struct {
	//gorm.Model
	GrupoID          uint64       `gorm:"type:bigint; primary_key; type: bigserial ;column:grupoid"`
	CodigoGrupo      string       `gorm:"type:varchar(25); unique; not null;column:codigogrupo" validate:"required"`
	DescricaoGrupo   string       `gorm:"type:varchar(30);column:descricaogrupo" validate:"required"`
	DataCriacaoGrupo sql.NullTime `gorm:"type:date;column:datacriacaogrupo" validate:"required"`
	// Se o softdele valer 0 o registro esta ativo, se valer 1, ele esta removido
	SoftDelete int `gorm:"type:smallint;column:softdelete" validate:"number,gte=0,lte=1"`
	//Define se o grupo é criado como grupo originalmente (valor g) ou se é criado originalmente de um usuário (valor u)
	TipoOrigemGrupo string `gorm:"type:char;not null;default:'g';column:tipoorigemgrupo" validate:"required"`
	//Declarando chave estrangeira []Nome da struct
	// FKACLUsuarioGrupo    []ACLUsuarioGrupo
	// FKACLGrupoAcessaMenu []ACLGrupoAcessaMenu
}

//ACLGrupoJSON : zz
type ACLGrupoJSON struct {
	CodigoGrupo      string `json:"CodigoGrupo" validate:"required"`
	DescricaoGrupo   string `json:"DescricaoGrupo" validate:"required"`
	DataCriacaoGrupo string `json:"DataCriacaoGrupo" `
	SoftDelete       int    `json:"SoftDelete" validate:"number,gte=0,lte=1"`
	DataTeste        string `json:"DataTeste"`
	TipoOrigemGrupo  string `json:"TipoOrigemGrupo" validate:"required"`
}

//ACLCodigoGrupoJSON : usado para receber parâmetro de busca que neste caso é o Código do Grupo
type ACLCodigoGrupoJSON struct {
	CodigoGrupo string `json:"CodigoGrupo" validate:"required"`
}

//****************************************
//Utils

// DataParaString : recebe um campo sql.NullTime e retorna "" (vazio) se ele tiver
// o valor do campo Valid como false ou a data formatada se o valor do campo Valid for true
func DataParaString(data sql.NullTime) string {
	if data.Valid == false {
		return ""
	}
	return data.Time.Format("02/01/2006")

}

// StringParaData : recebe uma string e retorna um sql.NullTime
// **func StringParaData(dataPar *sql.NullTime, dataStringPar string) error {
func StringParaData(dataStringPar string) (sql.NullTime, error) {
	var (
		dataLocal sql.NullTime
		erroLocal error
	)

	now.TimeFormats = append(now.TimeFormats, "02/01/2006")

	if dataStringPar == "" {
		fmt.Println("[MAGINFSPD001|modelAclGrupos.go|StringParaData] Comprimento de dataString:" + strconv.Itoa(len(dataStringPar)))
		//dataPar = sql.NullTime{}
		dataLocal.Valid = false
	} else {
		dataLocal.Time, erroLocal = now.Parse(dataStringPar)
		dataLocal.Valid = true
	}

	return dataLocal, erroLocal
}

//*****************************************

// CriaNovoGrupo : cria um novo grupo no BD
func CriaNovoGrupo(ACLGrupoPAR ACLGrupo, bdPar bancoDeDados.BDCon) error {
	var (
		// achou int
		ACLGrupoLocal ACLGrupo
		erroRetorno   error
	)

	validate = validator.New()

	erroRetorno = validate.Struct(ACLGrupoPAR)

	log.Println("[MAGERRCNG006 | modelAclGrupos.go|CriaNovoGrupo N.06] Valor do ACLGrupoPAR:", ACLGrupoPAR)
	if erroRetorno != nil {
		log.Println("[MAGERRCNG005 | modelAclGrupos.go|CriaNovoGrupo N.05] Valor ACLGrupoPAR:", ACLGrupoPAR)
		log.Println("[MAGERRCNG004 | modelAclGrupos.go|CriaNovoGrupo N.04] Erro de campo obrigatório:", erroRetorno)
		return erroRetorno
	}

	// var Achou = 0
	// Achou = VerificaSeLoginJaExisteNoBD(ACLUser.Login, bd)

	Achou := bdPar.BD.Where("codigogrupo = ?", ACLGrupoPAR.CodigoGrupo).First(&ACLGrupoLocal)

	if Achou.RowsAffected >= 1 {
		msgErro = "[MAGERRCNG003 | modelAclGrupos.go|CriaNovoGrupo N.03] Erro de insert: Grupo já existe: " + ACLGrupoPAR.CodigoGrupo
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()
	result := tx.Create(&ACLGrupoPAR)
	log.Println("[MAGERRCNG002 | modelAclGrupos.go|CriaNovoGrupo N.02] linhas afetadas:", result.RowsAffected)

	if result.RowsAffected == 0 {
		msgErro = "[MAGERRCNG001 | modelAclGrupos.go|CriaNovoGrupo N.01] Erro ao criar o grupo"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		tx.Rollback()
		return erroRetorno
	}

	tx.Commit()

	return nil
}

//SoftdeleteGrupo : realiza o softdelete do grupo
func SoftdeleteGrupo(codigoGrupoPar string, bdPar bancoDeDados.BDCon) error {
	var (
		// achou int
		ACLGrupoLocal ACLGrupo
		erroRetorno   error
	)

	log.Println(" [MAGINFSDL004 | modelAclGrupos.go|SoftdeleteGrupo 04] Valor CodGrupoPar:" + codigoGrupoPar)

	registroNaoLocalizado := bdPar.BD.Where("codigogrupo = ?", codigoGrupoPar).First(&ACLGrupoLocal).RecordNotFound()
	if registroNaoLocalizado == true {
		msgErro = "[MAGERRSDL001 | modelAclGrupos.go|Softdelete 01] Grupo não localizado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()
	regAfetados := tx.Model(&ACLGrupoLocal).Where("CodigoGrupo = ?", codigoGrupoPar).Update("SoftDelete", 1)

	// log.Println(" [MAGINFSDL002 | modelAclGrupos.go|SoftdeleteGrupo 02]  Quantidade de grupos removidos:" + strconv.FormatInt(regAfetados.RowsAffected, 10))

	if regAfetados.RowsAffected == 0 {
		tx.Rollback()
		msgErro = "[MAUERRSDL003 | modelAclGrupos.go|SoftdeleteGrupo 03] Erro ao remover (softdelete) o grupo"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx.Commit()
	erroRetorno = nil

	return erroRetorno
}

//ReverteSoftdeleteGrupo : reverte o softdelete do grupo. Se o grupo for excluido com
// softdelete, ele volta a ficar acessível.
func ReverteSoftdeleteGrupo(codigoGrupoPar string, bdPar bancoDeDados.BDCon) error {
	var (
		// achou int
		ACLGrupoLocal ACLGrupo
		erroRetorno   error
	)

	log.Println(" [MAGINFRSD004 | modelAclGrupos.go|ReverteSoftdeleteGrupo 04] Valor CodGrupoPar:" + codigoGrupoPar)

	registroNaoLocalizado := bdPar.BD.Where("codigogrupo = ?", codigoGrupoPar).First(&ACLGrupoLocal).RecordNotFound()
	if registroNaoLocalizado == true {
		msgErro = "[MAGERRRSD001 | modelAclGrupos.go|ReverteSoftdeleteGrupo 01] Grupo não localizado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()
	regAfetados := tx.Model(&ACLGrupoLocal).Where("CodigoGrupo = ?", codigoGrupoPar).Update("SoftDelete", 0)

	// log.Println(" [MAGINFRSD002 | modelAclGrupos.go|ReverteSoftdeleteGrupo 02]  Quantidade de grupos removidos:" + strconv.FormatInt(regAfetados.RowsAffected, 10))

	if regAfetados.RowsAffected == 0 {
		tx.Rollback()
		msgErro = "[MAUERRRSD003 | modelAclGrupos.go|ReverteSoftdeleteGrupo 03] Erro ao ativar (reverter softdelete) o grupo"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx.Commit()
	erroRetorno = nil
	return erroRetorno
}

//ListaTodosOsGrupos : retorna todos os dados dos grupos cadastrados.
func ListaTodosOsGrupos(bdPar bancoDeDados.BDCon) (ACLGrupoJSONRetorno []ACLGrupoJSON, erroRetorno error) {
	var (
		// achou int
		ACLGrupoLocal     []ACLGrupo
		ACLGrupoJSONLocal ACLGrupoJSON
	)
	// sql := "codigogrupo,descricaogrupo, datacriacaogrupo," +
	// 	"softdelete, CASE WHEN datateste isnull THEN '' ELSE 'datateste' END as datateste"
	sqlQuery := "codigogrupo,descricaogrupo, datacriacaogrupo," +
		"softdelete"
	registroNaoLocalizado := bdPar.BD.Select(sqlQuery).Where("SOFTDELETE = ?", 0).Find(&ACLGrupoLocal).RecordNotFound()

	// log.Println(" [MAGINFRSD002 | modelAclGrupos.go|ReverteSoftdeleteGrupo 02]  Quantidade de grupos removidos:" + strconv.FormatInt(regAfetados.RowsAffected, 10))

	if registroNaoLocalizado == true {
		msgErro = "[MAUERRLTG001 | modelAclGrupos.go|ListaTodosOsGrupo 01] Nenhum grupo foi encontrado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return nil, erroRetorno
	}

	for _, reg := range ACLGrupoLocal {
		fmt.Println("[MAULTG002 | modelAclGrupos.go|ListaTodosOsGrupo 02] Valor de ACLGrupoLocal:" + reg.CodigoGrupo)
		ACLGrupoJSONLocal.CodigoGrupo = reg.CodigoGrupo
		ACLGrupoJSONLocal.DescricaoGrupo = reg.DescricaoGrupo
		ACLGrupoJSONLocal.DataCriacaoGrupo = DataParaString(reg.DataCriacaoGrupo)
		ACLGrupoJSONRetorno = append(ACLGrupoJSONRetorno, ACLGrupoJSONLocal)
	}
	// ACLGrupoJSONRetorno = ACLGrupoLocal

	return ACLGrupoJSONRetorno, nil
}
