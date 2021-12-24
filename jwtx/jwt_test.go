package jwtx

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	Init(
		SetSecret("hello"),
		SetIssuer("fjsd"),
		SetExpire(time.Hour*200),
		//SetMethod(jwt.SigningMethodEdDSA),
		SetSubject("hello ju"),
	)
	generate, err := Generate(map[string]interface{}{
		"hello": "hello",
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("token: %s", generate)
}
func TestParseToken(t *testing.T) {
	Init()

	generate, err := Generate(map[string]interface{}{
		"hello": "hello",
	})
	if err != nil {
		t.FailNow()
	}
	fmt.Println("token: ", generate)
	payload, err := ParseToken(generate)
	if err != nil {
		t.FailNow()
		return
	}
	fmt.Println(payload)
}
