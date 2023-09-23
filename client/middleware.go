package client

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/candbright/client-auth/repo"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var identityKey = "phone_number"
var defaultUsername = "admin"
var defaultPassword = "admin@123456"

var adminUser = &repo.User{
	Username:    defaultUsername,
	Password:    defaultPassword,
	PhoneNumber: "18888888888",
}

type MiddlewareConfig struct {
	Realm                string
	GetUserByPhoneNumber func(phoneNumber string) (repo.User, error)
	Unauthorized         func(c *gin.Context, code int, message string)
	NoRoute              func(*gin.Context)
}

func (client *client) AuthMiddleware(router gin.IRouter, config *MiddlewareConfig) gin.IRouter {
	if config.Realm == "" {
		config.Realm = "unknown"
	}
	if config.GetUserByPhoneNumber == nil {
		config.GetUserByPhoneNumber = client.GetUserByPhoneNumber
	}
	if config.Unauthorized == nil {
		config.Unauthorized = func(c *gin.Context, code int, message string) {
			c.AbortWithStatusJSON(code, map[string]interface{}{
				"code":    -1,
				"message": message,
			})
		}
	}
	if config.NoRoute == nil {
		config.NoRoute = func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]interface{}{
				"code":    -1,
				"message": "page not found",
			})
		}
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       config.Realm,
		Key:         []byte("secret key"),
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*repo.User); ok {
				return jwt.MapClaims{
					identityKey: v.PhoneNumber,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			phoneNumber := claims[identityKey].(string)
			if phoneNumber == adminUser.PhoneNumber {
				return adminUser
			}
			user, _ := config.GetUserByPhoneNumber(phoneNumber)
			return &user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVal repo.User
			if err := c.ShouldBind(&loginVal); err != nil {
				return "", errors.New("missing phone number")
			}
			username := loginVal.Username
			password := loginVal.Password
			phoneNumber := loginVal.PhoneNumber

			if username == defaultUsername && password == defaultPassword {
				return adminUser, nil
			}
			user, err := config.GetUserByPhoneNumber(phoneNumber)
			if err != nil {
				return nil, errors.New("incorrect phone number")
			}
			return &user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: config.Unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		panic("JWT Error:" + err.Error())
	}
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/logout", authMiddleware.LogoutHandler)
	// Refresh time can be longer than token timeout
	router.GET("/refresh_token", authMiddleware.RefreshHandler)
	if eng, ok := router.(*gin.Engine); ok {
		eng.NoRoute(authMiddleware.MiddlewareFunc(), config.NoRoute)
	}
	authGroup := router.Group("")
	authGroup.Use(authMiddleware.MiddlewareFunc())
	return authGroup
}
