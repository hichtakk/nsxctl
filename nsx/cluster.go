package nsx

func (a *Agent) ClusterInfo() {
	// res := a.Request("GET", "/api/v1/cluster/status", nil, nil)
	// clusterId := res.Body.(map[string]interface{})["cluster_id"].(string)
	// fmt.Println("Cluster ID:", clusterId)
	// for _, node := range res.Body.(map[string]interface{})["mgmt_cluster_status"].(map[string]interface{})["online_nodes"].([]interface{}) {
	// 	fmt.Println()
	// 	fmt.Println("Node IP:   ", node.(map[string]interface{})["mgmt_cluster_listen_ip_address"])
	// 	fmt.Println("Node ID:   ", node.(map[string]interface{})["uuid"])
	// }
}
