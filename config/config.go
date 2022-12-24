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
	if t.CurrentSite == "" {
		return NsxTSite{}, nil
	}
	for _, s := range t.Sites {
		if s.Name == t.CurrentSite {
			return s, nil
		}
	}
	return NsxTSite{}, fmt.Errorf("site '%s' not found", t.CurrentSite)
}

func (t *NsxTConfig) GetSite(name string) (NsxTSite, error) {
	for _, s := range t.Sites {
		if s.Name == name {
			return s, nil
		}
	}
	return NsxTSite{}, fmt.Errorf("site '%s' not found", name)
}

type NsxAlbConfig struct {
	CurrentSite string       `json:"current-site" mapstructure:"current-site"`
	Sites       []NsxAlbSite `json:"sites" mapstructure:"sites"`
}

func (a *NsxAlbConfig) GetCurrentSite() (NsxAlbSite, error) {
	if a.CurrentSite == "" {
		return NsxAlbSite{}, nil
	}
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
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
	Version  string `json:"version"`
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

func (a *NsxAlbSite) GetCredential() map[string]string {
	return map[string]string{
		"username": a.User,
		"password": a.getPassword(),
	}
}
