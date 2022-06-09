package client

import (
	"fmt"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetSite() []string {
	req := c.makeRequest("GET", "/policy/api/v1/infra/sites")
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		return make([]string, 0, 0)
	}
	//data := readResponseBody(res)
	/*
		        gateways := data.(map[string]interface{})["results"]
				for _, gateway := range gateways.([]interface{}) {
					//fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
					b, _ := json.MarshalIndent(gateway, "", "  ")
					fmt.Println(string(b))
				}
	*/
	//_dumpResponse(res)

	// tentative: return "default"
	return []string{"default"}
}

func (c *NsxtClient) GetEnforcementPoint(site_id string) *[]structs.EnforcementPoint {
	path := "/policy/api/v1/infra/sites/" + site_id + "/enforcement-points"
	eps := []structs.EnforcementPoint{}
	res := c.Request("GET", path, nil, nil)
	for _, ep := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		id := ep.(map[string]interface{})["id"].(string)
		path := ep.(map[string]interface{})["path"].(string)
		e := structs.EnforcementPoint{Id: id, Path: path}
		eps = append(eps, e)
	}
	return &eps
}

func (c *NsxtClient) GetEdgeClusterUnderEnforcementPoint(site_id string, ep_id string) []string {
	//log.Println("edge-cluster")
	req := c.makeRequest("GET", "/policy/api/v1/infra/sites/"+site_id+"/enforcement-points/"+ep_id+"/edge-clusters")
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	//_dumpRequest(req)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		_dumpResponse(res)
		return make([]string, 0, 0)
	}
	/*
		    	data := readResponseBody(res)
				clusters := data.(map[string]interface{})["results"]
				for _, gateway := range gateways.([]interface{}) {
					//fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
					b, _ := json.MarshalIndent(gateway, "", "  ")
					fmt.Println(string(b))
				}
	*/

	//_dumpResponse(res)

	// tentative: return default path
	return []string{"/infra/sites/default/enforcement-points/" + site_id}
}
