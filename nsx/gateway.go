package nsx

import (
	"fmt"
)

func (a *Agent) GetTier0Gateway(gwId string) Tier0Gateways {
	// path := "/policy/api/v1/infra/tier-0s"
	// if gwId != "" {
	// 	path = path + "/" + gwId
	// }
	// res := a.Request("GET", path, nil, nil)
	// var gateways interface{}
	gws := Tier0Gateways{}
	// if gwId != "" {
	// 	gw := Tier0Gateway{}
	// 	res.UnmarshalBody(&gw)
	// 	gws = append(gws, gw)
	// } else {
	// 	gateways = res.Body.(map[string]interface{})["results"]
	// 	for _, gateway := range gateways.([]interface{}) {
	// 		gw := Tier0Gateway{}
	// 		jsonStr, err := json.Marshal(gateway)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		json.Unmarshal(jsonStr, &gw)
	// 		gws = append(gws, gw)
	// 	}
	// }
	return gws
}

func (a *Agent) GetTier1Gateway(gwId string) Tier1Gateways {
	// path := "/policy/api/v1/infra/tier-1s"
	// if gwId != "" {
	// 	path = path + "/" + gwId
	// }
	// res := a.Request("GET", path, nil, nil)
	// var gateways interface{}
	gws := Tier1Gateways{}
	// if gwId != "" {
	// 	gw := Tier1Gateway{}
	// 	res.UnmarshalBody(&gw)
	// 	gws = append(gws, gw)
	// } else {
	// 	gateways = res.Body.(map[string]interface{})["results"]
	// 	for _, gateway := range gateways.([]interface{}) {
	// 		gw := Tier1Gateway{}
	// 		jsonStr, err := json.Marshal(gateway)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		json.Unmarshal(jsonStr, &gw)
	// 		gws = append(gws, gw)
	// 	}
	// }
	return gws
}

func (a *Agent) GetLocaleService(tier0_id string) []string {
	// req := a.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services")
	// res, err := a.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	_dumpResponse(res)
	// 	return []string{}
	// }
	locale_service_id := "default"

	return []string{locale_service_id}
}

func (a *Agent) GetInterface(tier0_id string, locale_service_id string) []map[string]string {
	// req := a.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services/"+locale_service_id+"/interfaces")
	// res, err := a.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()
	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	_dumpResponse(res)
	// 	return []map[string]string{}
	// }
	// data := readResponseBody(res)
	interfaces := []map[string]string{}
	// ifs := data.(map[string]interface{})["results"]
	// for _, if_data := range ifs.([]interface{}) {
	// 	if_data_map := map[string]string{
	// 		"id":           if_data.(map[string]interface{})["id"].(string),
	// 		"edge_path":    if_data.(map[string]interface{})["edge_path"].(string),
	// 		"name":         if_data.(map[string]interface{})["display_name"].(string),
	// 		"segment_path": if_data.(map[string]interface{})["segment_path"].(string),
	// 		"uid":          if_data.(map[string]interface{})["unique_id"].(string),
	// 	}
	// 	interfaces = append(interfaces, if_data_map)
	// }

	return interfaces
}

func (a *Agent) GetInterfaceStatistics(tier0_id string, locale_service_id string, if_data map[string]string) RouterStats {
	// Quary Params:
	//   enforcement_point_path
	//   edge_path : ex edge_path=/infra/sites/default/enforcement-points/default/edge-clusters/57d2c653-4d63-48d8-b188-40b4e45a9bc8/edge-nodes/2ed9af04-21c9-11e9-be65-000c2902dff7
	// req := a.makeRequest("GET", "/policy/api/v1/infra/tier-0s/"+tier0_id+"/locale-services/"+locale_service_id+"/interfaces/"+if_data["id"]+"/statistics")
	// params := url.Values{
	// 	"enforcement_point_path": []string{"/infra/sites/default/enforcement-points/default"},
	// 	"edge_path":              []string{if_data["edge_path"]},
	// }
	// req.URL.RawQuery = params.Encode()
	// res, err := a.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()
	// //sites := a.GetSite()
	// //for _, site := range sites {
	// //	for _, ep_id := range a.GetEnforcementPoint(site) {
	// //		a.GetEdgeClusterUnderEnforcementPoint(site, ep_id)
	// //	}
	// //}
	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	_dumpResponse(res)
	// 	return RouterStats{}
	// }
	// data := readResponseBody(res)
	// jsonStr, err := json.Marshal(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var ifStats RouterStats
	// json.Unmarshal(jsonStr, &ifStats)

	return ifStats
}

