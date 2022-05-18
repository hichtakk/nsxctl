package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetGateway(tier int16, gwId string) structs.Tier0Gateways {
	var path string
	if tier == 0 {
		path = "/policy/api/v1/infra/tier-0s"
	} else {
		path = "/policy/api/v1/infra/tier-1s"
	}
	if gwId != "" {
		path = path + "/" + gwId
	}
	res := c.Request("GET", path, nil, nil)
	var data interface{}
	err := json.Unmarshal([]byte(res), &data)
	if err != nil {
		log.Fatal(err)
	}
	var gateways interface{}
	gws := structs.Tier0Gateways{}
	if gwId != "" {
		gateways = data
		gw := structs.Tier0Gateway{}
		jsonStr, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(jsonStr, &gw)
		gws = append(gws, gw)
	} else {
		gateways = data.(map[string]interface{})["results"]
		for _, gateway := range gateways.([]interface{}) {
			gw := structs.Tier0Gateway{}
			jsonStr, err := json.Marshal(gateway)
			if err != nil {
				log.Println(err)
			}
			json.Unmarshal(jsonStr, &gw)
			gws = append(gws, gw)
		}
	}
	return gws
}

func (c *NsxtClient) GetLocaleService(tier0_id string) []string {
	req := c.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services")
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	//_dumpRequest(req)

	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		_dumpResponse(res)
		return []string{}
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

	locale_service_id := "default"

	return []string{locale_service_id}
}

func (c *NsxtClient) GetInterface(tier0_id string, locale_service_id string) []map[string]string {
	req := c.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services/"+locale_service_id+"/interfaces")
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	//_dumpRequest(req)

	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		_dumpResponse(res)
		return []map[string]string{}
	}
	data := readResponseBody(res)
	//_dumpResponse(res)
	interfaces := []map[string]string{}
	ifs := data.(map[string]interface{})["results"]
	for _, if_data := range ifs.([]interface{}) {
		if_data_map := map[string]string{
			"id":           if_data.(map[string]interface{})["id"].(string),
			"edge_path":    if_data.(map[string]interface{})["edge_path"].(string),
			"name":         if_data.(map[string]interface{})["display_name"].(string),
			"segment_path": if_data.(map[string]interface{})["segment_path"].(string),
			"uid":          if_data.(map[string]interface{})["unique_id"].(string),
		}
		interfaces = append(interfaces, if_data_map)
	}
	//log.Println(interfaces)

	return interfaces
}

func (c *NsxtClient) GetInterfaceStatistics(tier0_id string, locale_service_id string, if_data map[string]string) structs.RouterStats {
	// Quary Params:
	//   enforcement_point_path
	//   edge_path : ex edge_path=/infra/sites/default/enforcement-points/default/edge-clusters/57d2c653-4d63-48d8-b188-40b4e45a9bc8/edge-nodes/2ed9af04-21c9-11e9-be65-000c2902dff7
	req := c.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services/"+locale_service_id+"/interfaces/"+if_data["id"]+"/statistics")
	params := url.Values{
		"enforcement_point_path": []string{"/infra/sites/default/enforcement-points/default"},
		"edge_path":              []string{if_data["edge_path"]},
	}
	req.URL.RawQuery = params.Encode()
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	//sites := c.GetSite()
	//for _, site := range sites {
	//	for _, ep_id := range c.GetEnforcementPoint(site) {
	//		c.GetEdgeClusterUnderEnforcementPoint(site, ep_id)
	//	}
	//}
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		_dumpResponse(res)
		return structs.RouterStats{}
	}
	data := readResponseBody(res)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	var ifStats structs.RouterStats
	json.Unmarshal(jsonStr, &ifStats)

	return ifStats
}

func (c *NsxtClient) GetGatewayInterfaceStats(gw structs.Tier0Gateway) map[int]structs.RouterStats {
	locale_service_id := c.GetLocaleService(gw.Id)
	interfaces := c.GetInterface(gw.Id, locale_service_id[0])
	mutex := &sync.Mutex{}
	result := make(map[int]structs.RouterStats)
	var wg sync.WaitGroup
	for idx, inf := range interfaces {
		wg.Add(1)
		go func(idx int, i map[string]string) {
			defer wg.Done()
			mutex.Lock()
			result[idx] = c.GetInterfaceStatistics(gw.Id, locale_service_id[0], i)
			mutex.Unlock()
		}(idx, inf)
	}
	wg.Wait()
	return result
}
