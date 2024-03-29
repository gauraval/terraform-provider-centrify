package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMultiplexedAccount_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMultiplexedAccountRead,

		Schema:             getDSMultiplexedAccountSchema(),
		DeprecationMessage: "dataresource centrifyvault_multiplexedaccount is deprecated will be removed in the future, use centrify_multiplexedaccount instead",
	}
}

func dataSourceMultiplexedAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMultiplexedAccountRead,

		Schema: getDSMultiplexedAccountSchema(),
	}
}

func getDSMultiplexedAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the multiplexed account",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description of the multiplexed account",
		},
		"account1_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account2_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account1": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account2": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"accounts": {
			Type:     schema.TypeSet,
			Computed: true,
			MinItems: 2,
			MaxItems: 2,
			Set:      schema.HashString,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"active_account": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceMultiplexedAccountRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding multiplexed account")
	client := m.(*restapi.RestClient)
	object := vault.NewMultiplexedAccount(client)
	object.Name = d.Get("name").(string)

	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving multiplexed account with name '%s': %s", object.Name, err)
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
