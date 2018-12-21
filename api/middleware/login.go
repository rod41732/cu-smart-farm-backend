package middleware

import (
	"log"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"
)

/// UserAuth is middleware for authenticating user
var UserAuth *jwt.GinJWTMiddleware

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// we use user struct just to distinguish it from string when passing data
type User struct {
	Username string
}

var identityKey = "username"

func Initialize() {
	var err error
	UserAuth, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "CUSmartFarm",
		Key:         common.SignKey,
		Timeout:     time.Hour * 99999,
		MaxRefresh:  time.Hour * 99999,
		IdentityKey: "username",
		// ------------------------ creation of JWT token --------------------
		// handle auth via request and return some info when success
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := common.SHA256(loginVals.Password)

			// auth using mongo
			mdb, err := common.Mongo()
			if common.PrintError(err) {
				return nil, jwt.ErrFailedAuthentication
			}
			var user map[string]interface{}
			col := mdb.DB("CUSmartFarm").C("users")
			err = col.Find(bson.M{
				"username": username,
				"password": password,
			}).One(&user)

			if err != nil || user == nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &User{
				Username: username,
			}, nil
		},
		// if success (can convert data to *User) => put data in to claims
		// which can be retreived by extract claims
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		// --------------- end of creation ----------------------------
		// ---------------- handling JWT in request -------------------
		// extracts claims which is set from PayloadFunc
		// and will set into c context via c.Set(identityKey)
		// which can retrieve on endPoint
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Username: claims[identityKey].(string),
			}
		},
		// handle whether we should allow
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*User); ok {
				return true
			}
			return false
		},
		// is called when unauthorized
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:    "header: Authorization, query: token, cookie: token",
		TokenHeadName:  "Bearer",
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false, //non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		CookieDomain:   "127.0.0.1",
		CookieName:     "token", // default jwt
	})

	if err != nil {
		log.Fatal(err)
	}
}
