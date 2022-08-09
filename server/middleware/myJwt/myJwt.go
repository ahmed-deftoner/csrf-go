package myjwt

import (
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/id_rsa"
	pubKeyPath  = "keys/id_rsa.pub"
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

func CreateNewTokens() {

}

func CheckAndRefreshTokens() {

}

func CreateAuthTokenString() {

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
