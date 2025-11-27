package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiclient "terraform-provider-semaphoreui/semaphoreui/client"
	"terraform-provider-semaphoreui/semaphoreui/client/integration"
	"terraform-provider-semaphoreui/semaphoreui/client/project"
	"terraform-provider-semaphoreui/semaphoreui/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectIntegrationExtractValueResource{}
	_ resource.ResourceWithConfigure   = &projectIntegrationExtractValueResource{}
	_ resource.ResourceWithImportState = &projectIntegrationExtractValueResource{}
)

func NewProjectIntegrationExtractValueResource() resource.Resource {
	return &projectIntegrationExtractValueResource{}
}

type projectIntegrationExtractValueResource struct {
	client *apiclient.SemaphoreUI
}

func (r *projectIntegrationExtractValueResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = client
}

func (r *projectIntegrationExtractValueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_integration_extract_value"
}

func (r *projectIntegrationExtractValueResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProjectIntegrationExtractValueSchema().GetResource(ctx)
}

func convertProjectIntegrationExtractValueModelToExtractValue(ev ProjectIntegrationExtractValueModel) *models.IntegrationExtractValue {
	return &models.IntegrationExtractValue{
		Name:         ev.Name.ValueString(),
		ValueSource:  ev.ValueSource.ValueString(),
		BodyDataType: ev.BodyDataType.ValueString(),
		Key:          ev.Key.ValueString(),
		Variable:     ev.Variable.ValueString(),
		VariableType: ev.VariableType.ValueString(),
	}
}

func convertProjectIntegrationExtractValueModelToRequest(ev ProjectIntegrationExtractValueModel) *models.IntegrationExtractValueRequest {
	return &models.IntegrationExtractValueRequest{
		Name:         ev.Name.ValueString(),
		ValueSource:  ev.ValueSource.ValueString(),
		BodyDataType: ev.BodyDataType.ValueString(),
		Key:          ev.Key.ValueString(),
		Variable:     ev.Variable.ValueString(),
		VariableType: ev.VariableType.ValueString(),
	}
}

func convertIntegrationExtractValueResponseToModel(response *models.IntegrationExtractValue, projectID int64) ProjectIntegrationExtractValueModel {
	return ProjectIntegrationExtractValueModel{
		ID:            types.Int64Value(response.ID),
		ProjectID:     types.Int64Value(projectID),
		IntegrationID: types.Int64Value(response.IntegrationID),
		Name:          types.StringValue(response.Name),
		ValueSource:   types.StringValue(response.ValueSource),
		BodyDataType:  types.StringValue(response.BodyDataType),
		Key:           types.StringValue(response.Key),
		Variable:      types.StringValue(response.Variable),
		VariableType:  types.StringValue(response.VariableType),
	}
}

// getExtractValueByID retrieves an extract value by ID from the list of extract values.
func getExtractValueByID(client *apiclient.SemaphoreUI, projectID int64, integrationID int64, extractValueID int64) (*ProjectIntegrationExtractValueModel, error) {
	response, err := client.Integration.GetProjectProjectIDIntegrationsIntegrationIDValues(&integration.GetProjectProjectIDIntegrationsIntegrationIDValuesParams{
		ProjectID:     projectID,
		IntegrationID: integrationID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read integration extract values: %s", err.Error())
	}
	for _, ev := range response.Payload {
		if ev.ID == extractValueID {
			model := convertIntegrationExtractValueResponseToModel(ev, projectID)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("integration extract value with ID %d not found", extractValueID)
}

// getExtractValueByName retrieves an extract value by name from the list of extract values.
func getExtractValueByName(client *apiclient.SemaphoreUI, projectID int64, integrationID int64, name string) (*ProjectIntegrationExtractValueModel, error) {
	response, err := client.Integration.GetProjectProjectIDIntegrationsIntegrationIDValues(&integration.GetProjectProjectIDIntegrationsIntegrationIDValuesParams{
		ProjectID:     projectID,
		IntegrationID: integrationID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read integration extract values: %s", err.Error())
	}
	for _, ev := range response.Payload {
		if ev.Name == name {
			model := convertIntegrationExtractValueResponseToModel(ev, projectID)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("integration extract value with name %s not found", name)
}

func (r *projectIntegrationExtractValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ProjectIntegrationExtractValueModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Project.PostProjectProjectIDIntegrationsIntegrationIDValues(&project.PostProjectProjectIDIntegrationsIntegrationIDValuesParams{
		ProjectID:                 plan.ProjectID.ValueInt64(),
		IntegrationID:             plan.IntegrationID.ValueInt64(),
		IntegrationExtractedValue: convertProjectIntegrationExtractValueModelToExtractValue(plan),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SemaphoreUI Integration Extract Value",
			"Could not create integration extract value, unexpected error: "+err.Error(),
		)
		return
	}

	// POST doesn't return the created object, so we need to read it back by name
	model, err := getExtractValueByName(r.client, plan.ProjectID.ValueInt64(), plan.IntegrationID.ValueInt64(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Extract Value",
			"Created integration extract value but could not read it back: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *projectIntegrationExtractValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ProjectIntegrationExtractValueModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := getExtractValueByID(r.client, state.ProjectID.ValueInt64(), state.IntegrationID.ValueInt64(), state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Extract Value",
			err.Error(),
		)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *projectIntegrationExtractValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ProjectIntegrationExtractValueModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Integration.PutProjectProjectIDIntegrationsIntegrationIDValuesExtractvalueID(&integration.PutProjectProjectIDIntegrationsIntegrationIDValuesExtractvalueIDParams{
		ProjectID:               plan.ProjectID.ValueInt64(),
		IntegrationID:           plan.IntegrationID.ValueInt64(),
		ExtractvalueID:          plan.ID.ValueInt64(),
		IntegrationExtractValue: convertProjectIntegrationExtractValueModelToRequest(plan),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SemaphoreUI Integration Extract Value",
			"Could not update integration extract value, unexpected error: "+err.Error(),
		)
		return
	}

	model, err := getExtractValueByID(r.client, plan.ProjectID.ValueInt64(), plan.IntegrationID.ValueInt64(), plan.ID.ValueInt64())
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

func (r *projectIntegrationExtractValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ProjectIntegrationExtractValueModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Integration.DeleteProjectProjectIDIntegrationsIntegrationIDValuesExtractvalueID(&integration.DeleteProjectProjectIDIntegrationsIntegrationIDValuesExtractvalueIDParams{
		ProjectID:      state.ProjectID.ValueInt64(),
		IntegrationID:  state.IntegrationID.ValueInt64(),
		ExtractvalueID: state.ID.ValueInt64(),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing SemaphoreUI Integration Extract Value",
			"Could not remove integration extract value, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *projectIntegrationExtractValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fields, err := parseImportFields(req.ID, []string{"project", "integration", "extractvalue"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Integration Extract Value Import ID",
			"Could not parse import ID: "+err.Error(),
		)
		return
	}

	model, err := getExtractValueByID(r.client, fields["project"], fields["integration"], fields["extractvalue"])
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
