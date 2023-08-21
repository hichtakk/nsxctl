package nsxalb

import (
	"bytes"

	"github.com/hichtakk/nsxctl/client"
	"github.com/hichtakk/nsxctl/config"
)

func NewNsxAlbAgent(site *config.NsxAlbSite, debug bool) *Agent {
	c, err := client.NewNsxAlbClient(site.Endpoint, site.Proxy, debug)
	if err != nil {
		return nil
	}
	return &Agent{c, site}
}

type Agent struct {
	client *client.NsxAlbClient
	site   *config.NsxAlbSite
}

func (a *Agent) Login() error {
	cred := a.site.GetCredential()
	return a.client.Login(cred)
}

func (a *Agent) Logout() error {
	return a.client.Logout()
}

func (a *Agent) ExecGet(path string, params map[string]string) (*client.Response, error) {
	req, err := a.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	res := a.client.Call(req)
	// switch res.StatusCode {
	// case 200, 201:
	// case 404:
	// case 500:
	// }
	return res, nil
}

func (a *Agent) ExecPost(path string, data []byte) (*client.Response, error) {
	req, err := a.client.NewRequest("POST", path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	res := a.client.Call(req)
	return res, nil
}

func (a *Agent) ExecPut(path string, data []byte) (*client.Response, error) {
	req, err := a.client.NewRequest("PUT", path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	res := a.client.Call(req)
	return res, nil
}

func (a *Agent) ExecPatch(path string, data []byte) (*client.Response, error) {
	req, err := a.client.NewRequest("PATCH", path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	res := a.client.Call(req)
	return res, nil
}

func (a *Agent) ExecDelete(path string, data []byte) (*client.Response, error) {
	req, err := a.client.NewRequest("DELETE", path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	res := a.client.Call(req)
	return res, nil
}
