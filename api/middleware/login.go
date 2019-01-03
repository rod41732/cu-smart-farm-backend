package middleware

import (
	"log"
	"time"

	"github.com/rod41732/cu-smart-farm-backend/storage"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/common"
	"gopkg.in/mgo.v2/bson"
)

//UserAuth is middleware for authenticating user
var UserAuth *jwt.GinJWTMiddleware

var identityKey = "user"

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// for marshalled db
type userData struct {
	Username string   `json:"username"`
	Devices  []string `json:"devices"`
}

func Initialize() {
	var err error
	UserAuth, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "CUSmartFarm",
		Key:         common.SignKey,
		Timeout:     time.Hour * 99999,
		MaxRefresh:  time.Hour * 99999,
		IdentityKey: identityKey,
		// ------------------------ creation of JWT token --------------------
		// handle auth via request and return user when success
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := common.SHA256(loginVals.Password)

			mdb, err := common.Mongo()
			if common.PrintError(err) {
				return nil, jwt.ErrFailedAuthentication
			}

			query := mdb.DB("CUSmartFarm").C("users").Find(bson.M{
				"username": username,
				"password": password,
			})
			if cnt, err := query.Count(); cnt == 0 || err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			_ = storage.GetUserStateInfo(username) // create user
			return username, nil
		},
		// this create JWT, and can be retrieved with extract claims
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(string); ok {
				return jwt.MapClaims{
					"username": v,
				}
			}
			return jwt.MapClaims{}
		},
		// --------------- end of creation ----------------------------
		// ---------------- handling JWT in request -------------------
		// extracts claims which is set from PayloadFunc
		// and will set into c context via c.Set(identityKey)
		// which can retrieve on endPoint
		IdentityHandler: func(c *gin.Context) interface{} { // return user object if success
			claims := jwt.ExtractClaims(c)
			// token := jwt.GetToken(c)
			mdb, err := common.Mongo()
			if err != nil {
				return nil
			}

			// Need to check with current token (Limit 1 device login)
			var user map[string]interface{}
			err = mdb.DB("CUSmartFarm").C("users").Find(bson.M{
				"username": claims["username"].(string),
				// "token":    token, // disable token check
			}).One(&user)
			if err != nil {
				return nil
			}
			return claims["username"].(string)
		},
		// handle whether we should allow
		Authorizator: func(data interface{}, c *gin.Context) bool {
			_, ok := data.(string)
			return ok
		},
		// is called when unauthorized
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
			c.Abort()
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
