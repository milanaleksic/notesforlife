module github.com/milanaleksic/notesforlife

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dropbox/dropbox-sdk-go-unofficial v5.4.0+incompatible
	github.com/golang/protobuf v1.3.1
	github.com/kolo/xmlrpc v0.0.0-20181023172212-16bdd962781d
	github.com/milanaleksic/minimal_consul_member v0.0.0-20180105162056-d210ecea4757
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/appengine v1.5.0
)

replace github.com/dropbox/dropbox-sdk-go-unofficial => ./vendor/github.com/dropbox/dropbox-sdk-go-unofficial
