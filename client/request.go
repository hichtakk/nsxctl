package client

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type NsxtClient struct {
	BaseUrl    string
	BasicAuth  bool
	Token      string
	httpClient *http.Client
	Debug      bool
}

func (c *NsxtClient) makeRequest(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, c.BaseUrl+path, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
	return req
}

func (c *NsxtClient) Request(method string, path string, req_data []byte) {
	if !(strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/policy/")) {
		fmt.Println("path must start with \"/api/\" or \"/policy/\"")
		return
	}
	req, _ := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(req_data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("StatusCode: %d\n", res.StatusCode)
	var data interface{}
	if len(res_body) > 0 {
		err = json.Unmarshal(res_body, &data)
		if err != nil {
			log.Println(err)
			return
		}
		j, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(j))
	} else {
		fmt.Println("no response body")
	}
}

func NewNsxtClient(basicAuth bool, debug bool) *NsxtClient {
	httpClient := newHttpClient()
	nsxtClient := &NsxtClient{BasicAuth: false, Token: "", httpClient: httpClient, Debug: debug}
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

func (c *NsxtClient) GetTlsFingerprint(server string, port uint) string {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", server, port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	fingerprint := sha256.Sum256(cert.Raw)
	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}
	conn.Close()
	return buf.String()
}

/*
 * functions for debugging
 */
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
