resource "centrifyvault_manualset" "my_accounts" {
    name = "My Accounts"
    type = "VaultAccount"
    description = "This Set contains my accounts."

    permission {
        principal_id = data.centrifyvault_user.admin.id
        principal_name = data.centrifyvault_user.admin.id
        principal_type = "User"
        rights = ["Grant","View","Edit","Delete"]
    }

    permission {
        principal_id = data.centrifyvault_role.system_admin.id
        principal_name = data.centrifyvault_role.system_admin.name
        principal_type = "Role"
        rights = ["Grant","View"]
    }

    member_permission {
        principal_id = data.centrifyvault_role.system_admin.id
        principal_name = data.centrifyvault_role.system_admin.name
        principal_type = "Role"
        rights = ["Grant","View","Checkout","Login","FileTransfer","Edit","Delete","UpdatePassword","WorkspaceLogin","RotatePassword"]
    }
}