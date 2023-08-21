package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdConfig() *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "configuration",
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Println("config")
		},
	}
	configCmd.AddCommand(
		NewCmdConfigView(),
		NewCmdConfigSetSite(),
		NewCmdConfigUseSite(),
		NewCmdConfigRemoveSite(),
		NewCmdConfigGetSites(),
	)

	return configCmd
}

func NewCmdConfigView() *cobra.Command {
	configViewCmd := &cobra.Command{
		Use:   "view",
		Short: "view configuration",
		Run: func(cmd *cobra.Command, args []string) {
			file, _ := ioutil.ReadFile(configfile)
			fmt.Println(string(file))
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return configViewCmd
}

func NewCmdConfigSetSite() *cobra.Command {
	var endpoint string
	var user string
	var password string
	var init bool
	var alb bool
	setSiteCmd := &cobra.Command{
		Use:   "set-site ${SITE_NAME}",
		Short: "add nsx-t/alb site configuration",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("argument for site name is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// name := args[0]
			// if _, err := os.Stat(configfile); os.IsNotExist(err) {
			// 	if !init {
			// 		log.Fatal("config file not found. If you need to create new config, add '--init' flag.")
			// 	}
			// 	conf.NsxT.Sites = make([]config.NsxTSite, 0)
			// 	conf.NsxAlb.Sites = make([]config.NsxAlbSite, 0)
			// 	file, err := json.MarshalIndent(conf, "", "  ")
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	err = ioutil.WriteFile(configfile, file, 0644)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// } else {
			// 	if init {
			// 		log.Fatal("config file already exist. You can't use '--init' flag.")
			// 	}
			// }
			// file, _ := ioutil.ReadFile(configfile)
			// json.Unmarshal(file, &conf)
			// if alb == false {
			// 	site := config.NsxTSite{
			// 		Name:     name,
			// 		Endpoint: endpoint,
			// 		User:     user,
			// 	}
			// 	site.SetPassword(password)
			// 	for _, s := range conf.NsxT.Sites {
			// 		if s.Name == name {
			// 			log.Fatalf("name '%s' is already exist in config", name)
			// 		}
			// 	}
			// 	conf.NsxT.Sites = append(conf.NsxT.Sites, site)
			// 	if len(conf.NsxT.Sites) == 1 {
			// 		conf.NsxT.CurrentSite = name
			// 	}
			// } else {
			// 	site := config.NsxAlbSite{
			// 		Name:     name,
			// 		Endpoint: endpoint,
			// 		User:     user,
			// 	}
			// 	site.SetPassword(password)
			// 	for _, s := range conf.NsxAlb.Sites {
			// 		if s.Name == name {
			// 			log.Fatalf("name '%s' is already exist in config", name)
			// 		}
			// 	}
			// 	conf.NsxAlb.Sites = append(conf.NsxAlb.Sites, site)
			// 	if len(conf.NsxAlb.Sites) == 1 {
			// 		conf.NsxAlb.CurrentSite = name
			// 	}
			// }
			// writeConfig()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	setSiteCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "endpoint for the new site")
	setSiteCmd.Flags().StringVarP(&user, "user", "u", "", "user for the new site")
	setSiteCmd.Flags().StringVarP(&password, "password", "p", "", "password for the new site user")
	setSiteCmd.Flags().BoolVarP(&init, "init", "", false, "create new configfile")
	setSiteCmd.Flags().BoolVarP(&alb, "alb", "", false, "NSX ALB config")
	setSiteCmd.MarkFlagRequired("endpoint")
	setSiteCmd.MarkFlagRequired("user")
	setSiteCmd.MarkFlagRequired("password")

	return setSiteCmd
}

func NewCmdConfigUseSite() *cobra.Command {
	var alb bool
	useSiteCmd := &cobra.Command{
		Use:   "use-site ${SITE_NAME}",
		Short: "set current nsx-t/alb site",
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// if len(args) != 0 {
			// 	return nil, cobra.ShellCompDirectiveNoFileComp
			// }
			// file, err := ioutil.ReadFile(configfile)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// json.Unmarshal(file, &conf)
			site_names := []string{}
			// if alb == false {
			// 	for _, s := range conf.NsxT.Sites {
			// 		site_names = append(site_names, s.Name)
			// 	}
			// } else {
			// 	for _, s := range conf.NsxAlb.Sites {
			// 		site_names = append(site_names, s.Name)
			// 	}
			// }
			return site_names, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			// name := args[0]
			// file, err := ioutil.ReadFile(configfile)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// json.Unmarshal(file, &conf)
			// siteExist := false
			// if alb == false {
			// 	for _, s := range conf.NsxT.Sites {
			// 		if s.Name == name {
			// 			siteExist = true
			// 		}
			// 	}
			// 	if siteExist == false {
			// 		log.Fatalf("site '%s' not found", name)
			// 	}
			// 	conf.NsxT.CurrentSite = name
			// } else {
			// 	for _, s := range conf.NsxAlb.Sites {
			// 		if s.Name == name {
			// 			siteExist = true
			// 		}
			// 	}
			// 	if siteExist == false {
			// 		log.Fatalf("site '%s' not found", name)
			// 	}
			// 	conf.NsxAlb.CurrentSite = name
			// }
			writeConfig()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	useSiteCmd.Flags().BoolVarP(&alb, "alb", "", false, "NSX ALB config")

	return useSiteCmd
}

func NewCmdConfigRemoveSite() *cobra.Command {
	var alb bool
	removeSiteCmd := &cobra.Command{
		Use:   "remove-site ${SITE_NAME}",
		Short: "remove specified nsx-t/alb site configuration",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("argument for site name is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// name := args[0]
			// if _, err := os.Stat(configfile); os.IsNotExist(err) {
			// 	log.Fatalf("config '%s' not found.", configfile)
			// }
			// file, _ := ioutil.ReadFile(configfile)
			// json.Unmarshal(file, &conf)
			// siteExist := false
			// if alb == false {
			// 	var newSites []config.NsxTSite
			// 	for _, s := range conf.NsxT.Sites {
			// 		if s.Name == name {
			// 			siteExist = true
			// 		} else {
			// 			newSites = append(newSites, s)
			// 		}
			// 	}
			// 	if siteExist == false {
			// 		log.Fatalf("site '%s' not found", name)
			// 	}
			// 	conf.NsxT.Sites = newSites
			// 	writeConfig()
			// } else {
			// 	var newSites []config.NsxAlbSite
			// 	for _, s := range conf.NsxAlb.Sites {
			// 		if s.Name == name {
			// 			siteExist = true
			// 		} else {
			// 			newSites = append(newSites, s)
			// 		}
			// 	}
			// 	if siteExist == false {
			// 		log.Fatalf("site '%s' not found", name)
			// 	}
			// 	conf.NsxAlb.Sites = newSites
			// 	writeConfig()
			// }
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	removeSiteCmd.Flags().BoolVarP(&alb, "alb", "", false, "NSX ALB config")

	return removeSiteCmd
}

func InitConfig(configPath string) {
	viper.SetConfigType("json")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore
			fmt.Println("config file not found")
		} else {
			fmt.Println("read config error")
			log.Fatal(err)
		}
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func writeConfig() error {
	file, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(configfile, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func NewCmdConfigGetSites() *cobra.Command {
	var alb bool
	getSitesCmd := &cobra.Command{
		Use:   "get-sites",
		Short: "get nsx-t/alb sites",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := ioutil.ReadFile(configfile)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(file, &conf)

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"Current", "Name", "Endpoint", "User"}, "\t") + "\n"))

			// if alb == false {
			// 	for _, s := range conf.NsxT.Sites {
			// 		current := ""
			// 		if s.Name == conf.NsxT.CurrentSite {
			// 			current = "*"
			// 		}
			// 		w.Write([]byte(strings.Join([]string{current, s.Name, s.Endpoint, s.User}, "\t") + "\n"))
			// 	}
			// } else {
			// 	for _, s := range conf.NsxAlb.Sites {
			// 		current := ""
			// 		if s.Name == conf.NsxAlb.CurrentSite {
			// 			current = "*"
			// 		}
			// 		w.Write([]byte(strings.Join([]string{current, s.Name, s.Endpoint, s.User}, "\t") + "\n"))
			// 	}
			// }
			w.Flush()
		},
	}
	getSitesCmd.Flags().BoolVarP(&alb, "alb", "", false, "NSX ALB config")

	return getSitesCmd
}
