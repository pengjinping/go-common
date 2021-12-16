package oauth2

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"git.kuainiujinke.com/oa/oa-common-golang/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
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

func ParseToken(r *http.Request) (*jwt.Token, error) {

	return request.ParseFromRequest(r, TokenExtractor{}, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
}
