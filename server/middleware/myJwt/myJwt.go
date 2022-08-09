package myjwt

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/ahmed-deftoner/csrf-go/db/models"
	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/id_rsa"
	pubKeyPath  = "keys/id_rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func JWTInit() error {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

func CreateNewTokens(uuid string, role string) (authTokenString, refreshTokenString, csrfSecret string, err error) {
	csrfSecret, err = models.GenerateCSRFSecret()
	if err != nil {
		return
	}

	refreshTokenString, err = createRefreshTokenString(uuid, role, csrfSecret)

	authTokenString, err = createAuthTokenString(uuid, role, csrfSecret)
	if err != nil {
		return
	}

	return
}

func CheckAndRefreshTokens() {

}

func createAuthTokenString(uuid string, role string, csrfSecret string) (authTokenString string, err error) {

}

func createRefreshTokenString(uuid string, role string, csrfString string) (refreshTokenString string, err error) {

}

func updateRefreshTokenExp(oldRefreshTokenString string) (newRefreshTokenString string, err error) {

}

func updateAuthTokenString(refreshTokenString string, oldAuthTokenString string) (newAuthTokenString, csrfSecret string, err error) {

}

func RevokeRefreshToken(refreshTokenString string) error {

}

func updateRefreshTokenCsrf(oldRefreshTokenString string, newCsrfString string) (newRefreshTokenString string, err error) {

}

func GrabUUID(authTokenString string) (string, error) {

}
