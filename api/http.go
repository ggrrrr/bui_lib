package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"

	"github.com/gorilla/mux"
)

var (
	apiName    string
	hostName   string
	serverName string
	isStarted  *atomic.Value = &atomic.Value{}
	isReady    *atomic.Value
	ctx        context.Context

	httpServer *http.Server
	httpRouter *mux.Router

	// Channel used to notify that the HTTP server has stoped
	httpShutDown = make(chan struct{}, 1)
	// Channel used to send shutothat the HTTP server has stoped
	// interrupt    = make(chan os.Signal, 1)
	buildInfo  *debug.BuildInfo
	modules    []ModulesInfo
	listenAddr string
)

func Ready() {
	log.Printf("ready status: OK")
	isReady.Store(true)
}

func Name() string {
	return apiName
}

func HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	log.Printf("handler.register: %s %s", path, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
	return httpRouter.HandleFunc(path, f)
}

type UserAgent struct {
	Mobile     bool
	Engine     string
	Browser    string
	OS         string
	Bot        bool
	OriginHost string
}

func (o UserAgent) String() string {
	mob := ""
	if o.Mobile {
		mob = "/mobile"
	}
	return fmt.Sprintf("OS:%s/%s%s", o.OS, o.Browser, mob)
}

func preHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addCorsHeader(w)
		for k, v := range r.Header {
			log.Printf("preHandler: \t%v: %v", k, v)

		}
		userAgent := parseRequestHeaders(r)
		ctx := SetUserAgent(r.Context(), userAgent)
		log.Printf("preHandler %v[%s]: %+v", r.RequestURI, r.Method, userAgent.String())
		if r.Method == "OPTIONS" {
			//handle preflight in here
			// log.Printf("corsHandler[%s]: %+v", r.Method, r.Header)
			// 	// log.Printf("corsHandler: %+v", r.Header)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			// h.ServeHTTP(w, r)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Create(c context.Context, ready bool) error {
	if isStarted.Load() != nil {
		return fmt.Errorf("aready initialized")
	}
	ctx = c
	isStarted.Store(false)
	isReady = &atomic.Value{}
	isReady.Store(ready)
	readBuildInfo()
	httpRouter = mux.NewRouter()

	// httpRouter.HandleFunc("/version", versionHandler)
	// httpRouter.Handle("/version", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(versionHandler)))
	HandleFunc("/version", versionHandler)
	HandleFunc("/healthz", readyz(isReady))
	HandleFunc("/readyz", readyz(isReady))

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://foo.com", "http://foo.com:8080"},
	// 	AllowCredentials: true,
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: true,
	// })

	// mycor := cors.New(
	// 	cors.Options{
	// 		AllowedOrigins: []string{"*"},
	// 		// AllowedMethods:   []string{"GET", "POST","OPTIONS"},
	// 		AllowedHeaders:   []string{"*"},
	// 		AllowCredentials: true,
	// 		Debug:            false,
	// 	},
	// )

	httpRouter.NotFoundHandler = handle404()
	// handlerCors1 := mycor.Handler(httpRouter)
	// handlerCors := preHandler(handlerCors1)
	handlerCors := preHandler(httpRouter)
	httpServer = &http.Server{
		Addr:    listenAddr,
		Handler: handlerCors,
		// Handler: httpRouter,
		// Handler: handlers.CompressHandler(httpRouter),
	}
	return nil
}

func SetResponseHeader(w http.ResponseWriter) {
	// http.Header
	w.Header().Set(string(HTTP_CONTENT_TYPE), string(HTTP_CT_JSON))
	// w.Header().Set("Server", "somecoolserver")
}

func GetAuthorizationBearer(r *http.Request) (string, error) {
	// http.Header
	authHeader := r.Header.Get(string(HTTP_CT_AUTH))
	parts := strings.Split(authHeader, string(HTTP_AUTH_BEARER))
	if len(parts) != 2 {
		log.Printf("GetAuthorizationBearer: %v", r.Header)
		return "", fmt.Errorf("header %s %s not found", HTTP_CT_AUTH, HTTP_AUTH_BEARER)
	}
	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return "", fmt.Errorf("header %s %s empty", HTTP_CT_AUTH, HTTP_AUTH_BEARER)
	}
	return token, nil
}
