module github.com/milanaleksic/notesforlife

go 1.11

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/aktau/github-release v0.7.2 // indirect
	github.com/dropbox/dropbox-sdk-go-unofficial v5.4.0+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/kolo/xmlrpc v0.0.0-20181023172212-16bdd962781d
	github.com/milanaleksic/minimal_consul_member v0.1.0
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2 // indirect
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	google.golang.org/appengine v1.5.0 // indirect
)

replace github.com/dropbox/dropbox-sdk-go-unofficial => ./vendor/github.com/dropbox/dropbox-sdk-go-unofficial
