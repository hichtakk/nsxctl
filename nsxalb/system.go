package nsxalb

import (
	"encoding/json"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxAlbClient) GetLicensingLedger() structs.LicensingLedger {
	resp := c.Request("GET", "/api/licensing/ledger/details", map[string]string{}, nil)
	var result structs.LicensingLedger
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &result)
	return result
}

func (c *NsxAlbClient) GetSystemConfiguration() structs.SystemConfiguration {
	resp := c.Request("GET", "/api/systemconfiguration", map[string]string{}, nil)
	var result structs.SystemConfiguration
	resByte, _ := resp.BodyBytes()
	json.Unmarshal(resByte, &result)
	return result
}
