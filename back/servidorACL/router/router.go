package router

import (
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	bancoDeDados "github.com/eajardini/ProjetoGoACL/GoACL/back/bancodedados"
	"github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib"
	mensagensErros "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/lib/model"
	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	bd                   bancoDeDados.BDCon
	msgErros             mensagensErros.LIBErroMSGRetorno
	r                    *gin.Engine
	libMSG               mensagensErros.LIBErroMSGSGBD
	authMiddlewareGlobal *jwt.GinJWTMiddleware
)

// ConfiguraBD :
func ConfiguraBD() {
	bd.ConfiguraStringDeConexao("../config/ConfigBancoDados.toml")
	bd.IniciaConexao()
}

// ConfiguraMSGErros : Carrega para a memória todas as mensagens de erros que
// estão cadastradas no BD
func ConfiguraMSGErros(bdPar bancoDeDados.BDCon) {
	var modulo string

	bd.AbreConexao()
	defer bd.FechaConexao()

	msgErros = libMSG.CarregaTodosAsMensagensDeErro(bd)
	if msgErros.Erro != nil {
		modulo = "[router.go|ConfiguraMSGErros|ERR01] "
		log.Fatal(modulo + msgErros.Erro.Error())
		os.Exit(0)
	}
	modulo = "[router.go|ConfiguraMSGErros|INF01] "
	log.Println(modulo + "Mensagens configuradas corretamente!")
}

//ConfiguraGin :
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
		port = ":20100" //acl
	}
	r = gin.Default()
	ConfiguraGin(r)
	ConfiguraBD()
	ConfiguraMSGErros(bd)

	//***********
	//**** JWT
	//***********
	authMiddlewareGlobal = lib.ConfiguraJWT()

	//************
	//**** Rotas
	//************
	ConfiguraACL(authMiddlewareGlobal)

	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
