package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"

	"github.com/mark3labs/mcp-go/mcp"

	"google.golang.org/api/iterator"
)

func listAssetsHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	projectID := strings.TrimSpace(request.GetString("project_id", ""))
	parent := strings.TrimSpace(request.GetString("parent", ""))

	resolvedParent, err := resolveParent(projectID, parent)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	contentType, err := parseContentType(
		request.GetString("content_type", ""),
		defaultAssetContent,
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	readTime, err := parseOptionalTimestamp(
		request.GetString("read_time", ""),
		"read_time",
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	relationshipTypes := compactStrings(
		request.GetStringSlice("relationship_types", nil),
	)
	if err := validateRelationshipRequest(contentType, relationshipTypes); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pageSize := normalizePageSize(
		request.GetInt("page_size", defaultListPageSize),
		defaultListPageSize,
		maxListPageSize,
	)
	pageToken := strings.TrimSpace(request.GetString("page_token", ""))

	req := &assetpb.ListAssetsRequest{
		Parent:            resolvedParent,
		ReadTime:          readTime,
		AssetTypes:        compactStrings(request.GetStringSlice("asset_types", nil)),
		ContentType:       contentType,
		RelationshipTypes: relationshipTypes,
	}

	log.Printf("list_gcp_assets parent=%s content_type=%s", resolvedParent, contentType.String())

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Asset client: %w", err)
	}
	defer closeAssetClient(client)

	it := client.ListAssets(ctx, req)
	pager := iterator.NewPager(it, pageSize, pageToken)

	var assetsPage []*assetpb.Asset
	nextPageToken, err := pager.NextPage(&assetsPage)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}

	resp := &assetpb.ListAssetsResponse{
		Assets:        assetsPage,
		NextPageToken: nextPageToken,
	}
	if raw, ok := it.Response.(*assetpb.ListAssetsResponse); ok {
		resp.ReadTime = raw.GetReadTime()
	}

	return newProtoToolResult(resp)
}

func searchResourcesHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	scope, err := request.RequireString("scope")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pageSize := normalizePageSize(
		request.GetInt("page_size", defaultSearchPageSize),
		defaultSearchPageSize,
		maxSearchPageSize,
	)
	pageToken := strings.TrimSpace(request.GetString("page_token", ""))

	req := &assetpb.SearchAllResourcesRequest{
		Scope:      strings.TrimSpace(scope),
		Query:      strings.TrimSpace(request.GetString("query", "")),
		AssetTypes: compactStrings(request.GetStringSlice("asset_types", nil)),
		OrderBy:    strings.TrimSpace(request.GetString("order_by", "")),
		ReadMask:   buildFieldMask(request.GetStringSlice("read_mask", nil)),
	}

	log.Printf("search_gcp_resources scope=%s", req.GetScope())

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Asset client: %w", err)
	}
	defer closeAssetClient(client)

	it := client.SearchAllResources(ctx, req)
	pager := iterator.NewPager(it, pageSize, pageToken)

	var results []*assetpb.ResourceSearchResult
	nextPageToken, err := pager.NextPage(&results)
	if err != nil {
		return nil, fmt.Errorf("failed to search resources: %w", err)
	}

	resp := &assetpb.SearchAllResourcesResponse{
		Results:       results,
		NextPageToken: nextPageToken,
	}

	return newProtoToolResult(resp)
}

func searchIAMPoliciesHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	scope, err := request.RequireString("scope")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pageSize := normalizePageSize(
		request.GetInt("page_size", defaultSearchPageSize),
		defaultSearchPageSize,
		maxSearchPageSize,
	)
	pageToken := strings.TrimSpace(request.GetString("page_token", ""))

	req := &assetpb.SearchAllIamPoliciesRequest{
		Scope:      strings.TrimSpace(scope),
		Query:      strings.TrimSpace(request.GetString("query", "")),
		AssetTypes: compactStrings(request.GetStringSlice("asset_types", nil)),
		OrderBy:    strings.TrimSpace(request.GetString("order_by", "")),
	}

	log.Printf("search_gcp_iam_policies scope=%s", req.GetScope())

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Asset client: %w", err)
	}
	defer closeAssetClient(client)

	it := client.SearchAllIamPolicies(ctx, req)
	pager := iterator.NewPager(it, pageSize, pageToken)

	var results []*assetpb.IamPolicySearchResult
	nextPageToken, err := pager.NextPage(&results)
	if err != nil {
		return nil, fmt.Errorf("failed to search IAM policies: %w", err)
	}

	resp := &assetpb.SearchAllIamPoliciesResponse{
		Results:       results,
		NextPageToken: nextPageToken,
	}

	return newProtoToolResult(resp)
}

func assetHistoryHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	parent, err := request.RequireString("parent")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	assetNames := compactStrings(request.GetStringSlice("asset_names", nil))
	if len(assetNames) == 0 {
		return mcp.NewToolResultError("asset_names is required"), nil
	}
	if len(assetNames) > maxHistoryAssetNames {
		return mcp.NewToolResultError(
			fmt.Sprintf("asset_names supports at most %d entries", maxHistoryAssetNames),
		), nil
	}

	contentType, err := parseContentType(
		request.GetString("content_type", ""),
		defaultHistoryContent,
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	relationshipTypes := compactStrings(
		request.GetStringSlice("relationship_types", nil),
	)
	if err := validateRelationshipRequest(contentType, relationshipTypes); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	readTimeWindow, err := parseTimeWindow(
		request.GetString("start_time", ""),
		request.GetString("end_time", ""),
		false,
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req := &assetpb.BatchGetAssetsHistoryRequest{
		Parent:            strings.TrimSpace(parent),
		AssetNames:        assetNames,
		ContentType:       contentType,
		ReadTimeWindow:    readTimeWindow,
		RelationshipTypes: relationshipTypes,
	}

	log.Printf("get_gcp_asset_history parent=%s assets=%d", req.GetParent(), len(assetNames))

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Asset client: %w", err)
	}
	defer closeAssetClient(client)

	resp, err := client.BatchGetAssetsHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset history: %w", err)
	}

	return newProtoToolResult(resp)
}

func queryAssetsHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	parent, err := request.RequireString("parent")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	statement := strings.TrimSpace(request.GetString("statement", ""))
	jobReference := strings.TrimSpace(request.GetString("job_reference", ""))

	if (statement == "" && jobReference == "") || (statement != "" && jobReference != "") {
		return mcp.NewToolResultError(
			"provide exactly one of statement or job_reference",
		), nil
	}

	readTime, err := parseOptionalTimestamp(
		request.GetString("read_time", ""),
		"read_time",
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	readTimeWindow, err := parseTimeWindow(
		request.GetString("start_time", ""),
		request.GetString("end_time", ""),
		true,
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if readTime != nil && readTimeWindow != nil {
		return mcp.NewToolResultError(
			"read_time cannot be combined with start_time or end_time",
		), nil
	}

	timeout, err := parseOptionalDuration(
		request.GetString("timeout", ""),
		"timeout",
	)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req := &assetpb.QueryAssetsRequest{
		Parent:    strings.TrimSpace(parent),
		PageToken: strings.TrimSpace(request.GetString("page_token", "")),
		Timeout:   timeout,
	}

	if pageSize := request.GetInt("page_size", 0); pageSize > 0 {
		if pageSize > 1000 {
			pageSize = 1000
		}
		req.PageSize = int32(pageSize)
	}

	if statement != "" {
		req.Query = &assetpb.QueryAssetsRequest_Statement{Statement: statement}
	} else {
		req.Query = &assetpb.QueryAssetsRequest_JobReference{
			JobReference: jobReference,
		}
	}

	if readTime != nil {
		req.Time = &assetpb.QueryAssetsRequest_ReadTime{ReadTime: readTime}
	}
	if readTimeWindow != nil {
		req.Time = &assetpb.QueryAssetsRequest_ReadTimeWindow{
			ReadTimeWindow: readTimeWindow,
		}
	}

	log.Printf("query_gcp_assets parent=%s", req.GetParent())

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Asset client: %w", err)
	}
	defer closeAssetClient(client)

	resp, err := client.QueryAssets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}

	return newProtoToolResult(resp)
}

func effectiveIAMPoliciesHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	scope, err := request.RequireString("scope")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	names := compactStrings(request.GetStringSlice("names", nil))
	if len(names) == 0 {
		return mcp.NewToolResultError("names is required"), nil
	}
	if len(names) > maxEffectiveIAMNames {
		return mcp.NewToolResultError(
			fmt.Sprintf("names supports at most %d entries", maxEffectiveIAMNames),
		), nil
	}

	req := &assetpb.BatchGetEffectiveIamPoliciesRequest{
		Scope: strings.TrimSpace(scope),
		Names: names,
	}

	log.Printf("get_gcp_effective_iam_policies scope=%s names=%d", req.GetScope(), len(names))

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Asset client: %w", err)
	}
	defer closeAssetClient(client)

	resp, err := client.BatchGetEffectiveIamPolicies(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get effective IAM policies: %w", err)
	}

	return newProtoToolResult(resp)
}
