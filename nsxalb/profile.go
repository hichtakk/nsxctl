package nsxalb

func (a *Agent) GetHealthMonitors() []HealthMonitor {
	// res := a.Request("GET", "/api/healthmonitor", map[string]string{}, nil)
	var results []HealthMonitor
	// str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	// json.Unmarshal(str, &results)

	return results
}
