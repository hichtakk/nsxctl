package nsxalb

func (a *Agent) GetGslbSites() Gslb {
	// res := c.Request("GET", "/api/gslb", map[string]string{}, nil)
	var results []Gslb
	// str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	// json.Unmarshal(str, &results)

	return results[0]
}

func (a *Agent) GetGslbServices() GslbServices {
	// res := c.Request("GET", "/api/gslbservice", map[string]string{}, nil)
	var results GslbServices
	// str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	// json.Unmarshal(str, &results)

	return results
}
