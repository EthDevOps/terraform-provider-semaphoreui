package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
)

type (
	ProjectIntegrationExtractValueModel struct {
		ID            types.Int64  `tfsdk:"id"`
		ProjectID     types.Int64  `tfsdk:"project_id"`
		IntegrationID types.Int64  `tfsdk:"integration_id"`
		Name          types.String `tfsdk:"name"`
		ValueSource   types.String `tfsdk:"value_source"`
		BodyDataType  types.String `tfsdk:"body_data_type"`
		Key           types.String `tfsdk:"key"`
		Variable      types.String `tfsdk:"variable"`
		VariableType  types.String `tfsdk:"variable_type"`
	}
)

func ProjectIntegrationExtractValueSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The project integration extract value",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to manage integration extract values. Extract values define how to extract data from webhook payloads and pass them as variables to template executions.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to read an integration extract value.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The extract value ID.",
				},
				Resource: &schemaR.Int64Attribute{
					Computed:      true,
					PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
				},
				DataSource: &schemaD.Int64Attribute{
					Required: true,
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
			"integration_id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The integration ID that this extract value belongs to.",
					Required:            true,
				},
				Resource: &schemaR.Int64Attribute{
					PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The display name of the extract value.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"value_source": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Where to extract the value from. Valid values are `body` or `header`.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("body", "header"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"body_data_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The data type of the body. Valid values are `json`, `xml`, or `string`.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("json", "xml", "string"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"key": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The key to extract from the body or header.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"variable": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The variable name to store the extracted value.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"variable_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of variable to set. Valid values are `environment` or `task`.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("environment", "task"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
