package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
)

type Config struct {
	NsxT   NsxTConfig   `json:"nsx-t" mapstructure:"nsx-t"`
	NsxAlb NsxAlbConfig `json:"nsx-alb,omitempty" mapstructure:"nsx-alb"`
}

type NsxTConfig struct {
	CurrentSite string     `json:"current-site" mapstructure:"current-site"`
	Sites       []NsxTSite `json:"sites"`
}

func (t *NsxTConfig) GetCurrentSite() (NsxTSite, error) {
	for _, s := range t.Sites {
		if s.Name == t.CurrentSite {
			return s, nil
		}
	}
	return NsxTSite{}, fmt.Errorf("site '%s' not found", t.CurrentSite)
}

type NsxAlbConfig struct {
	CurrentSite string       `json:"current-site" mapstructure:"current-site"`
	Sites       []NsxAlbSite `json:"sites" mapstructure:"sites"`
}

type NsxTSite struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (t *NsxTSite) getPassword() string {
	pswd, err := base64.StdEncoding.DecodeString(t.Password)
	if err != nil {
		log.Fatal(err)
	}
	return string(pswd)
}

func (t *NsxTSite) SetPassword(password string) {
	t.Password = base64.StdEncoding.EncodeToString([]byte(password))
}

func (t *NsxTSite) GetCredential() url.Values {
	return url.Values{
		"j_username": {t.User},
		"j_password": {t.getPassword()},
	}
}

type NsxAlbSite struct {
	NsxTSite
	Version string `json:"version"`
}
