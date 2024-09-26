package common

import (
	"math/rand"
	"time"
)

// 初始化随机数种子
func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

// 获取随机化字符串
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
