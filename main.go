// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"google.golang.org/api/iterator"
)

type AssetSummary struct {
	Name      string `json:"name"`
	AssetType string `json:"asset_type"`
	Location  string `json:"location"`
}

func main() {
	s := server.NewMCPServer(
		"Cloud Asset Inventory Server",
		"1.0.0",
	)

	tool := mcp.NewTool("list_gcp_assets",
		mcp.WithDescription("Lists Google Cloud assets within a specific project."),

		mcp.WithString("project_id",
			mcp.Required(),
			mcp.Description("The Google Cloud project ID to query for assets."),
		),

		mcp.WithArray("asset_types",
			mcp.Description("Optional. An array of strings for asset types to filter by (e.g., 'storage.googleapis.com/Bucket')."),
			mcp.Items(map[string]any{"type": "string"}),
		),
	)

	s.AddTool(tool, listAssetsHandler)

	log.Println("Starting MCP server for Cloud Asset Inventory...")
	log.Println("Tool 'list_gcp_assets' is available.")

	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server exited: %v", err)
	}
}

func listAssetsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID, err := request.RequireString("project_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	assetTypes := request.GetStringSlice("asset_types", nil)

	log.Printf("Handling request for project '%s', asset types: %v", projectID, assetTypes)

	client, err := asset.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Cloud Asset client: %w", err)
	}
	defer client.Close()

	req := &assetpb.ListAssetsRequest{
		Parent:      fmt.Sprintf("projects/%s", projectID),
		AssetTypes:  assetTypes,
		ContentType: assetpb.ContentType_RESOURCE,
	}

	it := client.ListAssets(ctx, req)

	var summaries []AssetSummary
	for {
		asset, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed while iterating over assets: %w", err)
		}
		summaries = append(summaries, AssetSummary{
			Name:      asset.Name,
			AssetType: asset.AssetType,
			Location:  asset.GetResource().GetLocation(),
		})
	}

	log.Printf("Found %d assets. Returning summary as a JSON string.", len(summaries))

	jsonData, err := json.Marshal(summaries)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results to json: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
