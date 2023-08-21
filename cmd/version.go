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

// show system version
func NewCmdShowVersion() *cobra.Command {
	aliases := []string{"v"}
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: aliases,
		Short:   fmt.Sprintf("show version of NSX/ALB [%s]", strings.Join(aliases, ",")),
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			if !alb {
				// NSX version
				ver := nsxAgent.GetVersion()
				w.Write([]byte(strings.Join([]string{"Version"}, "\t") + "\n"))
				w.Write([]byte(ver))
			} else {
				// NSX ALB version
				sysconf := albAgent.GetSystemConfiguration()
				tier := sysconf.LicenseTier
				license := albAgent.GetLicensingLedger()
				version := albAgent.GetVersion()
				usage := "-/-/-"
				for _, l := range license.TierUsages {
					if l.Tier == tier {
						usage = fmt.Sprintf("%v/%v/%v", l.Usage["consumed"], l.Usage["available"], l.Usage["remaining"])
						break
					}
				}
				w.Write([]byte(strings.Join([]string{"Version", "Tier", "TierUsage(Consumed/Capacity/Remain)"}, "\t") + "\n"))
				w.Write([]byte(strings.Join([]string{version, sysconf.LicenseTier, usage}, "\t")))
			}
			w.Write([]byte("\n"))
			w.Flush()
		},
	}
	versionCmd.PersistentFlags().BoolVarP(&alb, "alb", "", false, "show NSX ALB version")

	return versionCmd
}
