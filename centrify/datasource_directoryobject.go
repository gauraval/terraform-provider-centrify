package centrify

import (
	"fmt"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDirectoryObject_deprecated() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDirectoryObjectRead,

		Schema:             getDSDirectoryObjectSchema(),
		DeprecationMessage: "dataresource centrifyvault_directoryobject is deprecated will be removed in the future, use centrify_directoryobject instead",
	}
}

func dataSourceDirectoryObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDirectoryObjectRead,

		Schema: getDSDirectoryObjectSchema(),
	}
}

func getDSDirectoryObjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the directory object",
		},
		"object_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Type of the directory object",
			ValidateFunc: validation.StringInSlice([]string{
				"User",
				"Group",
			}, false),
		},
		"system_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "UPN of the directory object",
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Display name of the directory object",
		},
		"distinguished_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Distinguished name of the directory object",
		},
		"forest": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Forest name of the directory object",
		},
		"directory_services": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of UUID of directory services",
		},
	}
}

func dataSourceDirectoryObjectRead(d *schema.ResourceData, m interface{}) error {
	logger.Infof("Finding Directory Object")
	client := m.(*restapi.RestClient)
	object := vault.NewDirectoryObjects(client)
	object.QueryName = d.Get("name").(string)
	object.ObjectType = d.Get("object_type").(string)
	object.DirectoryServices = flattenSchemaSetToStringSlice(d.Get("directory_services"))

	err := object.Read()
	if err != nil {
		return fmt.Errorf("error retrieving directory services: %s", err)
	}

	var results []vault.DirectoryObject
	// Further narrow down with Distinguished if specified
	dn := d.Get("distinguished_name").(string)
	if dn != "" {
		for _, v := range object.DirectoryObjects {
			if dn == v.DistinguishedName {
				results = append(results, v)
			}
		}
	} else {
		results = object.DirectoryObjects
	}

	if len(results) == 0 {
		return fmt.Errorf("query returns 0 object for directory object %s", object.QueryName)
	}
	if len(results) > 1 {
		return fmt.Errorf("search directory object %s, but returns too many objects (found %d, expected 1)", object.QueryName, len(results))

	}

	var result = results[0]
	d.SetId(result.ID)
	d.Set("ID", result.ID)
	d.Set("display_name", result.DisplayName)
	d.Set("name", result.Name)
	d.Set("forest", result.Forest)
	d.Set("system_name", result.SystemName)

	return nil
}
