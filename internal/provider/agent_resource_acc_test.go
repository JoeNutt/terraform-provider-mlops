package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAgentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentResourceConfig("test-acc-agent", "alpine:latest"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mlops_agent.test", "name", "test-acc-agent"),
					resource.TestCheckResourceAttr("mlops_agent.test", "image", "alpine:latest"),
					resource.TestCheckResourceAttr("mlops_agent.test", "llm_provider", "terraform-provider-mlops"),
					resource.TestCheckResourceAttr("mlops_agent.test", "llm_model", "llama3-8b-8192"),
				),
			},
		},
	})
}

func testAccAgentResourceConfig(name, image string) string {
	return `
provider "mlops" {
  groq_api_key = "dummy-key"
}

resource "mlops_agent" "test" {
  name         = "` + name + `"
  image        = "` + image + `"
  llm_provider = "terraform-provider-mlops"
  llm_model    = "llama3-8b-8192"
}
`
}
