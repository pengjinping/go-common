package web

import (
	"fmt"
	"net/http"
	"oa-common/config"
	"oa-common/logger"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func ParseToken(r *http.Request, secret interface{}) (*jwt.Token, error) {
	var token *jwt.Token
	var err error
	token, err = doParseToken(r, secret)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"X-Client-Token"},
	Filter:    stripBearerPrefixFromTokenString,
}

func doParseToken(r *http.Request, secret interface{}) (*jwt.Token, error) {
	return request.ParseFromRequest(r, AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		}, request.WithParser(newParser()))
}

func newParser() *jwt.Parser {
	return &jwt.Parser{
		UseJSONNumber: true,
	}
}

// api请求认证
func ApiAuthMiddleware() gin.HandlerFunc {
	sk := []byte(config.GetString("jwt.signing-key"))
	return func(c *gin.Context) {
		tok, err := ParseToken(c.Request, sk)

		if err != nil {
			logger.Error(c, err.Error())
			FailWithMessage(http.StatusUnauthorized, "登录token无效", c)
			c.Abort()
			return
		}
		if tok == nil || !tok.Valid {
			FailWithMessage(http.StatusUnauthorized, "登录token无效", c)
			c.Abort()
			return
		}
		claims := tok.Claims.(jwt.MapClaims)
		//todo 如果在blacklist里，禁止登录

		c.Set("claims", claims)
		c.Set("email", claims["email"])
		c.Set("userId", fmt.Sprintf("%v", claims["sub"]))
		c.Next()
	}

}
