package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/milanaleksic/notesforlife/dokuwiki"
	"github.com/milanaleksic/notesforlife/dropbox"
)

func main() {
	token := flag.String("token", "", "API token for Dropbox")
	path := flag.String("path", "", "Path to track")
	username := flag.String("username", "", "Dokuwiki username")
	password := flag.String("password", "", "Dokuwiki password")
	wikiLocation := flag.String("wikiLocation", "", "Dokuwiki location")
	flag.Parse()

	wiki := dokuwiki.NewClient(fmt.Sprintf("%s/lib/exe/xmlrpc.php", *wikiLocation))
	err := wiki.Login(*username, *password)
	if err != nil {
		log.Fatalf("Failed to login to wiki: %+v", err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	client := dropbox.NewClient(*token, *path)
	go client.Process()
	app := &mainApp{
		client:        client,
		signalChannel: signalChannel,
		wiki:          wiki,
	}
	app.mainLoop()
}

type mainApp struct {
	contentsData  string
	client        *dropbox.Client
	signalChannel chan os.Signal
	wiki          *dokuwiki.Client
}

func (m *mainApp) mainLoop() {
	contentsDataLocal, err := m.wiki.GetPage("dropbox_sync")
	if err != nil {
		log.Fatalf("Failed to fetch page contents from wiki: %+v", err)
	}
	m.contentsData = contentsDataLocal
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
			break
		case <-m.signalChannel:
			os.Exit(0)
			break
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
	pagename := niceName(f.Name)
	if !strings.Contains(m.contentsData, fmt.Sprintf("(%s)", pagename)) {
		log.Printf("Adding link to %s (%s) to the contents page", pagename, f.Name)
		m.contentsData = m.contentsData + fmt.Sprintf("\n- [%s](%s)", f.Name, pagename)
		success, err := m.wiki.PutPage("dropbox_sync", m.contentsData)
		if err != nil {
			log.Fatalf("Failed to store page in wiki: %+v", err)
		}
		if !success {
			log.Println("Failed to store page in wiki")
		}
	}
	currentPage, err := m.wiki.GetPage(pagename)
	if err != nil {
		log.Fatalf("Failed to get page %s from wiki: %+v", pagename, err)
	}
	if string(f.Data) == currentPage {
		log.Printf("Same contents detected for %v, skipping", f.Name)
	} else {
		success, err := m.wiki.PutPage(pagename, string(f.Data))
		if err != nil {
			log.Fatalf("Failed to store page %s in wiki: %+v", pagename, err)
		}
		if !success {
			log.Println("Failed to store page in wiki")
		}
	}
}
