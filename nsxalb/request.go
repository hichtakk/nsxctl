package nsxalb

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hichtakk/nsxctl/client"
)

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

func (c *NsxAlbClient) Request(method string, path string, query_param map[string]string, req_data []byte) *client.Response {
	// validate path
	err := func() error {
		var match bool
		match = false
		for _, v := range []string{"/login", "/logout", "/api/"} {
			if strings.HasPrefix(path, v) {
				match = true
			}
		}
		if match == false {
			return fmt.Errorf("path must start with \"/api/\"")
		}
		return nil
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return &client.Response{}
	}
	req, _ := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(req_data))
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-CSRFToken", c.Token)
	}
	req.Header.Set("Referer", c.BaseUrl)
	if len(query_param) > 0 {
		params := req.URL.Query()
		for k, v := range query_param {
			params.Add(k, v)
		}
		req.URL.RawQuery = params.Encode()
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return &client.Response{}
	}
	defer res.Body.Close()
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return &client.Response{}
	}
	var data interface{}
	if len(res_body) > 0 {
		err = json.Unmarshal(res_body, &data)
		if err != nil {
			//log.Fatal("response json decode error")
			return &client.Response{res, nil, nil}
		}
		r := &client.Response{res, data, nil}
		return r
	} else {
		return &client.Response{res, nil, nil}
	}
}
