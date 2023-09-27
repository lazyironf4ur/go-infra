package util

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"strconv"
)

// AES 规定有3种长度的key: 16, 24, 32分别对应AES-128, AES-192, or AES-256
var default_key = defaultKey(32) 

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

type AesCrypt struct {
	Key []byte
}

// 参照 https://www.cnblogs.com/ssgeek/p/15721816.html#23-des%E5%AD%90%E5%AF%86%E9%92%A5%E7%94%9F%E6%88%90实现
func (aesc AesCrypt) Encrypt (plainText string) (cipherByte []byte, err error){
	key := default_key
	// if aesc.Key != "" {
	// 	key = aesc.Key
	// }

	plainByte := []byte(plainText)
	keyByte := []byte(key)
	// 创建加密算法aes
	c, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	cipherByte = make([]byte, len(plainByte))
	cfb.XORKeyStream(cipherByte, plainByte)
	return
}

func (aesc AesCrypt) Decrypt(cipherByte []byte) (plainText string, err error) {
	key := default_key
	// if aesc.Key != "" {
	// 	key = aesc.Key
	// }
	// 转换成字节数据, 方便加密
	keyByte := []byte(key)
	// 创建加密算法aes
	c, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plainByte := make([]byte, len(cipherByte))
	cfbdec.XORKeyStream(plainByte, cipherByte)
	plainText = strconv.QuoteToASCII(string(plainByte))
	return
}

func defaultKey (size int) []byte {
	key := make([]byte, size)
	for i := 0; i < size; i++ {
		key[i] = byte(rand.Intn(127))
	}
	return key
}	