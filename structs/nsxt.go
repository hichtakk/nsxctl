package structs

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/netip"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type EnforcementPoint struct {
	Id   string
	Path string
}

type EdgeClusterMember struct {
	Index int    `json:"member_index"`
	Id    string `json:"transport_node_id"`
	Name  string
}

type EdgeCluster struct {
	Id             string              `json:"id"`
	Name           string              `json:"display_name"`
	Type           string              `json:"deployment_type"`
	Members        []EdgeClusterMember `json:"members"`
	MemberNodeType string              `json:"member_node_type"`
}

type EdgeClusters []EdgeCluster

func (ecs *EdgeClusters) Print(edgeName map[string]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name", "EdgeNode"}, "\t") + "\n"))
	for _, ec := range *ecs {
		edge := []string{}
		for _, e := range ec.Members {
			edge = append(edge, edgeName[e.Id])
		}
		edgeStr := strings.Join(edge, ",")
		w.Write([]byte(strings.Join([]string{ec.Id, ec.Name, edgeStr}, "\t") + "\n"))
	}
	w.Flush()
}

func (ecs *EdgeClusters) GetClusterById(Id string) *EdgeCluster {
	for _, ec := range *ecs {
		if ec.Id == Id {
			return &ec
		}
	}
	return &EdgeCluster{}
}

type RouteEntry struct {
	Type       string `json:"route_type"`
	Network    string `json:"network"`
	NextHop    string `json:"next_hop"`
	Ad         uint8  `json:"admin_distance"`
	RouterId   string `json:"lr_component_id"`
	RouterType string `json:"lr_component_type"`
}

type EdgeRoute struct {
	NodePath string       `json:"edge_node"`
	Entries  []RouteEntry `json:"route_entries"`
}

func (er *EdgeRoute) Print() {
	for _, e := range er.Entries {
		var routeType string
		switch e.Type {
		case "t0c":
			routeType = "C"
		case "t0s":
			routeType = "S"
		case "b":
			routeType = "B"
		case "t0n":
			routeType = "N"
		case "t1c":
			routeType = "c"
		case "t1s":
			routeType = "s"
		case "t1n":
			routeType = "n"
		case "t1l":
			routeType = "l"
		case "t1ls":
			routeType = "ln"
		case "t1d":
			routeType = "d"
		case "t1ipsec":
			routeType = "p"
		case "isr":
			routeType = "i"
		}
		if routeType == "C" {
			fmt.Printf("%v> %v is directly connected\n", routeType, e.Network)
		} else if routeType == "i" && e.NextHop == "" {
			fmt.Printf("%v> %v [%v] blackhole\n", routeType, e.Network, e.Ad)
		} else {
			fmt.Printf("%v> %v [%v] via %v\n", routeType, e.Network, e.Ad, e.NextHop)
		}
	}
	fmt.Println()
}

func (er *EdgeRoute) GetEdgeClusterId() string {
	path := strings.Split(er.NodePath, "/")
	return path[7]
}

func (er *EdgeRoute) GetEdgeClusterNodeIdx() int {
	path := strings.Split(er.NodePath, "/")
	idx, _ := strconv.Atoi(path[9])
	return idx
}

func (er *EdgeRoute) GetEntries(version int) RouteEntries {
	var entries []RouteEntry
	var bitLen int
	if version == 6 {
		bitLen = 128
	} else {
		bitLen = 32
	}
	// filter whether IPv4 or IPv6
	for _, e := range er.Entries {
		eip, _ := netip.ParsePrefix(e.Network)
		if eip.Addr().BitLen() != bitLen {
			continue
		}
		entries = append(entries, e)
	}
	// check addresssing order
	nthSmall := make([]int, len(entries))
	for idx, en := range entries {
		small := 0
		for _, e := range entries {
			en_prefix, _ := netip.ParsePrefix(en.Network)
			e_prefix, _ := netip.ParsePrefix(e.Network)
			if e_prefix.Addr().Compare(en_prefix.Addr()) < 0 {
				small += 1
			}
		}
		nthSmall[idx] = small
	}
	sorted_entries := make([]RouteEntry, len(entries))
	for idx, se := range nthSmall {
		sorted_entries[se] = entries[idx]
	}
	return RouteEntries(sorted_entries)
}

type RouteEntries []RouteEntry

