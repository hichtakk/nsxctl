package cmd

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/hichtakk/nsxctl/config"
	"github.com/hichtakk/nsxctl/nsx"
	"github.com/hichtakk/nsxctl/nsxalb"
	"github.com/spf13/cobra"
)

var (
	nsxAgent   *nsx.Agent
	albAgent   *nsxalb.Agent
	conf       config.Config
	configfile string
	useSite    string
	useAlbSite string
	debug      bool
)

const version = "v0.1.1"

var revision = ""

func newCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nsxctl",
		Short: "NSX/ALB command-line client",
		Long:  "NSX/ALB command-line client",
	}
	rootCmd.AddCommand(
		NewCmdShow(),
		// NewCmdCreate(),
		// NewCmdDelete(),
		// NewCmdConfig(),
		// NewCmdClusterInfo(),
		// NewCmdTop(),
		NewCmdExec(),
		NewCmdVersion(),
		// NewCmdApply(),
		//NewCmdTest(),
	)
	homedir := os.Getenv("HOME")
	if homedir == "" && runtime.GOOS == "windows" {
		homedir = os.Getenv("USERPROFILE")
	}
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", homedir+"/.config/nsxctl.json", "path to nsxctl config file")
	rootCmd.PersistentFlags().StringVarP(&useSite, "site", "", "", "specify NSX site name (not applicable for alb-* subcommands)")
	rootCmd.PersistentFlags().StringVarP(&useAlbSite, "alb-site", "", "", "specify ALB site name (applies only to the alb-* subcommand)")
	rootCmd.RegisterFlagCompletionFunc("site", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		site_names := []string{}
		file, _ := os.ReadFile(configfile)
		json.Unmarshal(file, &conf)
		for _, s := range conf.Nsx.Sites {
			site_names = append(site_names, s.Name)
		}
		return site_names, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.RegisterFlagCompletionFunc("alb-site", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		site_names := []string{}
		file, _ := os.ReadFile(configfile)
		json.Unmarshal(file, &conf)
		for _, s := range conf.NsxAlb.Sites {
			site_names = append(site_names, s.Name)
		}
		return site_names, cobra.ShellCompDirectiveNoFileComp
	})

	return rootCmd
}

func GetCmdRoot() *cobra.Command {
	return newCmd()
}

func Login() error {
	file, _ := os.ReadFile(configfile)
	json.Unmarshal(file, &conf)
	var site config.NsxSite
	var err error
	if useSite != "" {
		site, err = conf.Nsx.GetSite(useSite)
	} else {
		site, err = conf.Nsx.GetCurrentSite()
	}
	if err != nil {
		return err
	}
	nsxAgent = nsx.NewNsxAgent(&site, debug)
	err = nsxAgent.Login()
	if err != nil {
		return err
	}
	return nil
}

func LoginALB() error {
	file, _ := os.ReadFile(configfile)
	json.Unmarshal(file, &conf)
	var albsite config.NsxAlbSite
	var err error
	if useAlbSite != "" {
		albsite, err = conf.NsxAlb.GetSite(useAlbSite)
	} else {
		albsite, err = conf.NsxAlb.GetCurrentSite()
	}
	if err != nil {
		return err
	}
	albAgent = nsxalb.NewNsxAlbAgent(&albsite, debug)
	err = albAgent.Login()
	if err != nil {
		return err
	}
	return nil
}
