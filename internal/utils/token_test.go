package utils

import (
	"fmt"
	"strings"
	"testing"
)

type Token struct {
	OpenId string
	Key    string
}

func TestToken(t *testing.T) {

	token := "4gFntQLPDu8 lmDRfGPwMT8mquCM4UJPpJH1gSQjLKCvVUHPgaNNHEVcjfUu9HogHqVMyoUmaYr9JOE0TbKUOeV3CLOkryVFaE7F/OpJXwtsaw8o3IeElffabqMB6DUp"
	s := strings.Replace(token, " ", "+", 10)
	r := Token{}
	if err := ParaseToken[Token](s, "cg20TmaktPuWGu4A", &r); err != nil {
		panic(err)
	}
	//4gFntQLPDu8 lmDRfGPwMT8mquCM4UJPpJH1gSQjLKCvVUHPgaNNHEVcjfUu9HogHqVMyoUmaYr9JOE0TbKUOeV3CLOkryVFaE7F/OpJXwtsaw8o3IeElffabqMB6DUp
	//SzzyMBffCVvYq7QiqnD1I4qe7KeUKONp

	fmt.Println(r)
}

func TestRadmon(t *testing.T) {
	s := GenerateRandomString(32)
	fmt.Println(s)
}
