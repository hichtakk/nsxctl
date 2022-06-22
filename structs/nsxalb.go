package structs

type VSResult struct {
	VirtualServices []VirtualService `json:"results"`
}

type VirtualService struct {
	Type    string            `json:"type"`
	UUID    string            `json:"uuid"`
	Tenant  string            `json:"tenant_ref"`
	Name    string            `json:"name"`
	Address map[string]string `json:"ip_address"`
	Ports   []VSPort          `json:"services"`
}

type VSPort struct {
	Ssl       bool `json:"enable_ssl"`
	Port      uint `json:"port"`
	PortRange uint `json:"port_range_end"`
}