func (a *Agent) GetGatewayInterfaceStats(gw Tier0Gateway) map[int]RouterStats {
	// locale_service_id := a.GetLocaleService(gw.Id)
	// interfaces := a.GetInterface(gw.Id, locale_service_id[0])
	// mutex := &sync.Mutex{}
	result := make(map[int]RouterStats)
	// var wg sync.WaitGroup
	// for idx, inf := range interfaces {
	// 	wg.Add(1)
	// 	go func(idx int, i map[string]string) {
	// 		defer wg.Done()
	// 		mutex.Lock()
	// 		result[idx] = a.GetInterfaceStatistics(gw.Id, locale_service_id[0], i)
	// 		mutex.Unlock()
	// 	}(idx, inf)
	// }
	// wg.Wait()
	return result
}

func (a *Agent) GetRoutingTable(tier0Id string) []EdgeRoute {
	// var path string
	// path = "/policy/api/v1/infra/tier-0s/" + tier0Id + "/routing-table"
	// res := a.Request("GET", path, nil, nil)
	entries := []EdgeRoute{}
	// for _, ed := range res.Body.(map[string]interface{})["results"].([]interface{}) {
	// 	str, _ := json.Marshal(ed)
	// 	var edgeRoutes EdgeRoute
	// 	json.Unmarshal(str, &edgeRoutes)
	// 	entries = append(entries, edgeRoutes)
	// }

	return entries
}

func (a *Agent) GetBgpConfig(tier0Id string, locale string) BgpConfig {
	// var path string
	// path = "/policy/api/v1/infra/tier-0s/" + tier0Id + "/locale-services/" + locale + "/bgp"
	// res := a.Request("GET", path, nil, nil)
	bgpConfig := BgpConfig{}
	// body, _ := res.BodyBytes()
	// json.Unmarshal(body, &bgpConfig)

	return bgpConfig
}

func (a *Agent) GetBgpNeighbors(tier0Id string, locale string) []BgpNeighbor {
	// var path string
	// path = "/policy/api/v1/infra/tier-0s/" + tier0Id + "/locale-services/" + locale + "/bgp/neighbors"
	// res := a.Request("GET", path, nil, nil)
	neighbors := []BgpNeighbor{}
	// for _, nb := range res.Body.(map[string]interface{})["results"].([]interface{}) {
	// 	str, _ := json.Marshal(nb)
	// 	var neighbor BgpNeighbor
	// 	json.Unmarshal(str, &neighbor)
	// 	neighbors = append(neighbors, neighbor)
	// }

	return neighbors
}

func (a *Agent) GetBgpNeighborsAdvRoutes(path string) []EdgeBgpAdvRoute {
	// path = "/policy/api/v1" + path + "/advertised-routes"
	// res := a.Request("GET", path, nil, nil)
	edges := []EdgeBgpAdvRoute{}
	// for _, nb := range res.Body.(map[string]interface{})["results"].([]interface{}) {
	// 	es := nb.(map[string]interface{})["edge_node_routes"].([]interface{})
	// 	for _, e := range es {
	// 		str, _ := json.Marshal(e)
	// 		var edge EdgeBgpAdvRoute
	// 		json.Unmarshal(str, &edge)
	// 		edges = append(edges, edge)
	// 	}
	// }

	return edges
}

func (a *Agent) GetGatewayAggregateInfo(gw_realization_id string) []map[string]string {
	// path := "/api/v1/ui-controller/l3/logical-routers/" + gw_realization_id + "/status/aggregate-info?source=realtime"
	// res := a.Request("GET", path, nil, nil)
	// if res == nil {
	// 	return nil
	// }

	// per_node_status := res.Body.(map[string]interface{})["status"].(map[string]interface{})["per_node_status"]
	// if per_node_status == nil {
	// 	//if a tier-0 is created but interface isn't set, per_node_status will be empty
	// 	return nil
	// }

	result := []map[string]string{}
	// for _, st := range per_node_status.([]interface{}) {
	// 	str, _ := json.Marshal(st)
	// 	var status map[string]string
	// 	json.Unmarshal(str, &status)
	// 	result = append(result, status)
	// }
	return result
}

func (a *Agent) GetTier0GatewayFromName(name string) (Tier0Gateway, error) {
	// gws := a.GetTier0Gateway("")
	// for _, gw := range gws {
	// 	if gw.Name == name {
	// 		return gw, nil
	// 	}
	// }
	return Tier0Gateway{}, fmt.Errorf("Error: Tier-0 gateway '%s' is not found", name)
}

func (a *Agent) GetTier1GatewayFromName(name string) (Tier1Gateway, error) {
	// gws := a.GetTier1Gateway("")
	// for _, gw := range gws {
	// 	if gw.Name == name {
	// 		return gw, nil
	// 	}
	// }
	return Tier1Gateway{}, fmt.Errorf("Error: Tier-1 gateway '%s' is not found", name)
}
