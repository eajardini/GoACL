package bancodedados

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/BurntSushi/toml"
)

//BDCon : Variável de conexão de Banco de Dados
type BDCon struct {
	BD *gorm.DB
}

//SysConfig :zz
type SysConfig struct {
	Principal    principalConfig
	BancoDeDados map[string]bancoDeDadosConfig
}

//principalConfig : zz
type principalConfig struct {
	Modo string
}

//BancoDeDados : zz
type bancoDeDadosConfig struct {
	SGBD     string
	Host     string
	User     string
	Port     string
	Ssl      string
	Database string
	Password string
}

//sysConfig : zz
var (
	sysConfig     SysConfig
	modo          string
	sgbd          string
	stringConexao string
)

// AbreArquivoDeConfiguracaoDoBancoDeDados : abre o arquivo ConfigBancoDados.toml e seta a
// estrutura de Config.
func (bdcon *BDCon) AbreArquivoDeConfiguracaoDoBancoDeDados(caminhoDoArquivoToml string) error {

	var err error

	if _, err = toml.DecodeFile(caminhoDoArquivoToml, &sysConfig); err != nil {
		fmt.Println("[bancodedados:AbreArquivoDeConfiguracaoDoBancoDeDados N.01] Erro ao abrio o arquivo TOML de configuração")
		log.Fatal(err)
	}
	return err
}

// SetaStringDeConexao : zz
func (bdcon *BDCon) SetaStringDeConexao(modo string) {

	//log.Println("[bancodados, SetaStringDeConexao] Modo de Abertura: ", modo)

	sgbd = sysConfig.BancoDeDados[modo].SGBD
	stringConexao = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		sysConfig.BancoDeDados[modo].User, sysConfig.BancoDeDados[modo].Password, sysConfig.BancoDeDados[modo].Database,
		sysConfig.BancoDeDados[modo].Host, sysConfig.BancoDeDados[modo].Ssl)
	fmt.Println("[string de Coneção:]", stringConexao)
	// fmt.Println("Host:", sysConfig.BancoDeDados[sysConfig.Principal.Modo].Database)
}

//ConfiguraStringDeConexao : este método abre o arquivo de configuração do banco de dados
//	e depois seta a string de conexão
func (bdcon *BDCon) ConfiguraStringDeConexao(caminhoDoArquivoToml string) error {
	err := bdcon.AbreArquivoDeConfiguracaoDoBancoDeDados(caminhoDoArquivoToml)
	bdcon.SetaStringDeConexao(sysConfig.Principal.Modo)

	return err
}

//IniciaConexao : Faz a conexão inicial. A string de conexão já deve estar
//								configurada previamente pelo método ConfiguraStringDeConexao()
func (bdcon *BDCon) IniciaConexao() error {
	db, err := gorm.Open(sgbd, stringConexao)
	if err != nil {
		log.Println("[bancodados.go| IniciaConexao N.02] Erro:", err.Error())
		log.Println("[bancodados.go| IniciaConexao N.03] Erro:String conexão", stringConexao)
		log.Fatal("[bancodados.go| IniciaConexao N.04] Erro ao Conectar ao Banco de Dados!!")

	}

	err = db.DB().Ping()

	if err != nil {
		log.Fatal("[bancodados.go| IniciaConexao N.05] Erro ao Pingar ao Banco de Dados!!")
	} else {
		log.Println("[bancodados.go| IniciaConexao N.06] Banco de Dados conectado corretamente!!")
	}

	db.Close()
	return err
}

//AbreConexao :zz
func (bdcon *BDCon) AbreConexao() error {
	db, err := gorm.Open(sgbd, stringConexao)
	if err != nil {
		log.Println("[BDDERRACN007  | bancodedados.go|AbreConexao N.07]: Erro ao abrir o Banco de Dados. Motivo: ", err)
		log.Fatalln(err)
	}
	//seta a variavel BD para fazer as consultas SQL

	db.SingularTable(true)
	bdcon.BD = db

	return err

}

//FechaConexao :zz
func (bdcon *BDCon) FechaConexao() error {
	err := bdcon.BD.Close()
	return err
}

//Insert : faz o insert no Banco de Dados
func (bdcon *BDCon) Insert(sgbd string, sqlinsert interface{}) {
	fmt.Println(sqlinsert)
}

// ExecutaMigrate : zz
func (bdcon *BDCon) ExecutaMigrate(schemaSQL []byte) error {
	_, err := bdcon.BD.DB().Exec(string(schemaSQL))
	if err != nil {
		log.Println("[bancodados.go|ExecutaMigrate N.08] Erro ao realizar o Migrate")
		log.Println("[bancodados.go|ExecutaMigrate N.09] Descarregando o arquivo: ", string(schemaSQL))
		log.Fatalln(err)
	}

	return err
}
