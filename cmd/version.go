package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version of nsxctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nsxctl version: %s, revision: %s\n", version, revision)
		},
	}

	return versionCmd
}

// show NSX-T version
func NewCmdShowVersion() *cobra.Command {
	aliases := []string{"v"}
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: aliases,
		Short:   fmt.Sprintf("show version of NSX-T/ALB [%s]", strings.Join(aliases, ",")),
		Run: func(cmd *cobra.Command, args []string) {
			if alb != true {
				fmt.Println(nsxtclient.GetVersion())
			} else {
				c := albclient.GetSystemConfiguration()
				tier := c.LicenseTier
				license := albclient.GetLicensingLedger()
				usage := "-/-/-"
				for _, l := range license.TierUsages {
					if l.Tier == tier {
						usage = fmt.Sprintf("%v/%v/%v", l.Usage["consumed"], l.Usage["available"], l.Usage["remaining"])
						break
					}
				}
				w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
				w.Write([]byte(strings.Join([]string{"Version", "Tier", "TierUsage(Consumed/Capacity/Remain)"}, "\t") + "\n"))
				w.Write([]byte(strings.Join([]string{albclient.Version, c.LicenseTier, usage}, "\t") + "\n"))
				w.Flush()
			}
		},
	}
	versionCmd.PersistentFlags().BoolVarP(&alb, "alb", "", false, "show NSX ALB version")

	return versionCmd
}
