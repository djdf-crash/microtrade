package middlewares

import (
	"time"

	"db"

	"crypto/rsa"

	"crypto/rand"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	verifyKey *rsa.PublicKey
	secretKey *rsa.PrivateKey
)

func init() {
	initKeys()
}

func initKeys() {
	secretKey, _ = rsa.GenerateKey(rand.Reader, 1024)

}

func AuthMiddleware() *jwt.GinJWTMiddleware {

	AuthMiddleware := &jwt.GinJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte(secretKey.D.String()),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: Authenticator,
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
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
	return AuthMiddleware
}

func Authenticator(username string, password string, ctx *gin.Context) (userName string, ok bool) {
	var userDb db.User

	db.InitDB()
	defer db.DB.Close()

	if username != "" || password != "" {
		db.DB.Where("username=?", username).Find(&userDb)

		if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(password)); err != nil {
			return "", false
		}

	}

	return userName, true
}
