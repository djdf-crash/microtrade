package handlers

import (
	"db"
	"middlewares"
	"net/http"
	"time"

	"utils"

	"errors"

	"config"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func resetPasswordReq(ctx *gin.Context) {
	var resetPassword ResetPasswordReq

	if err := ctx.ShouldBindJSON(&resetPassword); err == nil {
		hashPassword, _, err := getPasswordHash(resetPassword.Email)
		if err != nil {
			respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
			return
		}

		tokenReset := utils.NewToken(resetPassword.Email, 24*time.Hour, hashPassword, middlewares.AuthMiddleware.Key)
		fullPath := "http://localhost" + config.AppConfig.Port + "/token/" + tokenReset

		bodyMessage := "Please click " + fullPath + " for reset you password"
		err = utils.SendEmail(config.AppConfig.SendEmail.Server, config.AppConfig.SendEmail.Port, config.AppConfig.SendEmail.Sender,
			config.AppConfig.SendEmail.PasswordSender, resetPassword.Email, bodyMessage)
		if err != nil {
			respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
			return
		}

		respondWithMessage(http.StatusBadRequest, "Message send for email "+resetPassword.Email, ctx)

	} else {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
	}
}

func confirmPasswordReq(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		//respondWithMessage(http.StatusNotFound, "Invalid token", ctx)
		ctx.Redirect(http.StatusTemporaryRedirect, "/#/404/")
	}

	_, err := utils.VerifyToken(token, getPasswordHash, middlewares.AuthMiddleware.Key)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/#/404/")
		//respondWithMessage(http.StatusNotFound, "", ctx)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/#/confirm/"+token)
	//respondWithMessage(http.StatusOK, "login:"+login, ctx)
}

func getPasswordHash(login string) ([]byte, time.Time, error) {
	user, ok := db.CheckUserByEmail(login)
	if !ok {
		return nil, time.Now(), errors.New("User " + login + " not found")
	}

	return []byte(user.Password), user.LastLogin, nil
}

func changePasswordUser(ctx *gin.Context) {
	var changePassword ChangePasswordReq

	if err := ctx.ShouldBindJSON(&changePassword); err == nil {

		userEmail := ctx.GetString("userID")
		if userEmail == "" {
			respondWithMessage(http.StatusBadRequest, "User not found", ctx)
			return
		}

		user := db.FindUserByName(userEmail)

		if changePassword.NewPassword == "" {
			respondWithMessage(http.StatusBadRequest, "New password is empty", ctx)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassword.Password)); err != nil {
			respondWithMessage(http.StatusBadRequest, "Password is not correctly", ctx)
			return
		}

		newHash, _ := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)

		user.Password = string(newHash)
		user.LastLogin = time.Now()

		db.UpdateUser(&user)

	} else {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
	}

	respondWithMessage(http.StatusOK, "Successful", ctx)
}

func registerUser(ctx *gin.Context) {
	var userRegister Register
	var user db.User

	if err := ctx.ShouldBindJSON(&userRegister); err == nil {
		if userRegister.Password != userRegister.ConfirmPassword {
			respondWithMessage(http.StatusBadRequest, "Password and confirm password not equals", ctx)
			return
		}

		if _, ok := db.CheckUserByEmail(userRegister.Email); !ok {

			hash, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
			if err != nil {
				return
			}

			user.Email = userRegister.Email
			user.Password = string(hash)
			user.LastLogin = time.Now()

			if err := db.AddUser(&user); err != nil {
				respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
			} else {

				jwtToken, expire, _ := middlewares.AuthMiddleware.TokenGenerator(user.Email)

				respondWithMessage(http.StatusCreated, "token:"+jwtToken+"; expire:"+expire.Format(time.RFC3339), ctx)
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

func staticFilesGet(urlPrefix string, fs static.ServeFileSystem) gin.HandlerFunc {

	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			if c.Request.URL.Path != "/" {
				c.Writer.Header().Set("Cache-Control", "max-age=604800")
			}
			c.Abort()
		}
	}

}

func respondWithMessage(code int, message string, ctx *gin.Context) {

	response := ResponseMessage{
		Message{
			code,
			message,
		},
	}

	if ctx.Request.Method == http.MethodGet {
		ctx.Writer.WriteHeader(code)
	} else {
		ctx.JSON(code, &response)
	}

	ctx.Abort()
}
