package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccProjectIntegrationExtractValueExists(resourceName string) resource.TestCheckFunc {
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
		if rs.Primary.Attributes["integration_id"] == "" {
			return fmt.Errorf("no IntegrationID is set")
		}

		id, _ := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		projectId, _ := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		integrationId, _ := strconv.ParseInt(rs.Primary.Attributes["integration_id"], 10, 64)

		model, err := getExtractValueByID(testClient(), projectId, integrationId, id)
		if err != nil {
			return fmt.Errorf("error reading project integration extract value: %s", err.Error())
		}

		if model.Name.ValueString() != rs.Primary.Attributes["name"] {
			return fmt.Errorf("extract value name mismatch: %s != %s", model.Name.ValueString(), rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccProjectIntegrationExtractValueDependencyConfig(nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration" "test" {
  project_id  = semaphoreui_project.test.id
  name        = "Integration %[2]s"
  template_id = semaphoreui_project_template.test.id
}
`, testAccProjectIntegrationDependencyConfig(nameSuffix), nameSuffix)
}

func testAccProjectIntegrationExtractValueConfig(nameSuffix string, name string, variable string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration_extract_value" "test" {
  project_id     = semaphoreui_project.test.id
  integration_id = semaphoreui_project_integration.test.id
  name           = "%[2]s"
  value_source   = "body"
  body_data_type = "json"
  key            = "$.ref"
  variable       = "%[3]s"
  variable_type  = "environment"
}
`, testAccProjectIntegrationExtractValueDependencyConfig(nameSuffix), name, variable)
}

func testAccProjectIntegrationExtractValueImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}

		return fmt.Sprintf("project/%[1]s/integration/%[2]s/extractvalue/%[3]s",
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["integration_id"],
			rs.Primary.Attributes["id"]), nil
	}
}

func TestAcc_ProjectIntegrationExtractValueResource_basic(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectIntegrationExtractValueConfig(nameSuffix, fmt.Sprintf("ExtractValue %s", nameSuffix), "GIT_REF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccProjectIntegrationExtractValueExists("semaphoreui_project_integration_extract_value.test"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "name", fmt.Sprintf("ExtractValue %s", nameSuffix)),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "value_source", "body"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "body_data_type", "json"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "key", "$.ref"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "variable", "GIT_REF"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "variable_type", "environment"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_extract_value.test", "id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_extract_value.test", "project_id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_extract_value.test", "integration_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "semaphoreui_project_integration_extract_value.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccProjectIntegrationExtractValueImportID("semaphoreui_project_integration_extract_value.test"),
			},
			// Update testing
			{
				Config: testAccProjectIntegrationExtractValueConfig(nameSuffix, fmt.Sprintf("ExtractValue %s", nameSuffix), "GIT_BRANCH"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccProjectIntegrationExtractValueExists("semaphoreui_project_integration_extract_value.test"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "name", fmt.Sprintf("ExtractValue %s", nameSuffix)),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_extract_value.test", "variable", "GIT_BRANCH"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_extract_value.test", "id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_extract_value.test", "project_id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_extract_value.test", "integration_id"),
				),
			},
			// Delete testing
			{
				Config: testAccProjectIntegrationExtractValueDependencyConfig(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceNotExists("semaphoreui_project_integration_extract_value.test"),
				),
			},
		},
	})
}
