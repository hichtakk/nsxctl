package nsxalb

import (
	"encoding/json"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxAlbClient) ShowVirtualService() []structs.VirtualServiceInventory {
	resp := c.Request("GET", "/api/virtualservice-inventory/?include_name=true", map[string]string{}, nil)
	var results structs.VSResult
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &results)
	var vss []structs.VirtualServiceInventory
	for _, vs := range results.VirtualServiceInventories {
		vss = append(vss, vs)
	}

	return vss
}
