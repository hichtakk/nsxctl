package nsxalb

import (
	"encoding/json"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxAlbClient) GetGslbSites() structs.Gslb {
	res := c.Request("GET", "/api/gslb", map[string]string{}, nil)
	var results []structs.Gslb
	str, _ := json.Marshal(res.Body.(map[string]interface{})["results"].([]interface{}))
	json.Unmarshal(str, &results)

	return results[0]
}
