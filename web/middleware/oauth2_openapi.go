package middleware

import (
	"fmt"

	"git.kuainiujinke.com/oa/oa-common-golang/logger"
	"git.kuainiujinke.com/oa/oa-common-golang/utils/oauth2"
	"git.kuainiujinke.com/oa/oa-common-golang/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func OAuthCodeVerify() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := oauth2.ParseToken(c.Request)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		if err != nil || !token.Valid {
			errmsg = "the requested JWT is invalid. " + errmsg
			logger.Error(c, errmsg)
			web.FailWithMessage(errmsg, c)
			c.Abort()
			return
		}

		// 写入授权认证信息
		claims := token.Claims.(jwt.MapClaims)
		c.Set("employee_id", fmt.Sprintf("%s", claims["sub"]))
		c.Set("grant_type", oauth2.GrantTypeAuthorizationCode)
		c.Set("scopes", claims["scopes"])

		c.Next()
	}

}

func OAuthClientVerify() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := oauth2.ParseToken(c.Request)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		if err != nil || !token.Valid {
			errmsg = "the requested JWT is invalid. " + errmsg
			logger.Error(c, errmsg)
			web.FailWithMessage(errmsg, c)
			c.Abort()
			return
		}

		// 写入授权认证信息
		claims := token.Claims.(jwt.MapClaims)
		// todo: client_id
		c.Set("client_id", fmt.Sprintf("%s", claims["sub"]))
		c.Set("grant_type", oauth2.GrantTypeClientCredentials)
		c.Set("scopes", claims["scopes"])

		c.Next()
	}

}