func (res *RouteEntries) Print() {
	for _, e := range *res {
		var routeType string
		switch e.Type {
		case "t0c":
			routeType = "C"
		case "t0s":
			routeType = "S"
		case "b":
			routeType = "B"
		case "t0n":
			routeType = "N"
		case "t1c":
			routeType = "c"
		case "t1s":
			routeType = "s"
		case "t1n":
			routeType = "n"
		case "t1l":
			routeType = "l"
		case "t1ls":
			routeType = "ln"
		case "t1d":
			routeType = "d"
		case "t1ipsec":
			routeType = "p"
		case "isr":
			routeType = "i"
		}
		if routeType == "C" {
			fmt.Printf("%v> %v is directly connected\n", routeType, e.Network)
		} else if routeType == "i" && e.NextHop == "" {
			fmt.Printf("%v> %v [%v] blackhole\n", routeType, e.Network, e.Ad)
		} else {
			fmt.Printf("%v> %v [%v] via %v\n", routeType, e.Network, e.Ad, e.NextHop)
		}
	}
	fmt.Println()
}

type BgpConfig struct {
	Id              string `json:"id"`
	Name            string `json:"display_name"`
	Ecmp            bool   `json:"ecmp"`
	Enabled         bool   `json:"enabled"`
	GracefulRestart bool   `json:"graceful_restart"`
	InterSrRouting  bool   `json:"inter_sr_ibgp"`
	Asn             string `json:"local_as_num"`
}

type BgpAdvRouteEntry struct {
	AsPath    string `json:"as_path"`
	LocalPref int    `json:"local_pref"`
	Med       int    `json:"med"`
	Network   string `json:"network"`
	NextHop   string `json:"next_hop"`
	Weight    int    `json:"weight"`
}

type BgpAdvRouteEntries []BgpAdvRouteEntry

func (ar *BgpAdvRouteEntries) Print() {
	fmt.Printf("%-8s	 %-8s	%-8s	%-8s	%-5s\n", "Network", "Next Hop", "Metric", "Local Pref", "Path")
	for _, e := range *ar {
		fmt.Printf("%-8s	%8s	%6d		%10d %5s\n", e.Network, e.NextHop, e.Med, e.LocalPref, e.AsPath)
	}
}

type EdgeBgpAdvRoute struct {
	Source  string             `json:"source_address"`
	EdgeId  string             `json:"transport_node_id"`
	Entries []BgpAdvRouteEntry `json:"routes"`
}

func (ar *EdgeBgpAdvRoute) Print() {
	for _, e := range ar.Entries {
		fmt.Printf("%-8s	%-8s	%-8s	%-8s %-5s\n", "Network", "Next Hop", "Metric", "Local Pref", "Path")
		fmt.Printf("%-8s	%8s	%8d	%8d %5s\n", e.Network, e.NextHop, e.Med, e.LocalPref, e.AsPath)
	}
}

func (ar *EdgeBgpAdvRoute) GetEntries() BgpAdvRouteEntries {
	var entries []BgpAdvRouteEntry
	entries = ar.Entries

	// check addresssing order
	nthSmall := make([]int, len(entries))
	for idx, en := range entries {
		small := 0
		for _, e := range entries {
			en_prefix, _ := netip.ParsePrefix(en.Network)
			e_prefix, _ := netip.ParsePrefix(e.Network)
			if e_prefix.Addr().Compare(en_prefix.Addr()) < 0 {
				small += 1
			}
		}
		nthSmall[idx] = small
	}
	sorted_entries := make([]BgpAdvRouteEntry, len(entries))
	for idx, se := range nthSmall {
		sorted_entries[se] = entries[idx]
	}
	return BgpAdvRouteEntries(sorted_entries)
}

type ComputeManager struct {
	Id     string
	Name   string
	Type   string
	Server string
	Detail string
	Status *ComputeManagerStatus
}

func (cm *ComputeManager) Print() {
	fmt.Println("Name:    ", cm.Name)
	fmt.Println("ID:      ", cm.Id)
	fmt.Println("Type:    ", cm.Type)
	fmt.Println("FQDN/IP: ", cm.Server)
	fmt.Println("Version: ", cm.Detail)
	if cm.Status != nil {
		s := fmt.Sprintf("Status:   %s, %s", cm.Status.Connection, cm.Status.Registration)
		fmt.Println(s)
	}
}

type ComputeManagerStatus struct {
	Connection   string
	Registration string
}

type TransportZones []TransportZone

func (tzs *TransportZones) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name", "Type"}, "\t") + "\n"))
	for _, tz := range *tzs {
		w.Write([]byte(strings.Join([]string{tz.Id, tz.Name, tz.Type}, "\t") + "\n"))
	}
	w.Flush()
}

