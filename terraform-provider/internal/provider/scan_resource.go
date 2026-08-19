package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface compliance.
var (
	_ resource.Resource              = &scanResource{}
	_ resource.ResourceWithConfigure = &scanResource{}
)

type scanResource struct{ client *Client }

func NewScanResource() resource.Resource { return &scanResource{} }

type scanModel struct {
	Target         types.String `tfsdk:"target"`
	Effort         types.String `tfsdk:"effort"`
	FailOnSeverity types.String `tfsdk:"fail_on_severity"`
	ID             types.String `tfsdk:"id"`
	FindingsCount  types.Int64  `tfsdk:"findings_count"`
	HighCount      types.Int64  `tfsdk:"high_count"`
	MediumCount    types.Int64  `tfsdk:"medium_count"`
	SARIF          types.String `tfsdk:"sarif"`
	ReportURL      types.String `tfsdk:"report_url"`
}

func (r *scanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scan"
}

func (r *scanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runs an LLM pentest scan against `target`. The apply **fails** when a finding at " +
			"or above `fail_on_severity` is present — use it to gate a deploy on a clean scan.",
		Attributes: map[string]schema.Attribute{
			"target": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL/host to scan. Only scan systems you are authorized to test.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"effort": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "`low` | `medium` | `high` | `xhigh` | `max`. Default `high`.",
			},
			"fail_on_severity": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "`none` | `medium` | `high`. Default `high`.",
			},
			"id":             schema.StringAttribute{Computed: true, MarkdownDescription: "Scan ID."},
			"findings_count": schema.Int64Attribute{Computed: true},
			"high_count":     schema.Int64Attribute{Computed: true},
			"medium_count":   schema.Int64Attribute{Computed: true},
			"sarif":          schema.StringAttribute{Computed: true, MarkdownDescription: "SARIF 2.1.0 report JSON."},
			"report_url":     schema.StringAttribute{Computed: true},
		},
	}
}

func (r *scanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *scanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var m scanModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &m)...)
	if resp.Diagnostics.HasError() {
		return
	}

	effort := valueOr(m.Effort, "high")
	failOn := valueOr(m.FailOnSeverity, "high")

	started, err := r.client.StartScan(ctx, m.Target.ValueString(), effort)
	if err != nil {
		resp.Diagnostics.AddError("Scan failed to start", err.Error())
		return
	}
	result, err := r.client.WaitForScan(ctx, started.ID)
	if err != nil {
		resp.Diagnostics.AddError("Scan failed", err.Error())
		return
	}

	applyResult(&m, result, effort, failOn)
	// Save state first so the report stays inspectable even if the gate trips.
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
	if gate := severityGate(result, failOn); gate != "" {
		resp.Diagnostics.AddError("Scan gate failed", gate)
	}
}

func (r *scanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Scans are point-in-time; retain prior state.
	var m scanModel
	resp.Diagnostics.Append(req.State.Get(ctx, &m)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

func (r *scanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// `target` forces replacement; only effort/fail_on can change in place.
	var m scanModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &m)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}

func (r *scanResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// A scan is not infrastructure — nothing to destroy.
}

func valueOr(s types.String, def string) string {
	if s.IsNull() || s.IsUnknown() || s.ValueString() == "" {
		return def
	}
	return s.ValueString()
}

func counts(res *ScanResult) (high, medium int64) {
	for _, f := range res.Findings {
		switch f.Severity {
		case "error":
			high++
		case "warning":
			medium++
		}
	}
	return
}

func applyResult(m *scanModel, res *ScanResult, effort, failOn string) {
	high, med := counts(res)
	m.ID = types.StringValue(res.ID)
	m.Effort = types.StringValue(effort)
	m.FailOnSeverity = types.StringValue(failOn)
	m.FindingsCount = types.Int64Value(int64(len(res.Findings)))
	m.HighCount = types.Int64Value(high)
	m.MediumCount = types.Int64Value(med)
	m.SARIF = types.StringValue(res.SARIF)
	m.ReportURL = types.StringValue(res.ReportURL)
}

func severityGate(res *ScanResult, failOn string) string {
	high, med := counts(res)
	switch failOn {
	case "high":
		if high > 0 {
			return fmt.Sprintf("%d high-severity finding(s). Report: %s", high, res.ReportURL)
		}
	case "medium":
		if high+med > 0 {
			return fmt.Sprintf("%d finding(s) at or above medium. Report: %s", high+med, res.ReportURL)
		}
	}
	return ""
}
