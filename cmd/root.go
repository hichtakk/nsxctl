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
	debug      bool
)

func newCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nsxctl",
		Short: "NSX-T command-line client",
		Long:  "modern NSX-T client",
	}
	rootCmd.AddCommand(
		NewCmdShow(),
		NewCmdCreate(),
		NewCmdDelete(),
		NewCmdConfig(),
	)
	homedir := os.Getenv("HOME")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", homedir+"/.config/nsxctl.json", "path to nsxctl config file")

	return rootCmd
}

func GetCmdRoot() *cobra.Command {
	return newCmd()
}
