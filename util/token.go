package util

import "encoding/base64"

type TokenHelper struct {

}


func (tokenHelper TokenHelper) GenerateToken(plainText string) string{
	return base64.StdEncoding.EncodeToString([]byte(plainText))
}

func (tokenHelper TokenHelper) ParseToken(token string) string {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		panic(err)
	}
	return string(data)
}