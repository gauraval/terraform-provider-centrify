package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOidcWebApp_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOidcWebAppRead,

		Schema:             getDSOidcWebAppSchema(),
		DeprecationMessage: "dataresource centrifyvault_webapp_oidc is deprecated will be removed in the future, use centrify_webapp_oidc instead",
	}
}

func dataSourceOidcWebApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOidcWebAppRead,

		Schema: getDSOidcWebAppSchema(),
	}
}

func getDSOidcWebAppSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the Web App",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Description of the Web App",
		},
		"application_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Application ID. Specify the name or 'target' that the mobile application uses to find this application.",
		},
		"oauth_profile": getOidcProfileSchema(),
		"template_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Template name of the Web App",
		},
		"oidc_script": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		// Policy menu
		"default_profile_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default authentication profile ID",
		},
		// Account Mapping menu
		"username_strategy": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Account mapping",
		},
		"username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "All users share the user name. Applicable if 'username_strategy' is 'Fixed'",
		},
		"user_map_script": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Account mapping script. Applicable if 'username_strategy' is 'UseScript'",
		},
		// Workflow
		"workflow_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"workflow_approver": getWorkflowApproversSchema(),
		"challenge_rule":    getChallengeRulesSchema(),
		"policy_script": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Use script to specify authentication rules (configured rules are ignored)",
		},
	}
}

func dataSourceOidcWebAppRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding Oidc webapp")
	client := m.(*restapi.RestClient)
	object := vault.NewOidcWebApp(client)
	object.Name = d.Get("name").(string)
	object.ApplicationID = d.Get("application_id").(string)

	// We can't use simple Query method because it doesn't return all attributes
	err := object.GetByName()
	if err != nil {
		return fmt.Errorf("error retrieving Oauth webapp with name '%s': %s", object.Name, err)
	}
	d.SetId(object.ID)

	schemamap, err := vault.GenerateSchemaMap(object)
	if err != nil {
		return err
	}
	//logger.Debugf("Generated Map: %+v", schemamap)
	for k, v := range schemamap {
		switch k {
		case "oauth_profile":
			d.Set(k, []interface{}{v})
		case "challenge_rule":
			d.Set(k, v.(map[string]interface{})["rule"])
		case "workflow_settings":
			if v.(string) != "" {
				wfschema, err := convertWorkflowSchema(v.(string))
				if err != nil {
					return err
				}
				d.Set("workflow_approver", wfschema)
				d.Set(k, v)
			}
		default:
			d.Set(k, v)
		}
	}

	return nil
}
