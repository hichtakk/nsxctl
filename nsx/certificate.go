package nsx

import (
	"bytes"
)

func (a *Agent) GetApiCertificate() Certificates {
	// path := "/api/v1/trust-management/certificates"
	// res := a.Request("GET", path, map[string]string{"type": "api_certificate"}, nil)
	certs := Certificates{}
	// for _, c := range res.Body.(map[string]interface{})["results"].([]interface{}) {
	// 	str, _ := json.Marshal(c)
	// 	var cert Certificate
	// 	json.Unmarshal(str, &cert)
	// 	certs = append(certs, cert)
	// }

	return certs
}

func (a *Agent) GetTlsFingerprint(server string, port uint) string {
	// conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", server, port), &tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	// 	panic("failed to connect: " + err.Error())
	// }
	// cert := conn.ConnectionState().PeerCertificates[0]
	// fingerprint := sha256.Sum256(cert.Raw)
	var buf bytes.Buffer
	// for i, f := range fingerprint {
	// 	if i > 0 {
	// 		fmt.Fprintf(&buf, ":")
	// 	}
	// 	fmt.Fprintf(&buf, "%02X", f)
	// }
	// conn.Close()
	return buf.String()
}

type Certificate struct {
	Name string `json:"display_name"`
	Id   string `json:"id"`
	Pem  string `json:"pem_encoded"`
}

type Certificates []Certificate
