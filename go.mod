module github.com/milanaleksic/notesforlife

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dropbox/dropbox-sdk-go-unofficial v5.4.0+incompatible
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/kolo/xmlrpc v0.0.0-20181023172212-16bdd962781d
	github.com/milanaleksic/minimal_consul_member v0.1.0
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	google.golang.org/appengine v1.5.0 // indirect
)

replace github.com/dropbox/dropbox-sdk-go-unofficial => ./vendor/github.com/dropbox/dropbox-sdk-go-unofficial
