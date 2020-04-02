package router

import (
	aclpkg "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
)

//ConfiguraACL : contem todas as rotas para o controler ACL
func ConfiguraACL() {
	acl := r.Group("/acl")
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.POST("/NovoUsuario", aclpkg.NovoUsuario)
	}
}
