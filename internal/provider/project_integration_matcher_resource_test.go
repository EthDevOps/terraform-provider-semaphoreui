package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccProjectIntegrationMatcherExists(resourceName string) resource.TestCheckFunc {
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

		model, err := getMatcherByID(testClient(), projectId, integrationId, id)
		if err != nil {
			return fmt.Errorf("error reading project integration matcher: %s", err.Error())
		}

		if model.Name.ValueString() != rs.Primary.Attributes["name"] {
			return fmt.Errorf("matcher name mismatch: %s != %s", model.Name.ValueString(), rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccProjectIntegrationMatcherDependencyConfig(nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration" "test" {
  project_id  = semaphoreui_project.test.id
  name        = "Integration %[2]s"
  template_id = semaphoreui_project_template.test.id
}
`, testAccProjectIntegrationDependencyConfig(nameSuffix), nameSuffix)
}

func testAccProjectIntegrationMatcherConfig(nameSuffix string, name string, value string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration_matcher" "test" {
  project_id     = semaphoreui_project.test.id
  integration_id = semaphoreui_project_integration.test.id
  name           = "%[2]s"
  match_type     = "body"
  method         = "equals"
  body_data_type = "json"
  key            = "$.event"
  value          = "%[3]s"
}
`, testAccProjectIntegrationMatcherDependencyConfig(nameSuffix), name, value)
}

func testAccProjectIntegrationMatcherImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}

		return fmt.Sprintf("project/%[1]s/integration/%[2]s/matcher/%[3]s",
			rs.Primary.Attributes["project_id"],
			rs.Primary.Attributes["integration_id"],
			rs.Primary.Attributes["id"]), nil
	}
}

func TestAcc_ProjectIntegrationMatcherResource_basic(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectIntegrationMatcherConfig(nameSuffix, fmt.Sprintf("Matcher %s", nameSuffix), "push"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccProjectIntegrationMatcherExists("semaphoreui_project_integration_matcher.test"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "name", fmt.Sprintf("Matcher %s", nameSuffix)),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "match_type", "body"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "method", "equals"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "body_data_type", "json"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "key", "$.event"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "value", "push"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_matcher.test", "id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_matcher.test", "project_id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_matcher.test", "integration_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "semaphoreui_project_integration_matcher.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccProjectIntegrationMatcherImportID("semaphoreui_project_integration_matcher.test"),
			},
			// Update testing
			{
				Config: testAccProjectIntegrationMatcherConfig(nameSuffix, fmt.Sprintf("Matcher %s", nameSuffix), "pull_request"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccProjectIntegrationMatcherExists("semaphoreui_project_integration_matcher.test"),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "name", fmt.Sprintf("Matcher %s", nameSuffix)),
					resource.TestCheckResourceAttr("semaphoreui_project_integration_matcher.test", "value", "pull_request"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_matcher.test", "id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_matcher.test", "project_id"),
					resource.TestCheckResourceAttrSet("semaphoreui_project_integration_matcher.test", "integration_id"),
				),
			},
			// Delete testing
			{
				Config: testAccProjectIntegrationMatcherDependencyConfig(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccResourceNotExists("semaphoreui_project_integration_matcher.test"),
				),
			},
		},
	})
}
