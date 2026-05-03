package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure agentResource implements the resource.Resource interface.
var _ resource.Resource = &agentResource{}

type agentResource struct {
	data *MlopsProviderData
}

func NewAgentResource() resource.Resource {
	return &agentResource{}
}

func (r *agentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (r *agentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*MlopsProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MlopsProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.data = data
}

func (r *agentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"image": schema.StringAttribute{
				Required: true,
			},
			"llm_provider": schema.StringAttribute{
				Required: true,
			},
			"llm_model": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

type agentResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Image       types.String `tfsdk:"image"`
	LlmProvider types.String `tfsdk:"llm_provider"`
	LlmModel    types.String `tfsdk:"llm_model"`
}

func (r *agentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.data == nil {
		resp.Diagnostics.AddError(
			"Provider Not Configured",
			"The provider hasn't been configured before creating the resource.",
		)
		return
	}

	var plan agentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.data.DockerClient.DeployAgent(
		ctx,
		plan.Name.ValueString(),
		plan.Image.ValueString(),
		r.data.ApiKey,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deploying Agent",
			fmt.Sprintf("Could not deploy agent: %s", err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *agentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *agentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *agentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.data == nil {
		resp.Diagnostics.AddError(
			"Provider Not Configured",
			"The provider hasn't been configured before deleting the resource.",
		)
		return
	}

	var state agentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.data.DockerClient.RemoveAgent(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing Agent",
			fmt.Sprintf("Could not remove agent: %s", err.Error()),
		)
		return
	}
}
