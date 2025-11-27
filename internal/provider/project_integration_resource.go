package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiclient "terraform-provider-semaphoreui/semaphoreui/client"
	"terraform-provider-semaphoreui/semaphoreui/client/project"
	"terraform-provider-semaphoreui/semaphoreui/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectIntegrationResource{}
	_ resource.ResourceWithConfigure   = &projectIntegrationResource{}
	_ resource.ResourceWithImportState = &projectIntegrationResource{}
)

func NewProjectIntegrationResource() resource.Resource {
	return &projectIntegrationResource{}
}

type projectIntegrationResource struct {
	client *apiclient.SemaphoreUI
}

func (r *projectIntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *projectIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_integration"
}

func (r *projectIntegrationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProjectIntegrationSchema().GetResource(ctx)
}

func convertProjectIntegrationModelToIntegrationRequest(integration ProjectIntegrationModel) *models.IntegrationRequest {
	return &models.IntegrationRequest{
		ProjectID:  integration.ProjectID.ValueInt64(),
		Name:       integration.Name.ValueString(),
		TemplateID: integration.TemplateID.ValueInt64(),
	}
}

func convertIntegrationResponseToProjectIntegrationModel(response *models.Integration) ProjectIntegrationModel {
	return ProjectIntegrationModel{
		ID:         types.Int64Value(response.ID),
		ProjectID:  types.Int64Value(response.ProjectID),
		Name:       types.StringValue(response.Name),
		TemplateID: types.Int64Value(response.TemplateID),
	}
}

// getIntegrationByID retrieves an integration by ID from the list of integrations.
func getIntegrationByID(client *apiclient.SemaphoreUI, projectID int64, integrationID int64) (*ProjectIntegrationModel, error) {
	response, err := client.Project.GetProjectProjectIDIntegrations(&project.GetProjectProjectIDIntegrationsParams{
		ProjectID: projectID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read project integrations: %s", err.Error())
	}
	for _, integration := range response.Payload {
		if integration.ID == integrationID {
			model := convertIntegrationResponseToProjectIntegrationModel(integration)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("project integration with ID %d not found", integrationID)
}

// GetIntegrationByName retrieves an integration by name from the list of integrations.
func GetIntegrationByName(client *apiclient.SemaphoreUI, projectID int64, name string) (*ProjectIntegrationModel, error) {
	response, err := client.Project.GetProjectProjectIDIntegrations(&project.GetProjectProjectIDIntegrationsParams{
		ProjectID: projectID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read project integrations: %s", err.Error())
	}
	for _, integration := range response.Payload {
		if integration.Name == name {
			model := convertIntegrationResponseToProjectIntegrationModel(integration)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("project integration with name %s not found", name)
}

func (r *projectIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ProjectIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.Project.PostProjectProjectIDIntegrations(&project.PostProjectProjectIDIntegrationsParams{
		ProjectID:   plan.ProjectID.ValueInt64(),
		Integration: convertProjectIntegrationModelToIntegrationRequest(plan),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SemaphoreUI Project Integration",
			"Could not create project integration, unexpected error: "+err.Error(),
		)
		return
	}
	model := convertIntegrationResponseToProjectIntegrationModel(response.Payload)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *projectIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ProjectIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := getIntegrationByID(r.client, state.ProjectID.ValueInt64(), state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Project Integration",
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
func (r *projectIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ProjectIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Project.PutProjectProjectIDIntegrationsIntegrationID(&project.PutProjectProjectIDIntegrationsIntegrationIDParams{
		ProjectID:     plan.ProjectID.ValueInt64(),
		IntegrationID: plan.ID.ValueInt64(),
		Integration:   convertProjectIntegrationModelToIntegrationRequest(plan),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SemaphoreUI Project Integration",
			"Could not update project integration, unexpected error: "+err.Error(),
		)
		return
	}

	model, err := getIntegrationByID(r.client, plan.ProjectID.ValueInt64(), plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Project Integration",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ProjectIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Project.DeleteProjectProjectIDIntegrationsIntegrationID(&project.DeleteProjectProjectIDIntegrationsIntegrationIDParams{
		ProjectID:     state.ProjectID.ValueInt64(),
		IntegrationID: state.ID.ValueInt64(),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing SemaphoreUI Project Integration",
			"Could not remove project integration, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *projectIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fields, err := parseImportFields(req.ID, []string{"project", "integration"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Project Integration Import ID",
			"Could not parse import ID: "+err.Error(),
		)
		return
	}

	model, err := getIntegrationByID(r.client, fields["project"], fields["integration"])
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Project Integration",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
