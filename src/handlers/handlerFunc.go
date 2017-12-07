package handlers

import (
	"db"
	"middlewares"
	"net/http"
	"time"

	"utils"

	"errors"

	"config"

	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v8"
)

func loginReq(ctx *gin.Context) {
	tokenString, _, err := middlewares.AuthMiddleware.LoginHandler(ctx)
	if err != nil {
		RespondWithMessage(http.StatusBadRequest, 201, fmt.Sprintf(utils.LoginError[201], err.Error()), ctx)
		return
	}
	RespondWithMessage(http.StatusCreated, 109, fmt.Sprintf(utils.UserRegisterError[109], tokenString), ctx)
}

func resetPasswordReq(ctx *gin.Context) {
	var resetPassword ResetPasswordReq

	if err := ctx.ShouldBindJSON(&resetPassword); err == nil {
		hashPassword, _, err := getPasswordHash(resetPassword.Email)
		if err != nil {
			RespondWithMessage(http.StatusBadRequest, 304, utils.PasswordResetRequestError[304], ctx)
			return
		}

		tokenReset := utils.NewToken(resetPassword.Email, 24*time.Hour, hashPassword, middlewares.AuthMiddleware.Key)
		fullPath := "http://localhost" + config.AppConfig.Port + "/token/" + tokenReset

		bodyMessage := "Please click " + fullPath + " for reset you password"
		err = utils.SendEmail(config.AppConfig.SendEmail.Server, config.AppConfig.SendEmail.Port, config.AppConfig.SendEmail.Sender,
			config.AppConfig.SendEmail.PasswordSender, resetPassword.Email, bodyMessage)
		if err != nil {
			RespondWithMessage(http.StatusBadRequest, 302, utils.PasswordResetRequestError[302], ctx)
			return
		}

		RespondWithMessage(http.StatusBadRequest, 303, fmt.Sprintf(utils.PasswordResetRequestError[302], resetPassword.Email), ctx)

	} else {
		checkErrors(err, ctx)
	}
}

func confirmPasswordReq(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		//RespondWithMessage(http.StatusNotFound, "Invalid token", ctx)
		ctx.Redirect(http.StatusTemporaryRedirect, "/#/404/")
	}

	_, err := utils.VerifyToken(token, getPasswordHash, middlewares.AuthMiddleware.Key)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/#/404/")
		//RespondWithMessage(http.StatusNotFound, "", ctx)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/#/confirm/"+token)
	//RespondWithMessage(http.StatusOK, "login:"+login, ctx)
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
			RespondWithMessage(http.StatusBadRequest, 501, utils.PasswordChangeError[501], ctx)
			return
		}

		if changePassword.NewPassword == "" {
			RespondWithMessage(http.StatusBadRequest, 502, utils.PasswordChangeError[502], ctx)
			return
		}

		user := db.FindUserByName(userEmail)

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassword.Password)); err != nil {
			RespondWithMessage(http.StatusBadRequest, 503, utils.PasswordChangeError[503], ctx)
			return
		}

		newHash, _ := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)

		user.Password = string(newHash)
		user.LastLogin = time.Now()

		db.UpdateUser(&user)

		RespondWithMessage(http.StatusOK, 200, "Successful", ctx)

	} else {
		checkErrors(err, ctx)
	}

}

func registerUser(ctx *gin.Context) {
	var userRegister Register
	var user db.User

	if err := ctx.ShouldBindJSON(&userRegister); err == nil {
		if userRegister.Password != userRegister.ConfirmPassword {
			RespondWithMessage(http.StatusBadRequest, 107, utils.UserRegisterError[107], ctx)
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
				RespondWithMessage(http.StatusBadRequest, 108, utils.UserRegisterError[108], ctx)
			} else {

				jwtToken, _, _ := middlewares.AuthMiddleware.TokenGenerator(user.Email)

				RespondWithMessage(http.StatusCreated, 109, fmt.Sprintf(utils.UserRegisterError[109], jwtToken), ctx)
			}
		} else {
			RespondWithMessage(http.StatusBadRequest, 101, utils.UserRegisterError[101], ctx)
			return
		}

	} else {

		checkErrors(err, ctx)
	}
}
func refreshToken(ctx *gin.Context) {
	middlewares.AuthMiddleware.RefreshHandler(ctx)
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

func checkErrors(e interface{}, ctx *gin.Context) {
	switch e.(type) {

	case validator.ValidationErrors:
		validatorErrors := e.(validator.ValidationErrors)
		errorBindingValidation(validatorErrors, ctx)
	default:
		RespondWithMessage(http.StatusBadRequest, -2, utils.CommonError[-2], ctx)
	}
}

func errorBindingValidation(validatorErrors validator.ValidationErrors, ctx *gin.Context) {
	var mess string
	var codeError int
	for _, err := range validatorErrors {
		switch err.Name {
		case "Email":
			codeError = 104
			mess = utils.UserRegisterError[104]
		case "Password":
			codeError = 105
			mess = utils.UserRegisterError[105]
		case "ConfirmPassword":
			codeError = 106
			mess = utils.UserRegisterError[106]
		}
		break
	}
	RespondWithMessage(http.StatusBadRequest, codeError, mess, ctx)
}

func RespondWithMessage(codeResponse int, codeError int, message string, ctx *gin.Context) {

	response := map[string]interface{}{
		"code":    codeError,
		"message": message,
	}

	if ctx.Request.Method == http.MethodGet {
		ctx.Writer.WriteHeader(codeResponse)
	} else {
		ctx.JSON(codeResponse, &response)
	}

	ctx.Abort()
}
