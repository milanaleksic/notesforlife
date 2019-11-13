# notes for life

[![Build Status](https://semaphoreci.com/api/v1/milanaleksic/notesforlife/branches/master/badge.svg)](https://semaphoreci.com/milanaleksic/notesforlife)

> Please read and try to understand all of the design decisions listed below before you decide to use the application.

This application allows syncing of notes stored in Dropbox
into Dokuwiki.

The reason why this application is made is because, according to me, no easy or friendly way
exists that would allow day-to-day notes to be placed into Dokuwiki, but there are massive
number of very beautiful and easy-to-use Markdown mobile apps. I myself use:
- *https://stackedit.io* while on a laptop 
- *iA Writer* on my old Android tablet and 
- *Byword* on my iPhone. 

This application connects to Dropbox via their (unofficial Go) API client and constantly polls 
for changes. The changes are then pushed into Dokuwiki via XML-RPC. "Notes for life" will 
create a single page "dropboxsync" that will list all the files stored on Dropbox.

## Preparation

### Dropbox API token

On the page https://www.dropbox.com/developers/apps/create you will need to create a 
"Full Dropbox" app (not the *Business* one)  and press "Generated access token" to get one.

### Dokuwiki

You need to install some Markdown syntax plugin otherwise pushing Markdown contents
into Dokuwiki will produce pages with garbled look.

By default, Dokuwiki does not allow RPC connections/integrations. You need to do at least 
these actions:

- have "fresh enough" Dokuwiki installation
- allow XML-RPC
- put user you wish to use to push contents into `remoteuser` option

Please read in more details on [https://www.dokuwiki.org/devel:xmlrpc#get_it_working](https://www.dokuwiki.org/devel:xmlrpc#get_it_working)

## Running

The app is meant to be executed regularly or as a service on a real computer/server.

It is configured via a `notesforlife.toml` file which should contain:

- For Dropbox:
  + API token generated via their website, as explained above
  + path which should be watched for changes and pushed into Dokuwiki
- For Dokuwiki:
  + username
  + password
  + location (URL), e.g. "https://www.mydokuwikisite.com"

## Current state

Some MVP design choices:
- NFL will sync everything under a certain path in Dropbox (no extension filters), so be careful!
- re-runs should be idempotent (no multiple writes of identical contents to Dokuwiki should happen)
  + but, you can still reorder links in the contents page
- sync is only from the Dropbox into Dokuwiki (Dropbox is thus the master data store),
no binary sync is available
- I use my own BADUC system to do a continuous deployment into my cluster at home, to turn off that
integration just set `consulLocation` to an empty string otherwise you can't start the application
- empty files will put the link into the Dokuwiki contents listing page, but no page 
will be made since Dokuwiki does not allow empty pages.
- if you are using newer Go (>=1.11) then you need to turn off the modules (don't use `GO111MODULE=on`)
