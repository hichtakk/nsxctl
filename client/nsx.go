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

func NewNsxClient(endpoint string, proxy string, debug bool) (*NsxClient, error) {
	httpClient, err := newHttpClient(proxy)
	if err != nil {
		return nil, err
	}
	client := &NsxClient{
		BaseUrl:    endpoint,
		httpClient: httpClient,
		Debug:      debug,
	}

	return client, nil
}

type NsxClient struct {
	BaseUrl    string
	BasicAuth  bool
	Token      string
	httpClient HttpClient
	Debug      bool
	Proxy      string
}

// login to NSX and create session
func (c *NsxClient) Login(cred url.Values) error {
	path := "/api/session/create"
	if c.Debug {
		log.Println("-- LOGIN SESSION --")
	}
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	body := strings.NewReader(cred.Encode())
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return err
	}
	c.SetCustomHeader(req, header)
	res := c.Call(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("authentication failed")
	}
	// get session id from response header
	c.Token = res.Header.Get("X-Xsrf-Token")
	if c.Token == "" {
		log.Fatal("token not found")
	}
	if c.Debug {
		log.Println("-- LOGIN SUCCEEDED --")
		// roles, _ := res.JsonGet("/roles")
		// for _, role := range roles.([]interface{}) {
		// 	name, _ := jp.Get(role, "/role")
		// 	permissions, _ := jp.Get(role, "/permissions")
		// 	fmt.Printf("role: %s, permission: %v\n", name.(string), permissions)
		// }
	}

	return nil
}

// logout from NSX
func (c *NsxClient) Logout() error {
	path := "/api/session/destroy"
	if c.Debug {
		log.Println("-- LOGOUT SESSION --")
	}
	// create request object
	req, err := c.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	res := c.Call(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("StatusCode=%d\n", res.StatusCode)
		return nil
	}
	if c.Debug {
		log.Println("-- LOGOUT SESSION SUCCESS --")
	}

	return nil
}

// call NSX REST API and return response
func (c *NsxClient) Call(req *http.Request) *Response {
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

func (c *NsxClient) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	// validate path
	acceptablePrefixes := []string{"/api/", "/policy/"}
	err := validatePrefix(path, acceptablePrefixes)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(method, c.BaseUrl+path, body)
	c.SetDefaultHeader(req)

	return req, nil
}

func (c *NsxClient) SetDefaultHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
}

func (c *NsxClient) SetCustomHeader(req *http.Request, header map[string]string) {
	h := http.Header{}
	for k, v := range header {
		h.Set(k, v)
	}
	req.Header = h
}

func (c *NsxClient) SetHttpClient(hc HttpClient) {
	c.httpClient = hc
}
