package nsxalb

import (
	"net/http"
	"net/http/cookiejar"
)

type NsxAlbClient struct {
	BaseUrl     string
	BasicAuth   bool
	Token       string
	httpClient  *http.Client
	Debug       bool
	Version     string
	FullVersion string
}

func NewNsxAlbClient(basicAuth bool, debug bool) *NsxAlbClient {
	httpClient := newHttpClient()
	nsxAlbClient := &NsxAlbClient{BasicAuth: false, Token: "", httpClient: httpClient, Debug: debug}
	if basicAuth != true {
		jar, _ := cookiejar.New(nil)
		nsxAlbClient.httpClient.Jar = jar
	}

	return nsxAlbClient
}
