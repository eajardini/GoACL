package modelacl

import (
	"errors"
	"log"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/jinzhu/gorm"
)

//ACLUsuarioGrupo : zz
type ACLUsuarioGrupo struct {
	//gorm.Model
	UserGrupoID uint64 `gorm:"type:bigint; primary_key; type: bigserial ;column:usergrupoid"`
	UsuarioID   uint64 `gorm:"type:bigint; column:usuarioid; unique_index:idx_usuario_grupo" validate:"required"`
	GrupoID     uint64 `gorm:"type:bigint; column:grupoid;  unique_index:idx_usuario_grupo" validate:"required"`
}

//ACLUsuarioGrupoJSON : zz
type ACLUsuarioGrupoJSON struct {
	Login       string `json:"Login" validate:"required"`
	CodigoGrupo string `json:"CodigoGrupo" validate:"required"`
}

//RespUsuario :zz
type RespUsuario struct {
	UsuarioID int64
}

// TODO: fazer rotina para remover usuário no grupo
// 		 : fazer rotina para atualizar direito de acesso do crupo aos menus

// InsereUsuarioEmGrupo : inclui usuario em um grupo
func InsereUsuarioEmGrupo(loginPar string, codigoGrupoPar string, bdPar bancoDeDados.BDCon) error {
	var (
		ACLUsuarioGrupoLocal ACLUsuarioGrupo
		// rowUsuarioIDLocal    *sql.Rows
		erroLocal error
		// respUsuario RespUsuario
		mensagemDeErro            string
		numerosDeRegistrosCriados *gorm.DB
	)

	erroLocal = bdPar.BD.Table("acl_usuario").Select("usuarioid").Where("login = ?", loginPar).Row().Scan(&ACLUsuarioGrupoLocal.UsuarioID)
	if erroLocal != nil {
		mensagemDeErro = "[MSUERRIUG001 | modelAclSetaUsuarioEmGrupo.go|InsereUsuarioEmGrupo001] Não foi possível localizar o ID do usuário" + erroLocal.Error()
		erroLocal = errors.New(mensagemDeErro)
		return erroLocal
	}

	erroLocal = bdPar.BD.Table("acl_grupo").Select("grupoid").Where("codigogrupo = ?", codigoGrupoPar).Row().Scan(&ACLUsuarioGrupoLocal.GrupoID)
	if erroLocal != nil {
		mensagemDeErro = "[MSUERRIUG002  modelAclSetaUsuarioEmGrupo.go|InsereUsuarioEmGrupo002] Não foi possível localizar o ID do grupo" + erroLocal.Error()
		erroLocal = errors.New(mensagemDeErro)
		return erroLocal
	}

	tx := bdPar.BD.Begin()
	numerosDeRegistrosCriados = tx.Create(&ACLUsuarioGrupoLocal)

	if numerosDeRegistrosCriados.RowsAffected == 0 {
		mensagemDeErro = "[MSUERRIUG003 | modelAclSetaUsuarioEmGrupo.go|InsereUsuarioEmGrupo003] Problema em adicionar usuário ao grupo" + erroLocal.Error()
		log.Println(mensagemDeErro)
		erroLocal = errors.New(mensagemDeErro)
		tx.Rollback()
		return erroLocal
	}

	tx.Commit()
	return erroLocal
}
