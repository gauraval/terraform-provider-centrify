package platform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	logger "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/logging"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	jsoniter "github.com/json-iterator/go"
)

// Define this for convience usage
type keyValue map[string]interface{}

var (
	subArgs = make(map[string]interface{})
	// Right reppresents a struct of valid rights
	Right = struct {
		Grant, View, Edit, Delete, Add, Run, Login, Checkout, Retrieve, ManageSession, AgentAuth, OfflineRescue, AddAccount, UnlockAccount, RequestZoneRole, FileTransfer, UpdatePassword, WorkspaceLogin, RotatePassword, RetrieveSecret, ManagementAssignment string
	}{
		Grant:                "Grant",
		View:                 "View",
		Edit:                 "Edit",
		Delete:               "Delete",
		Add:                  "Add",
		Run:                  "Run",
		Login:                "Login",
		Checkout:             "Checkout",
		Retrieve:             "Retrieve",
		ManageSession:        "ManageSession",
		AgentAuth:            "AgentAuth",
		OfflineRescue:        "OfflineRescue",
		AddAccount:           "AddAccount",
		UnlockAccount:        "UnlockAccount",
		RequestZoneRole:      "RequestZoneRole",
		FileTransfer:         "FileTransfer",
		UpdatePassword:       "UpdatePassword",
		WorkspaceLogin:       "WorkspaceLogin",
		RotatePassword:       "RotatePassword",
		RetrieveSecret:       "RetrieveSecret",
		ManagementAssignment: "ManagementAssignment",
	}

	// ValidPermissionMap represents a struct of valid permissions
	ValidPermissionMap = struct {
		Generic, Set, WinNix, System, Database, Domain, Account, DBAccount, DomainAccount, CloudAccount, MultiplexAccount, Secret, SSHKey, Service, Application, Folder map[string]string
	}{
		Generic: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Edit: Right.Edit, Right.Delete: Right.Delete},
		// Set defines valid permissions for Set
		Set: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Edit: Right.Edit, Right.Delete: Right.Delete},
		// WinNix defines valid permissions for Windows and Unix
		WinNix: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.ManageSession: Right.ManageSession, Right.Edit: Right.Edit, Right.Delete: Right.Delete, Right.AgentAuth: Right.AgentAuth, Right.OfflineRescue: Right.OfflineRescue, Right.AddAccount: Right.AddAccount, Right.UnlockAccount: Right.UnlockAccount, Right.ManagementAssignment: "ManagePrivilegeElevationAssignment", Right.RequestZoneRole: Right.RequestZoneRole},
		// System defines valid permissions for other system types other than Windows and Unix
		System: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.ManageSession: Right.ManageSession, Right.Edit: Right.Edit, Right.Delete: Right.Delete, Right.AgentAuth: Right.AgentAuth, Right.OfflineRescue: Right.OfflineRescue, Right.AddAccount: Right.AddAccount, Right.UnlockAccount: Right.UnlockAccount},
		// Database defines valid permissions for database
		Database: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Edit: Right.Edit, Right.Delete: Right.Delete},
		// Domain defines valid permissions for domain
		Domain: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Edit: Right.Edit, Right.Delete: Right.Delete, Right.UnlockAccount: Right.UnlockAccount, Right.AddAccount: Right.AddAccount},
		// Account defines valid permissions for account in general
		Account: map[string]string{Right.Grant: "Owner", Right.View: Right.View, Right.Checkout: "Naked", Right.Login: Right.Login, Right.FileTransfer: Right.FileTransfer, Right.Edit: "Manage", Right.Delete: Right.Delete, Right.UpdatePassword: Right.UpdatePassword, Right.WorkspaceLogin: "UserPortalLogin", Right.RotatePassword: Right.RotatePassword},
		// DBAcount defines valid permissions for database account
		DBAccount: map[string]string{Right.Grant: "Owner", Right.View: Right.View, Right.Checkout: "Naked", Right.Edit: "Manage", Right.Delete: Right.Delete, Right.UpdatePassword: Right.UpdatePassword, Right.RotatePassword: Right.RotatePassword},
		// DomainAccount defines valid permissions for domain account
		DomainAccount: map[string]string{Right.Grant: "Owner", Right.View: Right.View, Right.Checkout: "Naked", Right.Login: Right.Login, Right.FileTransfer: Right.FileTransfer, Right.Edit: "Manage", Right.Delete: Right.Delete, Right.UpdatePassword: Right.UpdatePassword, Right.RotatePassword: Right.RotatePassword},
		// CloudAccount defines valid permissions for cloud provider account
		CloudAccount: map[string]string{Right.Grant: "Owner", Right.View: Right.View, Right.Checkout: "Naked", Right.Login: Right.Login, Right.Edit: "Manage", Right.Delete: Right.Delete, Right.UpdatePassword: Right.UpdatePassword, Right.RotatePassword: Right.RotatePassword},
		// MultiplexAccount defines valid permissions for multiplex account
		MultiplexAccount: map[string]string{Right.Grant: Right.Grant, Right.Edit: Right.Edit, Right.Delete: Right.Delete},
		// Secret defines valid permissions for secret
		Secret: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Edit: Right.Edit, Right.Delete: Right.Delete, Right.RetrieveSecret: "Retrieve"},
		// SSHKey defines valid permissions for ssh key
		SSHKey: map[string]string{Right.Grant: "Owner", Right.View: Right.View, Right.Retrieve: "Checkout", Right.Edit: "Manage", Right.Delete: Right.Delete},
		// Service defines valid permissions for service
		Service: map[string]string{Right.Grant: Right.Grant, Right.Edit: Right.Edit, Right.Delete: Right.Delete},
		// App defines valid permissions for application
		Application: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Run: "Execute"},
		// Folder defines valid permissions for secret folder
		Folder: map[string]string{Right.Grant: Right.Grant, Right.View: Right.View, Right.Edit: Right.Edit, Right.Delete: Right.Delete, Right.Add: Right.Add},
	}
)

