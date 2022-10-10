package nsxalb

import (
	"encoding/json"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxAlbClient) GetHealthMonitors() []structs.HealthMonitor {
	res := c.Request("GET", "/api/healthmonitor", map[string]string{}, nil)
	var results []structs.HealthMonitor
	str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	json.Unmarshal(str, &results)

	return results
}
