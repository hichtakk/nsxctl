package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if alb != true {
				return Login()
			}
			return LoginALB()
		},
		PersistentPostRun: func(c *cobra.Command, args []string) {
			if alb != true {
				nsxtclient.Logout()
				return
			}
			albclient.Logout()
		},
	}
	showCmd.AddCommand(
		NewCmdShowVersion(),
		NewCmdShowGateway(),
		NewCmdShowComputeManager(),
		NewCmdShowTransportZone(),
		NewCmdShowTransportNode(),
		NewCmdShowTransportNodeProfile(),
		NewCmdShowEnforcementPoint(),
		NewCmdShowIpPool(),
		NewCmdShowIpBlock(),
		NewCmdShowEdge(),
		NewCmdShowEdgeCluster(),
		NewCmdShowSegment(),
		NewCmdShowRoutingTable(),
		NewCmdShowBgpAdvRoutes(),
		NewCmdShowAlbCloud(),
		NewCmdShowAlbVirtualService(),
		NewCmdShowAlbServiceEngine(),
	)

	return showCmd
}
