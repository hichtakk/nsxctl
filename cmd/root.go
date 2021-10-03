package cmd

import (
	"os"

	c "github.com/hichtakk/nsxctl/client"
	"github.com/hichtakk/nsxctl/config"
	"github.com/spf13/cobra"
)

var (
	nsxtclient *c.NsxtClient
	conf       config.Config
	configfile string
)

func newCmd() *cobra.Command {
	var debug bool
	rootCmd := &cobra.Command{
		Use:   "nsxctl",
		Short: "NSX-T command-line client",
		Long:  "modern NSX-T client",
	}
	rootCmd.AddCommand(
		NewCmdShow(),
		NewCmdConfig(),
	)
	homedir := os.Getenv("HOME")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", homedir+"/.config/nsxctl.json", "path to nsxctl config file")

	return rootCmd
}

func InitCmd(c *c.NsxtClient) *cobra.Command {
	nsxtclient = c
	return newCmd()
}

func GetCmdRoot() *cobra.Command {
	return newCmd()
}
