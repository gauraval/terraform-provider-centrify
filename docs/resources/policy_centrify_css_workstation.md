---
subcategory: "Policy Configuration"
---

# centrify_css_workstation attribute

**centrify_css_workstation** is a sub attribute in settings attribute within **centrify_policy** Resource.

## Example Usage

```terraform
resource "centrify_policy" "test_policy" {
    name = "Test Policy"
    description = "Test Policy"
    link_type = "Role"
    policy_assignment = [
        data.centrify_role.system_admin.id,
    ]
    
    settings {
        centrify_css_workstation {
            authentication_enabled = true
            default_profile_id = data.centrify_authenticationprofile.newdevice_auth_pf.id
        }
    }
}
```

More examples can be found [here](https://github.com/centrify/terraform-provider-centrify/blob/main/examples/centrify_policy/policy_centrify_css_workstation.tf)

## Argument Reference

Optional:

- `authentication_enabled` - (Boolean) Enable authentication policy controls.
- `challenge_rule` (Block List) Authentication rules. Refer to [challenge_rule](./attribute_challengerule.md) attribute for details.
- `default_profile_id` - (String) Default Profile (used if no conditions matched)
