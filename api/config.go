package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ggrrrr/bui_lib/build"
	"github.com/ggrrrr/bui_lib/config"
	"github.com/spf13/viper"
)

// This is used to set/get user agent metadata to the ctx
type ctxKeyUserAgent string

const (
	LISTEN_ADDR = "listen.addr"

	ctxKeyUA = ctxKeyUserAgent("")
)

var (
	envParamsDefaults = []config.ParamValue{
		{
			Name:     LISTEN_ADDR,
			DefValue: ":8000",
			Info:     "listen address: <IP>:<PORT>, 127.0.0.1:8080, 0.0.0:80",
		},
	}
)

func GetUserAgent(ctx context.Context) *UserAgent {
	v, ok := ctx.Value(ctxKeyUA).(UserAgent)
	if ok {
		return &v
	}
	return &UserAgent{}
}

func SetUserAgent(ctx context.Context, ua UserAgent) context.Context {
	return context.WithValue(ctx, ctxKeyUA, ua)
}

func Configure() error {
	config.Configure(envParamsDefaults)
	addr := viper.GetString(LISTEN_ADDR)
	if addr == "" {
		fmt.Println(config.Help())
		return config.ErrorParamNotSet(LISTEN_ADDR)
	}
	listenAddr = addr
	if strings.Index(addr, ":") == -1 {
		fmt.Printf("%v\n", config.Help())
		return config.ErrorParamInvalid(LISTEN_ADDR, addr)
	}
	host, herr := os.Hostname()
	if herr != nil {
		log.Printf("unable to get hostname: %v", herr)
		host = herr.Error()
	}
	hostName = host
	serverName = fmt.Sprintf("%s/%s", apiName, build.Release)

	return nil
}
