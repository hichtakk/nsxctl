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

func (c *NsxAlbClient) GetPools() []structs.PoolInventory {
	path := "/api/pool-inventory?include_name=true"
	resp := c.Request("GET", path, nil, nil)
	var results structs.PoolResult
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &results)
	var pools []structs.PoolInventory
	for _, p := range results.PoolInventories {
		pools = append(pools, p)
	}

	return pools
}

func (c *NsxAlbClient) GetPool(id string) structs.Pool {
	path := "/api/pool/" + id + "?include_name=true"
	resp := c.Request("GET", path, nil, nil)
	var pool structs.Pool
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &pool)

	return pool
}
