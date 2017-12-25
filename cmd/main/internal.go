package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type healthzController struct {
	version   string
	startTime time.Time
}

func newHealthzController(startTime time.Time, version string) *healthzController {
	return &healthzController{
		version:   version,
		startTime: startTime,
	}
}

func internalHandlers(internalPort int) {
	if internalPort == -1 {
		log.Println("Not starting the internal healthz server")
		return
	}
	startTime := time.Now()
	go func() {
		http.Handle("/healthz", newHealthzController(startTime, Version))
		if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", internalPort), nil); err != nil {
			log.Fatalf("Internal handling server couldn't be started on port %d, err=%v", internalPort, err)
		}
	}()
}

func (w *healthzController) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if response, err := json.Marshal(&struct {
		Self          string
		SelfVersion   string
		TimeUpSeconds int
	}{
		Self:          "ok",
		SelfVersion:   w.version,
		TimeUpSeconds: int(time.Since(w.startTime).Seconds()),
	}); err != nil {
		log.Printf("Could not write healthz status, err=%s", err)
	} else if _, err := wr.Write(response); err != nil {
		log.Printf("Could not write healthz status, err=%s", err)
	}
}
