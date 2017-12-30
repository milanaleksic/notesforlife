package main

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

var config struct {
	System struct {
		InternalPort   int
		ConsulLocation string
	}
	Dropbox struct {
		APIToken string
		Path     string
	}
	Dokuwiki struct {
		URL      string
		Username string
		Password string
	}
}

func init() {
	if _, err := toml.DecodeFile(path.Join(path.Dir(os.Args[0]), "notesforlife.toml"), &config); err != nil {
		if _, err := toml.DecodeFile("notesforlife.toml", &config); err != nil {
			if _, err := toml.DecodeFile("../application.toml", &config); err != nil {
				log.Fatalf("Failure while parsing configuration: %v", err)
			}
		}
	}
}
