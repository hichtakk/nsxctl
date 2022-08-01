package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if alb != true {
				if err := Login(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return nil
			}
			if err := LoginALB(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nil
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
		NewCmdShowAlbGslb(),
	)

	return showCmd
}
