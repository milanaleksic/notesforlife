package minimal_consul_member

import (
	"encoding/json"
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
