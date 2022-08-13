package db

import "github.com/ahmed-deftoner/csrf-go/db/models"

var users = map[string]models.User{}

func DBInit() {

}

func StoreUser(username string, password string, role string) (uuid string, err error) {

}

func DeleteUser(uuid string) {

}

func FetchUserById(uuid string) (models.User, error) {

}

func FetchUserByUsername(username string) (models.User, string, error) {

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
