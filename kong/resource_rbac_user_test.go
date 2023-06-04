package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccKongRBACUser(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRBACUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRBACUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRBACUserExists("kong_rbacUser.rbacUser"),
					resource.TestCheckResourceAttr("kong_rbacUser.rbacUser", "name", "myrbacUser"),
				),
			},
			{
				Config: testUpdateRBACUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRBACUserExists("kong_rbacUser.rbacUser"),
					resource.TestCheckResourceAttr("kong_rbacUser.rbacUser", "name", "yourrbacUser"),
				),
			},
		},
	})
}

func TestAccKongRBACUserImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRBACUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRBACUserConfig,
			},
			{
				ResourceName:      "kong_rbacUser.rbacUser",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongRBACUserDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.RBACUsers

	rbacUsers := getResourcesByType("kong_rbacUser", state)

	if len(rbacUsers) > 1 {
		return fmt.Errorf("expecting max 1 rbacUser resource, found %v", len(rbacUsers))
	}

	if len(rbacUsers) == 0 {
		return nil
	}

	response, err := client.ListAll(context.Background())

	if err != nil {
		return fmt.Errorf("error thrown when trying to list rbacUsers: %v", err)
	}

	if response != nil {
		for _, element := range response {
			if *element.ID == rbacUsers[0].Primary.ID {
				return fmt.Errorf("rbacUser %s still exists, %+v", rbacUsers[0].Primary.ID, response)
			}
		}
	}

	return nil
}

func testAccCheckKongRBACUserExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.RBACUsers
		rbacUsers, err := client.ListAll(context.Background())

		if !kong.IsNotFoundErr(err) && err != nil {
			return err
		}

		if rbacUsers == nil {
			return fmt.Errorf("rbacUser with id %v not found", rs.Primary.ID)
		}

		if len(rbacUsers) != 2 {
			return fmt.Errorf("expected two rbacUsers (default & just-created), found %v", len(rbacUsers))
		}

		return nil
	}
}

const testRBACUserConfig = `
resource "kong_rbacUser" "rbacUser" {
	name			= "myrbacUser"
}
`
const testUpdateRBACUserConfig = `
resource "kong_rbacUser" "rbacUser" {
	name			= "yourrbacUser"
}
