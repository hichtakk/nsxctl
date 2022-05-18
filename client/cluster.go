package client

import (
	"fmt"
)

func (c *NsxtClient) ClusterInfo() {
	res := c.Request("GET", "/api/v1/cluster/status", nil, nil)
	fmt.Println(res)
}
