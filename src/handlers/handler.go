package handlers

import (
	"db"
	"middlewares"
	"net/http"

	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	jwtGO "gopkg.in/dgrijalva/jwt-go.v3"
)

type Register struct {
	Email           string `json:"email" binding:"required,emailValidator"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// Login form structure.
type Login struct {
	Email    string `json:"email" binding:"required,emailValidator"`
	Password string `json:"password" binding:"required"`
}

type ResponseMessage struct {
	Code    int
	Message string
}

func LoginHandler(ctx *gin.Context) {
	gwtMiddleware := middlewares.AuthMiddleware()
	initLogin(ctx, gwtMiddleware)
}

func initLogin(c *gin.Context, mw *jwt.GinJWTMiddleware) {
	// Initial middleware default setting.
	mw.MiddlewareInit()

	var loginVals Login

	if c.ShouldBindWith(&loginVals, binding.JSON) != nil {
		mw.Unauthorized(c, http.StatusBadRequest, "Missing Username or Password")
		return
	}

	if mw.Authenticator == nil {
		mw.Unauthorized(c, http.StatusInternalServerError, "Missing define authenticator func")
		return
	}

	userID, ok := mw.Authenticator(loginVals.Email, loginVals.Password, c)

	if !ok {
		mw.Unauthorized(c, http.StatusUnauthorized, "Incorrect Username / Password")
		return
	}

	// Create the token
	token := jwtGO.New(jwtGO.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwtGO.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(loginVals.Email) {
			claims[key] = value
		}
	}

	if userID == "" {
		userID = loginVals.Email
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["id"] = userID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	tokenString, err := token.SignedString(mw.Key)

	if err != nil {
		mw.Unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func LogoutHandler(ctx *gin.Context) {
	resp := map[string]string{"Logout": "OK"}
	ctx.JSON(http.StatusOK, resp)
}

func RegisterHandler(ctx *gin.Context) {
	registerUser(ctx)
}
func registerUser(ctx *gin.Context) {
	var userRegister Register
	var user db.Users

	if err := ctx.ShouldBindJSON(&userRegister); err == nil {
		if userRegister.Password != userRegister.ConfirmPassword {
			respondWithMessage(http.StatusBadRequest, "Password and confirm password not equals", ctx)
			return
		}

		if !db.CheckUserByEmail(userRegister.Email) {

			hash, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
			if err != nil {
				return
			}

			user.Email = userRegister.Email
			user.Password = string(hash)

			if err := db.AddUser(&user); err != nil {
				respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
			} else {
				respondWithMessage(http.StatusCreated, "User registered", ctx)
			}
		} else {
			respondWithMessage(http.StatusBadRequest, "User name is exist", ctx)
			return
		}

	} else {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
		return
	}
}

func RefreshTokenHandler(ctx *gin.Context) {
	gwtMiddleware := middlewares.AuthMiddleware()
	gwtMiddleware.RefreshHandler(ctx)
}

func respondWithMessage(code int, message string, ctx *gin.Context) {
	response := ResponseMessage{
		code,
		message,
	}

	ctx.JSON(code, &response)
	ctx.Abort()
}
