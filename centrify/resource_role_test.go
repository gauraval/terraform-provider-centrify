package centrify

import (
	"fmt"
	"regexp"
	"testing"

	vault "github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/platform"
	"github.com/centrify/terraform-provider-centrify/cloud-golang-sdk/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceRoleCreation(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "centrify_role.testrole"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBasicDataExists(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccBasicDataExists(rName string) string {
	return fmt.Sprintf(`resource "centrify_role" "testrole" {
		name = %[1]q
	}`, rName)
}

func testAccCheckRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*restapi.RestClient)
	object := vault.NewUser(client)
	for _, res := range s.RootModule().Resources {
		if res.Type != "centrify_role" {
			continue
		}
		object.ID = res.Primary.ID
		err := object.Read()
		if err == nil {
			return fmt.Errorf("Role Still Exists")
		}

		if err != nil {
			notFoundErr := "not found"
			expectedErr := regexp.MustCompile(notFoundErr)
			if !expectedErr.Match([]byte(err.Error())) {
				return fmt.Errorf("expected %s, got %s", notFoundErr, err)
			}
		}
	}
	return nil

}
