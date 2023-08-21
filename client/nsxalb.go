package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func NewNsxAlbClient(endpoint string, proxy string, debug bool) (*NsxAlbClient, error) {
	httpClient, err := newHttpClient(proxy)
	if err != nil {
		return nil, err
	}
	client := &NsxAlbClient{
		BaseUrl:    endpoint,
		httpClient: httpClient,
		Debug:      debug,
	}

	return client, nil
}

type NsxAlbClient struct {
	BaseUrl     string
	BasicAuth   bool
	Token       string
	httpClient  HttpClient
	Debug       bool
	Version     string
	FullVersion string
	Proxy       string
}

// login to NSX ALB and create session
func (c *NsxAlbClient) Login(cred url.Values) error {
	path := "/login"
	if c.Debug {
		log.Println("-- LOGIN SESSION --")
	}
	// convert url.Values(map[string][string]) to map[string]string
	cr := map[string]string{}
	for k, v := range cred {
		cr[k] = v[0]
	}
	credJson, _ := json.Marshal(cr)
	body := strings.NewReader(string(credJson))
	// create request object
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return err
	}
	// override default header
	header := map[string]string{
		"Content-Type": "application/json",
	}
	c.SetCustomHeader(req, header)
	// do request
	res := c.Call(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		if c.Debug {
			log.Printf("StatusCode=%d\n", res.StatusCode)
			log.Println(req)
			log.Println(req.Header)
		}
		return fmt.Errorf("authentication failed")
	}
	// get session id from cookie
	c.Token = c.getCookie(c.BaseUrl+path, "csrftoken")
	if c.Debug {
		// log.SetOutput(os.Stderr)
		log.Println("-- LOGIN SUCCEEDED --")
		// log.Println(res.Header)
	}
	// extract
	version, err := res.JsonGet("/version/Version")
	if err != nil {
		return err
	}
	c.Version = version.(string)

	return nil
}

func (c *NsxAlbClient) Logout() error {
	path := "/logout"
	if c.Debug {
		log.Println("-- LOGOUT SESSION --")
	}
	// create request object
	req, err := c.NewRequest("POST", path, nil)
	if err != nil {
		log.Println("request preparation error")
		log.Println(err.Error())
		return err
	}
	res := c.Call(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("StatusCode=%d\n", res.StatusCode)
		return fmt.Errorf("request failed")
	}
	if c.Debug {
		log.Println("-- LOGOUT SESSION SUCCESS --")
	}
	return nil
}

// call NSX REST API and return response
func (c *NsxAlbClient) Call(req *http.Request) *Response {
	res, err := c.httpClient.Do(req)
	if c.Debug {
		log.Println("-- REQUEST --")
		log.Println(req)
		log.Println("-- RESPONSE --")
		log.Println(res)
	}
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	r := &Response{res, nil}
	if r.ContentLength != 0 {
		// As some response return context-type with encoding like "application/json;charset=UTF-8", check prefix
		contentType := res.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			// decode response json body then store to Content field
			err = json.NewDecoder(r.Body).Decode(&r.Content)
			if err != nil {
				log.Println(r.Body)
				log.Printf("parse response body failed: %v", err.Error())
			}
		} else if strings.HasPrefix(contentType, "text/html") {
			buf := &bytes.Buffer{}
			buf.ReadFrom(r.Body)
			r.Content = buf.String()
		}
	}
	if c.Debug {
		log.Println("-- CONTENT --")
		log.Println(r.Content)
	}

	return r
}

func (c *NsxAlbClient) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	// validate path
	acceptablePrefixes := []string{"/api/", "/login", "/logout"}
	err := validatePrefix(path, acceptablePrefixes)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(method, c.BaseUrl+path, body)
	c.SetDefaultHeader(req)

	return req, nil
}

func (c *NsxAlbClient) SetDefaultHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Avi-Version", c.Version)
	//req.Header.Set("Authorization", "Basic xxxxxxxxxxx")
	req.Header.Set("X-CSRFToken", c.Token)
	switch req.Method {
	case "GET":
		req.Header.Set("Accept-Encoding", "application/json")
	case "POST":
		req.Header.Set("Referer", c.BaseUrl)
	}
}

func (c *NsxAlbClient) SetCustomHeader(req *http.Request, header map[string]string) {
	h := http.Header{}
	for k, v := range header {
		h.Set(k, v)
	}
	req.Header = h
}

func (c *NsxAlbClient) getCookie(u string, key string) string {
	url, _ := url.Parse(u)
	value := ""
	// cast interface to http.Client type
	hc := c.httpClient.(*http.Client)
	cookies := hc.Jar.Cookies(url)
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == key {
			value = cookies[i].Value
		}
	}
	return value
}

func (c *NsxAlbClient) SetHttpClient(hc HttpClient) {
	c.httpClient = hc
}
