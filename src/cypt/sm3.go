package cypt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm3"
)

// sm3 hash算法
func Sm3(data string) (string, error) {
	h := sm3.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	return string(sum), nil
}

// sm3 hash算法
func Sm3Hmac111(data, key string) (string, error) {
	dataByte := []byte(data)
	keyByte, err2 := hex.DecodeString(key)
	if err2 != nil {
		return "", err2
	}
	hash := hmac.New(sm3.New, keyByte)
	_, err := hash.Write(dataByte)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return str, nil
}

// sm3 hash算法
func Sm3Hmac(data, key string) (string, error) {
	dataByte := []byte(data)
	keyByte := []byte(key)
	hash := hmac.New(sm3.New, keyByte)
	_, err := hash.Write(dataByte)
	if err != nil {
		return "", err
	}
	str := hex.EncodeToString(hash.Sum(nil))
	return str, nil
}
