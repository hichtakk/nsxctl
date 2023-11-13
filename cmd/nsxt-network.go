package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hichtakk/nsxctl/structs"
	"github.com/spf13/cobra"
)

func NewCmdShowIpPool() *cobra.Command {
	aliases := []string{"ipp"}
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ip address pools [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pools := nsxtclient.GetIpPool()
			pools.Print()
		},
	}

	return ipPoolCmd
}

func NewCmdShowIpBlock() *cobra.Command {
	aliases := []string{"ipb"}
	ipPoolCmd := &cobra.Command{
		Use:     "ip-block",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ip address blocks [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			blocks := nsxtclient.GetIpBlock()
			blocks.Print()
		},
	}

	return ipPoolCmd
}

func NewCmdCreateIpPool() *cobra.Command {
	//var transportType string
	aliases := []string{"ipp"}
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: aliases,
		Short:   fmt.Sprintf("create a new ip pool [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			nsxtclient.CreateIpPool(name)
		},
	}
	//ipPoolCmd.Flags().StringVarP(&transportType, "type", "t", "", "transport zone type [vlan, overlay]")
	//ipPoolCmd.MarkFlagRequired("transportType")

	return ipPoolCmd
}

func NewCmdCreateIpBlock() *cobra.Command {
	var cidr string
	aliases := []string{"ipb"}
	ipBlockCmd := &cobra.Command{
		Use:     "ip-block",
		Aliases: aliases,
		Short:   fmt.Sprintf("create a new ip block [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			nsxtclient.CreateIpBlock(name, cidr)
		},
	}
	ipBlockCmd.Flags().StringVarP(&cidr, "cidr", "", "", "CIDR block (10.0.0.0/16)")
	ipBlockCmd.MarkFlagRequired("cidr")

	return ipBlockCmd
}

func NewCmdDeleteIpPool() *cobra.Command {
	aliases := []string{"ipp"}
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: aliases,
		Short:   fmt.Sprintf("delete ip pool [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			nsxtclient.DeleteIpPool(name)
		},
	}

	return ipPoolCmd
}

func NewCmdDeleteIpBlock() *cobra.Command {
	aliases := []string{"ipb"}
	ipPoolCmd := &cobra.Command{
		Use:     "ip-block",
		Aliases: aliases,
		Short:   fmt.Sprintf("delete ip block [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			nsxtclient.DeleteIpBlock(name)
		},
	}

	return ipPoolCmd
}

func NewCmdCreateSegment() *cobra.Command {
	var transportzone string
	var vlan_ids string              // CSV format
	var gateway string
	var interface_address string     // CIDR format
	segmentCmd := &cobra.Command{
		Use: "segment",
		Short: fmt.Sprintf("create a new segment"),
		Args: cobra.ExactArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			if err := Login(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			segment_name := args[0]
			err := nsxtclient.CreateSegment(segment_name, transportzone, vlan_ids, gateway, interface_address)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	segmentCmd.Flags().StringVarP(&interface_address, "interface-address", "", "", "interface address (CIDR format) to connect specified gateway")
	segmentCmd.Flags().StringVarP(&gateway, "gateway", "", "", "gateway name to connect")
	segmentCmd.Flags().StringVarP(&transportzone, "transportzone", "", "", "transportzone name")
	segmentCmd.Flags().StringVarP(&vlan_ids, "vlan-ids", "", "", "vlan ids (CSV format)")
	segmentCmd.RegisterFlagCompletionFunc("gateway", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		Login()
		gw_names := []string{}
		for _, gw := range nsxtclient.GetGateways(-1) {
			gw_names = append(gw_names, gw.Name)
		}
		return gw_names, cobra.ShellCompDirectiveNoFileComp
	})
	segmentCmd.RegisterFlagCompletionFunc("transportzone", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		Login()
		sites := nsxtclient.GetSite()
		endpoints := nsxtclient.GetEnforcementPoint(sites[0])
		transport_zones := nsxtclient.GetPolicyTransportZone(sites[0], (*endpoints)[0].Id)
		var tz_names []string
		for _, tz := range *transport_zones {
			tz_names = append(tz_names, tz.Name)
		}
		return tz_names, cobra.ShellCompDirectiveNoFileComp
	})

	return segmentCmd
}

func NewCmdDeleteSegment() *cobra.Command {
	segmentCmd := &cobra.Command{
		Use: "segment",
		Short: fmt.Sprintf("delete a segment"),
		Args: cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			Login()
			segment_names := []string{}
			for _, seg := range nsxtclient.GetSegment() {
				segment_names = append(segment_names, seg.Name)
			}
			return segment_names, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			segment_name := args[0]
			err := nsxtclient.DeleteSegment(segment_name)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return segmentCmd
}

func NewCmdDeleteBridge() *cobra.Command {
	var vlan_ids string              // CSV format
	segmentCmd := &cobra.Command{
		Use: "bridge",
		Short: fmt.Sprintf("delete edge bridge from segment"),
		Args: cobra.ExactArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			if err := Login(); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			Login()
			segment_names := []string{}
			for _, seg := range nsxtclient.GetSegment() {
				if seg.BridgeProfiles != nil {
					segment_names = append(segment_names, seg.Name)
				}
			}
			return segment_names, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			segment_name := args[0]
			var segment structs.Segment
			for _, seg := range nsxtclient.GetSegment() {
				if seg.Name == segment_name {
					segment = seg
					break
				}
			}

			new_bps := []structs.BridgeProfileInfo{}
			if vlan_ids != "" {
				for _, bp := range segment.BridgeProfiles {
					new_vlans := []string{}
					for _, vlan := range bp.Vlans {
						found := false
						for _, vlan_to_delete := range strings.Split(vlan_ids, ",") {
							if vlan == vlan_to_delete {
								found = true
							}
						}
						if !found {
							new_vlans = append(new_vlans, vlan)
						}
					}
					if len(new_vlans) > 0 {
						bp.Vlans = new_vlans
						new_bps = append(new_bps, bp)
					}
				}
			}
			segment.BridgeProfiles = new_bps
			err := nsxtclient.UpdateSegment(segment)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	segmentCmd.Flags().StringVarP(&vlan_ids, "vlan-ids", "", "", "vlan ids (CSV format)")

	return segmentCmd
}
