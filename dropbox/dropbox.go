package dropbox

import (
	"io/ioutil"
	"log"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type ChangedFile struct {
	Data    []byte
	Name    string
	Initial bool
}

type Client struct {
	ChangedFile chan ChangedFile
	known       map[string]string
	dbx         files.Client
	path        string
}

func NewClient(token, path string) *Client {
	output := make(chan ChangedFile, 10)
	return &Client{
		ChangedFile: output,
		path:        path,
		known:       make(map[string]string),
		dbx:         files.New(dropbox.Config{Token: token}),
	}
}

func (c *Client) Process() {
	items, cursor := c.listCurrentItems()
	for _, item := range items {
		c.handleChangedFile(item, true)
	}
	c.longPoll(cursor)
}

func (c *Client) longPoll(cursor string) {
	activeCursor := cursor
	for {
		longPollArg := files.NewListFolderLongpollArg(activeCursor)
		resp, err := c.dbx.ListFolderLongpoll(longPollArg)
		if err != nil {
			log.Fatalf("Could not long poll file list: %v", err)
		}
		if resp.Changes {
			items, newCursor := c.listCurrentItems()
			activeCursor = newCursor
			for _, item := range items {
				c.handleChangedFile(item, false)
			}
		} else {
			log.Println("No changes detected")
		}
	}
}

func (c *Client) listCurrentItems() (items []*files.FileMetadata, cursor string) {
	items = make([]*files.FileMetadata, 0)
	listFolderArg := files.NewListFolderArg(c.path)
	listFolderArg.Recursive = true
	resp, err := c.dbx.ListFolder(listFolderArg)
	if err != nil {
		log.Fatalf("Could not fetch file list: %v", err)
	}
	for _, item := range resp.Entries {
		if _, ok := item.(*files.FileMetadata); !ok {
			continue
		}
		if c.known[item.(*files.FileMetadata).Name] != item.(*files.FileMetadata).ContentHash {
			c.known[item.(*files.FileMetadata).Name] = item.(*files.FileMetadata).ContentHash
			items = append(items, item.(*files.FileMetadata))
		}
	}

	hasMore := resp.HasMore
	cursor = resp.Cursor
	for hasMore {
		listFolderArg := files.NewListFolderContinueArg(cursor)
		respContinue, err := c.dbx.ListFolderContinue(listFolderArg)
		if err != nil {
			log.Fatalf("Could not continue fetching file list: %v", err)
		}
		hasMore = respContinue.HasMore
		cursor = respContinue.Cursor
		for _, item := range respContinue.Entries {
			if _, ok := item.(*files.FileMetadata); !ok {
				continue
			}
			if c.known[item.(*files.FileMetadata).Name] != item.(*files.FileMetadata).ContentHash {
				c.known[item.(*files.FileMetadata).Name] = item.(*files.FileMetadata).ContentHash
				items = append(items, item.(*files.FileMetadata))
			}
		}
	}
	return
}

func (c *Client) handleChangedFile(item *files.FileMetadata, initial bool) {
	log.Printf("File to be analyzed for contents: %+v, initial=%v", item, initial)
	downloadArg := files.NewDownloadArg(item.PathDisplay)
	res, content, err := c.dbx.Download(downloadArg)
	if err != nil {
		log.Fatalf("Could not download file: %+v", err)
	}
	defer func() { _ = content.Close() }()
	data, err := ioutil.ReadAll(content)
	if err != nil {
		log.Fatalf("Could not read data from remote: %+v", err)
	}

	c.ChangedFile <- ChangedFile{
		Name:    res.PathDisplay,
		Data:    data,
		Initial: initial,
	}
}
