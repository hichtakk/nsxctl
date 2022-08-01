package structs

import (
	"fmt"
	"os"
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

type SEResult struct {
	ServiceEngineInventories []ServiceEngineInventory `json:"results"`
}

type ServiceEngineInventory struct {
	Config  ServiceEngineConfig  `json:"config"`
	Health  map[string]int       `json:"health_score"`
	Runtime ServiceEngineRuntime `json:"runtime"`
}

func (sei *ServiceEngineInventory) Print(w *tabwriter.Writer) {
	id := sei.Config.UUID
	name := sei.Config.Name
	ip := sei.Config.Address["addr"]
	// health := sei.Health["health_score"]
	status := strings.Split(sei.Runtime.Status.State, "_")[1]
	// reason := strings.Join(sei.Runtime.Status.Reason, "\n")
	cloud := (strings.Split(sei.Config.CloudRef, "#"))[1]
	segroup := (strings.Split(sei.Config.SEGroupRef, "#"))[1]

	w.Write([]byte(strings.Join([]string{id, name, ip, cloud, segroup, status}, "\t") + "\n"))
}

type ServiceEngineConfig struct {
	CloudRef   string            `json:"cloud_ref"`
	Address    map[string]string `json:"mgmt_ip_address"`
	Name       string            `json:"name"`
	SEGroupRef string            `json:"se_group_ref"`
	TenantRef  string            `json:"tenant_ref"`
	UUID       string            `json:"uuid"`
}

type ServiceEngineRuntime struct {
	GatewayUp    bool                `json:"gateway_up"`
	MigrateState string              `json:"migrate_state"`
	PowerState   string              `json:"power_state"`
	Status       ServiceEngineStatus `json:"oper_status"`
}

type ServiceEngineStatus struct {
	State  string   `json:"state"`
	Reason []string `json:"reason"`
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
	Health  map[string]int `json:"heal`
	Runtime VSRuntime      `json:"runtime"`
}

type VSRuntime struct {
	PersentSEsUp int                  `json:"percent_ses_up"`
	VipSummary   []VipSummary         `json:"vip_summary"`
	Status       VirtualServiceStatus `json:"oper_status"`
}

type VipSummary struct {
	Id     string                   `json:"vip_id"`
	Status map[string]string        `json:"oper_status"`
	Se     []map[string]interface{} `json:"service_engine"`
}

type VirtualService struct {
	Type       string   `json:"type"`
	UUID       string   `json:"uuid"`
	Tenant     string   `json:"tenant_ref"`
	Name       string   `json:"name"`
	Ports      []VSPort `json:"services"`
	CloudRef   string   `json:"cloud_ref"`
	SeGroupRef string   `json:"se_group_ref"`
	Vips       []Vip    `json:"vip"`
}

type VirtualServiceStatus struct {
	State  string   `json:"state"`
	Reason []string `json:"reason"`
}

func (v *VirtualServiceInventory) Print(w *tabwriter.Writer, verbose bool) {
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
	status := strings.Split(v.Runtime.Status.State, "_")[1]
	// reason := strings.Join(v.Runtime.Status.Reason, "\n")
	var seNames []string
	for _, vip := range v.Runtime.VipSummary {
		for _, se := range vip.Se {
			name := strings.Split(se["url"].(string), "#")[1]
			primary := se["primary"].(bool)
			// standby := se["standby"].(bool)
			if primary {
				name = name + "(p)"
			}
			seNames = append(seNames, name)
		}
	}
	if verbose {
		w.Write([]byte(strings.Join([]string{v.Config.UUID, v.Config.Name, vips, ports, cloud[1], segroup[1], status, strings.Join(seNames, ", ")}, "\t") + "\n"))
	} else {
		w.Write([]byte(strings.Join([]string{v.Config.UUID, v.Config.Name, vips, ports, cloud[1], segroup[1], status}, "\t") + "\n"))
	}
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

type LicensingLedger struct {
	Id         string      `json:"uuid"`
	SeInfos    []SeInfo    `json:"se_infos"`
	TierUsages []TierUsage `json:"tier_usages"`
}

type SeInfo struct {
	Id           string `json:"uuid"`
	LastUpdated  int    `json:"last_updated"`
	ServiceCores int    `json:"service_cores"`
	TenantUUID   string `json:"tenant_uuid"`
	Tier         string `json:"tier"`
}

type TierUsage struct {
	Tier  string         `json:"tier"`
	Usage map[string]int `json:"usage"`
}

type SystemConfiguration struct {
	Id          string `json:"uuid"`
	LicenseTier string `json:"default_license_tier"`
	TierUsage   TierUsage
}

type Gslb struct {
	Id                string              `json:"uuid"`
	Domains           []map[string]string `json:"dns_configs"`
	Name              string              `json:"name"`
	ReplicationPolicy map[string]string   `json:"replication_policy"`
	GslbSites         []GslbSite          `json:"sites"`
	ThirdPartySites   []ThirdPartySite    `json:"third_party_sites"`
	LeaderUuid        string              `json:"leader_cluster_uuid"`
}

func (g *Gslb) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"Name", "Type", "IP Address"}, "\t") + "\n"))
	for _, s := range g.GslbSites {
		w.Write([]byte(strings.Join([]string{s.Name, s.Type, s.Address[0]["addr"]}, "\t") + "\n"))
	}
	w.Flush()
}

type GslbSite struct {
	Id      string              `json:"cluster_uuid"`
	Name    string              `json:"name"`
	Enabled bool                `json:"enabled"`
	Address []map[string]string `json:"ip_addresses"`
	Type    string              `json:"member_type"`
	DnsVs   []DnsVs             `json:"dns_vses"`
}

type ThirdPartySite struct {
	Id      string `json:"cluster_uuid"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type DnsVs struct {
	Id      string   `json:"dns_vs_uuid"`
	Domains []string `json:"domain_names"`
}
