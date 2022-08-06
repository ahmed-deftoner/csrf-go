package randomstrings

import "encoding/base64"

func GenerateRandomBytes(n int) ([]byte, error) {

}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
