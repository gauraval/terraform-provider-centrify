package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAuthenticationProfile_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthenticationProfileRead,

		Schema:             getDSAuthenticationProfileSchema(),
		DeprecationMessage: "dataresource centrifyvault_authenticationprofile is deprecated will be removed in the future, use centrify_authenticationprofile instead",
	}
}

func dataSourceAuthenticationProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthenticationProfileRead,

		Schema: getDSAuthenticationProfileSchema(),
	}
}

func getDSAuthenticationProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "UUID of the authentication profile",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the authentication profile",
		},
		"challenges": {
			Type:     schema.TypeList,
			MaxItems: 2,
			MinItems: 1,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Authentication mechanisms for challenges",
		},
		"additional_data": {
			Type:     schema.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"number_of_questions": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "Number of questions user must answer",
					},
				},
			},
		},
		"pass_through_duration": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Pass through duration of the authentication profile",
		},
	}
}

func dataSourceAuthenticationProfileRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding authentication profile")
	client := m.(*restapi.RestClient)
	object := vault.NewAuthenticationProfile(client)
	object.Name = d.Get("name").(string)

	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving authentication profile with name '%s': %s", object.Name, err)
	}
	d.SetId(object.ID)

	schemamap, err := vault.GenerateSchemaMap(object)
	if err != nil {
		return err
	}
	//logger.Debugf("Generated Map: %+v", schemamap)
	for k, v := range schemamap {
		switch k {
		case "additional_data":
			d.Set(k, flattenAdditionalData(object.AdditionalData))
		default:
			d.Set(k, v)
		}
	}

	return nil
}
