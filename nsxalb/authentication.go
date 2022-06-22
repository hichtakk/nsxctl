package nsxalb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func (c *NsxAlbClient) Login(cred map[string]string) {
	target_url := c.BaseUrl + "/login"
	credJson, _ := json.Marshal(cred)
	req, _ := http.NewRequest("POST", target_url, bytes.NewBuffer(credJson))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("StatusCode=%d\n", res.StatusCode)
		log.Println(req)
		log.Println(res)
		data := readResponseBody(res)
		log.Println(data)
		return
	}
	url, _ := url.Parse(target_url)
	cookies := c.httpClient.Jar.Cookies(url)
	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "csrftoken" {
			c.Token = cookies[i].Value
		}
	}
	data := readResponseBody(res)
	if c.Debug {
		log.SetOutput(os.Stderr)
		log.Println("login successful")
		log.Println(res.Header)
		log.Println(data)
	}
	versionData := data.(map[string]interface{})["version"]
	version := versionData.(map[string]interface{})["Version"]
	c.Version = version.(string)
}

func (c *NsxAlbClient) Logout() {
	target_url := c.BaseUrl + "/logout"
	req, _ := http.NewRequest("POST", target_url, nil)
	req.Header.Set("X-CSRFToken", c.Token)
	req.Header.Set("Referer", c.BaseUrl)
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("StatusCode=%d\n", res.StatusCode)
		return
	}
	if c.Debug {
		log.Printf("logout successful\n")
		log.Println(res.Header)
	}
}

func (c *NsxAlbClient) Cluster() {
	target_url := c.BaseUrl + "/api/pool"
	req, _ := http.NewRequest("GET", target_url, nil)
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
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
