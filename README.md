# Google Cloud Asset Inventory MCP server

This repository contains a Go-based MCP server for Google Cloud Asset
Inventory. It runs over `stdio` and exposes a set of read-only tools for
listing assets, searching resources, searching IAM policies, querying asset
history, issuing SQL queries, and retrieving effective IAM policies.

## Prerequisites

You need Google Cloud credentials that can call the Cloud Asset Inventory API.
The server uses Application Default Credentials, so the most common setup is
the Google Cloud CLI:

```sh
gcloud auth application-default login
```

You also need the Cloud Asset API enabled for the project, folder, or
organization you want to inspect.

This repository is pinned to Go `1.26.1`. If you use `mise`, run:

```sh
mise install
```

## Build and test

You can build and verify the server locally with the standard Go toolchain.

```sh
go build -o mcp-asset-server .
go test ./...
```

## Live smoke test

You can run a real end-to-end MCP smoke test against Cloud Asset Inventory by
passing a project ID to the helper script. The script uses your local
Application Default Credentials, builds the current server, initializes MCP
over `stdio`, and calls `list_gcp_assets`.

```sh
./scripts/smoke-test-mcp.sh your-project-id
```

The first argument is required and must be the target Google Cloud project ID.
You can optionally pass a second argument to change the requested page size.

## Available tools

The server exposes Cloud Asset Inventory read APIs as MCP tools. Each tool
returns JSON text.

- `list_gcp_assets`
  Lists assets for a `project_id` or `parent` scope and supports asset type
  filters, pagination, snapshot timestamps, and relationship snapshots.
- `search_gcp_resources`
  Wraps `searchAllResources` for free-text and fielded resource search across a
  project, folder, or organization.
- `search_gcp_iam_policies`
  Wraps `searchAllIamPolicies` for IAM binding search across a project, folder,
  or organization.
- `get_gcp_asset_history`
  Wraps `batchGetAssetsHistory` for historical snapshots of up to 100 named
  assets.
- `query_gcp_assets`
  Wraps `queryAssets` for BigQuery-compatible SQL queries and follow-up reads
  with a `job_reference`.
- `get_gcp_effective_iam_policies`
  Wraps `effectiveIamPolicies.batchGet` for effective policy lookups on up to
  20 full resource names.

## Example MCP client usage

You can run this server as a subprocess from any MCP-compatible client. The
following Python example shows one way to wire it into an ADK agent:

```python
from google.adk import Agent
from google.adk.tools.mcp_tool import MCPToolset
from google.adk.tools.mcp_tool.mcp_toolset import StdioServerParameters

root_agent = Agent(
    model="your-model",
    name="asset_inventory_agent",
    description="Answers questions with Cloud Asset Inventory data.",
    instruction="Use the MCP tools to inspect Google Cloud assets.",
    tools=[
        MCPToolset(
            connection_params=StdioServerParameters(
                command="./mcp-asset-server",
                args=["stdio"],
            ),
        ),
    ],
)
```

## API coverage notes

As of March 19, 2026, this MCP server matches a focused subset of the Cloud
Asset Inventory v1 REST API, not the full specification documented in the
[official Cloud Asset Inventory REST reference](https://docs.cloud.google.com/asset-inventory/docs/reference/rest).

The current MCP tools cover these v1 methods:

- `assets.list`
- `searchAllResources`
- `searchAllIamPolicies`
- `batchGetAssetsHistory`
- `queryAssets`
- `effectiveIamPolicies.batchGet`

The server does not yet expose the remaining v1 REST surface, including:

- `analyzeIamPolicy`
- `analyzeIamPolicyLongrunning`
- `analyzeMove`
- `analyzeOrgPolicies`
- `analyzeOrgPolicyGovernedAssets`
- `analyzeOrgPolicyGovernedContainers`
- `exportAssets`
- `feeds` CRUD operations
- `operations.get`
- `savedQueries` CRUD operations

It also does not yet mirror every optional request shape for the methods above.
For example, the current `query_gcp_assets` tool does not expose `outputConfig`
for writing query results externally.
