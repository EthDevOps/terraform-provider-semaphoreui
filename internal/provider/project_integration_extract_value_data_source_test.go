package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccProjectIntegrationExtractValueDataSourceConfigByID(nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration_extract_value" "test" {
  project_id     = semaphoreui_project.test.id
  integration_id = semaphoreui_project_integration.test.id
  name           = "ExtractValue %[2]s"
  value_source   = "body"
  body_data_type = "json"
  key            = "$.ref"
  variable       = "GIT_REF"
  variable_type  = "environment"
}

data "semaphoreui_project_integration_extract_value" "test" {
  project_id     = semaphoreui_project.test.id
  integration_id = semaphoreui_project_integration.test.id
  id             = semaphoreui_project_integration_extract_value.test.id
}
`, testAccProjectIntegrationExtractValueDependencyConfig(nameSuffix), nameSuffix)
}

func TestAcc_ProjectIntegrationExtractValueDataSource_basicID(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIntegrationExtractValueDataSourceConfigByID(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_extract_value.test", "name", fmt.Sprintf("ExtractValue %s", nameSuffix)),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_extract_value.test", "value_source", "body"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_extract_value.test", "body_data_type", "json"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_extract_value.test", "key", "$.ref"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_extract_value.test", "variable", "GIT_REF"),
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration_extract_value.test", "variable_type", "environment"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration_extract_value.test", "id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration_extract_value.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration_extract_value.test", "integration_id"),
				),
			},
		},
	})
}
