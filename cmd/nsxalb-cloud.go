package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hichtakk/nsxctl/structs"
	"github.com/spf13/cobra"
)

func NewCmdShowAlbCloud() *cobra.Command {
	aliases := []string{"ac", "cloud"}
	cloudCmd := &cobra.Command{
		Use:     "alb-cloud",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB Cloud [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			albclient.ShowCloud()
			//albclient.DownloadSeImage()
		},
	}

	return cloudCmd
}

func NewCmdShowAlbVirtualService() *cobra.Command {
	var verbose bool
	aliases := []string{"alb-vs", "vs"}
	cloudCmd := &cobra.Command{
		Use:     "alb-virtualservice",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB VirtualService [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			vss := albclient.ShowVirtualService()

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			if verbose {
				w.Write([]byte(strings.Join([]string{"ID", "Name", "VIP", "Port", "Network", "Cloud", "SEGroup", "VRF", "Status", "ServiceEngines"}, "\t") + "\n"))
			} else {
				w.Write([]byte(strings.Join([]string{"ID", "Name", "VIP", "Port", "Network", "Cloud", "SEGroup", "VRF", "Status"}, "\t") + "\n"))
			}
			for _, vs := range vss {
				vs.Print(w, verbose)
			}
			w.Flush()
		},
	}
	cloudCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display serviceengine placement")

	return cloudCmd
}

func NewCmdShowAlbServiceEngine() *cobra.Command {
	aliases := []string{"alb-se", "se"}
	cloudCmd := &cobra.Command{
		Use:     "alb-serviceengine",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB ServiceEngine [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ses := albclient.GetServiceEngine()

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"ID", "Name", "IP", "Cloud", "SEGroup", "Status"}, "\t") + "\n"))
			for _, se := range ses {
				se.Print(w)
			}
			w.Flush()
		},
	}

	return cloudCmd
}

func NewCmdShowAlbPool() *cobra.Command {
	aliases := []string{"pool"}
	cloudCmd := &cobra.Command{
		Use:     "alb-pool",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB Pool [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: GetPoolNames,
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			pools := albclient.GetPools()

			if len(args) > 0 {
				var pool structs.Pool
				for _, p := range pools {
					if p.Config.Name == args[0] {
						pool = albclient.GetPool(p.Config.UUID)
					}
				}
				if pool.Name == "" {
					log.Fatal("pool not found.")
					return
				}
				pool.Print()
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"ID", "Name", "VirtualService", "Ready", "Status", "Cloud"}, "\t") + "\n"))
			for _, p := range pools {
				p.Print(w)
			}
			w.Flush()
		},
	}

	return cloudCmd
}

func GetPoolNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	LoginALB()
	pool_names := []string{}
	pools := albclient.GetPools()
	for _, p := range pools {
		pool_names = append(pool_names, p.Config.Name)
	}
	return pool_names, cobra.ShellCompDirectiveNoFileComp
}
