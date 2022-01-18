package token

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/golang-jwt/jwt/v4"
)

func CreateClaimId(ctx context.Context) string {
	userAgent := fmt.Sprintf("%v", api.GetUserAgent(ctx))
	sha := sha256.New()
	sha.Write([]byte(userAgent))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func Verify(inputToken string, ctx context.Context) (*ApiClaims, error) {
	out, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		log.Printf("JWT error: %s !", err)
		return nil, err
	}
	jwtAlg := out.Header["alg"]
	if SignMethod != jwtAlg {
		log.Printf("JWT modified, ALG does not match. possible ATTACK JWT.alg: %s !", jwtAlg)
		return nil, fmt.Errorf("bad alg in the JWT")
	}
	// claimId := CreateClaimId(ctx)
	// if out.

	claims, ok := out.Claims.(jwt.MapClaims)
	if !out.Valid {
		return nil, fmt.Errorf("not valid")
	}
	log.Printf("%v ", out.Valid)
	log.Printf("%v ", claims)
	roles := (claims["roles"]).(string)
	sub := (claims["sub"]).(string)
	if !ok {
		return nil, fmt.Errorf("not able to get claims")
	}
	if claims["jti"] == "" {
		return nil, fmt.Errorf("missing claims.jti")
	}
	if claims["jti"] != CreateClaimId(ctx) {
		return nil, fmt.Errorf("wrong  claims.jti")
	}
	if sub == "" {
		return nil, fmt.Errorf("missing claims.sub")
	}
	if roles == "" {
		return nil, fmt.Errorf("missing claims.roles")
	}

	newApi := ApiClaims{
		Roles:            roles,
		RegisteredClaims: jwt.RegisteredClaims{Subject: sub},
	}
	return &newApi, nil
}