type TransportZone struct {
	Id   string `json:"id"`
	Name string `json:"display_name"`
	Type string `json:"tz_type"`
	Path string `json:"path"`
}

type TransportNodes []TransportNode

func (tns *TransportNodes) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name", "IP", "Tunnel"}, "\t") + "\n"))
	for _, tn := range *tns {
		ip := strings.Join(tn.EdgeNodeDeploymentInfo.IPAddress, ",")
		numTun := fmt.Sprintf("%v", len(tn.Tunnels))
		w.Write([]byte(strings.Join([]string{tn.Id, tn.Name, ip, numTun}, "\t") + "\n"))
	}
	w.Flush()
}

type TransportNode struct {
	Id                     string                 `json:"id"`
	Name                   string                 `json:"display_name"`
	HostSwitchSpec         HostSwitchSpec         `json:"host_switch_spec"`
	EdgeNodeDeploymentInfo EdgeNodeDeploymentInfo `json:"node_deployment_info"`
	ResourceType           string                 `json:"resource_type"`
	Tunnels                TransportNodeTunnels   `json:",omitempty"`
}

type TransportNodeTunnels []TransportNodeTunnel

type TransportNodeTunnel struct {
	Name            string `json:"name"`
	Status          string `json:"status"`
	Encapsulation   string `json:"encap"`
	EgressInterface string `json:"egress_interface"`
	LocalIp         string `json:"local_ip"`
	RemoteIp        string `json:"remote_ip"`
	RemoteNodeName  string `json:"remote_node_display_name"`
	RemoteNodeId    string `json:"remote_node_id"`
}

type TransportNodeProfiles []TransportNodeProfile

func (tnps *TransportNodeProfiles) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name"}, "\t") + "\n"))
	for _, tnp := range *tnps {
		w.Write([]byte(strings.Join([]string{tnp.Id, tnp.Name}, "\t") + "\n"))
	}
	w.Flush()
}

type TransportNodeProfile struct {
	Id             string         `json:"id"`
	Name           string         `json:"display_name"`
	HostSwitchSpec HostSwitchSpec `json:"host_switch_spec"`
}

type HostSwitchSpec struct {
	HostSwitches []HostSwitch `json:"host_switches"`
	ResourceType string       `json:"resource_type"`
}

type HostSwitch struct {
	Mode                   string                   `json:"host_switch_mode"`
	Name                   string                   `json:"host_switch_name"`
	UplinkProfiles         []map[string]string      `json:"host_switch_profile_ids"`
	Type                   string                   `json:"host_switch_type"`
	IPAssignment           IpAssignmentSpec         `json:"ip_assignment_spec"`
	Pnics                  []map[string]string      `json:"pnics"`
	TransportZoneEndpoints []TransportZoneEndpoints `json:"transport_zone_endpoints"`
}

type TransportZoneEndpoints struct {
	TransportZoneId         string              `json:"transport_zone_id"`
	TransportZoneProfileIds []map[string]string `json:"transport_zone_profile_ids"`
}

type IpAssignmentSpec struct {
	ResourceType string `json:"resource_type"`
	IpPoolId     string `json:"ip_pool_id,omitempty"`
}

type EdgeNodeDeploymentInfo struct {
	Name                 string               `json:"display_name"`
	IPAddress            []string             `json:"ip_addresses"`
	EdgeDeploymentConfig EdgeDeploymentConfig `json:"deployment_config"`
	NodeSettings         NodeSettings         `json:"node_settings"`
	ResourceType         string               `json:"resource_type"`
}

type EdgeDeploymentConfig struct {
	Size               string             `json:"form_factor"`
	Users              map[string]string  `json:"node_user_settings"`
	VMDeploymentConfig VMDeploymentConfig `json:"vm_deployment_config"`
}

type VMDeploymentConfig struct {
	ComputeId             string          `json:"compute_id"`
	DataNetworkIds        []string        `json:"data_network_ids"`
	DefaultGateway        []string        `json:"default_gateway_addresses"`
	ManagementNetworkId   string          `json:"management_network_id"`
	ManagementPortSubnets []Subnet        `json:"management_port_subnets"`
	ReservationInfo       ReservationInfo `json:"reservation_info"`
	StorageId             string          `json:"storage_id"`
	VcId                  string          `json:"vc_id"`
	PlacementType         string          `json:"placement_type"`
}

type Subnet struct {
	IPAddresses  []string `json:"ip_addresses"`
	PrefixLength int      `json:"prefix_length"`
}

