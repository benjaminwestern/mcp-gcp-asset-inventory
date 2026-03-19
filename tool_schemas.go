package main

import "github.com/mark3labs/mcp-go/mcp"

func newListAssetsTool() mcp.Tool {
	return mcp.NewTool(
		"list_gcp_assets",
		mcp.WithDescription(
			"Lists Cloud Asset Inventory assets for a project, folder, or "+
				"organization. Use either `project_id` or `parent`.",
		),
		mcp.WithString(
			"project_id",
			mcp.Description(
				"Optional project ID shortcut. If `parent` is omitted, the tool "+
					"uses `projects/{project_id}`.",
			),
		),
		mcp.WithString(
			"parent",
			mcp.Description(
				"Optional asset scope in the form `projects/{id}`, "+
					"`folders/{number}`, or `organizations/{number}`.",
			),
		),
		mcp.WithArray(
			"asset_types",
			mcp.Description(
				"Optional asset type filters or RE2 patterns such as "+
					"`storage.googleapis.com/Bucket` or `compute.googleapis.com.*`.",
			),
			mcp.WithStringItems(),
		),
		mcp.WithString(
			"content_type",
			mcp.Description(
				"Optional asset content type. Defaults to `RESOURCE` to preserve "+
					"the prior behavior of this tool.",
			),
			mcp.Enum(contentTypeEnumValues...),
		),
		mcp.WithString(
			"read_time",
			mcp.Description(
				"Optional RFC3339 timestamp for the asset snapshot, for example "+
					rfc3339Example+".",
			),
		),
		mcp.WithArray(
			"relationship_types",
			mcp.Description(
				"Optional relationship types. Only valid when `content_type` is "+
					"`RELATIONSHIP`.",
			),
			mcp.WithStringItems(),
		),
		mcp.WithNumber(
			"page_size",
			mcp.Description(
				"Optional page size. Defaults to 100 and is capped at 1000.",
			),
			mcp.Min(1),
			mcp.Max(maxListPageSize),
		),
		mcp.WithString(
			"page_token",
			mcp.Description("Optional page token from a previous response."),
		),
	)
}

func newSearchResourcesTool() mcp.Tool {
	return mcp.NewTool(
		"search_gcp_resources",
		mcp.WithDescription(
			"Searches resources with Cloud Asset Inventory "+
				"`searchAllResources`.",
		),
		mcp.WithString(
			"scope",
			mcp.Required(),
			mcp.Description(
				"Resource Manager scope such as `projects/{id}`, "+
					"`folders/{number}`, or `organizations/{number}`.",
			),
		),
		mcp.WithString(
			"query",
			mcp.Description(
				"Optional Cloud Asset search query. If omitted, the tool returns "+
					"all searchable resources in scope.",
			),
		),
		mcp.WithArray(
			"asset_types",
			mcp.Description("Optional asset type filters or RE2 patterns."),
			mcp.WithStringItems(),
		),
		mcp.WithNumber(
			"page_size",
			mcp.Description(
				"Optional page size. Defaults to 100 and is capped at 500.",
			),
			mcp.Min(1),
			mcp.Max(maxSearchPageSize),
		),
		mcp.WithString(
			"page_token",
			mcp.Description("Optional page token from a previous response."),
		),
		mcp.WithString(
			"order_by",
			mcp.Description(
				"Optional sort expression such as `location DESC, name`.",
			),
		),
		mcp.WithArray(
			"read_mask",
			mcp.Description(
				"Optional list of fields to return, for example `name`, "+
					"`location`, or `versionedResources`. Use `*` for all fields.",
			),
			mcp.WithStringItems(),
		),
	)
}

func newSearchIAMPoliciesTool() mcp.Tool {
	return mcp.NewTool(
		"search_gcp_iam_policies",
		mcp.WithDescription(
			"Searches IAM policies with Cloud Asset Inventory "+
				"`searchAllIamPolicies`.",
		),
		mcp.WithString(
			"scope",
			mcp.Required(),
			mcp.Description(
				"Resource Manager scope such as `projects/{id}`, "+
					"`folders/{number}`, or `organizations/{number}`.",
			),
		),
		mcp.WithString(
			"query",
			mcp.Description(
				"Optional IAM policy search query. If omitted, the tool returns "+
					"all searchable IAM bindings in scope.",
			),
		),
		mcp.WithArray(
			"asset_types",
			mcp.Description("Optional attached asset type filters or RE2 patterns."),
			mcp.WithStringItems(),
		),
		mcp.WithNumber(
			"page_size",
			mcp.Description(
				"Optional page size. Defaults to 100 and is capped at 500.",
			),
			mcp.Min(1),
			mcp.Max(maxSearchPageSize),
		),
		mcp.WithString(
			"page_token",
			mcp.Description("Optional page token from a previous response."),
		),
		mcp.WithString(
			"order_by",
			mcp.Description(
				"Optional sort expression such as `assetType DESC, resource`.",
			),
		),
	)
}

