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
	ProjectIntegrationMatcherModel struct {
		ID            types.Int64  `tfsdk:"id"`
		ProjectID     types.Int64  `tfsdk:"project_id"`
		IntegrationID types.Int64  `tfsdk:"integration_id"`
		Name          types.String `tfsdk:"name"`
		MatchType     types.String `tfsdk:"match_type"`
		Method        types.String `tfsdk:"method"`
		BodyDataType  types.String `tfsdk:"body_data_type"`
		Key           types.String `tfsdk:"key"`
		Value         types.String `tfsdk:"value"`
	}
)

func ProjectIntegrationMatcherSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The project integration matcher",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource allows you to manage integration matchers. Matchers define conditions that must be met for an integration webhook to trigger a template execution.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source allows you to read an integration matcher.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "The matcher ID.",
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
					MarkdownDescription: "The integration ID that this matcher belongs to.",
					Required:            true,
				},
				Resource: &schemaR.Int64Attribute{
					PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The display name of the matcher.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"match_type": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Where to look for the match. Valid values are `body` or `header`.",
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
			"method": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The comparison method. Valid values are `equals`, `unequals`, or `contains`.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("equals", "unequals", "contains"),
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
					MarkdownDescription: "The key to match against in the body or header.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"value": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The value to compare against.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
