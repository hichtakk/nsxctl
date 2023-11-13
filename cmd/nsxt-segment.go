package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hichtakk/nsxctl/structs"
	"github.com/spf13/cobra"
)

func NewCmdShowSegment() *cobra.Command {
	var verbose bool
	aliases := []string{"seg"}
	segmentCmd := &cobra.Command{
		Use:     "segment",
		Aliases: aliases,
		Short:   fmt.Sprintf("show segments [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			segments := nsxtclient.GetSegment()
			sort.Slice(segments, func(i, j int) bool {
				return segments[i].Name < segments[j].Name
			})

			if verbose {
				bridgeInfo := map[string]string{}
				for _, seg := range segments {
					bps := []string{}
					for _, bp := range seg.BridgeProfiles {
						res2 := nsxtclient.Request("GET", "/policy/api/v1" + bp.Path, nil, nil)
						var bridgeProfile structs.BridgeProfile
						res2.UnmarshalBody(&bridgeProfile)
						bps = append(bps, fmt.Sprintf("%s(%s)", strings.Join(bp.Vlans, ","), bridgeProfile.Name))
					}
					bridgeInfo[seg.Id] = strings.Join(bps, ":")
				}
				segments.Print(bridgeInfo)
			} else {
				segments.Print(nil)
			}
		},
	}
	segmentCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "display bridge info")

	return segmentCmd
}
