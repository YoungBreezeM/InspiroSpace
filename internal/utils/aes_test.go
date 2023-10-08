package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {

	a := struct {
		OpenId string `json:"openid"`
		OK     string `json:"ok"`
	}{
		OpenId: "hello",
		OK:     "world",
	}

	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	key := "ABCDABCDABCDABCDABCDABCDABCDABCD"
	b2, err := AesEncrypt(b, []byte(key))
	if err != nil {
		panic(err)
	}

	b3, err := AesDecrypt(b2, []byte(key))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b3))

}
