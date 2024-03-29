package platform

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
)

// Connector - Encapsulates a single Connector
type Connector struct {
	vaultObject

	MachineName           string `json:"MachineName,omitempty" schema:"machine_name,omitempty"`
	DnsHostName           string `json:"DnsHostName,omitempty" schema:"dns_host_name,omitempty"`
	Forest                string `json:"Forest,omitempty" schema:"forest,omitempty"`
	SSHService            string `json:"SSHService,omitempty" schema:"ssh_service,omitempty"`
	RDPService            string `json:"RDPService,omitempty" schema:"rdp_service,omitempty"`
	ADProxy               string `json:"ADProxy,omitempty" schema:"ad_proxy,omitempty"`
	AppGateway            string `json:"AppGateway,omitempty" schema:"app_gateway,omitempty"`
	HttpAPIService        string `json:"HttpAPIService,omitempty" schema:"http_api_service,omitempty"`
	LDAPProxy             string `json:"LDAPProxy,omitempty" schema:"ldap_proxy,omitempty"`
	RadiusService         string `json:"RadiusService,omitempty" schema:"radius_service,omitempty"`
	RadiusExternalService string `json:"RadiusExternalService,omitempty" schema:"radius_external_service,omitempty"`
	Online                bool   `json:"Online,omitempty" schema:"online,omitempty"`
	Version               string `json:"Version,omitempty" schema:"version,omitempty"`
	VpcIdentifier         string `json:"VpcIdentifier,omitempty" schema:"vpc_identifier,omitempty"`
	VmIdentifier          string `json:"VmIdentifier,omitempty" schema:"vm_identifier,omitempty"`
	Status                string `json:"-"` // Used to represent Online status
}

// NewConnector is a Connector constructor
func NewConnector(c *restapi.RestClient) *Connector {
	s := Connector{}
	s.client = c

	return &s
}

// Read function fetches a Connector from source, including attribute values. Returns error if any
func (o *Connector) Read() error {
	return nil
}

// Delete function deletes a Connector and returns a map that contains deletion result
func (o *Connector) Delete() (*restapi.GenericMapResponse, error) {
	return o.deleteObjectMapAPI("")
}

// Update function updates an existing Connector and returns a map that contains update result
func (o *Connector) Update() (*restapi.GenericMapResponse, error) {
	return nil, nil
}

// Query function returns a single Connector object in map format
func (o *Connector) Query() (map[string]interface{}, error) {
	query := "SELECT * FROM Proxy WHERE 1=1"
	if o.Name != "" {
		query += " AND Name='" + o.Name + "'"
	}
	if o.Status != "" {
		if o.Status == "Active" {
			query += " AND Online=1"
		} else {
			query += " AND Online=0"
		}
	}
	if o.Version != "" {
		query += " AND Version='" + o.Version + "'"
	}
	if o.VpcIdentifier != "" {
		query += " AND VpcIdentifier='" + o.VpcIdentifier + "'"
	}
	if o.VpcIdentifier != "" {
		query += " AND VmIdentifier='" + o.VmIdentifier + "'"
	}
	if o.MachineName != "" {
		query += " AND MachineName='" + o.MachineName + "'"
	}
	if o.DnsHostName != "" {
		query += " AND DnsHostName='" + o.DnsHostName + "'"
	}
	if o.Forest != "" {
		query += " AND Forest='" + o.Forest + "'"
	}

	return queryVaultObject(o.client, query)
}

// GetIDByName returns vault object ID by name
func (o *Connector) GetIDByName() (string, error) {
	if o.Name == "" {
		return "", fmt.Errorf("%s name must be provided", GetVarType(o))
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return "", fmt.Errorf("error retrieving %s: %s", GetVarType(o), err)
	}
	o.ID = result["ID"].(string)

	return o.ID, nil
}

// GetByName retrieves vault object from tenant by name
func (o *Connector) GetByName() error {
	if o.Name == "" {
		return fmt.Errorf("%s name must be provided", GetVarType(o))
	}

	result, err := o.Query()
	if err != nil {
		logger.Errorf(err.Error())
		return fmt.Errorf("error retrieving %s: %s", GetVarType(o), err)
	}
	mapToStruct(o, result)

	return nil
}
