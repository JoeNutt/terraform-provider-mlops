package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAgentResourceSchema(t *testing.T) {
	ctx := context.Background()

	p := &MlopsProvider{}
	resources := p.Resources(ctx)

	var agentResource resource.Resource
	for _, resFunc := range resources {
		res := resFunc()
		// In a real scenario, we'd check metadata.
		// For this test, we expect to find a resource that we can then check.
		// Since p.Resources(ctx) returns nil currently, agentResource will remain nil.
		agentResource = res
		break
	}

	if agentResource == nil {
		t.Fatal("ml_agent resource not implemented or not returned by provider")
	}

	resp := &resource.SchemaResponse{}
	agentResource.Schema(ctx, resource.SchemaRequest{}, resp)

	expectedAttributes := []struct {
		name     string
		required bool
	}{
		{"name", true},
		{"image", true},
		{"llm_provider", true},
		{"llm_model", true},
	}

	for _, tc := range expectedAttributes {
		t.Run(tc.name, func(t *testing.T) {
			attr, ok := resp.Schema.Attributes[tc.name]
			if !ok {
				t.Fatalf("Attribute %s not found in schema", tc.name)
			}

			sAttr, ok := attr.(schema.StringAttribute)
			if !ok {
				t.Fatalf("Attribute %s is not a StringAttribute, got %T", tc.name, attr)
			}

			if sAttr.Required != tc.required {
				t.Errorf("Attribute %s: expected required %v, got %v", tc.name, tc.required, sAttr.Required)
			}
		})
	}
}

func TestAgentResourceCreate(t *testing.T) {
	ctx := context.Background()
	res := NewAgentResource()

	planValue := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":         tftypes.String,
				"image":        tftypes.String,
				"llm_provider": tftypes.String,
				"llm_model":    tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"name":         tftypes.NewValue(tftypes.String, "test-agent"),
			"image":        tftypes.NewValue(tftypes.String, "alpine:latest"),
			"llm_provider": tftypes.NewValue(tftypes.String, "terraform-provider-mlops"),
			"llm_model":    tftypes.NewValue(tftypes.String, "llama3-8b-8192"),
		},
	)

	req := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw: planValue,
			Schema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"name":         schema.StringAttribute{Required: true},
					"image":        schema.StringAttribute{Required: true},
					"llm_provider": schema.StringAttribute{Required: true},
					"llm_model":    schema.StringAttribute{Required: true},
				},
			},
		},
	}
	resp := &resource.CreateResponse{}

	// This should fail because provider data (Docker client) is missing
	res.Create(ctx, req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected error due to missing provider data, but got none")
	}
}

func TestAgentResourceDelete(t *testing.T) {
	ctx := context.Background()
	res := NewAgentResource()

	stateValue := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":         tftypes.String,
				"image":        tftypes.String,
				"llm_provider": tftypes.String,
				"llm_model":    tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"name":         tftypes.NewValue(tftypes.String, "test-agent"),
			"image":        tftypes.NewValue(tftypes.String, "alpine:latest"),
			"llm_provider": tftypes.NewValue(tftypes.String, "terraform-provider-mlops"),
			"llm_model":    tftypes.NewValue(tftypes.String, "llama3-8b-8192"),
		},
	)

	req := resource.DeleteRequest{
		State: tfsdk.State{
			Raw: stateValue,
			Schema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"name":         schema.StringAttribute{Required: true},
					"image":        schema.StringAttribute{Required: true},
					"llm_provider": schema.StringAttribute{Required: true},
					"llm_model":    schema.StringAttribute{Required: true},
				},
			},
		},
	}
	resp := &resource.DeleteResponse{}

	// This should fail because provider data (Docker client) is missing
	res.Delete(ctx, req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Expected error due to missing provider data, but got none")
	}
}