type ReservationInfo struct {
	Cpu    CpuReservationInfo    `json:"cpu_reservation"`
	Memory MemoryReservationInfo `json:"memory_reservation"`
}

type CpuReservationInfo struct {
	MHz      uint64 `json:"reservation_in_mhz"`
	Priority string `json:"reservation_in_shares"`
}

type MemoryReservationInfo struct {
	Percentage int `json:"reservation_percentage"`
}

type NodeSettings struct {
	AllowSshRootLogin bool     `json:"allow_ssh_root_login"`
	DnsServers        []string `json:"dns_servers"`
	EnableSsh         bool     `json:"enable_ssh"`
	Hostname          string   `json:"hostname"`
	NtpServers        []string `json:"ntp_servers"`
	SearchDomains     []string `json:"search_domains"`
}

type PerNodeStatisticsRx struct {
	TotalBytes                           uint64 `json:"total_bytes"`
	TotalPackets                         uint64 `json:"total_packets"`
	DroppedPackets                       uint64 `json:"dropped_packets"`
	BlockedPackets                       uint64 `json:"blocked_packets"`
	DestinationUnsupportedDroppedPackets uint64 `json:"destination_unsupported_dropped_packets"`
	FirewallDroppedPackets               uint64 `json:"firewall_dropped_packets"`
	IpsecDroppedPackets                  uint64 `json:"ipsec_dropped_packets"`
	IpsecNoSaDroppedPackets              uint64 `json:"ipsec_no_sa_dropped_packets"`
	IpsecNoVtiDroppedPackets             uint64 `json:"ipsec_no_vti_dropped_packets"`
	Ipv6DroppedPackets                   uint64 `json:"ipv6_dropped_packets"`
	KniDroppedPackets                    uint64 `json:"kni_dropped_packets"`
	L4portUnsupportedDroppedPackets      uint64 `json:"l4port_unsupported_dropped_packets"`
	MalformedDroppedPackets              uint64 `json:"malformed_dropped_packets"`
	NoReceiverDroppedPackets             uint64 `json:"no_receiver_dropped_packets"`
	NoRouteDroppedPackets                uint64 `json:"no_route_dropped_packets"`
	ProtoUnsupportedDroppedPackets       uint64 `json:"proto_unsupported_dropped_packets"`
	RedirectDroppedPackets               uint64 `json:"redirect_dropped_packets"`
	RpfCheckDroppedPackets               uint64 `json:"rpf_check_dropped_packets"`
	TtlExceededDroppedPackets            uint64 `json:"ttl_exceeded_dropped_packets"`
}

type PerNodeStatisticsTx struct {
	TotalBytes                  uint64 `json:"total_bytes"`
	TotalPackets                uint64 `json:"total_packets"`
	DroppedPackets              uint64 `json:"dropped_packets"`
	BlockedPackets              uint64 `json:"blocked_packets"`
	FirewallDroppedPackets      uint64 `json:"firewall_dropped_packets"`
	IpsecDroppedPackets         uint64 `json:"ipsec_dropped_packets"`
	IpsecNoSaDroppedPackets     uint64 `json:"ipsec_no_sa_dropped_packets"`
	IpsecNoVtiDroppedPackets    uint64 `json:"ipsec_no_vti_dropped_packets"`
	DadDroppedPackets           uint64 `json:"dad_dropped_packets"`
	FragNeededDroppedPackets    uint64 `json:"frag_needed_dropped_packets"`
	IpSecPolBlockDroppedPackets uint64 `json:"ipsec_pol_block_dropped_packets"`
	IpSecPolErrDroppedPackets   uint64 `json:"ipsec_pol_err_dropped_packets"`
	NoArpDroppedPackets         uint64 `json:"no_arp_dropped_packets"`
	NoLinkedDroppedPackets      uint64 `json:"no_linked_dropped_packets"`
	NoMemDroppedPackets         uint64 `json:"no_mem_dropped_packets"`
	NonIpDroppedPackets         uint64 `json:"non_ip_dropped_packets"`
	ServiceInsertDroppedPackets uint64 `json:"service_insert_dropped_packets"`
}

type PerNodeStatistics struct {
	LastUpdate uint64              `json:"last_update_timestamp"`
	Rx         PerNodeStatisticsRx `json:"rx"`
	Tx         PerNodeStatisticsTx `json:"tx"`
}

type RouterStats struct {
	PortId            string              `json:"logical_router_port_id"`
	PerNodeStatistics []PerNodeStatistics `json:"per_node_statistics"`
}

