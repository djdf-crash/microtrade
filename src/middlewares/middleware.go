package middlewares

import (
	"db"
	"io/ioutil"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	signBytes []byte
)

var AuthMiddleware = &GinJWTMiddleware{
	Realm:            "test zone",
	SigningAlgorithm: "RS256",
	VerifyKey:        initVerifyKey(),
	SignKey:          initSignKey(),
	Timeout:          time.Hour * 24,
	MaxRefresh:       time.Hour * 168,
	Authenticator:    Authenticator,
	Authorizator:     Authorizator,
	//PayloadFunc:    PayloadFunc,
	Unauthorized: Unauthorized,
	Response:     Response,
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

func initVerifyKey() []byte {
	signBytes, _ = ioutil.ReadFile("./keys/public.rsa")
	return signBytes

}

func initSignKey() []byte {
	signBytes, _ = ioutil.ReadFile("./keys/secret.rsa")
	return signBytes

}

func Authenticator(email string, password string, ctx *gin.Context) (userName *db.User, ok bool) {

	var user *db.User

	if email != "" || password != "" {

		user = db.FindUserByName(email)

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return user, false
		}

		user.LastLogin = time.Now()
		db.UpdateUser(user)

		return user, true
	}

	return user, false

}

func Authorizator(email string, ctx *gin.Context) bool {

	//ctx.Request.URL.Path
	return true
}

//func PayloadFunc(userID string) map[string]interface{} {
//
//	//return map[string]interface{}
//	//user := db.FindUserByName(userID)
//	//md5 := md5.New()
//	//newHash := string(md5.Sum([]byte(user.Password)))
//	//return map[string]interface{}{
//	//	"hash": newHash,
//	//}
//}

func Response(codeHTTP, codeERR int, message string, ctx *gin.Context) {
	response := map[string]interface{}{
		"code":    codeERR,
		"message": message,
	}

	if ctx.Request.Method == http.MethodGet {
		ctx.Writer.WriteHeader(codeHTTP)
	} else {
		ctx.JSON(codeHTTP, &response)
	}
	ctx.Abort()
}

func Unauthorized(ctx *gin.Context, codeHTTP int, codeERR int, message string) {

	response := map[string]interface{}{
		"code":    codeERR,
		"message": message,
	}

	ctx.JSON(codeHTTP, &response)
	ctx.Abort()
}
