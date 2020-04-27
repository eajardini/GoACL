package router

import (
	"log"
	"os"
	"time"

	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	bd bancoDeDados.BDCon
	r  *gin.Engine
)

func ConfiguraBD() {
	bd.ConfiguraStringDeConexao("../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
}

func ConfiguraGin(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

}

// IniciaRouter : sobe o servidor
func IniciaRouter() {
	port := os.Getenv("PORT")

	if port == "" {
		port = ":8211" //acl
	}
	r = gin.Default()
	ConfiguraGin(r)
	ConfiguraBD()
	//************
	//**** Rotas
	//************
	ConfiguraACL()

	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
