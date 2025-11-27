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
	_ resource.Resource                = &projectIntegrationMatcherResource{}
	_ resource.ResourceWithConfigure   = &projectIntegrationMatcherResource{}
	_ resource.ResourceWithImportState = &projectIntegrationMatcherResource{}
)

func NewProjectIntegrationMatcherResource() resource.Resource {
	return &projectIntegrationMatcherResource{}
}

type projectIntegrationMatcherResource struct {
	client *apiclient.SemaphoreUI
}

func (r *projectIntegrationMatcherResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *projectIntegrationMatcherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_integration_matcher"
}

func (r *projectIntegrationMatcherResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProjectIntegrationMatcherSchema().GetResource(ctx)
}

func convertProjectIntegrationMatcherModelToMatcher(matcher ProjectIntegrationMatcherModel) *models.IntegrationMatcher {
	return &models.IntegrationMatcher{
		Name:         matcher.Name.ValueString(),
		MatchType:    matcher.MatchType.ValueString(),
		Method:       matcher.Method.ValueString(),
		BodyDataType: matcher.BodyDataType.ValueString(),
		Key:          matcher.Key.ValueString(),
		Value:        matcher.Value.ValueString(),
	}
}

func convertProjectIntegrationMatcherModelToMatcherRequest(matcher ProjectIntegrationMatcherModel) *models.IntegrationMatcherRequest {
	return &models.IntegrationMatcherRequest{
		Name:         matcher.Name.ValueString(),
		MatchType:    matcher.MatchType.ValueString(),
		Method:       matcher.Method.ValueString(),
		BodyDataType: matcher.BodyDataType.ValueString(),
		Key:          matcher.Key.ValueString(),
		Value:        matcher.Value.ValueString(),
	}
}

func convertIntegrationMatcherResponseToModel(response *models.IntegrationMatcher, projectID int64) ProjectIntegrationMatcherModel {
	return ProjectIntegrationMatcherModel{
		ID:            types.Int64Value(response.ID),
		ProjectID:     types.Int64Value(projectID),
		IntegrationID: types.Int64Value(response.IntegrationID),
		Name:          types.StringValue(response.Name),
		MatchType:     types.StringValue(response.MatchType),
		Method:        types.StringValue(response.Method),
		BodyDataType:  types.StringValue(response.BodyDataType),
		Key:           types.StringValue(response.Key),
		Value:         types.StringValue(response.Value),
	}
}

// getMatcherByID retrieves a matcher by ID from the list of matchers.
func getMatcherByID(client *apiclient.SemaphoreUI, projectID int64, integrationID int64, matcherID int64) (*ProjectIntegrationMatcherModel, error) {
	response, err := client.Integration.GetProjectProjectIDIntegrationsIntegrationIDMatchers(&integration.GetProjectProjectIDIntegrationsIntegrationIDMatchersParams{
		ProjectID:     projectID,
		IntegrationID: integrationID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read integration matchers: %s", err.Error())
	}
	for _, matcher := range response.Payload {
		if matcher.ID == matcherID {
			model := convertIntegrationMatcherResponseToModel(matcher, projectID)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("integration matcher with ID %d not found", matcherID)
}

// getMatcherByName retrieves a matcher by name from the list of matchers.
func getMatcherByName(client *apiclient.SemaphoreUI, projectID int64, integrationID int64, name string) (*ProjectIntegrationMatcherModel, error) {
	response, err := client.Integration.GetProjectProjectIDIntegrationsIntegrationIDMatchers(&integration.GetProjectProjectIDIntegrationsIntegrationIDMatchersParams{
		ProjectID:     projectID,
		IntegrationID: integrationID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read integration matchers: %s", err.Error())
	}
	for _, matcher := range response.Payload {
		if matcher.Name == name {
			model := convertIntegrationMatcherResponseToModel(matcher, projectID)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("integration matcher with name %s not found", name)
}

func (r *projectIntegrationMatcherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ProjectIntegrationMatcherModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Project.PostProjectProjectIDIntegrationsIntegrationIDMatchers(&project.PostProjectProjectIDIntegrationsIntegrationIDMatchersParams{
		ProjectID:          plan.ProjectID.ValueInt64(),
		IntegrationID:      plan.IntegrationID.ValueInt64(),
		IntegrationMatcher: convertProjectIntegrationMatcherModelToMatcher(plan),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SemaphoreUI Integration Matcher",
			"Could not create integration matcher, unexpected error: "+err.Error(),
		)
		return
	}

	// POST doesn't return the created object, so we need to read it back by name
	model, err := getMatcherByName(r.client, plan.ProjectID.ValueInt64(), plan.IntegrationID.ValueInt64(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Matcher",
			"Created integration matcher but could not read it back: "+err.Error(),
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
func (r *projectIntegrationMatcherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ProjectIntegrationMatcherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := getMatcherByID(r.client, state.ProjectID.ValueInt64(), state.IntegrationID.ValueInt64(), state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Matcher",
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
func (r *projectIntegrationMatcherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ProjectIntegrationMatcherModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Integration.PutProjectProjectIDIntegrationsIntegrationIDMatchersMatcherID(&integration.PutProjectProjectIDIntegrationsIntegrationIDMatchersMatcherIDParams{
		ProjectID:          plan.ProjectID.ValueInt64(),
		IntegrationID:      plan.IntegrationID.ValueInt64(),
		MatcherID:          plan.ID.ValueInt64(),
		IntegrationMatcher: convertProjectIntegrationMatcherModelToMatcherRequest(plan),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SemaphoreUI Integration Matcher",
			"Could not update integration matcher, unexpected error: "+err.Error(),
		)
		return
	}

	model, err := getMatcherByID(r.client, plan.ProjectID.ValueInt64(), plan.IntegrationID.ValueInt64(), plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Matcher",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectIntegrationMatcherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ProjectIntegrationMatcherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Integration.DeleteProjectProjectIDIntegrationsIntegrationIDMatchersMatcherID(&integration.DeleteProjectProjectIDIntegrationsIntegrationIDMatchersMatcherIDParams{
		ProjectID:     state.ProjectID.ValueInt64(),
		IntegrationID: state.IntegrationID.ValueInt64(),
		MatcherID:     state.ID.ValueInt64(),
	}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing SemaphoreUI Integration Matcher",
			"Could not remove integration matcher, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *projectIntegrationMatcherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	fields, err := parseImportFields(req.ID, []string{"project", "integration", "matcher"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Integration Matcher Import ID",
			"Could not parse import ID: "+err.Error(),
		)
		return
	}

	model, err := getMatcherByID(r.client, fields["project"], fields["integration"], fields["matcher"])
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SemaphoreUI Integration Matcher",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
