package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (c *NsxtClient) Login(cred url.Values) error {
	target_url := c.BaseUrl + "/api/session/create"
	req, _ := http.NewRequest("POST", target_url, strings.NewReader(cred.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("authentication failed")
	}
	c.Token = res.Header.Get("X-Xsrf-Token")
	if c.Token == "" {
		log.Fatal("token not found")
	}
	data := readResponseBody(res)
	if c.Debug {
		log.Println("login successful")
		log.Println(res.Header)
		roles := data.(map[string]interface{})["roles"]
		for _, v := range roles.([]interface{}) {
			fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
		}
	}

	return nil
}

func (c *NsxtClient) Logout() {
	if c.Token == "" {
		return
	}
	target_url := c.BaseUrl + "/api/session/destroy"
	req, _ := http.NewRequest("POST", target_url, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
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
