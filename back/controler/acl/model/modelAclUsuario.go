package modelacl

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/go-playground/validator/v10"
)

// TODO: fazer rotina que cria automaticamente um
// 		grupo com o nome do usuario para dar direito
//		pessoal do usuário poder acessar os grupos alem
//		daqueles padrões do do grupo que ele está inserido

var (
	// bd          bancoDeDados.BDCon
	msgErro string
	//erroRetorno error
	validate *validator.Validate
)

//ACLUsuario : zz
type ACLUsuario struct {
	//gorm.Model
	UsuarioID           uint64    `gorm:"type:bigint; primary_key; type: bigserial ;column:usuarioid" json:"usuarioid" `
	Login               string    `gorm:"type:varchar(25); unique; not null" json:"login" validate:"required"`
	Password            string    `gorm:"type:varchar(50)" json:"password" validate:"required"`
	Datacriacao         time.Time `gorm:"type:date" json:"datacriacao" validate:"required"`
	Datavalidade        time.Time `gorm:"type:date" json:"datavalidade"`
	Userbloqueado       int       `gorm:"type:integer" json:"userbloqueado" validate:"number,gte=0,lte=1"`
	Userativo           int       `gorm:"type:integer" json:"userativo" validate:"number,gte=0,lte=1"`
	LoginAntesDeApagado string    `gorm:"type:varchar(25); column:LoginAntesDeApagado" json:"loginAntesDeApagado"`
	//Declarando chave estrangeira []Nome da struct
	// FKACLUsuarioGrupo []ACLUsuarioGrupo
}

//ACLUsuarioJSON : zz
type ACLUsuarioJSON struct {
	Login         string `json:"login" validate:"required"`
	Password      string `json:"password" validate:"required"`
	Datacriacao   string `json:"datacriacao" validate:"required"`
	Datavalidade  string `json:"datavalidade" `
	Userbloqueado int    `json:"userbloqueado" validate:"number,gte=0,lte=1"`
	Userativo     int    `json:"userativo" validate:"number,gte=0,lte=1"`
}

//ACLUsuarioLoginJSON : usado para receber parâmetro de busca que neste caso é o login do usuário
type ACLUsuarioLoginJSON struct {
	Login string `json:"login" validate:"required"`
}

//FazAutenticacao : recebe o login e senha e confronta com os dados do BD.
// se coincidir, retorna true; senao retorn false.
func FazAutenticacao(userNamePar string, passwordPar string, bdPar bancoDeDados.BDCon) bool {
	var (
		qtdadeRegistrosAchadosLocal int
	)

	bdPar.BD.Table("acl_usuario").Where("login = ? and password = ?", userNamePar, passwordPar).Count(&qtdadeRegistrosAchadosLocal)

	if qtdadeRegistrosAchadosLocal == 0 {
		return false
	}

	return true

}

// VerificaSeLoginJaExisteNoBD : Faz a verificação se o login existe ou no BD. Se existir retorn 1 senão retorna 0
func VerificaSeLoginJaExisteNoBD(loginPar string, bdPar bancoDeDados.BDCon) int {
	var (
		// achou int
		ACLUser ACLUsuario
	)
	// bd.AbreConexao()
	// defer bd.FechaConexao()

	fmt.Println("[MAUINFVLE001 | modelAcleUsuario.go|VerificaSeLoginJaExisteNoBD_01] Valor do Login:", loginPar)

	// achou = 0
	// achou := bdPar.BD.Raw("SELECT count(usuarioid) FROM acl_usuario WHERE login = ?", loginPar)
	// bdPar.BD.Where("login = ?", loginPar).Find(&ACLUser).Count(&achou)
	achou := bdPar.BD.Where("login = ?", loginPar).First(&ACLUser)
	fmt.Println("[MAUINFVLE002 | modelAclUsuario.go|VerificaSeLoginJaExisteNoBD_02] Valor do COUNT na busca:", achou.RowsAffected)

	return int(achou.RowsAffected)
}

