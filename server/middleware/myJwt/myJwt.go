package myjwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/ahmed-deftoner/csrf-go/db"
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

func CheckAndRefreshTokens(oldAuthTokenString string, oldRefreshTokenString string, oldCsrfSecret string) (newAuthTokenString, newRefreshTokenString, newCsrfSecret string, err error) {
	if oldCsrfSecret == "" {
		log.Println("No CSRF token!")
		err = errors.New("Unauthorized")
		return
	}
	authToken, err := jwt.ParseWithClaims(oldAuthTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}
	if oldCsrfSecret != authTokenClaims.Csrf {
		log.Println("CSRF token doesn't match jwt!")
		err = errors.New("Unauthorized")
		return
	}
	if authToken.Valid {
		log.Println("Auth token is valid")

		newCsrfSecret = authTokenClaims.Csrf

		newRefreshTokenString, err = updateRefreshTokenExp(oldRefreshTokenString)
		newAuthTokenString = oldAuthTokenString
		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		log.Println("Auth token is not valid")
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			log.Println("Auth token is expired")

			newAuthTokenString, newCsrfSecret, err = updateAuthTokenString(oldRefreshTokenString, oldAuthTokenString)
			if err != nil {
				return
			}

			newRefreshTokenString, err = updateRefreshTokenExp(oldRefreshTokenString)
			if err != nil {
				return
			}

			newRefreshTokenString, err = updateRefreshTokenCsrf(newRefreshTokenString, newCsrfSecret)
			return
		} else {
			log.Println("Error in auth token")
			err = errors.New("Error in auth token")
			return
		}
	} else {
		log.Println("Error in auth token")
		err = errors.New("Error in auth token")
		return
	}

	err = errors.New("Unauthorized")
	return
}

func createAuthTokenString(uuid string, role string, csrfSecret string) (authTokenString string, err error) {
	AuthTokenExp := time.Now().Add(models.AuthTokenValidTime).Unix()
	authClaims := models.TokenClaims{
		jwt.StandardClaims{
			Subject:   uuid,
			ExpiresAt: AuthTokenExp,
		},
		role,
		csrfSecret,
	}
	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)

	authTokenString, err = authJwt.SignedString(signKey)
	return
}

func createRefreshTokenString(uuid string, role string, csrfString string) (refreshTokenString string, err error) {
	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()
	refreshJti, err := db.StoreRefreshToken()
	if err != nil {
		return
	}

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        refreshJti,
			Subject:   uuid,
			ExpiresAt: refreshTokenExp,
		},
		role,
		csrfString,
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	refreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func updateRefreshTokenExp(oldRefreshTokenString string) (newRefreshTokenString string, err error) {
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}
	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        oldRefreshTokenClaims.StandardClaims.Id, // jti
			Subject:   oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: refreshTokenExp,
		},
		oldRefreshTokenClaims.Role,
		oldRefreshTokenClaims.Csrf,
	}

	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func updateAuthTokenString(refreshTokenString string, oldAuthTokenString string) (newAuthTokenString, csrfSecret string, err error) {

}

func RevokeRefreshToken(refreshTokenString string) error {

}

func updateRefreshTokenCsrf(oldRefreshTokenString string, newCsrfString string) (newRefreshTokenString string, err error) {

}

func GrabUUID(authTokenString string) (string, error) {
	authToken, _ := jwt.ParseWithClaims(authTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return "", errors.New("Error fetching claims")
	})
	authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
	if !ok {
		return "", errors.New("Error fetching claims")
	}

	return authTokenClaims.StandardClaims.Subject, nil
}
