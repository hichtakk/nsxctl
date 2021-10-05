package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hichtakk/nsxctl/config"
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
			name := args[0]
			if _, err := os.Stat(configfile); os.IsNotExist(err) {
				if !init {
					log.Fatal("config file not found. If you need to create new config, add '--init' flag.")
				}
				conf.NsxT.Sites = make([]config.NsxTSite, 0)
				conf.NsxAlb.Sites = make([]config.NsxAlbSite, 0)
				file, err := json.MarshalIndent(conf, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(configfile, file, 0644)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				if init {
					log.Fatal("config file already exist. You can't use '--init' flag.")
				}
			}
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			site := config.NsxTSite{
				Name:     name,
				Endpoint: endpoint,
				User:     user,
			}
			site.SetPassword(password)
			for _, s := range conf.NsxT.Sites {
				if s.Name == name {
					log.Fatalf("name '%s' is already exist in config", name)
				}
			}
			conf.NsxT.Sites = append(conf.NsxT.Sites, site)
			if len(conf.NsxT.Sites) == 1 {
				conf.NsxT.CurrentSite = name
			}
			writeConfig()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	setSiteCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "endpoint for the new site")
	setSiteCmd.Flags().StringVarP(&user, "user", "u", "", "user for the new site")
	setSiteCmd.Flags().StringVarP(&password, "password", "p", "", "password for the new site user")
	setSiteCmd.Flags().BoolVarP(&init, "init", "", false, "create new configfile")
	setSiteCmd.MarkFlagRequired("endpoint")
	setSiteCmd.MarkFlagRequired("user")
	setSiteCmd.MarkFlagRequired("password")

	return setSiteCmd
}

func NewCmdConfigUseSite() *cobra.Command {
	useSiteCmd := &cobra.Command{
		Use:   "use-site ${SITE_NAME}",
		Short: "add nsx-t/alb site configuration",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				cmd.Help()
				return fmt.Errorf("argument for site name is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			file, err := ioutil.ReadFile(configfile)
			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal(file, &conf)
			siteExist := false
			for _, s := range conf.NsxT.Sites {
				if s.Name == name {
					siteExist = true
				}
			}
			if siteExist == false {
				log.Fatalf("site '%s' not found", name)
			}
			conf.NsxT.CurrentSite = name
			writeConfig()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return useSiteCmd
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
