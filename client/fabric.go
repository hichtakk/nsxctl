package client

import (
	"encoding/json"
	"fmt"
)

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
