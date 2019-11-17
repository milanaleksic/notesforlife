package dokuwiki

import (
	"errors"
	"log"

	"github.com/kolo/xmlrpc"
)

type Client struct {
	client *xmlrpc.Client
	dryRun bool
}

func NewClient(serverPath string, dryRun bool) *Client {
	client, err := xmlrpc.NewClient(serverPath, nil)
	if err != nil {
		log.Fatalf("Could not create XML-RPC client: %+v", err)
	}
	return &Client{
		client: client,
		dryRun: dryRun,
	}
}

func (d *Client) Login(username, password string) (err error) {
	var result bool
	err = d.client.Call("dokuwiki.login", []interface{}{username, password}, &result)
	if err != nil {
		return
	}
	if !result {
		return errors.New("credentials provided were not correct")
	}
	return
}

func (d *Client) GetVersion() (version string, err error) {
	err = d.client.Call("dokuwiki.getVersion", nil, &version)
	return
}

func (d *Client) GetPage(pagename string) (page string, err error) {
	err = d.client.Call("wiki.getPage", []interface{}{pagename}, &page)
	return
}

func (d *Client) PutPage(pagename, data string) (success bool, err error) {
	if d.dryRun {
		log.Printf("Since dry run is active, update of the page %s will not be done. Planned on saving: %s", pagename, data)
	} else {
		err = d.client.Call("wiki.putPage", []interface{}{pagename, data}, &success)
	}
	return
}
