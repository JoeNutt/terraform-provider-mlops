package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mlops": providerserver.NewProtocol6(New("test")()),
}

func TestProviderConfigure(t *testing.T) {
	ctx := context.Background()
	p := &MlopsProvider{}

	configValue := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"groq_api_key": tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"groq_api_key": tftypes.NewValue(tftypes.String, "test-key"),
		},
	)

	schema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"groq_api_key": schema.StringAttribute{
				Optional: true,
			},
		},
	}

	req := provider.ConfigureRequest{
		Config: tfsdk.Config{
			Raw:    configValue,
			Schema: schema,
		},
	}
	resp := &provider.ConfigureResponse{}

	p.Configure(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure returned errors: %v", resp.Diagnostics)
	}
}

func TestProviderSchema(t *testing.T) {
	ctx := context.Background()
	var p provider.Provider = &MlopsProvider{}
	res := &provider.SchemaResponse{}
	p.Schema(ctx, provider.SchemaRequest{}, res)

	testCases := []struct {
		name      string
		attribute string
		optional  bool
		sensitive bool
	}{
		{
			name:      "groq_api_key",
			attribute: "groq_api_key",
			optional:  true,
			sensitive: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if res.Schema.Attributes == nil {
				t.Fatalf("Schema attributes are nil")
			}

			attr, ok := res.Schema.Attributes[tc.attribute]
			if !ok {
				t.Fatalf("Attribute %s not found in schema", tc.attribute)
			}

			sAttr, ok := attr.(schema.StringAttribute)
			if !ok {
				t.Fatalf("Attribute %s is not a StringAttribute, got %T", tc.attribute, attr)
			}

			if sAttr.Optional != tc.optional {
				t.Errorf("Attribute %s: expected optional %v, got %v", tc.attribute, tc.optional, sAttr.Optional)
			}

			if sAttr.Sensitive != tc.sensitive {
				t.Errorf("Attribute %s: expected sensitive %v, got %v", tc.attribute, tc.sensitive, sAttr.Sensitive)
			}
		})
	}
}
