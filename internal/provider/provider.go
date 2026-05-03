package provider

import (
	"context"
	"fmt"

	"github.com/joseph/terraform-provider-groq/internal/infrastructure"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure MlopsProvider implements the provider.Provider interface.
var _ provider.Provider = &MlopsProvider{}

// MlopsProvider defines the provider implementation.
type MlopsProvider struct {
	version string
}

// MlopsProviderData is the data structure passed to resources.
type MlopsProviderData struct {
	ApiKey       string
	DockerClient *infrastructure.DockerClient
}

// Metadata returns the provider metadata.
func (p *MlopsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "terraform-provider-mlops"
}

// Schema defines the provider-level schema for configuration data.
func (p *MlopsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"groq_api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a Groq API client for data sources and resources.
func (p *MlopsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config struct {
		GroqApiKey types.String `tfsdk:"groq_api_key"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dockerClient, err := infrastructure.NewDockerClient()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Initialize Docker Client",
			fmt.Sprintf("An unexpected error occurred while initializing the Docker client: %s", err.Error()),
		)
		return
	}

	data := &MlopsProviderData{
		ApiKey:       config.GroqApiKey.ValueString(),
		DockerClient: dockerClient,
	}

	resp.ResourceData = data
	resp.DataSourceData = data
}

// DataSources defines the data sources implemented in the provider.
func (p *MlopsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *MlopsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAgentResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MlopsProvider{
			version: version,
		}
	}
}
