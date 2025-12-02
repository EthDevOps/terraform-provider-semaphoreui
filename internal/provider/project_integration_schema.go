package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
)

type (
	ProjectIntegrationModel struct {
		ID           types.Int64  `tfsdk:"id"`
		ProjectID    types.Int64  `tfsdk:"project_id"`
		Name         types.String `tfsdk:"name"`
		TemplateID   types.Int64  `tfsdk:"template_id"`
		Searchable   types.Bool   `tfsdk:"searchable"`
		AuthMethod   types.String `tfsdk:"auth_method"`
		AuthSecretID types.Int64  `tfsdk:"auth_secret_id"`
		AuthHeader   types.String `tfsdk:"auth_header"`
	}
)

func ProjectIntegrationSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The project integration",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to manage integrations (webhooks) for a project. Integrations enable external systems to trigger template executions via HTTP requests.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to read an integration.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The integration ID.",
				},
				Resource: &schemaR.Int64Attribute{
					Computed:      true,
					PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				},
				DataSource: &schemaD.Int64Attribute{
					Optional: true,
					Computed: true,
					Validators: []validator.Int64{
						int64validator.ExactlyOneOf(
							path.MatchRoot("id"),
							path.MatchRoot("name"),
						),
					},
				},
			},
			"project_id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The project ID that the integration belongs to.",
					Required:            true,
				},
				Resource: &schemaR.Int64Attribute{
					PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The display name of the integration.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(
							path.MatchRoot("id"),
							path.MatchRoot("name"),
						),
					},
				},
			},
			"template_id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The template ID that this integration will trigger.",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"searchable": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "When enabled, the integration uses matchers to route incoming webhooks via the project alias. When disabled, the integration has its own dedicated alias endpoint.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"auth_method": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The authentication method for the integration webhook. Valid values are `token`, `github`, `bitbucket`, `hmac`, `basic`. When not set, no authentication is required.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("token", "github", "bitbucket", "hmac", "basic"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"auth_secret_id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The ID of the project key containing the secret used for authentication. Required when `auth_method` is set.",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"auth_header": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The custom header name for authentication (e.g., `X-Webhook-Token`). Used with `token` authentication method.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
