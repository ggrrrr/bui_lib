package token

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ggrrrr/bui_lib/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

const (
	JWT_TTL_MIN     = "jwt.ttl.min"
	JWT_KEY_FILE    = "jwt.key.file"
	JWT_CRT_FILE    = "jwt.crt.file"
	JWT_CA_FILE     = "jwt.ca.file"
	JWT_SIGN_METHOD = "jwt.sign.method"
)

var (

	// openssl genrsa -out jwt.key 512
	// openssl rsa -in jwt.key -pubout > jwt.crt
	verifyKey   *rsa.PublicKey
	SignMethod  string
	TokenTTLMin int

	envParamsDefaults = []config.ParamValue{
		{
			Name:     JWT_TTL_MIN,
			Info:     "JWT ttl in minutes,",
			DefValue: "60",
		},
		{
			Name:     JWT_KEY_FILE,
			Info:     "JWT certificate priviate key file,",
			DefValue: "jwt.key",
		}, {
			Name:     JWT_CRT_FILE,
			Info:     "JWT certificate public key file,",
			DefValue: "jwt.crt",
		},
		{
			Name:     JWT_CA_FILE,
			Info:     "JWT certificate public ca file,",
			DefValue: "ca.crt",
		},
		{
			Name:     JWT_SIGN_METHOD,
			Info:     "JWT signing method RS256, RS512 etc...",
			DefValue: "RS256",
		},
	}
)

func Configure() error {
	var err error
	config.Configure(envParamsDefaults)
	crtFile := viper.GetString(JWT_CRT_FILE)
	SignMethod = viper.GetString(JWT_SIGN_METHOD)
	TokenTTLMin = viper.GetInt(JWT_TTL_MIN)
	// signAlg := jwt.GetAlgorithms()
	if err := jwt.GetSigningMethod(SignMethod); err == nil {
		log.Printf("JWT_SIGN_METHOD: SignMethod %v", err)
		return fmt.Errorf("unkown sign method: %s param: %v", SignMethod, JWT_SIGN_METHOD)
	}

	verifyBytes, err := ioutil.ReadFile(crtFile)
	if err != nil {
		log.Printf("crtFile: %v: %v", crtFile, err)
		return fmt.Errorf("unable to read file: %s param: %v %v", crtFile, JWT_CRT_FILE, err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	log.Printf("jwt.public: %v", verifyKey.Size())
	if err != nil {
		fmt.Println(config.Help())
		return fmt.Errorf("unable to use:%s: %v", crtFile, JWT_KEY_FILE)
	}
	return nil
}

type ApiClaims struct {
	Namespace string `json:"ns"`
	Roles     string `json:"roles"`
	jwt.RegisteredClaims
}
