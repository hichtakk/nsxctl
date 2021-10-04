package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *NsxtClient) GetComputeManager() {
	path := "/api/v1/fabric/compute-managers"
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
		//fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
		b, _ := json.MarshalIndent(cm, "", "  ")
		fmt.Println(string(b))
	}
}

func (c *NsxtClient) GetTransportNode() {
	req := c.makeRequest("GET", "/api/v1/transport-nodes")
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
	gateways := data.(map[string]interface{})["results"]
	for _, gateway := range gateways.([]interface{}) {
		//fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
		b, _ := json.MarshalIndent(gateway, "", "  ")
		fmt.Println(string(b))
	}
}

func (c *NsxtClient) GetTransportNodeProfile() {
	req := c.makeRequest("GET", "/api/v1/transport-node-profiles")
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
	gateways := data.(map[string]interface{})["results"]
	for _, gateway := range gateways.([]interface{}) {
		//fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
		b, _ := json.MarshalIndent(gateway, "", "  ")
		fmt.Println(string(b))
	}
}
