package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdShowSegment() *cobra.Command {
	aliases := []string{"seg"}
	ipPoolCmd := &cobra.Command{
		Use:     "segment",
		Aliases: aliases,
		Short:   fmt.Sprintf("show segments [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			segments := nsxtclient.GetSegment()
			sort.Slice(segments, func(i, j int) bool {
				return segments[i].Name < segments[j].Name
			})
			segments.Print()
		},
	}

	return ipPoolCmd
}
