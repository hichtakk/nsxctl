package nsxalb

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxAlbClient) ShowVirtualService() {
	resp := c.Request("GET", "/api/virtualservice", map[string]string{}, nil)
	var results structs.VSResult
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &results)
	for _, vs := range results.VirtualServices {
		name := vs.Name
		uuid := vs.UUID
		ip := vs.Address
		ports := vs.Ports
		var ports_str []string
		for _, p := range ports {
			port_num := fmt.Sprintf("%d", p.Port)
			if p.Port != p.PortRange {
				port_num = fmt.Sprintf("%s-%d", port_num, p.PortRange)
			}
			if p.Ssl == true {
				port_num = fmt.Sprintf("%s(ssl)", port_num)
			}
			ports_str = append(ports_str, port_num)
		}
		fmt.Printf("%s  %-10s  %s  %s\n", uuid, name, ip["addr"], strings.Join(ports_str, ","))
	}
}
