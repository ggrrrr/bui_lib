package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ggrrrr/bui_lib/api"
)

func main() {
	err := api.Configure()
	if err != nil {
		log.Fatalf(err.Error())
	}
	root := context.Background()
	err = api.Create(root, false)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// err = api.NewApi(root, false)
	go func() {
		err := api.Start()
		if err != nil {
			log.Printf("http error: %+v", err)
		}
	}()
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals)
	log.Printf("os.signal: %v", <-osSignals)
	defer api.Shutdown()

}
