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
	"os"
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

type Response struct {
	*http.Response
	Body interface{}
}

func (r *Response) BodyBytes() ([]byte, error) {
	return json.Marshal(r.Body)
}

func (r *Response) UnmarshalBody(strct interface{}) {
	bytes, _ := r.BodyBytes()
	json.Unmarshal(bytes, strct)
}

func (c *NsxtClient) makeRequest(method string, path string) *http.Request {
	req, _ := http.NewRequest(method, c.BaseUrl+path, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
	return req
}

func (c *NsxtClient) Request(method string, path string, query_param map[string]string, req_data []byte) *Response {
	// validate path
	err := func() error {
		var match bool
		match = false
		for _, v := range []string{"/api/", "/policy/"} {
			if strings.HasPrefix(path, v) {
				match = true
			}
		}
		if match == false {
			return fmt.Errorf("path must start with \"/api/\" or \"/policy/\"")
		}
		return nil
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return &Response{}
	}
	req, _ := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(req_data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
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
		return &Response{}
	}
	defer res.Body.Close()
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return &Response{}
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		fmt.Fprintf(os.Stderr, "StatusCode: %d\n", res.StatusCode)
		fmt.Fprintf(os.Stderr, "%s", res_body)
		return &Response{}
	}
	var data interface{}
	if len(res_body) > 0 {
		err = json.Unmarshal(res_body, &data)
		if err != nil {
			log.Println(err)
			return &Response{}
		}
		//j, _ := json.MarshalIndent(data, "", "  ")
		r := &Response{res, data}
		return r
		//return string(j)
	} else {
		fmt.Println("no response body")
	}

	return &Response{}
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
