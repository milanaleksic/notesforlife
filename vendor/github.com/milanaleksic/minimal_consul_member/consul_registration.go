package minimal_consul_member

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type agentServiceRegistration struct {
	ID                string   `json:",omitempty"`
	Name              string   `json:",omitempty"`
	Tags              []string `json:",omitempty"`
	Port              int      `json:",omitempty"`
	Address           string   `json:",omitempty"`
	EnableTagOverride bool     `json:",omitempty"`
	Checks            []*agentServiceCheck
}

type agentServiceCheck struct {
	Script            string `json:",omitempty"`
	DockerContainerID string `json:",omitempty"`
	Shell             string `json:",omitempty"` // Only supported for Docker.
	Interval          string `json:",omitempty"`
	Timeout           string `json:",omitempty"`
	TTL               string `json:",omitempty"`
	HTTP              string `json:",omitempty"`
	TCP               string `json:",omitempty"`
	Status            string `json:",omitempty"`
	Notes             string `json:",omitempty"`
	TLSSkipVerify     bool   `json:",omitempty"`
}

func registerOnConsul(internalPort int, consulLocation string) (err error) {
	if consulLocation == "" {
		return
	}
	registration, err := json.Marshal(&agentServiceRegistration{
		Name: "notes_for_life",
		Port: internalPort,
		Checks: []*agentServiceCheck{
			{
				Interval: "30s",
				HTTP:     fmt.Sprintf("http://127.0.0.1:%d/healthz", internalPort),
			},
		},
	})
	if err != nil {
		return errors.New(fmt.Sprintf("could not register to consul, err=%v", err))
	}
	request, err := http.NewRequest("PUT", consulLocation+"/v1/agent/service/register", bytes.NewBuffer(registration))
	if err != nil {
		return errors.New(fmt.Sprintf("could not register to consul, err=%v", err))
	}
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("could not register to consul, err=%v", err))
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("could not register to consul, received status code=%v", response.StatusCode))
	}
	return nil
}
