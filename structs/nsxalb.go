package structs

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

type Cloud struct {
	Id        string `json:"uuid"`
	Name      string `json:"name"`
	TenantRef string `json:"tenant_ref"`
}

type ServiceEngineGroup struct {
	Id           string `json:"uuid"`
	Name         string `json:"name"`
	HaMode       string `json:"ha_mode"`
	Buffer       int    `json:"buffer_se"`
	VsLimitPerSe int    `json:"max_vs_per_se"`
	TenantRef    string `json:"tenant_ref"`
	SE           []ServiceEngine
}

type ServiceEngine struct {
	Primary    bool              `json:"is_primary"`
	Secondary  bool              `json:"is_secondary"`
	Memory     int               `json:"memory"`
	Address    map[string]string `json:"mgmt_ip"`
	SeRef      string            `json:"se_ref"`
	Cpu        int               `json:"vcpus"`
	Interfaces []SeInterface     `json:"vip_intf_list"`
	VipMac     string            `json:"vip_intf_mac"`
	VipMask    int               `json:"vip_intf_mask"`
}

type SeInterface struct {
	Address map[string]string `json:"vip_intf_ip"`
	Mac     string            `json:"vip_intf_mac"`
	Vlan    int               `json:"vlan_id"`
}

type VSResult struct {
	VirtualServiceInventories []VirtualServiceInventory `json:"results"`
}

type VirtualServiceInventory struct {
	Config  VirtualService `json:"config"`
	Runtime VSRuntime      `json:"runtime"`
}

type VSRuntime struct {
	PersentSEsUp int          `json:"percent_ses_up"`
	VipSummary   []VipSummary `json:"vip_summary"`
}

type VipSummary struct {
	Id     string              `json:"vip_id"`
	Status map[string]string   `json:"oper_status"`
	Se     []map[string]string `json:"service_engine"`
}

type VirtualService struct {
	Type        string   `json:"type"`
	UUID        string   `json:"uuid"`
	Tenant      string   `json:"tenant_ref"`
	Name        string   `json:"name"`
	Ports       []VSPort `json:"services"`
	CloudRef    string   `json:"cloud_ref"`
	SeGroupRef  string   `json:"se_group_ref"`
	Vips        []Vip    `json:"vip"`
}

func (v *VirtualServiceInventory) Print(w *tabwriter.Writer) {
	ports := ""
	for i, p := range v.Config.Ports {
		summary := p.GetSummary()
		ports += summary
		if i != len(v.Config.Ports)-1 {
			ports += ","
		}
	}
	vips := ""
	for i, vip := range v.Config.Vips {
		vips += vip.Address["addr"]
		if i != len(v.Config.Vips)-1 {
			ports += ","
		}
	}
	cloud := strings.Split(v.Config.CloudRef, "#")
	segroup := strings.Split(v.Config.SeGroupRef, "#")
	w.Write([]byte(strings.Join([]string{v.Config.UUID, v.Config.Name, vips, ports, cloud[1], segroup[1]}, "\t") + "\n"))
}

func (v *VirtualService) GetCloudId() string {
	path := strings.Split(v.CloudRef, "/")

	return path[len(path)-1]
}

func (v *VirtualService) GetSegId() string {
	path := strings.Split(v.SeGroupRef, "/")

	return path[len(path)-1]
}

type Vip struct {
	Id      string            `json:"vip_id"`
	Address map[string]string `json:"ip_address"`
}

type VipRuntime struct {
	Se []ServiceEngine `json:"se_list"`
}

type VSPort struct {
	Ssl       bool `json:"enable_ssl"`
	Port      uint `json:"port"`
	PortRange uint `json:"port_range_end"`
}

func (p *VSPort) GetSummary() string {
	var summary string
	if p.Port != p.PortRange {
		summary = fmt.Sprintf("%v-%v", p.Port, p.PortRange)
	} else {
		summary = fmt.Sprintf("%v", p.Port)
	}
	if p.Ssl {
		summary += fmt.Sprintf("(SSL)")
	}

	return summary
}
