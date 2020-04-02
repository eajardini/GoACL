package modelacl

import (
	"errors"
	"log"
	"time"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
)

var (
	bd          bancoDeDados.BDCon
	msgErro     string
	erroRetorno error
)

//ACLUsuario : zz
type ACLUsuario struct {
	//gorm.Model
	UsuarioID     uint64    `gorm:"type:bigint; primary_key; type: serial ;column:usuarioid" json:"usuarioid" validate:"required"`
	Login         string    `gorm:"type:varchar(20)" json:"login" validate:"required"`
	Password      string    `gorm:"type:varchar(20)" json:"password" validate:"required"`
	Datacriacao   time.Time `gorm:"type:date" json:"datacriacao" validate:"required"`
	Datavalidade  time.Time `gorm:"type:date" json:"datavalidade" validate:"required"`
	Userbloqueado bool      `gorm:"type:boolean" json:"userbloqueado" validate:"required"`
	Userativo     bool      `gorm:"type:boolean" json:"userativo"`
}

//ACLUsuarioJSON : zz
type ACLUsuarioJSON struct {
	UsuarioID     string `json:"usuarioid" validate:"required"`
	Login         string `json:"login" validate:"required"`
	Password      string `json:"password" validate:"required"`
	Datacriacao   string `json:"datacriacao" validate:"required"`
	Datavalidade  string `json:"datavalidade" validate:"required"`
	Userbloqueado bool   `json:"userbloqueado" validate:"required"`
	Userativo     bool   `json:"userativo"`
}

// CriaNovoUsuario : cria um novo usuário no BD
func CriaNovoUsuario(ACLUser ACLUsuario) error {
	bd.AbreConexao()
	defer bd.FechaConexao()

	result := bd.BD.Create(&ACLUser)
	log.Println("[modelAclUsuario.go|CriaNovousuario N.02] linhas afetadas:", result.RowsAffected)

	if result.RowsAffected != 0 {
		msgErro = "[modelAclUsuario.go|CriaNovousuario N.01] Erro ao criar no usuário"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	return nil

}
