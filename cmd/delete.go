package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmdDelete is subcommand to delete resources.
func NewCmdDelete() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete resources",
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
	deleteCmd.AddCommand(
		NewCmdDeleteComputeManager(),
		NewCmdDeleteTransportZone(),
		NewCmdDeleteIpPool(),
		NewCmdDeleteIpBlock(),
		NewCmdDeleteSegment(),
		NewCmdDeleteBridge(),
	)

	return deleteCmd
}
