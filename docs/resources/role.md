---
subcategory: "Access"
---

# centrify_role (Resource)

This resource allows you to create/update/delete role.

## Example Usage

```terraform
resource "centrify_role" "test_role" {
    name = "Test Role"
    description = "Test role with basic admin right"
    adminrights = [
        "Privileged Access Service User",
    ]
}
```

More examples can be found [here](https://github.com/centrify/terraform-provider-centrify/tree/main/examples/centrify_role)

## Argument Reference

### Required

- `name` - (String) Name of the role.

### Optional

- `description` - (String) Description of an role.
- `adminrights` - (Set of String) List of administrative rights.
- `member` - (Block Set) (see [below reference for member](#reference-for-member))

## [Reference for `member`]

Required:

- `id` - (String) ID of the member.
- `type` - (String) Type of the member. Can be set to `User`, `Group` or `Role`.

## Import

Role can be imported using the resource `id`, e.g.

```shell
terraform import centrify_role.example xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
