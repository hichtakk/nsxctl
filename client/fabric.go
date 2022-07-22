package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

// DEPRECATED
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

func (c *NsxtClient) GetPolicyTransportZone(site string, ep string) *structs.TransportZones {
	/*
		path := "/policy/api/v1" + ep_path + "/transport-zones"
		tzs := []structs.TransportZone{}
		res := c.Request("GET", path, nil, nil)
		for _, tz := range res.Body.(map[string]interface{})["results"].([]interface{}) {
			id := tz.(map[string]interface{})["id"].(string)
			name := tz.(map[string]interface{})["display_name"].(string)
			tz_type := tz.(map[string]interface{})["tz_type"].(string)
			tzs = append(tzs, structs.TransportZone{Id: id, Name: name, Type: tz_type})
		}
		return &tzs
	*/
	path := "/policy/api/v1/infra/sites/" + site + "/enforcement-points/" + ep + "/transport-zones"
	res := c.Request("GET", path, nil, nil)
	zones := structs.TransportZones{}
	str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	json.Unmarshal(str, &zones)

	return &zones
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

func (c *NsxtClient) GetTransportNode(site string, ep string) structs.TransportNodes {
	path := "/policy/api/v1/infra/sites/" + site + "/enforcement-points/" + ep + "/host-transport-nodes"
	params := map[string]string{"node_types": "HostNode"}
	res := c.Request("GET", path, params, nil)
	nodes := []structs.TransportNode{}
	for _, n := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(n)
		var node structs.TransportNode
		json.Unmarshal(str, &node)
		nodes = append(nodes, node)
	}

	return nodes
}

func (c *NsxtClient) GetTransportNodeTunnels(node_id string) structs.TransportNodeTunnels {
	path := "/api/v1/transport-nodes/" + node_id + "/tunnels"
	res := c.Request("GET", path, nil, nil)
	tunnels := []structs.TransportNodeTunnel{}
	for _, t := range res.Body.(map[string]interface{})["tunnels"].([]interface{}) {
		str, _ := json.Marshal(t)
		var tunnel structs.TransportNodeTunnel
		json.Unmarshal(str, &tunnel)
		tunnels = append(tunnels, tunnel)
	}

	return structs.TransportNodeTunnels(tunnels)
}

func (c *NsxtClient) GetTransportNodeStatus(id string) string {
	path := "/api/v1/transport-nodes/" + id + "/status?source=realtime"
	res := c.Request("GET", path, nil, nil)
	if res.Error != nil {
		return ""
	}
	return res.Body.(map[string]interface{})["status"].(string)
}

func (c *NsxtClient) GetTransportNodeById(uuid string) *structs.TransportNode {
	path := "/api/v1/transport-nodes/" + uuid
	res := c.Request("GET", path, nil, nil)
	str, _ := json.Marshal(res.Body)
	var node structs.TransportNode
	json.Unmarshal(str, &node)
	return &node
}

func (c *NsxtClient) GetTransportNodeProfile() *structs.TransportNodeProfiles {
	path := "/policy/api/v1/infra/host-transport-node-profiles"
	res := c.Request("GET", path, nil, nil)
	str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	var profiles structs.TransportNodeProfiles
	json.Unmarshal(str, &profiles)
	return &profiles
}

func (c *NsxtClient) GetEdgeCluster() *[]structs.EdgeCluster {
	path := "/api/v1/edge-clusters/"
	res := c.Request("GET", path, nil, nil)
	edgeClusters := []structs.EdgeCluster{}
	for _, ec := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(ec)
		var edgeCluster structs.EdgeCluster
		json.Unmarshal(str, &edgeCluster)
		edgeClusters = append(edgeClusters, edgeCluster)
	}
	return &edgeClusters
}

func (c *NsxtClient) CreateEdge(name string, template_name string, address string, root_password string, admin_password string) {
	edges := c.GetEdge()
	if edges == nil {
		log.Fatal("template edge not found")
		return
	}

	var template_edge structs.TransportNode
	for _, e := range edges {
		if e.Name == template_name {
			template_edge = e
			break
		}
	}

	template_edge.Name = name
	template_edge.EdgeNodeDeploymentInfo.IPAddress[0] = address
	template_edge.EdgeNodeDeploymentInfo.EdgeDeploymentConfig.VMDeploymentConfig.ManagementPortSubnets[0].IPAddresses[0] = address
	template_edge.EdgeNodeDeploymentInfo.EdgeDeploymentConfig.Users["cli_password"] = admin_password
	template_edge.EdgeNodeDeploymentInfo.EdgeDeploymentConfig.Users["root_password"] = root_password

	jsonObj, err := json.Marshal(template_edge)
	if err != nil {
		log.Print(err)
		return
	}

	path := "/api/v1/transport-nodes"
	var resp *Response
	resp = c.Request("POST", path, nil, jsonObj)
	fmt.Println(resp)
}

func (c *NsxtClient) GetEdge() []structs.TransportNode {
	path := "/api/v1/transport-nodes?node_types=EdgeNode"
	res := c.Request("GET", path, nil, nil)
	edges := []structs.TransportNode{}
	for _, e := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(e)
		var edge structs.TransportNode
		json.Unmarshal(str, &edge)
		edges = append(edges, edge)
	}

	return edges
}
