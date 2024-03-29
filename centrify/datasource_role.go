package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRole_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRoleRead,

		Schema:             getDSRoleSchema(),
		DeprecationMessage: "dataresource centrifyvault_role is deprecated will be removed in the future, use centrify_role instead",
	}
}

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRoleRead,

		Schema: getDSRoleSchema(),
	}
}

func getDSRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the role",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description of an role",
		},
		"adminrights": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"member": {
			Type:     schema.TypeSet,
			Optional: true,
			Set:      customRoleMemberHash,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the member",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Name of the member",
					},
					"type": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Type of the member",
					},
				},
			},
		},
	}
}

func dataSourceRoleRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding role")
	client := m.(*restapi.RestClient)
	object := vault.NewRole(client)
	object.Name = d.Get("name").(string)

	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving role with name '%s': %s", object.Name, err)
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
