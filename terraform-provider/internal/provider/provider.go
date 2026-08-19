package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type disassemblyProvider struct{ version string }

// New returns the provider factory used by providerserver.Serve.
func New(version string) func() provider.Provider {
	return func() provider.Provider { return &disassemblyProvider{version: version} }
}

func (p *disassemblyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "disassembly"
	resp.Version = p.version
}

type providerModel struct {
	APIKey   types.String `tfsdk:"api_key"`
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *disassemblyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configure access to the Disassembly.AI pentest API.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "API key. Falls back to the `DISASSEMBLY_API_KEY` environment variable.",
			},
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "API base URL. Defaults to `https://api.disassembly.ai`.",
			},
		},
	}
}

func (p *disassemblyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("DISASSEMBLY_API_KEY")
	if !cfg.APIKey.IsNull() && cfg.APIKey.ValueString() != "" {
		apiKey = cfg.APIKey.ValueString()
	}
	endpoint := "https://api.disassembly.ai"
	if !cfg.Endpoint.IsNull() && cfg.Endpoint.ValueString() != "" {
		endpoint = cfg.Endpoint.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API key",
			"Set `api_key` on the provider or the DISASSEMBLY_API_KEY environment variable.",
		)
		return
	}

	client := &Client{HTTP: http.DefaultClient, Endpoint: endpoint, APIKey: apiKey}
	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *disassemblyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewScanResource}
}

func (p *disassemblyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
