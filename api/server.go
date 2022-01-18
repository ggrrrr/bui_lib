package api

import (
	"log"
)

func Shutdown() error {
	log.Print("Shutdown...")
	return httpServer.Shutdown(ctx)
}

func Start() error {
	if isStarted.Load() == true {
		log.Printf("httpServer.started!")
		return nil
	}

	isStarted.Store(true)
	log.Printf("httpServer:listenAndServe: %s", httpServer.Addr)
	err := httpServer.ListenAndServe()
	isStarted.Store(false)
	return err
}
