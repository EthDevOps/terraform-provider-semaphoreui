package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccProjectIntegrationMatcherDataSourceConfigByID(nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration_matcher" "test" {
  project_id     = semaphoreui_project.test.id
  integration_id = semaphoreui_project_integration.test.id
  name           = "Matcher %[2]s"
  match_type     = "body"
  method         = "equals"
  body_data_type = "json"
  key            = "$.event"
  value          = "push"
}

data "semaphoreui_project_integration_matcher" "test" {
  project_id     = semaphoreui_project.test.id
  integration_id = semaphoreui_project_integration.test.id
  id             = semaphoreui_project_integration_matcher.test.id
}
`, testAccProjectIntegrationMatcherDependencyConfig(nameSuffix), nameSuffix)
}

func TestAcc_ProjectIntegrationMatcherDataSource_basicID(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIntegrationMatcherDataSourceConfigByID(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_matcher.test", "name", fmt.Sprintf("Matcher %s", nameSuffix)),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_matcher.test", "match_type", "body"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_matcher.test", "method", "equals"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_matcher.test", "body_data_type", "json"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_matcher.test", "key", "$.event"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_matcher.test", "value", "push"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration_matcher.test", "id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration_matcher.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration_matcher.test", "integration_id"),
				),
			},
		},
	})
}
