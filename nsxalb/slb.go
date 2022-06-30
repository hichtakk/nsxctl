package nsxalb

import (
	"encoding/json"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxAlbClient) ShowVirtualService() []structs.VirtualService {
	resp := c.Request("GET", "/api/virtualservice", map[string]string{}, nil)
	var results structs.VSResult
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &results)
	var vss []structs.VirtualService
	for _, vs := range results.VirtualServices {
		vss = append(vss, vs)
	}

	return vss
}
