---
subcategory: "Policy Configuration"
---

# sshkey_set attribute

**sshkey_set** is a sub attribute in settings attribute within **centrify_policy** Resource.

## Example Usage

```terraform
resource "centrify_policy" "test_policy" {
    name = "Test Policy"
    description = "Test Policy"
    link_type = "Collection"
    policy_assignment = [
        format("SshKeys|%s", data.centrify_manualset.test_set.id),
    ]
    
    settings {
        sshkey_set {
            default_profile_id = data.centrify_authenticationprofile.newdevice_auth_pf.id
            challenge_rule {
                authentication_profile_id = data.centrify_authenticationprofile.newdevice_auth_pf.id
                rule {
                    filter = "IpAddress"
                    condition = "OpInCorpIpRange"
                }
            }
        }
    }
}
```

More examples can be found [here](https://github.com/centrify/terraform-provider-centrify/blob/main/examples/centrify_policy/policy_sshkey_set.tf)

## Argument Reference

Optional:

- `challenge_rule` (Block List) Authentication rules. Refer to [challenge_rule](./attribute_challengerule.md) attribute for details.
- `default_profile_id` - (String) Default Profile (used if no conditions matched).
