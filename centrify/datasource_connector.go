package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceConnector_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectorRead,

		Schema:             getDSConnectorSchema(),
		DeprecationMessage: "dataresource centrifyvault_connector is deprecated will be removed in the future, use centrify_connector instead",
	}
}

func dataSourceConnector() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectorRead,

		Schema: getDSConnectorSchema(),
	}
}

func getDSConnectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the Connector",
		},
		"machine_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns_host_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"forest": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ssh_service": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"rdp_service": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ad_proxy": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"app_gateway": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"http_api_service": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ldap_proxy": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"radius_service": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"radius_external_service": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"version": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"status": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive"}, false),
		},
		"online": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"vpc_identifier": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"vm_identifier": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func dataSourceConnectorRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding connector")
	client := m.(*restapi.RestClient)
	object := vault.NewConnector(client)
	object.Name = d.Get("name").(string)
	object.MachineName = d.Get("machine_name").(string)
	object.DnsHostName = d.Get("dns_host_name").(string)
	object.Forest = d.Get("forest").(string)
	object.Status = d.Get("status").(string)
	object.Version = d.Get("version").(string)
	object.VpcIdentifier = d.Get("vpc_identifier").(string)
	object.VmIdentifier = d.Get("vm_identifier").(string)

	// We can't use simple Query method because it doesn't return all attributes
	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving Connector with name '%s': %s", object.Name, err)
	}
	d.SetId(object.ID)

	schemamap, err := vault.GenerateSchemaMap(object)
	if err != nil {
		return err
	}
	//logger.Debugf("Generated Map: %+v", schemamap)
	for k, v := range schemamap {
		d.Set(k, v)
	}

	return nil
}
