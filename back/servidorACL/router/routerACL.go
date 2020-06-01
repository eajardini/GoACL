package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
)

//ConfiguraACL : contem todas as rotas para o controler ACL
func ConfiguraACL(authMiddlewarePar *jwt.GinJWTMiddleware) {
	acl := r.Group("/acl")
	acl.Use(authMiddlewarePar.MiddlewareFunc())
	{
		// rh.GET("/retornafotofuncionario/:idFuncionario", funcionarios.RetornaFotoFuncionario)
		acl.POST("/NovoUsuario", aclcontroler.NovoUsuario)
		acl.POST("/DesativaUsuario", aclcontroler.DesativaUsuario)
		acl.POST("/RemoveFisicamenteUsuario", aclcontroler.RemoveFisicamenteUsuario)
		acl.POST("/AtivaUsuario", aclcontroler.AtivaUsuario)
		acl.POST("/AlteraUsuario", aclcontroler.AlteraUsuario)
		acl.GET("/BuscaUsuarioPorLogin", aclcontroler.BuscaUsuarioPorLogin)
		acl.GET("/BuscaTodosUsuario", aclcontroler.BuscaTodosUsuario)
		acl.GET("/BuscaTodosUsuariosAtivos", aclcontroler.BuscaTodosUsuariosAtivos)

		//Grupos

		acl.POST("/NovoGrupo", aclcontroler.NovoGrupo)
		acl.POST("/ReverteSoftDeleteGrupo", aclcontroler.ReverteSoftDeleteGrupo)
		acl.DELETE("/SoftDeleteGrupo", aclcontroler.SoftDeleteGrupo)
		acl.GET("/ListaTodosOsGrupos", aclcontroler.ListaTodosOsGrupos)

		//Menu
		acl.GET("/MontaMenu", aclcontroler.MontaMenu)

	}
	r.POST("/login", authMiddlewarePar.LoginHandler)
}
