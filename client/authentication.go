package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (c *NsxtClient) Login(cred url.Values) {
	target_url := c.BaseUrl + "/api/session/create"
	req, _ := http.NewRequest("POST", target_url, strings.NewReader(cred.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		return
	}
	log.Println("login successful")
	c.Token = res.Header.Get("X-Xsrf-Token")
	if c.Token == "" {
		log.Fatal("token not found")
	}
	data := readResponseBody(res)
	roles := data.(map[string]interface{})["roles"]
	for _, v := range roles.([]interface{}) {
		fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
	}
}

func (c *NsxtClient) Logout() {
	target_url := c.BaseUrl + "/api/session/destroy"
	req, _ := http.NewRequest("POST", target_url, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		return
	}
	log.Printf("logout successful\n")
}
