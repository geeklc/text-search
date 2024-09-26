package cypt

import (
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm4"
)

// sm4加密
func Sm4Enc(plainText, key string) (string, error) {
	srcByte := []byte(plainText)
	keyByte, err2 := hex.DecodeString(key)
	if err2 != nil {
		return "", err2
	}
	encrypt, err := sm4.Sm4Ecb(keyByte, srcByte, true)
	if err != nil {
		return "", err
	}
	data := hex.EncodeToString(encrypt)
	return data, nil
}

// sm4解密
func Sm4Dec(cipherText, key string) (string, error) {
	srcByte, err1 := hex.DecodeString(cipherText)
	if err1 != nil {
		return "", err1
	}
	keyByte, err2 := hex.DecodeString(key)
	if err2 != nil {
		return "", err2
	}
	dec, err := sm4.Sm4Ecb(keyByte, srcByte, false)
	if err != nil {
		return "", err
	}
	data := string(dec)
	return data, nil
}
