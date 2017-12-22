# notes for life

> Please read and try to understand all of the deaign decisions listed below before you decide to use the application.

This application allows syncing of notes written in Markdown format and stored in Dropbox
into Dokuwiki (with active Markdown content plugin, of course).

The reason why this application is made is because, according to me, no easy or friendly way
exists that would allow day-to-day notes to be placed into Dokuwiki, but there are massive
number of very beautiful and easy-to-use Markdown mobile apps (I myself use *iA Writer* on my
old Android tablet and *Byword* on my iPhone). 

This application connects to Dropbox via their (unofficial Go) API client and constantly polls 
for changes. The changes are then pushed into Dokuwiki via XML-RPC. "Notes for life" will 
create a single page "dropboxsync" that will list all the files stored on Dropbox.

## Preparation

### Dropbox API token

On the page https://www.dropbox.com/developers/apps/create you will need to create a 
"Full Dropbox" app (not the *Business* one)  and press "Generated access token" to get one.

### Dokuwiki

By default, Dokuwiki does not allow RPC connections/integrations. You need to do at least 
these actions:

- install some Markdown syntax plugin otherwise page contents will look garbled
- have "fresh enough" Dokuwiki installation
- allow XML-RPC
- put user you wish to use to push contents into `remoteuser` option

Please read in more details on [https://www.dokuwiki.org/devel:xmlrpc#get_it_working](https://www.dokuwiki.org/devel:xmlrpc#get_it_working)

## Running

The app is meant to be executed regularly or as a service on a real computer/server.

You should run it like this:

```bash
notes_for_life \
  -token DROPBOX_API_TOKEN
  -path /Apps/Byword \
  -username DOKUWIKI_USERNAME \
  -password DOKUWIKI_PASSWORD \
  -wikiLocation "https://www.mydokuwikisite.com"
```

## Current state

Some MVP design choices:
- NFL will sync everything under a certain path in Dropbox (no extension filters), so be careful!
- there is no support for env variables, configuration files
- re-runs should be idempotent (no multiple writes of identical contents to Dokuwiki should happen)
  + but, you can still reorder links in the contents page
- sync is only from the dropbox into Dokuwiki (Dropbox is thus the master data store),
no binary sync is available
- empty files will put the link into the Dokuwiki contents listing page, but no page 
will be made since Dokuwiki does not allow empty pages.
