package db

import (
	"errors"

	"github.com/ahmed-deftoner/csrf-go/db/models"
)

var users = map[string]models.User{}

var refreshTokens map[string]string

func DBInit() {
	refreshTokens = make(map[string]string)
}

func StoreUser(username string, password string, role string) (uuid string, err error) {

}

func DeleteUser(uuid string) {

}

func FetchUserById(uuid string) (models.User, error) {

}

func FetchUserByUsername(username string) (models.User, string, error) {
	for k, v := range users {
		if v.Username == username {
			return v, k, nil
		}
	}

	return models.User{}, "", errors.New("User not found that matches given username")
}

func StoreRefreshToken() (jti string, err error) {

}

func DeleteRefreshToken(jti string) {

}

func CheckRefreshToken(jti string) bool {

}

func LogUserIn(username string, password string) (models.User, string, error) {

}

func generateBcryptHash(password string) (string, error) {

}

func checkPasswordAgainstHash(hash string, password string) error {

}
