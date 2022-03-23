package cmd

import (
	"fmt"
	"log"
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
			nsxtclient.GetSegment()
		},
	}

	return ipPoolCmd
}
