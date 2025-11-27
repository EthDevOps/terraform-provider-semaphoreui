package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccProjectIntegrationDataSourceConfigByID(nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration" "test" {
  project_id  = semaphoreui_project.test.id
  name        = "Test Integration %[2]s"
  template_id = semaphoreui_project_template.test.id
}

data "semaphoreui_project_integration" "test" {
  project_id = semaphoreui_project.test.id
  id         = semaphoreui_project_integration.test.id
}
`, testAccProjectIntegrationDependencyConfig(nameSuffix), nameSuffix)
}

func testAccProjectIntegrationDataSourceConfigByName(nameSuffix string) string {
	return fmt.Sprintf(`
%[1]s
resource "semaphoreui_project_integration" "test" {
  project_id  = semaphoreui_project.test.id
  name        = "Test Integration %[2]s"
  template_id = semaphoreui_project_template.test.id
}

data "semaphoreui_project_integration" "test" {
  project_id = semaphoreui_project.test.id
  name       = semaphoreui_project_integration.test.name
}
`, testAccProjectIntegrationDependencyConfig(nameSuffix), nameSuffix)
}

func TestAcc_ProjectIntegrationDataSource_basicID(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIntegrationDataSourceConfigByID(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration.test", "name", fmt.Sprintf("Test Integration %s", nameSuffix)),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration.test", "id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration.test", "template_id"),
				),
			},
		},
	})
}

func TestAcc_ProjectIntegrationDataSource_basicName(t *testing.T) {
	nameSuffix := acctest.RandString(8)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIntegrationDataSourceConfigByName(nameSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.semaphoreui_project_integration.test", "name", fmt.Sprintf("Test Integration %s", nameSuffix)),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration.test", "id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration.test", "project_id"),
					resource.TestCheckResourceAttrSet("data.semaphoreui_project_integration.test", "template_id"),
				),
			},
		},
	})
}
