package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetTier0Gateway(gwId string) structs.Tier0Gateways {
	path := "/policy/api/v1/infra/tier-0s"
	if gwId != "" {
		path = path + "/" + gwId
	}
	res := c.Request("GET", path, nil, nil)
	var gateways interface{}
	gws := structs.Tier0Gateways{}
	if gwId != "" {
		gw := structs.Tier0Gateway{}
		res.UnmarshalBody(&gw)
		gws = append(gws, gw)
	} else {
		gateways = res.Body.(map[string]interface{})["results"]
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

func (c *NsxtClient) GetTier1Gateway(gwId string) structs.Tier1Gateways {
	path := "/policy/api/v1/infra/tier-1s"
	if gwId != "" {
		path = path + "/" + gwId
	}
	res := c.Request("GET", path, nil, nil)
	var gateways interface{}
	gws := structs.Tier1Gateways{}
	if gwId != "" {
		gw := structs.Tier1Gateway{}
		res.UnmarshalBody(&gw)
		gws = append(gws, gw)
	} else {
		gateways = res.Body.(map[string]interface{})["results"]
		for _, gateway := range gateways.([]interface{}) {
			gw := structs.Tier1Gateway{}
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

	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		_dumpResponse(res)
		return []string{}
	}
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
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d\n", res.StatusCode)
		_dumpResponse(res)
		return []map[string]string{}
	}
	data := readResponseBody(res)
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

func (c *NsxtClient) GetRoutingTable(tier0Id string) []structs.EdgeRoute {
	var path string
	path = "/policy/api/v1/infra/tier-0s/" + tier0Id + "/routing-table"
	res := c.Request("GET", path, nil, nil)
	entries := []structs.EdgeRoute{}
	for _, ed := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(ed)
		var edgeRoutes structs.EdgeRoute
		json.Unmarshal(str, &edgeRoutes)
		entries = append(entries, edgeRoutes)
	}

	return entries
}

func (c *NsxtClient) GetBgpConfig(tier0Id string, locale string) structs.BgpConfig {
	var path string
	path = "/policy/api/v1/infra/tier-0s/" + tier0Id + "/locale-services/" + locale + "/bgp"
	res := c.Request("GET", path, nil, nil)
	bgpConfig := structs.BgpConfig{}
	body, _ := res.BodyBytes()
	json.Unmarshal(body, &bgpConfig)

	return bgpConfig
}

func (c *NsxtClient) GetBgpNeighbors(tier0Id string, locale string) []structs.BgpNeighbor {
	var path string
	path = "/policy/api/v1/infra/tier-0s/" + tier0Id + "/locale-services/" + locale + "/bgp/neighbors"
	res := c.Request("GET", path, nil, nil)
	neighbors := []structs.BgpNeighbor{}
	for _, nb := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(nb)
		var neighbor structs.BgpNeighbor
		json.Unmarshal(str, &neighbor)
		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

func (c *NsxtClient) GetBgpNeighborsAdvRoutes(path string) []structs.EdgeBgpAdvRoute {
	path = "/policy/api/v1" + path + "/advertised-routes"
	res := c.Request("GET", path, nil, nil)
	edges := []structs.EdgeBgpAdvRoute{}
	for _, nb := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		es := nb.(map[string]interface{})["edge_node_routes"].([]interface{})
		for _, e := range es {
			str, _ := json.Marshal(e)
			var edge structs.EdgeBgpAdvRoute
			json.Unmarshal(str, &edge)
			edges = append(edges, edge)
		}
	}

	return edges
}

func (c *NsxtClient) GetGatewayAggregateInfo(gw_realization_id string) []map[string]string {
	path := "/api/v1/ui-controller/l3/logical-routers/" + gw_realization_id + "/status/aggregate-info?source=realtime"
	res := c.Request("GET", path, nil, nil)
	if res == nil {
		return nil
	}

	per_node_status := res.Body.(map[string]interface{})["status"].(map[string]interface{})["per_node_status"]
	if per_node_status == nil {
		//if a tier-0 is created but interface isn't set, per_node_status will be empty
		return nil
	}

	result := []map[string]string{}
	for _, st := range per_node_status.([]interface{}) {
		str, _ := json.Marshal(st)
		var status map[string]string
		json.Unmarshal(str, &status)
		result = append(result, status)
	}
	return result
}