type Gateway interface {
	Print()
}

type Tier0Gateway struct {
	Id            string `json:"id"`
	HaMode        string `json:"ha_mode"`
	Name          string `json:"display_name"`
	FailoverMode  string `json:"failover_mode"`
	RealizationId string `json:"realization_id"`
	Path          string `json:"path"`
}

func (gw *Tier0Gateway) Print() {
	fmt.Printf("ID:   %v\n", gw.Id)
	fmt.Printf("Name: %v\n", gw.Name)
	fmt.Printf("HA Mode: %v\n", gw.HaMode)
	fmt.Printf("Failover Mode: %v\n", gw.FailoverMode)
}

type Tier0Gateways []Tier0Gateway

func (gws *Tier0Gateways) Print(output string) {
	if output == "json" {
	} else {
		fmt.Printf("%-8s	%-8s	%-8s	%-8s\n", "ID", "Name", "HA Mode", "Failover Mode")
		for _, gw := range *gws {
			fmt.Printf("%-8s	%8s	%8s	%8s\n", gw.Id, gw.Name, gw.HaMode, gw.FailoverMode)
		}
	}

}

type Tier1Gateway struct {
	Id            string `json:"id"`
	HaMode        string `json:"ha_mode"`
	Name          string `json:"display_name"`
	FailoverMode  string `json:"failover_mode"`
	RealizationId string `json:"realization_id"`
	Path          string `json:"path"`
}

func (gw *Tier1Gateway) Print() {
	fmt.Printf("ID:   %v\n", gw.Id)
	fmt.Printf("Name: %v\n", gw.Name)
	fmt.Printf("HA Mode: %v\n", gw.HaMode)
	fmt.Printf("Failover Mode: %v\n", gw.FailoverMode)
}

type Tier1Gateways []Tier1Gateway

func (gws *Tier1Gateways) Print(output string) {
	if output == "json" {
	} else {
		fmt.Printf("%-8s	%-8s	%-8s	%-8s\n", "ID", "Name", "HA Mode", "Failover Mode")
		for _, gw := range *gws {
			fmt.Printf("%-8s	%8s	%8s	%8s\n", gw.Id, gw.Name, gw.HaMode, gw.FailoverMode)
		}
	}
}

type BgpNeighbor struct {
	Name    string   `json:"display_name"`
	Id      string   `json:"id"`
	Address string   `json:"neighbor_address"`
	Path    string   `json:"path"`
	Asn     string   `json:"remote_as_num"`
	Source  []string `json:"source_addresses"`
}

type Segments []Segment

func (segs *Segments) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name", "Gateway", "Subnet", "State"}, "\t") + "\n"))
	for _, seg := range *segs {
		gw := seg.Connectivity
		if gw == "" {
			gw = "-"
		} else {
			gw = strings.Split(gw, "/")[2]
		}
		subnets := []string{}
		for _, subnet := range seg.Subnets {
			subnets = append(subnets, subnet.Gateway)
		}
		subnetStr := ""
		if len(subnets) > 0 {
			subnetStr = strings.Join(subnets, ",")
		} else {
			subnetStr = "-"
		}
		w.Write([]byte(strings.Join([]string{seg.Id, seg.Name, gw, subnetStr, seg.AdminState}, "\t") + "\n"))
	}
	w.Flush()
}

type Segment struct {
	Name              string                 `json:"display_name"`
	Id                string                 `json:"id"`
	AdminState        string                 `json:"admin_state,omitempty"`
	AdvancedConifg    map[string]interface{} `json:"advanced_config,omitempty"`
	Connectivity      string                 `json:"connectivity_path,omitempty"`
	ReplicationMode   string                 `json:"replication_mode,omitempty"`
	Subnets           []SegmentSubnet        `json:"subnets,omitempty"`
	TransportZonePath string                 `json:"transport_zone_path,omitempty"`
	Vlans             []string               `json:"vlan_ids"`
}

type SegmentSubnet struct {
	Gateway string `json:"gateway_address,omitempty"`
	Network string `json:"network,omitempty"`
}

type IpBlocks []IpBlock

func (bs *IpBlocks) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name", "CIDR"}, "\t") + "\n"))
	for _, b := range *bs {
		w.Write([]byte(strings.Join([]string{b.Id, b.Name, b.Cidr}, "\t") + "\n"))
	}
	w.Flush()
}

type IpBlock struct {
	Name string `json:"display_name"`
	Id   string `json:"id"`
	Cidr string `json:"cidr"`
	Path string `json:"path"`
}

