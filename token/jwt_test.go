package token_test

import (
	"context"
	"testing"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/ggrrrr/bui_lib/token/sign"
	"github.com/golang-jwt/jwt/v4"
)

func TestPasswd1(t *testing.T) {
	const userShit = "userShit"
	ctx := context.WithValue(context.Background(), api.CtxKeyUA, userShit)

	var err error
	// t.Logf("%v", os.Getenv("jwt.key.file"))
	if token.Configure() != nil {
		t.Fatal("cant config")
	}
	if sign.Configure() != nil {
		t.Fatal("cant config")
	}
	// sign.Configure()
	// viver.GetString
	err = token.Configure()
	if err != nil {
		t.Fatalf("%v", err)
	}

	claims := token.ApiClaims{
		Roles: "system",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "username",
		},
	}

	// jwtStd := jwt.StandardClaims{Id: "asd"}
	// _ = jwtStd
	// asd := token.ApiClaims{Groups: "Asd", StandardClaims: jwtStd}
	tokenString, err := sign.SignKey(claims, ctx)
	if err != nil {
		t.Fatalf("cant sign: %v", err)
	}
	t.Logf("tokenString: %v", tokenString)

	v, err := token.Verify(tokenString, ctx)
	if err != nil {
		t.Fatalf("cant verify: %v", err)
	}
	t.Logf("vasd: %+v", v)
	if v.Subject != "username" {
		t.Errorf("unable to no ID")
	}
	if v.Roles != "system" {
		t.Errorf("unable to no Roles")
	}
}
