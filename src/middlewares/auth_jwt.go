package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"db"

	"utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// GinJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type GinJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	VerifyKey []byte

	SignKey []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(userID string, password string, c *gin.Context) (*db.User, bool)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userID string, c *gin.Context) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(userID string) map[string]interface{}

	// User can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, int, string)

	Response func(codeHTTP int, codeERR int, message interface{}, c *gin.Context)

	// Set the identity handler function
	IdentityHandler func(jwt.MapClaims) string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time
}

// Login form structure.
type Login struct {
	Email    string `json:"email" binding:"required,emailValidator"`
	Password string `json:"password" binding:"required"`
}

// MiddlewareInit initialize jwt configs.
func (mw *GinJWTMiddleware) MiddlewareInit() error {

	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(userID string, c *gin.Context) bool {
			return true
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) string {
			return claims["id"].(string)
		}
	}

	if mw.Response == nil {
		return errors.New("func Response is nil")
	}

	if mw.Realm == "" {
		return errors.New("realm is required")
	}

	if mw.SignKey == nil {
		return errors.New("secret key is required")
	}

	return nil
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, 4, utils.ValidationReqError[4])
			return
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, -3, utils.CommonError[-3])
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	id := mw.IdentityHandler(claims)
	c.Set("JWT_TOKEN", token)
	c.Set("JWT_PAYLOAD", claims)
	c.Set("userID", id)

	if !mw.Authorizator(id, c) {
		mw.unauthorized(c, http.StatusUnauthorized, -1, utils.CommonError[-1])
		return
	}

	c.Next()
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) (map[string]string, error) {

	var tokens map[string]string

	// Initial middleware default setting.
	mw.MiddlewareInit()

	var loginVals Login

	if c.ShouldBindWith(&loginVals, binding.JSON) != nil {

		return tokens, errors.New("Missing Username or Password")
	}

	if mw.Authenticator == nil {
		return tokens, errors.New("Missing define authenticator func")
	}

	userID, ok := mw.Authenticator(loginVals.Email, loginVals.Password, c)

	if !ok {
		return tokens, errors.New("Incorrect Username / Password")
	}

	tokens, err := mw.TokenGenerator(userID)

	if err != nil {
		return tokens, errors.New("Create JWT Token faild")
	}

	userID.RefreshToken = tokens["token_refresh"]
	userID.LastLogin = time.Now()
	db.UpdateUser(userID)

	return tokens, nil

}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {

	var token *jwt.Token
	var err error

	tokenVal, ok := c.Get("JWT_TOKEN")
	if !ok {
		token, err = mw.parseToken(c)
		if err != nil {
			mw.unauthorized(c, http.StatusUnauthorized, -3, utils.CommonError[-3])
			return
		}
	} else {
		token = tokenVal.(*jwt.Token)
	}

	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, -3, utils.CommonError[-3])
		return
	}

	user := db.FindUserByName(claims["id"].(string))

	if !strings.EqualFold(token.Raw, user.RefreshToken) {
		mw.unauthorized(c, http.StatusUnauthorized, -5, utils.CommonError[-5])
		return
	}

	// Create the token
	tokensMap, _ := mw.TokenGenerator(user)

	user.RefreshToken = tokensMap["token_refresh"]

	db.UpdateUser(user)

	mw.Response(http.StatusOK, 200, tokensMap, c)
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}

// TokenGenerator method that clients can use to get a jwt token.
func (mw *GinJWTMiddleware) TokenGenerator(userID *db.User) (map[string]string, error) {

	claims := jwt.MapClaims{}

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID.Email) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().UTC().Add(mw.Timeout)
	claims["id"] = userID.Email
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	token := jwt.NewWithClaims(jwt.GetSigningMethod(mw.SigningAlgorithm), claims)

	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(mw.SignKey)

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		return map[string]string{}, errors.New("Create JWT Token faild")
	}

	tokenRefresh, err := mw.TokenRefreshGenerator(userID, userID.Email)

	if err != nil {
		return map[string]string{}, errors.New("Create JWT Token faild")
	}

	return map[string]string{
		"token":         tokenString,
		"token_refresh": tokenRefresh,
	}, nil
}

// TokenGenerator method that clients can use to get a jwt token.
func (mw *GinJWTMiddleware) TokenRefreshGenerator(user *db.User, userID string) (string, error) {

	newClaims := jwt.MapClaims{}

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(user.Email) {
			newClaims[key] = value
		}
	}

	//origIat := float64(claims["orig_iat"].(int64))

	expire := mw.TimeFunc().Add(mw.Timeout).Add(mw.MaxRefresh)
	newClaims["id"] = userID
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = mw.TimeFunc().Unix()

	token := jwt.NewWithClaims(jwt.GetSigningMethod(mw.SigningAlgorithm), &newClaims)

	signKey, _ := jwt.ParseRSAPrivateKeyFromPEM(mw.SignKey)

	tokenString, err := token.SignedString(signKey)

	return tokenString, err
}

func (mw *GinJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (mw *GinJWTMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)

	if token == "" {
		return "", errors.New("Query token empty")
	}

	return token, nil
}

func (mw *GinJWTMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)

	if cookie == "" {
		return "", errors.New("Cookie token empty")
	}

	return cookie, nil
}

func (mw *GinJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split(mw.TokenLookup, ":")
	switch parts[0] {
	case "header":
		token, err = mw.jwtFromHeader(c, parts[1])
	case "query":
		token, err = mw.jwtFromQuery(c, parts[1])
	case "cookie":
		token, err = mw.jwtFromCookie(c, parts[1])
	}

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(mw.VerifyKey)

		return verifyKey, nil
	})
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, codeHTTP int, codeERR int, message string) {

	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, codeHTTP, codeERR, message)

	return
}
