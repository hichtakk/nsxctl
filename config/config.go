package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
)

type Config struct {
	Nsx    NsxConfig    `json:"nsx" mapstructure:"nsx"`
	NsxAlb NsxAlbConfig `json:"nsx-alb,omitempty" mapstructure:"nsx-alb"`
}

type NsxConfig struct {
	CurrentSite string    `json:"current-site" mapstructure:"current-site"`
	Sites       []NsxSite `json:"sites"`
}

func (t *NsxConfig) GetCurrentSite() (NsxSite, error) {
	for _, s := range t.Sites {
		if s.Name == t.CurrentSite {
			return s, nil
		}
	}
	return NsxSite{}, fmt.Errorf("site '%s' not found on NSX configurations", t.CurrentSite)
}

func (t *NsxConfig) GetSite(name string) (NsxSite, error) {
	for _, s := range t.Sites {
		if s.Name == name {
			return s, nil
		}
	}
	return NsxSite{}, fmt.Errorf("site '%s' not found", name)
}

type NsxAlbConfig struct {
	CurrentSite string       `json:"current-site" mapstructure:"current-site"`
	Sites       []NsxAlbSite `json:"sites" mapstructure:"sites"`
}

func (a *NsxAlbConfig) GetCurrentSite() (NsxAlbSite, error) {
	for _, s := range a.Sites {
		if s.Name == a.CurrentSite {
			return s, nil
		}
	}
	return NsxAlbSite{}, fmt.Errorf("site '%s' not found", a.CurrentSite)
}

func (a *NsxAlbConfig) GetSite(name string) (NsxAlbSite, error) {
	for _, s := range a.Sites {
		if s.Name == name {
			return s, nil
		}
	}
	return NsxAlbSite{}, fmt.Errorf("site '%s' not found", name)
}

type NsxSite struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
	Proxy    string `json:"proxy,omitempty"`
}

func (t *NsxSite) getPassword() string {
	pswd, err := base64.StdEncoding.DecodeString(t.Password)
	if err != nil {
		log.Fatal(err)
	}
	return string(pswd)
}

func (t *NsxSite) SetPassword(password string) {
	t.Password = base64.StdEncoding.EncodeToString([]byte(password))
}

func (t *NsxSite) GetCredential() url.Values {
	return url.Values{
		"j_username": {t.User},
		"j_password": {t.getPassword()},
	}
}

type NsxAlbSite struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
	Version  string `json:"version"`
	Proxy    string `json:"proxy,omitempty"`
}

func (a *NsxAlbSite) getPassword() string {
	pswd, err := base64.StdEncoding.DecodeString(a.Password)
	if err != nil {
		log.Fatal(err)
	}
	return string(pswd)
}

func (a *NsxAlbSite) SetPassword(password string) {
	a.Password = base64.StdEncoding.EncodeToString([]byte(password))
}

func (a *NsxAlbSite) GetCredential() url.Values {
	return url.Values{
		"username": {a.User},
		"password": {a.getPassword()},
	}
}
