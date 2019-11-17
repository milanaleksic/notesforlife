package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/milanaleksic/minimal_consul_member"
	"github.com/milanaleksic/notesforlife/dokuwiki"
	"github.com/milanaleksic/notesforlife/dropbox"
)

var (
	// Version is the current version of the program and should be filled in during the build
	Version = "undefined"
)

const PageNameForContents = "dropbox_sync"

func main() {
	err := minimal_consul_member.Activate(config.System.InternalPort, config.System.ConsulLocation, Version, "notes_for_life")
	if err != nil {
		log.Fatalf("Failed to activate BADUC integration: %+v", err)
	}

	wiki := dokuwiki.NewClient(fmt.Sprintf("%s/lib/exe/xmlrpc.php", config.Dokuwiki.URL), config.System.DryRun)
	err = wiki.Login(config.Dokuwiki.Username, config.Dokuwiki.Password)
	if err != nil {
		log.Fatalf("Failed to login to wiki: %+v", err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	client := dropbox.NewClient(config.Dropbox.APIToken, config.Dropbox.Path)
	go client.Process()
	app := &mainApp{
		client:                client,
		signalChannel:         signalChannel,
		wiki:                  wiki,
		contentsData:          newFileContentDetails(),
	}
	app.mainLoop()
}

type mainApp struct {
	contentsData          *fileContentDetails
	client                *dropbox.Client
	signalChannel         chan os.Signal
	wiki                  *dokuwiki.Client
}

func (m *mainApp) mainLoop() {
	contentsDataRaw, err := m.wiki.GetPage(PageNameForContents)
	if err != nil {
		log.Fatalf("Failed to fetch page contents from wiki: %+v", err)
	}
	m.contentsData.buildContentsFromRaw(contentsDataRaw)
	for {
		select {
		case f := <-m.client.ChangedFile:
			if f.Initial {
				log.Printf("initial file: %+v, data length=%v", f.Name, len(f.Data))
				m.updateFileInWiki(f)
			} else {
				log.Printf("changed file: %+v, data length=%v", f.Name, len(f.Data))
				m.updateFileInWiki(f)
			}
		case <-m.signalChannel:
			os.Exit(0)
		}
	}
}

var titleReplacer = regexp.MustCompile("[^a-zA-Z0-9]")

func niceName(s string) string {
	return titleReplacer.ReplaceAllString(strings.ToLower(s), "_")
}

func (m *mainApp) updateFileInWiki(f dropbox.ChangedFile) {
	if strings.TrimSpace(string(f.Data)) == "" {
		return
	}
	pageName := niceName(f.Name)
	if details, ok := m.contentsData.knownPages[pageName]; !ok {
		log.Printf("Adding link to %s (%s) to the contents page", pageName, f.Name)
		m.contentsData.knownPages[pageName] = &contentsData{
			fileName:     f.Name,
			wikiPageName: pageName,
			status:       availableOnDropbox,
		}
	} else {
		details.status = availableOnDropbox
	}
	success, err := m.wiki.PutPage(PageNameForContents, m.contentsData.buildPageFromContents())
	if err != nil {
		log.Fatalf("Failed to store page in wiki: %+v", err)
	}
	if !success {
		log.Println("Failed to store page in wiki")
	}
	currentPage, err := m.wiki.GetPage(pageName)
	if err != nil {
		log.Fatalf("Failed to get page %s from wiki: %+v", pageName, err)
	}
	if string(f.Data) == currentPage {
		log.Printf("Same contents detected for %v, skipping", f.Name)
	} else {
		success, err := m.wiki.PutPage(pageName, string(f.Data))
		if err != nil {
			log.Fatalf("Failed to store page %s in wiki: %+v", pageName, err)
		}
		if !success {
			log.Println("Failed to store page in wiki")
		}
	}
}
