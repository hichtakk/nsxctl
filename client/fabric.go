package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetComputeManager() *[]structs.ComputeManager {
	path := "/api/v1/fabric/compute-managers"
	res := c.Request("GET", path, nil, nil)
	res_cms := res.Body.(map[string]interface{})["results"]
	cms := []structs.ComputeManager{}
	for _, res_cm := range res_cms.([]interface{}) {
		c := structs.ComputeManager{
			Id:     res_cm.(map[string]interface{})["id"].(string),
			Name:   res_cm.(map[string]interface{})["display_name"].(string),
			Type:   res_cm.(map[string]interface{})["origin_type"].(string),
			Server: res_cm.(map[string]interface{})["server"].(string),
			Detail: res_cm.(map[string]interface{})["origin_properties"].([]interface{})[0].(map[string]interface{})["value"].(string),
		}
		cms = append(cms, c)
	}
	return &cms
}

func (c *NsxtClient) GetComputeManagerStatus(cmId string) *structs.ComputeManagerStatus {
	path := "/api/v1/fabric/compute-managers/" + cmId + "/status"
	res := c.Request("GET", path, nil, nil)
	status := structs.ComputeManagerStatus{}
	status.Connection = res.Body.(map[string]interface{})["connection_status"].(string)
	status.Registration = res.Body.(map[string]interface{})["registration_status"].(string)
	return &status
}

func (c *NsxtClient) CreateComputeManager(name string, address string, thumbprint string, user string, password string, trust bool) {
	path := "/api/v1/fabric/compute-managers"
	reqData := make(map[string]interface{})
	reqData["display_name"] = name
	reqData["server"] = address
	reqData["origin_type"] = "vCenter"
	reqData["set_as_oidc_provider"] = trust
	reqData["credential"] = map[string]string{
		"credential_type": "UsernamePasswordLoginCredential",
		"username":        user,
		"password":        password,
		"thumbprint":      thumbprint,
	}
	reqJson, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", c.BaseUrl+path, bytes.NewBuffer(reqJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	data := readResponseBody(res)
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(b)
	}
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
	cms := data.(map[string]interface{})["results"]
	for _, cm := range cms.([]interface{}) {
		b, _ := json.MarshalIndent(cm, "", "  ")
		fmt.Println(string(b))
	}
}

func (c *NsxtClient) DeleteComputeManager(cmId string) {
	path := "/api/v1/fabric/compute-managers"
	req, _ := http.NewRequest("DELETE", c.BaseUrl+path+"/"+cmId, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(b)
	}
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
	_dumpRequest(req)
	_dumpResponse(res)
}

func (c *NsxtClient) GetTransportZone() {
	path := "/api/v1/transport-zones"
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

func (c *NsxtClient) CreateTransportZone(name string, transportType string) {
	path := "/api/v1/transport-zones"
	reqData := make(map[string]string)
	reqData["display_name"] = name
	reqData["transport_type"] = transportType
	reqJson, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", c.BaseUrl+path, bytes.NewBuffer(reqJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	data := readResponseBody(res)
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(b)
	}
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
	cms := data.(map[string]interface{})["results"]
	for _, cm := range cms.([]interface{}) {
		b, _ := json.MarshalIndent(cm, "", "  ")
		fmt.Println(string(b))
	}
}

func (c *NsxtClient) DeleteTransportZone(tzId string) {
	path := "/api/v1/transport-zones"
	req, _ := http.NewRequest("DELETE", c.BaseUrl+path+"/"+tzId, nil)
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(b)
	}
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		fmt.Println(res)
		return
	}
	_dumpRequest(req)
	_dumpResponse(res)
}

func (c *NsxtClient) PublishFQDN() {
	path := "/api/v1/configs/management"
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
	reqData := data.(map[string]interface{})
	reqData["publish_fqdns"] = false
	fmt.Println(reqData)
	reqJson, _ := json.Marshal(reqData)
	req, _ = http.NewRequest("PUT", c.BaseUrl+path, bytes.NewBuffer(reqJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", c.Token)
	res, err = c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		return
	}
	_dumpResponse(res)
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
