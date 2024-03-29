package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceUser_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema:             getDSUserSchema(),
		DeprecationMessage: "dataresource centrifyvault_user is deprecated will be removed in the future, use centrify_user instead",
	}
}

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: getDSUserSchema(),
	}
}

func getDSUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The username in loginid@suffix format",
		},
		"email": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Email address",
		},
		"displayname": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Display name",
		},
		"password_never_expire": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Password never expires",
		},
		"force_password_change_next": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Require password change at next login",
		},
		"oauth_client": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Is OAuth confidential client",
		},
		"send_email_invite": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Send email invite for user profile setup",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description of the user",
		},
		"office_number": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"home_number": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mobile_number": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"redirect_mfa_user_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Redirect multi factor authentication to a different user account (UUID value)",
		},
		"manager_username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Username of the manager",
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding user")
	client := m.(*restapi.RestClient)
	object := vault.NewUser(client)
	object.Name = d.Get("username").(string)

	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving usert with name '%s': %s", object.Name, err)
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
