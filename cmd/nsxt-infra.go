package cmd

import (
	"fmt"
	"log"
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
			eps := nsxtclient.GetEnforcementPoint("defaualt")
			fmt.Println(eps)
		},
	}

	return enforcementPointCmd
}
