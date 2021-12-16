package oauth2

import (
	"crypto/rsa"
	"fmt"
	"git.kuainiujinke.com/oa/oa-common-golang/logger"
	"git.kuainiujinke.com/oa/oa-common-golang/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

const (
	GrantTypeAuthorizationCode   = "AuthorizationCode"
	GrantTypeImplict             = "Implicit"
	GrantTypeClientCredentials   = "ClientCredentials"
	GrantTypePasswordCredentials = "PasswordCredentials"
)

func InitJWT(jwtPublicKeyPath string) {
	var err error
	var pubContent []byte
	if pubContent, err = ioutil.ReadFile(jwtPublicKeyPath); err != nil {
		logger.Error(&gin.Context{}, "can not read JWT public key file")
		return
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubContent)
	return
}

type TokenExtractor struct {
}

func (t TokenExtractor) ExtractToken(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.FormValue("access_token")
	}
	if token == "" {
		return "", request.ErrNoTokenInRequest
	}
	return strings.TrimSpace(strings.TrimPrefix(token, "Bearer")), nil
}

func parseToken(r *http.Request) (*jwt.Token, error) {

	return request.ParseFromRequest(r, TokenExtractor{}, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
}

func AuthCodeVerify() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := parseToken(c.Request)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		if err != nil || !token.Valid {
			errmsg = "the requested JWT is invalid" + errmsg
			//logger.Error(&gin.Context{}, errmsg)
			web.FailWithMessage(errmsg, c)
			c.Abort()
			return
		}

		// 写入授权认证信息
		claims := token.Claims.(jwt.MapClaims)
		c.Set("employee_id", fmt.Sprintf("%s", claims["sub"]))
		c.Set("grant_type", GrantTypeAuthorizationCode)
		c.Set("scopes", claims["scopes"])

		c.Next()
	}

}

func ClientVerify() gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := parseToken(c.Request)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		if err != nil || !token.Valid {
			errmsg = "the requested JWT is invalid" + errmsg
			logger.Error(&gin.Context{}, errmsg)
			web.FailWithMessage(errmsg, c)
			//c.JSON(http.StatusOK, gin.H{
			//	"code":    1,
			//	"message": errmsg,
			//	"data":    nil,
			//})
			c.Abort()
			return
		}

		// 写入授权认证信息
		claims := token.Claims.(jwt.MapClaims)
		// todo: client_id
		c.Set("client_id", fmt.Sprintf("%s", claims["sub"]))
		c.Set("grant_type", GrantTypeClientCredentials)
		c.Set("scopes", claims["scopes"])

		c.Next()
	}

}
