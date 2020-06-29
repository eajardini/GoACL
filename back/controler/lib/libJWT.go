package lib

// Site: https://github.com/appleboy/gin-jwt
// fazer requisicao inicial:http -v --json POST localhost:20100/login username=admin password=admin
// fazer requisicao de reflash :

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	aclcontroler "github.com/eajardini/ProjetoGoACL/GoACL/back/controler/acl"
	"github.com/gin-gonic/gin"
)

//LoginJWT :usado para autenticação do JWT
type LoginJWT struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//RetornoTokenJWT : usado para pegar o retorno das funções do JWT
type RetornoTokenJWT struct {
	Code   int    `json:"code"` //Code=200 se usuer autenticado; Code=401 se user não autenticado
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

var identityKey = "id"

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

//UserRetorno : zz
type UserRetorno struct {
	ID        string
	UserName  string
	FirstName string
	LastName  string
}

// GetClaims :zz
func GetClaims(c *gin.Context) jwt.MapClaims {
	return jwt.ExtractClaims(c)
}

//ConfiguraJWT :zz
func ConfiguraJWT() *jwt.GinJWTMiddleware {

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		// Realm:       "test zone",
		Key:         []byte("BorodinChave"),
		Timeout:     time.Minute * 2,
		MaxRefresh:  time.Minute * 2,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var (
				loginVals        LoginJWT
				autenticadoLocal bool //Se estiver autenticado, retornar true, caso contrario, retornar false
			)

			if err := c.ShouldBind(&loginVals); err != nil {
				log.Println("[libJWT.go|autheticador|INFO003] Não fiz Bind")
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			log.Println("[libJWT.go|autheticador|INFO003] username e password:" + userID + " " + password)

			autenticadoLocal = aclcontroler.FazAutenticacao(userID, password)
			if autenticadoLocal == true {
				return &User{
					UserName:  userID,
					LastName:  "",
					FirstName: "",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		panic(0)
	}
	return authMiddleware

}
