package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	apiclient "terraform-provider-semaphoreui/semaphoreui/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &projectIntegrationExtractValueDataSource{}
)

func NewProjectIntegrationExtractValueDataSource() datasource.DataSource {
	return &projectIntegrationExtractValueDataSource{}
}

type projectIntegrationExtractValueDataSource struct {
	client *apiclient.SemaphoreUI
}

func (d *projectIntegrationExtractValueDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.SemaphoreUI)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.SemaphoreUI, got %T. Please report this issue to the provider developers.",
		)
		return
	}
	d.client = client
}

// Metadata returns the data source type name.
func (d *projectIntegrationExtractValueDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_integration_extract_value"
}

// Schema defines the schema for the data source.
func (d *projectIntegrationExtractValueDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProjectIntegrationExtractValueSchema().GetDataSource(ctx)
}

func (d *projectIntegrationExtractValueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ProjectIntegrationExtractValueModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := getExtractValueByID(d.client, config.ProjectID.ValueInt64(), config.IntegrationID.ValueInt64(), config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Extract Value",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