func newAssetHistoryTool() mcp.Tool {
	return mcp.NewTool(
		"get_gcp_asset_history",
		mcp.WithDescription(
			"Gets historical asset snapshots with Cloud Asset Inventory "+
				"`batchGetAssetsHistory`.",
		),
		mcp.WithString(
			"parent",
			mcp.Required(),
			mcp.Description(
				"Parent scope in the form `projects/{id}`, `projects/{number}`, "+
					"or `organizations/{number}`.",
			),
		),
		mcp.WithArray(
			"asset_names",
			mcp.Required(),
			mcp.Description(
				"Full resource names to inspect. A single request supports up to "+
					"100 assets.",
			),
			mcp.WithStringItems(),
		),
		mcp.WithString(
			"content_type",
			mcp.Description(
				"Optional history content type. Defaults to `RESOURCE`.",
			),
			mcp.Enum(contentTypeEnumValues...),
		),
		mcp.WithString(
			"start_time",
			mcp.Description(
				"Optional RFC3339 window start time, for example "+
					rfc3339Example+".",
			),
		),
		mcp.WithString(
			"end_time",
			mcp.Description(
				"Optional RFC3339 window end time, for example "+
					rfc3339Example+". If omitted, the API uses the current time.",
			),
		),
		mcp.WithArray(
			"relationship_types",
			mcp.Description(
				"Optional relationship types. Only valid when `content_type` is "+
					"`RELATIONSHIP`.",
			),
			mcp.WithStringItems(),
		),
	)
}

func newQueryAssetsTool() mcp.Tool {
	return mcp.NewTool(
		"query_gcp_assets",
		mcp.WithDescription(
			"Runs Cloud Asset Inventory SQL queries with `queryAssets`. Provide "+
				"exactly one of `statement` or `job_reference`.",
		),
		mcp.WithString(
			"parent",
			mcp.Required(),
			mcp.Description(
				"Parent scope in the form `projects/{id}`, `folders/{number}`, "+
					"or `organizations/{number}`.",
			),
		),
		mcp.WithString(
			"statement",
			mcp.Description(
				"Optional BigQuery-compatible SQL statement for a new query job.",
			),
		),
		mcp.WithString(
			"job_reference",
			mcp.Description(
				"Optional job reference from a previous `query_gcp_assets` call.",
			),
		),
		mcp.WithNumber(
			"page_size",
			mcp.Description(
				"Optional maximum number of rows to return. The API caps results "+
					"at 1000 rows per page.",
			),
			mcp.Min(1),
			mcp.Max(1000),
		),
		mcp.WithString(
			"page_token",
			mcp.Description("Optional page token from a previous response."),
		),
		mcp.WithString(
			"timeout",
			mcp.Description(
				"Optional Go duration such as `30s` or `2m`.",
			),
		),
		mcp.WithString(
			"read_time",
			mcp.Description(
				"Optional RFC3339 point-in-time query timestamp, for example "+
					rfc3339Example+".",
			),
		),
		mcp.WithString(
			"start_time",
			mcp.Description(
				"Optional RFC3339 query window start time. When set, the query "+
					"uses a time window instead of `read_time`.",
			),
		),
		mcp.WithString(
			"end_time",
			mcp.Description(
				"Optional RFC3339 query window end time. If omitted, the API uses "+
					"the current time.",
			),
		),
	)
}

func newEffectiveIAMPoliciesTool() mcp.Tool {
	return mcp.NewTool(
		"get_gcp_effective_iam_policies",
		mcp.WithDescription(
			"Gets effective IAM policies for full resource names with "+
				"`effectiveIamPolicies.batchGet`.",
		),
		mcp.WithString(
			"scope",
			mcp.Required(),
			mcp.Description(
				"Resource Manager scope such as `projects/{id}`, "+
					"`folders/{number}`, or `organizations/{number}`.",
			),
		),
		mcp.WithArray(
			"names",
			mcp.Required(),
			mcp.Description(
				"Full resource names to inspect. A single request supports up to "+
					"20 names.",
			),
			mcp.WithStringItems(),
		),
	)
}