type IpPools []IpPool

func (ps *IpPools) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join([]string{"ID", "Name", "Usage(allocated/available)"}, "\t") + "\n"))
	for _, p := range *ps {
		usage := fmt.Sprintf("%v/%v", p.Usage["allocated_ip_allocations"], p.Usage["available_ips"])
		w.Write([]byte(strings.Join([]string{p.Id, p.Name, usage}, "\t") + "\n"))
	}
	w.Flush()
}

type IpPool struct {
	Name  string         `json:"display_name"`
	Id    string         `json:"id"`
	Usage map[string]int `json:"pool_usage"`
}

type DfwPolicies struct {
	Policies []DfwPolicy `json:"results"`
	Count    int         `json:"result_count"`
}

type DfwPolicy struct {
	Path       string    `json:"path"`
	Name       string    `json:"display_name"`
	Id         string    `json:"id"`
	Seq        int64     `json:"sequence_number"`
	Scope      []string  `json:"scope"`
	Stateful   bool      `json:"stateful"`
	TcpStrict  bool      `json:"tcp_strict"`
	RuleCount  int       `json:"rule_count"`
	Category   string    `json:"category"`
	Rules      []DfwRule `json:"rules"`
}

type DfwRules struct {
	Rules    []DfwRule   `json:"results"`
	Count    int         `json:"result_count"`
}

type DfwRule struct {
	Name                   string    `json:"display_name"`
	Id                     string    `json:"id"`
	RuleId                 int       `json:"rule_id"`
	Sources                []string  `json:"source_groups"`
	SourcesExcluded        bool      `json:"sources_excluded"`
	Destinations           []string  `json:"destination_groups"`
	DestinationsExcluded   bool      `json:"destinations_excluded"`
	Services               []string  `json:"services"`
	Profiles               []string  `json:"profiles"`
	Scope                  []string  `json:"scope"`
	Action                 string    `json:"action"`
	Direction              string    `json:"direction"`
	IpProtocol             string    `json:"ip_protocol"`
	Logged                 bool      `json:"logged"`
}

func (r *DfwRule) Print(w *tabwriter.Writer, policy DfwPolicy) {
	cr := len(r.Sources)
	cd := len(r.Destinations)
	cs := len(r.Services)
	cp := len(r.Profiles)
	ca := len(r.Scope)
	i := 0
	for (cr > i) || (cd > i) || (cs > i) || (cp > i) || (ca > i) || (i < 1) {
		src := ""
		dest := ""
		srv := ""
		prof := ""
		scope := ""
		if cr > i {
			src = r.Sources[i]
			if strings.HasPrefix(r.Sources[i], "/infra/") {
				src = filepath.Base(r.Sources[i])
			}
		}
		if cd > i {
			dest = r.Destinations[i]
			if strings.HasPrefix(r.Destinations[i], "/infra/") {
				dest = filepath.Base(r.Destinations[i])
			}
		}
		if cs > i {
			srv = filepath.Base(r.Services[i])
		}
		if cp > i {
			prof = filepath.Base(r.Profiles[i])
		}
		if ca > i {
			scope = filepath.Base(r.Scope[i])
		}
		if i == 0 {
			w.Write([]byte(strings.Join([]string{policy.Name, r.Name, strconv.Itoa(r.RuleId), src, dest, srv, prof, scope, r.Action, r.Direction, r.IpProtocol, strconv.FormatBool(r.Logged)}, "\t")+ "\n"))
		} else {
			w.Write([]byte(strings.Join([]string{"", "", "", src, dest, srv, prof, scope, "", "", "", ""}, "\t")+ "\n"))
		}
		i = i + 1
	}
}

func GetSummary(paths []string) string {
	var basenames []string
	for _, p := range paths {
		basenames = append(basenames, filepath.Base(p))
	}
	return strings.Join(basenames, "\n")
}

func (r *DfwRule) PrintCsv(w *csv.Writer, policy DfwPolicy) {
	src := GetSummary(r.Sources)
	dest := GetSummary(r.Destinations)
	srv := GetSummary(r.Services)
	prof := GetSummary(r.Profiles)
	scope := GetSummary(r.Scope)
	err := w.Write([]string{policy.Name, r.Name, strconv.Itoa(r.RuleId), src, dest, srv, prof, scope, r.Action, r.Direction, r.IpProtocol, strconv.FormatBool(r.Logged)})
	if err != nil {
		log.Fatal(err)
		return
	}
}
