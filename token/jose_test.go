package token_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/config"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func Test(t *testing.T) {

	token.Configure()
	signKeyBytes, err := ioutil.ReadFile(viper.GetString(token.JWT_KEY_FILE))
	// verifyKeyBytes := viper.GetString(token.JWT_CRT_FILE)
	if err != nil {
		fmt.Println(config.Help())
		t.Fatal(err)
	}

	signKeyPem, err := LoadPrivateKey(signKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	signKey := jose.SigningKey{Algorithm: jose.RS256, Key: signKeyPem}

	signer, err := jose.NewSigner(signKey, (&jose.SignerOptions{}))
	// signer, err := jose.NewSigner(signKey, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		t.Fatal(err)
	}

	cl := jwt.Claims{
		Subject:   "subject",
		Issuer:    "issuer",
		NotBefore: jwt.NewNumericDate(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)),
		Audience:  jwt.Audience{"leela", "fry"},
	}

	jwt, err := jwt.Signed(signer).Claims(cl).CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: %v", jwt)
	// verifyKey, err := LoadPrivateKey([]byte(verifyKeyBytes))
	// t.Logf("verifyKey %v, %+v", err, verifyKey)
	const userShit = "userShit"
	ctx := context.WithValue(context.Background(), api.CtxKeyUA, userShit)
	asd, err := token.Verify(jwt, ctx)
	t.Logf("verifyKey %v, %v", err, asd)

	// sonWebSig, err := jose.ParseSigned(jwt)
	// // t.Logf("verifyKey %v, %+v", err, sonWebSig)
	// t.Logf("verifyKey %v, %v", err, sonWebSig.Signatures[0].Header)
	// t.Logf("verifyKey %v, %v", err, sonWebSig.Signatures[0].Protected)
	// t.Logf("verifyKey %v, %v", err, sonWebSig.Signatures[0].Signature)
	// t.Logf("verifyKey %v, %v", err, sonWebSig.Signatures[0].Unprotected)
	// t.Logf("verifyKey %v, %s", err, sonWebSig.UnsafePayloadWithoutVerification())
}

func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	// jwk, err3 := LoadJSONWebKey(input, false)
	// if err3 == nil {
	// 	return jwk, nil
	// }

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s', '%s' and '%s'", err0, err1, err2, err2)
}

func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s' ", err0, err1)
}
