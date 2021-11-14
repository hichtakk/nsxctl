package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *NsxtClient) GetIpPool() {
	path := "/policy/api/v1/infra/ip-pools"
	req, _ := http.NewRequest("GET", c.BaseUrl+path, nil)
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
	data := readResponseBody(res)
	cms := data.(map[string]interface{})["results"]
	for _, cm := range cms.([]interface{}) {
		b, _ := json.MarshalIndent(cm, "", "  ")
		fmt.Println(string(b))
	}
}

func (c *NsxtClient) GetIpBlock() {
	path := "/policy/api/v1/infra/ip-blocks"
	req, _ := http.NewRequest("GET", c.BaseUrl+path, nil)
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
	data := readResponseBody(res)
	blks := data.(map[string]interface{})["results"]
	for _, blk := range blks.([]interface{}) {
		b, _ := json.MarshalIndent(blk, "", "  ")
		fmt.Println(string(b))
	}
}

func (c *NsxtClient) CreateIpPool(name string) {
	path := "/policy/api/v1/infra/ip-pools"
	reqData := make(map[string]string)
	reqData["display_name"] = name
	reqJson, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("PATCH", c.BaseUrl+path+"/"+name, bytes.NewBuffer(reqJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// no content returned when request succeeded
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
}

func (c *NsxtClient) CreateIpBlock(name string, cidr string) {
	path := "/policy/api/v1/infra/ip-blocks"
	reqData := make(map[string]string)
	reqData["display_name"] = name
	reqData["cidr"] = cidr
	reqJson, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("PATCH", c.BaseUrl+path+"/"+name, bytes.NewBuffer(reqJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// no content returned when request succeeded
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
}

func (c *NsxtClient) DeleteIpPool(name string) {
	path := "/policy/api/v1/infra/ip-pools"
	req, _ := http.NewRequest("DELETE", c.BaseUrl+path+"/"+name, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// no content returned when request succeeded
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
}

func (c *NsxtClient) DeleteIpBlock(name string) {
	path := "/policy/api/v1/infra/ip-blocks"
	req, _ := http.NewRequest("DELETE", c.BaseUrl+path+"/"+name, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// no content returned when request succeeded
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
}
