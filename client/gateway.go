package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetGateway(gwId string, tier int) structs.Gateways {
	path := fmt.Sprintf("/policy/api/v1/infra/tier-%ds", tier)
	if gwId != "" {
		path = path + "/" + gwId
	}
	res := c.Request("GET", path, nil, nil)
	var gateways interface{}
	gws := structs.Gateways{}
	if gwId != "" {
		gw := structs.Gateway{}
		res.UnmarshalBody(&gw)
		gw.Tier = tier
		gws = append(gws, gw)
	} else {
		gateways = res.Body.(map[string]interface{})["results"]
		for _, gateway := range gateways.([]interface{}) {
			gw := structs.Gateway{}
			jsonStr, err := json.Marshal(gateway)
			if err != nil {
				log.Println(err)
			}
			json.Unmarshal(jsonStr, &gw)
			gw.Tier = tier
			gws = append(gws, gw)
		}
	}
	return gws
}

func (c *NsxtClient) GetGateways (tier int) structs.Gateways {
	var gws structs.Gateways
	if tier == -1 {
		gws = append(gws, c.GetGateway("", 0)...)
		gws = append(gws, c.GetGateway("", 1)...)
	} else {
		gws = append(gws, c.GetGateway("", tier)...)
	}
	return gws
}

func (c *NsxtClient) GetLocaleService(gw_id string, tier int) []structs.LocaleService {
	path := fmt.Sprintf("/policy/api/v1/infra/tier-%ds/%s/locale-services", tier, gw_id)
	res := c.Request("GET", path, nil, nil)

	var localeServices []structs.LocaleService
	for _, obj := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		var localeService structs.LocaleService
		jsonStr, err := json.Marshal(obj)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(jsonStr, &localeService)
		localeServices = append(localeServices, localeService)
	}
	return localeServices
}

func (c *NsxtClient) GetInterface(gw_id string, tier int, locale_service_id string) []structs.GatewayInterface {
	path := fmt.Sprintf("/policy/api/v1/infra/tier-%ds/%s/locale-services/%s/interfaces", tier, gw_id, locale_service_id)
	res := c.Request("GET", path, nil, nil)

	var interfaces []structs.GatewayInterface
	for _, if_data := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		var gw_if structs.GatewayInterface
		jsonStr, err := json.Marshal(if_data)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(jsonStr, &gw_if)
		interfaces = append(interfaces, gw_if)
	}
	return interfaces
}

func (c *NsxtClient) GetInterfaceStatistics(tier0_id string, locale_service_id string, if_data structs.GatewayInterface) structs.RouterStats {
	// Quary Params:
	//   enforcement_point_path
	//   edge_path : ex edge_path=/infra/sites/default/enforcement-points/default/edge-clusters/57d2c653-4d63-48d8-b188-40b4e45a9bc8/edge-nodes/2ed9af04-21c9-11e9-be65-000c2902dff7
	req := c.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services/"+locale_service_id+"/interfaces/"+if_data.Id+"/statistics")
	params := url.Values{
		"enforcement_point_path": []string{"/infra/sites/default/enforcement-points/default"},
		"edge_path":              []string{if_data.EdgePath},
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

func (c *NsxtClient) GetGatewayInterfaceStats(gw structs.Gateway) map[int]structs.RouterStats {
	locale_service := c.GetLocaleService(gw.Id, gw.Tier)
	interfaces := c.GetInterface(gw.Id, gw.Tier, locale_service[0].Id)
	mutex := &sync.Mutex{}
	result := make(map[int]structs.RouterStats)
	var wg sync.WaitGroup
	for idx, inf := range interfaces {
		wg.Add(1)
		go func(idx int, i structs.GatewayInterface) {
			defer wg.Done()
			mutex.Lock()
			result[idx] = c.GetInterfaceStatistics(gw.Id, locale_service[0].Id, i)
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

func (c *NsxtClient) GetGatewayAggregateInfo(gw_realization_id string) structs.GatewayAggregateInfo {
	path := "/api/v1/ui-controller/l3/logical-routers/" + gw_realization_id + "/status/aggregate-info?source=realtime"
	res := c.Request("GET", path, nil, nil)
	resBytes, err := res.BodyBytes()
	if err != nil {
		log.Fatal(err)
	}

	var gwAggregateInfo structs.GatewayAggregateInfo
	json.Unmarshal(resBytes, &gwAggregateInfo)
	return gwAggregateInfo
}

func (c *NsxtClient) GetGatewayFromName(name string, tier int) (structs.Gateway, error) {
	for _, gw := range c.GetGateways(tier) {
		if gw.Name == name {
			return gw, nil
		}
	}
	return structs.Gateway{}, fmt.Errorf("Error: Tier-0/1 gateway '%s' is not found", name)
}
