package nsxalb

import (
	"encoding/json"
	"log"
	"strings"

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

func (c *NsxAlbClient) GetEvhHostname(vsId string) string {
	resp := c.Request("GET", "/api/virtualservice/" + vsId, nil, nil)
	resByte, _ := resp.BodyBytes()
	obj := &struct{
		Markers []map[string]any `json:"markers"`
	}{}
	err := json.Unmarshal(resByte, obj)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range obj.Markers {
		if v, ok := m["key"]; ok {
			if v.(string) == "Host" {
				if v, ok = m["values"]; ok {
					hosts := []string{}
					for _, s := range v.([]interface{}) {
						hosts = append(hosts, s.(string))
					}
					return strings.Join(hosts, ",")
				}
			}
		}
	}
	return ""
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
