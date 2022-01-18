package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type HTTP_HEADER_ string
type HTTP_CONTENT_TYPE_ string

const (
	HTTP_CONTENT_TYPE HTTP_HEADER_       = "Content-Type"
	HTTP_CT_JSON      HTTP_CONTENT_TYPE_ = "application/json"
	HTTP_CT_AUTH      HTTP_CONTENT_TYPE_ = "Authorization"
	HTTP_AUTH_BEARER  HTTP_CONTENT_TYPE_ = "Bearer"
)

type ModulesInfo struct {
	Path    string `json:"path"`    // module path
	Version string `json:"version"` // module version
	Sum     string `json:"sum"`     // checksum
	Replace string `json:"rep"`     // replaced by this module
}

type ApiVersion struct {
	ApiName    string        `json:"apiName"`
	ApiVersion string        `json:"apiVersion"`
	HostName   string        `json:"hostName"`
	BuildInfo  interface{}   `json:"buildInfo"`
	Modules    []ModulesInfo `json:"Modules"`
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	SetResponseHeader(w)
	err := json.NewEncoder(w).Encode(ApiVersion{
		ApiName:    apiName,
		ApiVersion: buildInfo.Main.Version,
		HostName:   hostName,
		BuildInfo:  fmt.Sprintf("version: %s build: %s", Version(), BuildInfo()),
		Modules:    modules,
	})
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

}

func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handle404() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		newHttpError(
			w, http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			fmt.Errorf("path: %s", r.URL))
	})
}
