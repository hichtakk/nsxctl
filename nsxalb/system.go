package nsxalb

import (
	"encoding/json"
	"strings"

	jp "github.com/mattn/go-jsonpointer"
)

func (a *Agent) GetLicensingLedger() *LicensingLedger {
	result := &LicensingLedger{}
	path := "/api/licensing/ledger/details"
	req, _ := a.client.NewRequest("GET", path, nil)
	res := a.client.Call(req)
	// convert interface{} to []byte
	b, _ := json.Marshal(res.Content)
	// convert []byte to struct
	err := json.Unmarshal(b, result)
	if err != nil {
		return nil
	}

	return result
}

func (a *Agent) GetSystemConfiguration() *SystemConfiguration {
	result := &SystemConfiguration{}
	path := "/api/systemconfiguration"
	req, _ := a.client.NewRequest("GET", path, nil)
	res := a.client.Call(req)
	// convert interface{} to []byte
	b, _ := json.Marshal(res.Content)
	// convert []byte to struct
	err := json.Unmarshal(b, result)
	if err != nil {
		return nil
	}

	return result
}

func (a *Agent) GetVersion() string {
	// this API endpoint is not on swagger
	path := "/api/cluster/runtime"
	req, _ := a.client.NewRequest("GET", path, nil)
	res := a.client.Call(req)
	version, err := jp.Get(res.Content, "/node_info/version")
	if err != nil {
		return "-- error --"
	}
	// as 'version' includes datetime, extract version and build number
	verslice := strings.Split(version.(string), " ")

	return verslice[0]
}
