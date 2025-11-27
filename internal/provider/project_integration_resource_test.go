package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccProjectIntegrationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ProjectID is set")
		}

		id, _ := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		projectId, _ := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)

		model, err := getIntegrationByID(testClient(), projectId, id)
		if err != nil {
			return fmt.Errorf("error reading project integration: %s", err.Error())
		}

		if model.Name.ValueString() != rs.Primary.Attributes["name"] {
			return fmt.Errorf("integration name mismatch: %s != %s", model.Name.ValueString(), rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccProjectIntegrationDependencyConfig(nameSuffix string) string {
	return fmt.Sprintf(`
resource "semaphoreui_project" "test" {
  name = "test-%[1]s"
}

resource "semaphoreui_project_key" "test" {
  project_id = semaphoreui_project.test.id
  name       = "None-%[1]s"
  none       = {}
}

resource "semaphoreui_project_repository" "test" {
  project_id = semaphoreui_project.test.id
  name       = "Repo-%[1]s"
  url        = "git@github.com:example/test.git"
  branch     = "main"
  ssh_key_id = semaphoreui_project_key.test.id
}

resource "semaphoreui_project_inventory" "test" {
  project_id = semaphoreui_project.test.id
  name       = "Inventory-%[1]s"
  ssh_key_id = semaphoreui_project_key.test.id
  file = {
    path          = "path/to/inventory"
    repository_id = semaphoreui_project_repository.test.id
  }
}

resource "semaphoreui_project_environment" "test" {
  project_id = semaphoreui_project.test.id
  name       = "Env-%[1]s"
}

resource "semaphoreui_project_template" "test" {
  project_id     = semaphoreui_project.test.id
  environment_id = semaphoreui_project_environment.test.id
  inventory_id   = semaphoreui_project_inventory.test.id
  repository_id  = semaphoreui_project_repository.test.id
  name           = "Template-%[1]s"
  playbook       = "playbook.yml"
}
`, nameSuffix)
}

func testAccProjectIntegrationConfig(nameSuffix string, name string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration" "test" {
  project_id  = semaphoreui_project.test.id
  name        = "%[2]s"
  template_id = semaphoreui_project_template.test.id
}
`, testAccProjectIntegrationDependencyConfig(nameSuffix), name)
}

func testAccProjectIntegrationImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}

		return fmt.Sprintf("project/%[1]s/integration/%[2]s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAcc_ProjectIntegrationResource_basic(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectIntegrationConfig(nameSuffix, fmt.Sprintf("Test Integration %s", nameSuffix)),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccProjectIntegrationExists("semaphoreui_project_integration.test"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration.test", "name", fmt.Sprintf("Test Integration %s", nameSuffix)),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration.test", "id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration.test", "project_id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration.test", "template_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "semaphoreui_project_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccProjectIntegrationImportID("semaphoreui_project_integration.test"),
			},
			// Update testing
			{
				Config: testAccProjectIntegrationConfig(nameSuffix, fmt.Sprintf("Updated Integration %s", nameSuffix)),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccProjectIntegrationExists("semaphoreui_project_integration.test"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration.test", "name", fmt.Sprintf("Updated Integration %s", nameSuffix)),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration.test", "id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration.test", "project_id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration.test", "template_id"),
				),
			},
			// Delete testing
			{
				Config: testAccProjectIntegrationDependencyConfig(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceNotExists("semaphoreui_project_integration.test"),
				),
			},
		},
	})
}
