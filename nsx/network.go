package nsx

func (a *Agent) GetIpPool() IpPools {
	// path := "/policy/api/v1/infra/ip-pools"
	// res := c.Request("GET", path, nil, nil)
	ipps := IpPools{}
	// str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	// json.Unmarshal(str, &ipps)

	return ipps
}

func (a *Agent) GetIpBlock() IpBlocks {
	// path := "/policy/api/v1/infra/ip-blocks"
	// res := c.Request("GET", path, nil, nil)
	ipbs := IpBlocks{}
	// for _, b := range res.Body.(map[string]interface{})["results"].([]interface{}) {
	// 	str, _ := json.Marshal(b)
	// 	var block IpBlock
	// 	json.Unmarshal(str, &block)
	// 	ipbs = append(ipbs, block)
	// }

	return ipbs
}

func (a *Agent) CreateIpPool(name string) {
	// path := "/policy/api/v1/infra/ip-pools"
	// reqData := make(map[string]string)
	// reqData["display_name"] = name
	// reqJson, _ := json.Marshal(reqData)
	// req, _ := http.NewRequest("PATCH", c.BaseUrl+path+"/"+name, bytes.NewBuffer(reqJson))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Xsrf-Token", c.Token)
	// res, err := c.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()

	// // no content returned when request succeeded
	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	fmt.Println(res)
	// 	return
	// }
}

func (a *Agent) CreateIpBlock(name string, cidr string) {
	// path := "/policy/api/v1/infra/ip-blocks"
	// reqData := make(map[string]string)
	// reqData["display_name"] = name
	// reqData["cidr"] = cidr
	// reqJson, _ := json.Marshal(reqData)
	// req, _ := http.NewRequest("PATCH", c.BaseUrl+path+"/"+name, bytes.NewBuffer(reqJson))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Xsrf-Token", c.Token)
	// res, err := c.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()

	// // no content returned when request succeeded
	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	fmt.Println(res)
	// 	return
	// }
}

func (a *Agent) DeleteIpPool(name string) {
	// path := "/policy/api/v1/infra/ip-pools"
	// req, _ := http.NewRequest("DELETE", c.BaseUrl+path+"/"+name, nil)
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Xsrf-Token", c.Token)
	// res, err := c.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()

	// // no content returned when request succeeded
	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	fmt.Println(res)
	// 	return
	// }
}

func (a *Agent) DeleteIpBlock(name string) {
	// path := "/policy/api/v1/infra/ip-blocks"
	// req, _ := http.NewRequest("DELETE", c.BaseUrl+path+"/"+name, nil)
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Xsrf-Token", c.Token)
	// res, err := c.httpClient.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer res.Body.Close()

	// // no content returned when request succeeded
	// if res.StatusCode != 200 {
	// 	fmt.Printf("StatusCode=%d\n", res.StatusCode)
	// 	fmt.Println(res)
	// 	return
	// }
}
