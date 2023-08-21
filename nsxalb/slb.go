package nsxalb

func (a *Agent) ShowVirtualService() []VirtualServiceInventory {
	// resp := c.Request("GET", "/api/virtualservice-inventory/?include_name=true", map[string]string{}, nil)
	// var results VSResult
	// resByte, _ := resp.BodyBytes()
	// json.Unmarshal(resByte, &results)
	var vss []VirtualServiceInventory
	// for _, vs := range results.VirtualServiceInventories {
	// 	vss = append(vss, vs)
	// }

	return vss
}

func (a *Agent) GetPools() []PoolInventory {
	// path := "/api/pool-inventory?include_name=true"
	// resp := c.Request("GET", path, nil, nil)
	// var results PoolResult
	// resByte, _ := resp.BodyBytes()
	// json.Unmarshal(resByte, &results)
	var pools []PoolInventory
	// for _, p := range results.PoolInventories {
	// 	pools = append(pools, p)
	// }

	return pools
}

func (a *Agent) GetPool(id string) Pool {
	// path := "/api/pool/" + id + "?include_name=true"
	// resp := c.Request("GET", path, nil, nil)
	var pool Pool
	// resByte, _ := resp.BodyBytes()
	// json.Unmarshal(resByte, &pool)

	return pool
}
