package main

import (
	"log"

	asset "cloud.google.com/go/asset/apiv1"

	"github.com/mark3labs/mcp-go/server"
)

func closeAssetClient(client *asset.Client) {
	if err := client.Close(); err != nil {
		log.Printf("failed to close Cloud Asset client: %v", err)
	}
}

func registerTools(s *server.MCPServer) {
	s.AddTool(newListAssetsTool(), listAssetsHandler)
	s.AddTool(newSearchResourcesTool(), searchResourcesHandler)
	s.AddTool(newSearchIAMPoliciesTool(), searchIAMPoliciesHandler)
	s.AddTool(newAssetHistoryTool(), assetHistoryHandler)
	s.AddTool(newQueryAssetsTool(), queryAssetsHandler)
	s.AddTool(newEffectiveIAMPoliciesTool(), effectiveIAMPoliciesHandler)
}
