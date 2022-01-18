package api

import (
	"log"
	"net/http"
	"net/url"

	"github.com/mssola/user_agent"
)

func addCorsHeader(res http.ResponseWriter) {
	// AllowedHeaders:   []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", "Authorization"},

	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	// headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET,POST")
}

func parseRequestHeaders(r *http.Request) UserAgent {
	uaHeader := r.Header.Get("User-Agent")
	ua := user_agent.New(uaHeader)
	e1, _ := ua.Engine()
	b1, _ := ua.Browser()
	o1 := ua.OSInfo().Name
	userAgent := UserAgent{
		Mobile:  ua.Mobile(),
		Engine:  e1,
		Browser: b1,
		OS:      o1,
		Bot:     ua.Bot(),
	}

	if len(r.Header["Origin"]) > 0 {
		url, err := url.ParseRequestURI(r.Header["Origin"][0])
		userAgent.OriginHost = url.Host
		log.Printf("\tOrigin: \t%v", r.Header["Origin"])
		log.Printf("\turl: \t%+v: %v", url, err)
	}

	return userAgent
}
