package minimal_consul_member

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Activate will register an application check on a certain internal port
// and activate a tiny server on 127.0.0.1:<internalPort>/healthz for checkups
func Activate(internalPort int, consulLocation string, applicationVersion string) (err error) {
	if internalPort == -1 || consulLocation == "" {
		log.Println("Not starting the internal healthz server")
		return
	}
	startTime := time.Now()
	go func() {
		http.Handle("/healthz", newHealthzController(startTime, applicationVersion))
		if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", internalPort), nil); err != nil {
			log.Fatalf("Internal healthz server couldn't be started on port %d, err=%v", internalPort, err)
		}
	}()
	if err = registerOnConsul(internalPort, consulLocation); err != nil {
		return err
	}
	return nil
}
