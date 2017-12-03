package middlewares

import (
	"db"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	signBytes []byte
)

var AuthMiddleware = &GinJWTMiddleware{
	Realm:            "test zone",
	SigningAlgorithm: "HS256",
	Key:              initKeys(),
	Timeout:          time.Hour,
	MaxRefresh:       time.Hour,
	Authenticator:    Authenticator,
	Authorizator:     Authorizator,
	Unauthorized:     Unauthorized,
	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup: "header:Authorization",
	// TokenLookup: "query:token",
	// TokenLookup: "cookie:token",

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName: "Bearer",

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc: time.Now,
}

func initKeys() []byte {
	signBytes, _ = ioutil.ReadFile("./keys/secret.rsa")
	return signBytes

}

func Authenticator(email string, password string, ctx *gin.Context) (userName string, ok bool) {

	if email != "" || password != "" {

		user := db.FindUserByName(email)

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return "", false
		}

	}

	return userName, true
}

func Authorizator(email string, ctx *gin.Context) bool {

	//if email == "test" {
	//	return true
	//}

	return true
}

func Unauthorized(ctx *gin.Context, code int, message string) {

	//if mw.Realm == "" {
	//	mw.Realm = "gin jwt"
	//}
	//
	//ctx.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	ctx.Abort()

	ctx.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