func init() {
	subArgs["Caching"] = -1
	//subArgs["PageSize"] = 10000
	//subArgs["Limit"] = 10000
}

// mapToStruct takes map as input and populate struct attribute accordingly
func mapToStruct(i interface{}, m map[string]interface{}) error {
	jsonString, _ := json.Marshal(m)
	err := json.Unmarshal(jsonString, i)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal map: %v", err)
	}

	return nil
}

// generateRequestMap takes struct object and convert it to map
//func generateRequestMap(i vaultObjectInterface) (map[string]interface{}, error) {
func generateRequestMap(i interface{}) (map[string]interface{}, error) {
	var mapData = make(map[string]interface{})
	dataBytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}

	return mapData, nil
}

// GenerateSchemaMap converts object into map according to object's json schema definition
func GenerateSchemaMap(i interface{}) (map[string]interface{}, error) {
	var mapData = make(map[string]interface{})
	schemaJSON := jsoniter.Config{TagKey: "schema", OnlyTaggedField: true}.Froze()
	dataBytes, err := schemaJSON.Marshal(i)
	if err != nil {
		panic(err)
	}
	err = schemaJSON.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}

	return mapData, nil
}

// flattenNestedMap converts nested map to flat map. It is used by Policy object.
// It is assumed that json tag of each nested struct element is unique
func flattenNestedMap(flatMap map[string]interface{}, nestedMap interface{}) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}:
			if err := flattenNestedMap(flatMap, v); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nestedMap.(type) {
	case map[string]interface{}:
		for k, v := range nestedMap.(map[string]interface{}) {
			assign(k, v)
		}
	default:
		return fmt.Errorf("Not a valid input, must be a map")
	}

	return nil
}

func flattenSettings(flatten map[string]interface{}, nestedMap interface{}) error {
	if nestedMap != nil {
		for k1, v1 := range nestedMap.(map[string]interface{}) {
			// This is first level that deals with CentrifyServices, CentrifyClient, CentrifyCSSServer, etc.
			switch v1.(type) {
			case map[string]interface{}:
				// this is second level that deals with actual attributes but also may be map such as
				// 		/Core/Css/WinClient/AuthenticationRules
				// 		/Core/PasswordReset/ADAdminPass
				for k2, v2 := range v1.(map[string]interface{}) {
					flatten[k2] = v2
				}
			default:
				flatten[k1] = v1
			}
		}
	}
	return nil
}

// RedRockQuery issues RedRock API query
func RedRockQuery(client *restapi.RestClient, query string, args map[string]interface{}) ([]interface{}, error) {
	var queryArg = make(map[string]interface{})
	queryArg["Script"] = query
	if args == nil {
		queryArg["Args"] = subArgs
	} else {
		queryArg["Args"] = args
	}

	logger.Debugf("Query arguments: %+v", queryArg)
	resp, err := client.CallGenericMapAPI("/RedRock/query", queryArg)
	//logger.Debugf("Query response: %+v", resp)
	if err != nil {
		logger.ErrorTracef(err.Error())
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf(resp.Message)
	}

	var results = resp.Result["Results"].([]interface{})

	return results, nil
}

