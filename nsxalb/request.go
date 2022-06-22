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

func (c *NsxAlbClient) Request(method string, path string, req_data []byte) string {
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
		return ""
	}

	req, _ := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(req_data))
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-CSRFToken", c.Token)
	}
	req.Header.Set("Referer", c.BaseUrl)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		fmt.Printf("StatusCode: %d\n", res.StatusCode)
		fmt.Printf("%s\n", res_body)
	}
	var data interface{}
	if len(res_body) > 0 {
		err = json.Unmarshal(res_body, &data)
		if err != nil {
			log.Println(err)
			return ""
		}
		j, _ := json.MarshalIndent(data, "", "  ")
		return string(j)
	} else {
		fmt.Println("no response body")
	}

	return ""
}
