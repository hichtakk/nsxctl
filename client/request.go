package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"time"
)

type NsxtClient struct {
	BaseUrl    string
	BasicAuth  bool
	Token      string
	httpClient *http.Client
}

func (c *NsxtClient) makeRequest(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, c.BaseUrl+path, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
	return req
}

// functions for debugging

func _dumpRequest(req *http.Request) {
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s\n\n", dump)
}

func _dumpResponse(res *http.Response) {
	dump, _ := httputil.DumpResponse(res, true)
	fmt.Printf("%s\n\n", dump)
}

func _dumpCookie(c *NsxtClient, target_url string) {
	set_cookie_url, _ := url.Parse(target_url)
	cookies := c.httpClient.Jar.Cookies(set_cookie_url)
	fmt.Printf("%v\n\n", cookies)
}

func NewNsxtClient(basicAuth bool) *NsxtClient {
	httpClient := newHttpClient()
	nsxtClient := &NsxtClient{BasicAuth: false, Token: "", httpClient: httpClient}
	if basicAuth != true {
		jar, _ := cookiejar.New(nil)
		nsxtClient.httpClient.Jar = jar
	}

	return nsxtClient
}

func newHttpClient() *http.Client {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   time.Duration(30) * time.Second,
	}
	return client
}

func readResponseBody(res *http.Response) interface{} {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
