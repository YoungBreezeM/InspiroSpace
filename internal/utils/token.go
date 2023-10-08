package utils

import (
	"encoding/json"
	"math/rand"
)

func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func CreateToken(data interface{}, key string) (token string, err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}

	b2, err := EncryptByAes(b, []byte(key))
	if err != nil {
		return
	}

	token = b2

	return
}

func ParaseToken[T any](token string, key string, s *T) (err error) {

	b2, err := DecryptByAes(token, []byte(key))
	if err != nil {
		return
	}

	if err = json.Unmarshal(b2, s); err != nil {
		return
	}

	return
}
