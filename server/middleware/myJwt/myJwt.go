package myjwt

import "errors"

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

func JWTInit() error {
	return errors.New("error")
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
