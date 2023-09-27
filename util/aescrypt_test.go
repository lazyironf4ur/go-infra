package util

import (
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	a := AesCrypt{}
	cipherByte, err := a.Encrypt("hahahaha")
	if err != nil {
		panic(err)
	}

	fmt.Printf("加密后字符串：%s\n", string(cipherByte))

	plainText, err2 := a.Decrypt(cipherByte)
	if err2 != nil {
		panic(err)
	}

	fmt.Printf("解密后字符串: %s\n", plainText)
}