// CriaNovoUsuario : cria um novo usuário no BD
func CriaNovoUsuario(ACLUser ACLUsuario, bdPar bancoDeDados.BDCon) error {
	var (
		// achou int
		ACLUserLocal ACLUsuario
		erroRetorno  error
	)

	validate = validator.New()

	erroRetorno = validate.Struct(ACLUser)

	// log.Println("[modelAclUsuario.go|CriaNovousuario N.06] Valor do ACLUser:", ACLUser)
	if erroRetorno != nil {
		log.Println("[modelAclUsuario.go|CriaNovousuario| ERRO 005] Valor ACLUser:", ACLUser)
		log.Println("[modelAclUsuario.go|CriaNovousuario| ERRO 004] Erro de campo obrigatório:", erroRetorno)
		return erroRetorno
	}

	// var Achou = 0
	// Achou = VerificaSeLoginJaExisteNoBD(ACLUser.Login, bd)

	Achou := bdPar.BD.Where("login = ?", ACLUser.Login).First(&ACLUserLocal)

	if Achou.RowsAffected >= 1 {
		msgErro = "[MAUERRCNU003 | modelAclUsuario.go|CriaNovousuario N.03] Erro de insert: Usuário já existe: " + ACLUser.Login
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()
	result := tx.Create(&ACLUser)
	// log.Println("[MAUINFCNU002 | modelAclUsuario.go|CriaNovousuario N.02] linhas afetadas:", result.RowsAffected)

	if result.RowsAffected == 0 {
		msgErro = "[MAUERRCNU001 | modelAclUsuario.go|CriaNovousuario N.01] Erro ao criar no usuário"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		tx.Rollback()
		return erroRetorno
	}

	tx.Commit()

	return nil
}

//** Remove usuário

// RemoveUsuarioPorLogin : atribui valor 0 no campo UsuarioAtivo
// e assim apagar logicamente (soft delete) o usuário
func RemoveUsuarioPorLogin(loginPar string, bdPar bancoDeDados.BDCon) error {

	var (
		ACLUser     ACLUsuario
		erroRetorno error
	)

	// achou := VerificaSeLoginJaExisteNoBD(loginPar, bdPar)
	achou := bdPar.BD.Where("login = ?", loginPar).First(&ACLUser)
	if achou.RowsAffected == 0 {
		msgErro = "[MAUINFRPL001 | modelAclUsuario.go|RemoveUsuarioPorLogin 01] Usuário não localizado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()
	regAfetados := tx.Model(&ACLUser).Where("login = ?", loginPar).Update("userativo", 0)

	log.Println(" [MAUINFRPL002 | modelAclUsuario.go|RemoveUsuarioPorLogin 02]  Quantidade de usuários removidos:" + strconv.FormatInt(regAfetados.RowsAffected, 10))

	if regAfetados.RowsAffected == 0 {
		tx.Rollback()
		msgErro = "[MAUERRRPL003 | modelAclUsuario.go|RemoveUsuarioPorLogin 03] Erro ao remover o usuário"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx.Commit()
	erroRetorno = nil

	return erroRetorno
}

//RemoveFisicamenteUsuarioPorLogin : remove o login por meio do camando delete
//Usado apenas para fins didádicos. O Login não pode ser removdo da base.
//Se eu quiser liberar um login, eu uso a função LiberarLogin
func RemoveFisicamenteUsuarioPorLogin(loginPar string, bdPar bancoDeDados.BDCon) error {

	var (
		ACLUser     ACLUsuario
		erroRetorno error
	)

	achou := bdPar.BD.Where("login = ?", loginPar).First(&ACLUser)
	if achou.RowsAffected == 0 {
		msgErro = "[MAUINFRFL001 | modelAclUsuario.go|RemoveFisicamenteUsuarioPorLogin 01] Usuário não localizado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()

	// regAfetados := bdPar.BD.Unscoped().Where("login = ?", loginPar).Delete(&ACLUser)
	regAfetados := tx.Exec("Delete from acl_usuario WHERE login = ?", loginPar)

	log.Println(" [MAUINFRFL002 | modelAclUsuario.go|RemoveFisicamenteUsuarioPorLogin 02]  Quantidade de usuários removidos fisicamente:" + strconv.FormatInt(regAfetados.RowsAffected, 10))

	if regAfetados.RowsAffected == 0 {
		tx.Rollback()
		msgErro = "[MAUERRRPF003 | modelAclUsuario.go|RemoveFisicamenteUsuarioPorLogin 03] Erro ao remover fisicamente o usuário"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	// Não esta removendo de forma correta o usuário. Talvez colocar código sql direto.

	tx.Commit()
	erroRetorno = nil

	return erroRetorno
}

//Ativar usuário

// AtivaUsuarioPorLogin : atribui valor 1 no campo UsuarioAtivo
// e assim ativa usuário apagado logicamente
func AtivaUsuarioPorLogin(loginPar string, bdPar bancoDeDados.BDCon) error {

	var (
		ACLUser     ACLUsuario
		erroRetorno error
	)

	// achou := VerificaSeLoginJaExisteNoBD(loginPar, bdPar)
	achou := bdPar.BD.Where("login = ?", loginPar).First(&ACLUser)
	if achou.RowsAffected == 0 {
		msgErro = "[MAUINFAPL001 | modelAclUsuario.go|AtivaUsuarioPorLogin 01] Usuário não localizado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx := bdPar.BD.Begin()
	regAfetados := tx.Model(&ACLUser).Where("login = ?", loginPar).Update("userativo", 1)

	log.Println(" [MAUINFAPL002 | modelAclUsuario.go|AtivaUsuarioPorLogin 02]  Quantidade de usuários removidos:" + strconv.FormatInt(regAfetados.RowsAffected, 10))

	if regAfetados.RowsAffected == 0 {
		tx.Rollback()
		msgErro = "[MAUERRAPL003 | modelAclUsuario.go|AtivaUsuarioPorLogin 03] Erro ao remover o usuário"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}

	tx.Commit()
	erroRetorno = nil

	return erroRetorno
}

// AlteraUsuario : altera os dados do usuário no BD
func AlteraUsuario(ACLUser ACLUsuario, bdPar bancoDeDados.BDCon) error {
	var (
		// achou int
		ACLUserLocal ACLUsuario
		erroRetorno  error
	)

	validate = validator.New()

	erroRetorno = validate.Struct(ACLUser)

	log.Println("[MAUINFAUS006 | modelAclUsuario.go|AlteraUsuario N.06] Valor do ACLUser:", ACLUser)
	if erroRetorno != nil {
		log.Println("[MAUINFAUS005 | modelAclUsuario.go|AlteraUsuario N.05] Valor ACLUser:", ACLUser)
		log.Println("[MAUINFAUS004 | modelAclUsuario.go|AlteraUsuario N.04] Erro de campo obrigatório:", erroRetorno)
		return erroRetorno
	}

	// var Achou = 0
	// Achou = VerificaSeLoginJaExisteNoBD(ACLUser.Login, bd)

	tx := bdPar.BD.Begin()
	Achou := tx.Where("login = ?", ACLUser.Login).First(&ACLUserLocal)

	if Achou.RowsAffected == 0 {
		msgErro = "[MAUINFAUS003 |  modelAclUsuario.go|AlteraUsuario N.03] Usuário não existe: " + ACLUser.Login
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		return erroRetorno
	}
	UsuarioID := ACLUserLocal.UsuarioID
	ACLUserLocal = ACLUser
	ACLUserLocal.UsuarioID = UsuarioID

	// result := tx.Update(&ACLUser).Where("login = ?", ACLUser.Login)
	result := tx.Save(&ACLUserLocal)
	log.Println("[MAUINFAUS002 |  modelAclUsuario.go|AlteraUsuario N.02] linhas afetadas:", result.RowsAffected)

	if result.RowsAffected == 0 {
		msgErro = "[MAUINFAUS001 |  modelAclUsuario.go|AlteraUsuario N.01] Erro ao alterar no usuário"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		tx.Rollback()
		return erroRetorno
	}

	tx.Commit()

	return nil
}

//BuscaTodosUsuario : Busca todos os dados dos usuários
func BuscaTodosUsuario(bdPar bancoDeDados.BDCon) (ACLRetorno []ACLUsuario, err error) {
	var (
		erroRetorno error
	)
	tx := bdPar.BD.Begin()
	registrosEncontrados := tx.Find(&ACLRetorno)
	if registrosEncontrados.RowsAffected == 0 {
		msgErro = "[MAUERRBTU001 |  modelAclUsuario.go|BuscaTodosUsuario N.01] Nenhum usuário encontrado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		tx.Rollback()
		return nil, erroRetorno
	}

	tx.Commit()

	return ACLRetorno, nil

}

//BuscaTodosUsuariosAtivos : Busca todos os dados dos usuários ativos
func BuscaTodosUsuariosAtivos(bdPar bancoDeDados.BDCon) (ACLRetorno []ACLUsuario, err error) {

	var (
		erroRetorno error
	)

	tx := bdPar.BD.Begin()
	registrosEncontrados := tx.Where("userativo = ?", 1).Find(&ACLRetorno)
	if registrosEncontrados.RowsAffected == 0 {
		msgErro = "[MAUERRBTA001 |  modelAclUsuario.go|BuscaTodosUsuariosAtivos N.01] Nenhum usuário ativo encontrado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		tx.Rollback()
		return nil, erroRetorno
	}

	tx.Commit()

	return ACLRetorno, nil

}

//BuscaUsuarioPorLogin : Busca os dados do usuário ativos por login
func BuscaUsuarioPorLogin(loginPAR string, bdPar bancoDeDados.BDCon) (ACLRetorno ACLUsuario, err error) {

	var (
		ACLUserLocal ACLUsuario
		erroRetorno  error
	)

	tx := bdPar.BD.Begin()
	registrosEncontrados := tx.Where("login = ? AND userativo = ?", loginPAR, 1).Find(&ACLUserLocal)
	if registrosEncontrados.RowsAffected == 0 {
		log.Println("[MAUERRBTL002 |  modelAclUsuario.go|BuscaUsuarioPorLogin N.02] loginPAR:", loginPAR)
		msgErro = "[MAUERRBTL001 |  modelAclUsuario.go|BuscaUsuarioPorLogin N.01] Nenhum usuário ativo encontrado"
		log.Println(msgErro)
		erroRetorno = errors.New(msgErro)
		tx.Rollback()
		return ACLRetorno, erroRetorno
	}

	tx.Commit()

	log.Println("[MAUERRBTL003 |  modelAclUsuario.go|BuscaUsuarioPorLogin N.03] Valor do ACLUserLocal.login: ", ACLUserLocal.Login)
	ACLRetorno = ACLUserLocal

	return ACLUserLocal, nil

}

//TODO: liberar Login
