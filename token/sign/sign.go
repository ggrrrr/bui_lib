package sign

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ggrrrr/bui_lib/config"
	"github.com/ggrrrr/bui_lib/token"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var signKey *rsa.PrivateKey

func Configure() error {
	keyFile := viper.GetString(token.JWT_KEY_FILE)

	signKeyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println(config.Help())
		return fmt.Errorf("unable to read file: %s param: %v %v", keyFile, token.JWT_KEY_FILE, err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyBytes)
	log.Printf("jwt.private: %+v", signKey.Size())
	if err != nil {
		fmt.Println(config.Help())
		return fmt.Errorf("unable to use:%s: %v", keyFile, token.JWT_KEY_FILE)
	}
	return nil
}

func SignKey(claims token.ApiClaims, ctx context.Context) (string, error) {
	if claims.Subject == "" {
		return "", fmt.Errorf("claim.Subject is missing")
	}
	ttl := token.TokenTTL
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(ttl))
	claims.ID = token.CreateClaimId(ctx)
	log.Printf("Sign: %+v", claims)
	token := jwt.NewWithClaims(jwt.GetSigningMethod(token.SignMethod), claims)
	// token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	tokenString, err := token.SignedString(signKey)
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// ss, err := token.SignedString(signKey)
	// log.Printf("SignKey: %v %v\n", tokenString, err)
	return tokenString, err
}
