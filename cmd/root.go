package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/hichtakk/nsxctl/client"
	c "github.com/hichtakk/nsxctl/client"
	"github.com/hichtakk/nsxctl/config"
	ac "github.com/hichtakk/nsxctl/nsxalb"
	"github.com/spf13/cobra"
)

var (
	nsxtclient *c.NsxtClient
	albclient  *ac.NsxAlbClient
	conf       config.Config
	configfile string
	useSite    string
	useAlbSite    string
	debug      bool
)

const version = "v0.1.1"

var revision = ""

func newCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nsxctl",
		Short: "NSX-T/ALB command-line client",
		Long:  "modern NSX-T/ALB client",
	}
	rootCmd.AddCommand(
		NewCmdShow(),
		NewCmdCreate(),
		NewCmdDelete(),
		NewCmdConfig(),
		NewCmdClusterInfo(),
		NewCmdTop(),
		NewCmdExec(),
		NewCmdVersion(),
		NewCmdApply(),
		//NewCmdTest(),
	)
	homedir := os.Getenv("HOME")
	if homedir == "" && runtime.GOOS == "windows" {
		homedir = os.Getenv("USERPROFILE")
	}
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", homedir+"/.config/nsxctl.json", "path to nsxctl config file")
	rootCmd.PersistentFlags().StringVarP(&useSite, "site", "", "", "specify NSX-T site name (not applicable for alb-* subcommands)")
	rootCmd.PersistentFlags().StringVarP(&useAlbSite, "alb-site", "", "", "specify ALB site name (applies only to the alb-* subcommand)")
	rootCmd.RegisterFlagCompletionFunc("site", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		site_names := []string{}
		file, _ := ioutil.ReadFile(configfile)
		json.Unmarshal(file, &conf)
		for _, s := range conf.NsxT.Sites {
			site_names = append(site_names, s.Name)
		}
		return site_names, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.RegisterFlagCompletionFunc("alb-site", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		site_names := []string{}
		file, _ := ioutil.ReadFile(configfile)
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
	file, _ := ioutil.ReadFile(configfile)
	json.Unmarshal(file, &conf)
	nsxtclient = client.NewNsxtClient(false, debug)

	var site config.NsxTSite
	var err error
	if useSite != "" {
		site, err = conf.NsxT.GetSite(useSite)
	} else {
		site, err = conf.NsxT.GetCurrentSite()
	}
	if err != nil {
		return err
	}
	nsxtclient.BaseUrl = site.Endpoint
	err = nsxtclient.Login(site.GetCredential())
	if err != nil {
		return err
	}
	return nil
}

func LoginALB() error {
	file, _ := ioutil.ReadFile(configfile)
	json.Unmarshal(file, &conf)
	albclient = ac.NewNsxAlbClient(false, debug)

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
	albclient.BaseUrl = albsite.Endpoint
	err = albclient.Login(albsite.GetCredential())
	if err != nil {
		return err
	}
	return nil
}

/*
func NewCmdTest() *cobra.Command {
	var testCmd = &cobra.Command{
		Use:   "try",
		Short: "test",
		PersistentPreRunE: func(cb *cobra.Command, args []string) error {
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			albclient = ac.NewNsxAlbClient(false, debug)
			albsite, err := conf.NsxAlb.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			albclient.BaseUrl = albsite.Endpoint
			albclient.Login(albsite.GetCredential())
			return nil
		},
		Run: func(cb *cobra.Command, args []string) {
			albclient.GenerateSeImage()
			albclient.Logout()
		},
	}

	return testCmd
}
*/
