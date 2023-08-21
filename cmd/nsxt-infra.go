package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdShowEnforcementPoint() *cobra.Command {
	aliases := []string{"ep"}
	enforcementPointCmd := &cobra.Command{
		Use:     "infra-enforcementpoint",
		Aliases: aliases,
		Short:   fmt.Sprintf("show infra-enforcementpoint [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// eps := nsxtclient.GetEnforcementPoint("default")
			// for _, ep := range *eps {
			// 	fmt.Println(ep.Id, ep.Path)
			// }
		},
	}

	return enforcementPointCmd
}
