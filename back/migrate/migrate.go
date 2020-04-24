package main

import (
	"log"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	modelacl "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl/model"
)

var (
	bd bancoDeDados.BDCon
)

func main() {
	bd.ConfiguraStringDeConexao("../config/ConfigBancoDados.toml")
	bd.IniciaConexao()

	bd.AbreConexao()
	defer bd.FechaConexao()

	//Realizando a migração das tabelas dos models
	bd.BD.SingularTable(true)
	bd.BD.AutoMigrate(&modelacl.ACLUsuario{}, &modelacl.ACLGrupo{}, &modelacl.ACLUsuarioGrupo{})

	log.Println("[migrate.go|main N.01] Tabelas criado com sucesso!")
}