func queryVaultObject(client *restapi.RestClient, query string) (map[string]interface{}, error) {
	results, err := RedRockQuery(client, query, nil)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		errmsg := "Query returns 0 object"
		//logger.Errorf(errmsg)
		logger.ErrorTracef(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
	if len(results) > 1 {
		errmsg := fmt.Sprintf("Query returns too many objects (found %d, expected 1)", len(results))
		//logger.Errorf(errmsg)
		logger.ErrorTracef(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
	var result = results[0].(map[string]interface{})
	var row = result["Row"].(map[string]interface{})

	return row, nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

// Insert an element into slice at position i
func insert(a []map[string]interface{}, c map[string]interface{}, i int) []map[string]interface{} {
	return append(a[:i], append([]map[string]interface{}{c}, a[i:]...)...)
}

// Find the intersection of two iterable values.
func intersect(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		if contains(b, el) {
			set = append(set, el)
		}
	}

	return set
}

// contains test if e is an element of a
func contains(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}

// FlattenSliceToString converts ["value1", "value2"] to "value1,value2"
func FlattenSliceToString(input []string) string {
	var str string
	for i, v := range input {
		str = str + v
		// Append "," if it is not the last element
		if i < len(input)-1 {
			str = str + ","
		}
	}

	return str
}

func validatePermissions(permissions []Permission, valid []string) error {
	if permissions != nil {
		for _, v := range permissions {
			rights := strings.Split(v.Rights, ",")
			if len(intersect(rights, valid)) != len(rights) {
				errmsg := fmt.Sprintf("%v can only contain %v", rights, valid)
				logger.ErrorTracef(errmsg)
				return fmt.Errorf(errmsg)
				//return fmt.Errorf("%v can only contain %v", rights, valid)
			}
		}
	}

	return nil
}

// ConvertToValidList converts provide list of rights to actual values that can be used for API call
// Converts []string{"a1", "b1"} to []string{"a2", "b2"} from map[string]string{"a1": "a2", "b1": "b2"}
func ConvertToValidList(input []string, validMap map[string]string) ([]string, error) {
	var converted []string
	for _, k := range input {
		v := validMap[k]
		if v != "" {
			converted = append(converted, v)
		} else {
			errmsg := fmt.Sprintf("Invalid right %s", k)
			logger.ErrorTracef(errmsg)
			return nil, fmt.Errorf(errmsg)
			//return nil, fmt.Errorf("Invalid right %s", k)
		}
	}

	return converted, nil
}

// ResolvePermissions given a list of Permissions, resolve PrincipalID and convert the given rights to actual rights
func ResolvePermissions(c *restapi.RestClient, perms []Permission, validPerms map[string]string) error {
	var err error
	for i, p := range perms {
		// Resolove PrincipalID
		switch strings.ToLower(p.PrincipalType) {
		case "user":
			user := NewUser(c)
			user.Name = p.PrincipalName
			perms[i].PrincipalID, err = user.GetIDByName()
			if err != nil {
				return err
			}
		case "role":
			role := NewRole(c)
			role.Name = p.PrincipalName
			perms[i].PrincipalID, err = role.GetIDByName()
			if err != nil {
				return err
			}
		default:
			errmsg := fmt.Sprintf("Invalid PrincipalType %s", p.PrincipalType)
			logger.ErrorTracef(errmsg)
			return fmt.Errorf(errmsg)
			//return fmt.Errorf("Invalid PrincipalType %s", p.PrincipalType)
		}

		// Convert rights
		var rights []string
		if p.Rights != "" {
			rights, err = ConvertToValidList(strings.Split(p.Rights, ","), validPerms)
		} else if p.RightList != nil {
			rights, err = ConvertToValidList(p.RightList, validPerms)
		}
		if err != nil {
			return err
		}
		perms[i].Rights = FlattenSliceToString(rights)
	}
	logger.Debugf("Resolved permissions: %+v", perms)

	return nil
}

// ResolvePermissions2 detects if PrincipalID is set, if not then resolve it
func ResolvePermissions2(c *restapi.RestClient, perms []Permission, validPerms map[string]string) error {
	var err error
	for i, p := range perms {
		// Resolove PrincipalID
		if perms[i].PrincipalID == "" {
			switch strings.ToLower(p.PrincipalType) {
			case "user":
				user := NewUser(c)
				user.Name = p.PrincipalName
				perms[i].PrincipalID, err = user.GetIDByName()
				if err != nil {
					return err
				}
			case "role":
				role := NewRole(c)
				role.Name = p.PrincipalName
				perms[i].PrincipalID, err = role.GetIDByName()
				if err != nil {
					return err
				}
			default:
				errmsg := fmt.Sprintf("Invalid PrincipalType %s", p.PrincipalType)
				logger.ErrorTracef(errmsg)
				return fmt.Errorf(errmsg)
				//return fmt.Errorf("Invalid PrincipalType %s", p.PrincipalType)
			}
		}

		// Convert rights
		var rights []string
		if p.Rights != "" {
			rights, err = ConvertToValidList(strings.Split(p.Rights, ","), validPerms)
		} else if p.RightList != nil {
			rights, err = ConvertToValidList(p.RightList, validPerms)
		}
		if err != nil {
			return err
		}
		perms[i].Rights = FlattenSliceToString(rights)
	}
	logger.Debugf("Resolved permissions: %+v", perms)

	return nil
}

func noFoundError() string {
	errmsg := "Query returns 0 object"
	logger.Errorf(errmsg)
	return errmsg
}

func foundTooManyError(no int) string {
	errmsg := fmt.Sprintf("Query returns too many objects (found %d, expected 1)", no)
	logger.Errorf(errmsg)
	return errmsg
}

func queryError(no int) error {
	if no == 0 {
		return fmt.Errorf(noFoundError())
	}
	if no > 1 {
		return fmt.Errorf(foundTooManyError(no))
	}
	return nil
}

// GetVarType returns variable type name as string
func GetVarType(myvar interface{}) string {
	valueOf := reflect.ValueOf(myvar)
	var varType string
	if valueOf.Type().Kind() == reflect.Ptr {
		varType = fmt.Sprintf(reflect.Indirect(valueOf).Type().Name())
	} else {
		varType = fmt.Sprintf(valueOf.Type().Name())
	}
	return varType
}

func FlattenWorkflowApprovers(approvers []WorkflowApprover) string {
	approvers_str := ""

	for _, v := range approvers {
		fields := reflect.TypeOf(v)
		values := reflect.ValueOf(v)
		num := fields.NumField()
		approver_str := ""
		for i := 0; i < num; i++ {
			field := fields.Field(i)
			value := values.Field(i)
			//logger.Debugf("Field name %s value is %v", field.Name, value)
			if field.Name != "DirectoryService" && field.Name != "DirectoryName" {
				//logger.Debugf("Field name %s is being flatenning", field.Name)
				switch field.Name {
				case "OptionsSelector":
					if value.Bool() {
						convertred_str := "\"" + field.Name + "\":true"
						if approver_str == "" {
							approver_str = approver_str + convertred_str
						} else {
							approver_str = approver_str + "," + convertred_str
						}
					}
				case "BackupApprover":
					if !value.IsNil() {
						originalValue := value.Elem()
						ba_num := originalValue.NumField()
						if ba_num > 0 {
							backupapprover_str := ""
							for j := 0; j < ba_num; j++ {
								ba_field := originalValue.Type().Field(j)
								ba_value := originalValue.Field(j)
								if ba_field.Name != "DirectoryService" && ba_field.Name != "DirectoryName" {
									if backupapprover_str == "" {
										backupapprover_str = backupapprover_str + "\"" + ba_field.Name + "\":\"" + ba_value.String() + "\""
									} else {
										backupapprover_str = backupapprover_str + "," + "\"" + ba_field.Name + "\":\"" + ba_value.String() + "\""
									}
								}
							}

							convertred_str := "\"" + field.Name + "\":{" + backupapprover_str + "}"
							if approver_str == "" {
								approver_str = approver_str + convertred_str
							} else {
								approver_str = approver_str + "," + convertred_str
							}
						}
					}
				default:
					if value.String() != "" {
						convertred_str := "\"" + field.Name + "\":\"" + value.String() + "\""
						if approver_str == "" {
							approver_str = approver_str + convertred_str
						} else {
							approver_str = approver_str + "," + convertred_str
						}
					}
				}
			}
		}
		if approvers_str == "" {
			approvers_str = "[{" + approver_str + "}"
		} else {
			approvers_str = approvers_str + ",{" + approver_str + "}"
		}
	}
	// Finally, close the list
	if approvers_str != "" {
		approvers_str = approvers_str + "]"
	}

	return approvers_str
}

func ResolveWorkflowApprovers(c *restapi.RestClient, approvers []WorkflowApprover) error {
	for i, v := range approvers {
		// Resolve guid of user that is from Active Directory, LDAP, Google Directory or Federated Directory
		if v.Type != "Manager" && v.Guid == "" {
			// Get directory service
			dirs := NewDirectoryServices(c)
			dir, err := dirs.GetByName(v.DirectoryService, v.DirectoryName)
			if err != nil {
				return err
			}
			//logger.Debugf("Resolved Directory Service: %+v", dir)
			// Get directory object
			objs := NewDirectoryObjects(c)
			obj, err := objs.GetByName(v.Type, v.Name, *dir)
			if err != nil {
				return err
			}
			//logger.Debugf("Resolved Directory Object: %+v", obj)
			if v.Type == "Role" {
				v.Guid = obj.RoleID
			} else {
				v.Guid = obj.ID
			}
		} else if v.Type == "Manager" && v.BackupApprover != nil && v.BackupApprover.Guid == "" {
			backup := v.BackupApprover
			// Get directory service
			dirs := NewDirectoryServices(c)
			dir, err := dirs.GetByName(backup.DirectoryService, backup.DirectoryName)
			if err != nil {
				return err
			}
			// Get directory object
			objs := NewDirectoryObjects(c)
			obj, err := objs.GetByName(backup.Type, backup.Name, *dir)
			if err != nil {
				return err
			}
			if backup.Type == "Role" {
				v.BackupApprover.Guid = obj.RoleID
			} else {
				v.BackupApprover.Guid = obj.ID
			}
		}
		approvers[i] = v
	}
	return nil
}

func GetAllZoneRoles(c *restapi.RestClient, domainid string) (map[string]ZoneRole, error) {
	var queryArgs = make(map[string]interface{})
	queryArgs["PageNumber"] = 1
	queryArgs["PageSize"] = 100000
	queryArgs["Limit"] = 100000
	queryArgs["SortBy"] = ""
	queryArgs["direction"] = false
	queryArgs["Caching"] = -1

	var requestArg = make(map[string]interface{})
	requestArg["DomainId"] = domainid
	requestArg["Args"] = subArgs

	// Attempt to read from an upstream API
	resp, err := c.CallSliceAPI("/ZoneRoleWorkflow/GetAllRoles", requestArg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	if resp.Success {
		var zonerolemap = make(map[string]ZoneRole)
		results := resp.Result
		for _, v := range results {
			zonerole := &ZoneRole{}
			mapToStruct(zonerole, v.(map[string]interface{}))
			if zonerole != nil {
				zonerolemap[zonerole.Name] = *zonerole
			}
		}
		return zonerolemap, nil
	} else {
		errmsg := fmt.Sprintf("%s %s", resp.Message, resp.Exception)
		logger.Errorf(errmsg)
		return nil, fmt.Errorf(errmsg)
	}
}

func FlattenZoneRoles(zoneroles []ZoneRole) string {
	zoneroles_str := ""

	for _, v := range zoneroles {
		fields := reflect.TypeOf(v)
		values := reflect.ValueOf(v)
		num := fields.NumField()
		zonerole_str := ""

		for i := 0; i < num; i++ {
			field := fields.Field(i)
			value := values.Field(i)
			//logger.Debugf("Field name %s value is %v", field.Name, value)
			switch field.Name {
			case "Windows", "Unix":
				if value.Bool() {
					convertred_str := "\"" + field.Name + "\":true"
					if zonerole_str == "" {
						zonerole_str = zonerole_str + convertred_str
					} else {
						zonerole_str = zonerole_str + "," + convertred_str
					}
				}

			default:
				if value.String() != "" {
					convertred_str := "\"" + field.Name + "\":\"" + value.String() + "\""
					if zonerole_str == "" {
						zonerole_str = zonerole_str + convertred_str
					} else {
						zonerole_str = zonerole_str + "," + convertred_str
					}
				}
			}
		}
		if zoneroles_str == "" {
			zoneroles_str = "[{" + zonerole_str + "}"
		} else {
			zoneroles_str = zoneroles_str + ",{" + zonerole_str + "}"
		}
	}
	// Finally, close the list
	if zoneroles_str != "" {
		zoneroles_str = zoneroles_str + "]"
	}

	return zoneroles_str
}

func resolveZoneRoles(c *restapi.RestClient, zoneroles []ZoneRole, domainid string) error {
	if domainid == "" {
		errmsg := "missing domain id"
		logger.ErrorTracef(errmsg)
		return fmt.Errorf(errmsg)
	}
	// Retrieve all zone roles
	var allzoneroles = make(map[string]ZoneRole)
	var err error
	if len(zoneroles) > 0 {
		allzoneroles, err = GetAllZoneRoles(c, domainid)
		if err != nil {
			return err
		}
	}

	// Fill zone role attributes
	for i, v := range zoneroles {
		// Resolve attributes for each zone role
		if v.Name != "" {
			zoneroles[i] = allzoneroles[v.Name]
		}
	}

	return nil
}
