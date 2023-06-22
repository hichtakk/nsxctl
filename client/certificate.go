package client

import (
	"encoding/json"
)

func (c *NsxtClient) GetApiCertificate() Certificates {
	path := "/api/v1/trust-management/certificates"
	res := c.Request("GET", path, map[string]string{"type": "api_certificate"}, nil)
	certs := Certificates{}
	for _, c := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(c)
		var cert Certificate
		json.Unmarshal(str, &cert)
		certs = append(certs, cert)
	}

	return certs
}

type Certificate struct {
	Name string `json:"display_name"`
	Id   string `json:"id"`
	Pem  string `json:"pem_encoded"`
}

type Certificates []Certificate
