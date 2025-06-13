# Google Cloud Asset Inventory - MCP Server

This project provides a basic Model Context Protocol (MCP) server written in Go. It exposes the Google Cloud Asset Inventory API as a tool, allowing language models and other MCP-compatible clients to list cloud assets in a specified Google Cloud project.

This server is designed to be run as a local subprocess by an MCP client, communicating over standard input and standard output (`stdio`).

## Prerequisites

Before you begin, ensure you have the following installed and configured:

1. **Go**: A recent version of the Go programming language (version 1.21 or newer is recommended). You can download it from the [official Go website](https://go.dev/dl/).
2. **Google Cloud SDK**: You must have the `gcloud` command-line tool installed and authenticated. The server uses Application Default Credentials for authorisation, which is most easily configured by logging in with the SDK.

    Run the following command to authenticate:

    ```sh
    gcloud auth application-default login
    ```

3. **Enabled APIs**: Ensure the **Cloud Asset API** is enabled for the Google Cloud project you intend to query. You can manage APIs in the Google Cloud Console.

## Getting Started

First, retrieve the necessary Go package dependencies for the project. Open your terminal in the project directory and run:

```sh
go mod tidy
```

## How to Build the Server

You can compile the server into a single executable file. The build process will vary slightly depending on your target operating system.

### Standard Build (for your current OS)

To build the executable for your current operating system and architecture, run the following command. This will create a binary named `mcp-asset-server` (or `mcp-asset-server.exe` on Windows) in your project directory.

```sh
go build -o mcp-asset-server .
```

### Cross-Compilation (Building for other Operating Systems)

Go makes it straightforward to build an executable for a different operating system from the one you are developing on. This is achieved by setting the `GOOS` (target operating system) and `GOARCH` (target architecture) environment variables.

#### Building for Linux (amd64)

From macOS or Windows, you can build a Linux binary using this command:

```sh
GOOS=linux GOARCH=amd64 go build -o mcp-asset-server-linux .
```

#### Building for macOS (amd64 / Apple Silicon)

* For Intel-based Macs (amd64):

    ```sh
    GOOS=darwin GOARCH=amd64 go build -o mcp-asset-server-macos-amd64 .
    ```

* For Apple Silicon Macs (M1/M2/M3, arm64):

    ```sh
    GOOS=darwin GOARCH=arm64 go build -o mcp-asset-server-macos-arm64 .
    ```

#### Building for Windows (amd64)

To build a Windows executable (which will have a `.exe` suffix), use this command:

```sh
GOOS=windows GOARCH=amd64 go build -o mcp-asset-server.exe .
```

## Usage with Python ADK

The primary use for this MCP server is to be loaded as a toolset by an agent, such as one built with the `google-adk` Python library. The agent will manage the lifecycle of the server, running it as a subprocess and communicating with it via `stdio`.

First, ensure you have the `google-adk` library installed:

```sh
pip install google-adk
```

Next, you can define an agent in Python that uses your compiled binary. In this example, the `MCPToolset` is configured to run the macOS `arm64` binary.

**Example `.env` file:**

```plaintext
GOOGLE_GENAI_USE_VERTEXAI=FALSE
GOOGLE_API_KEY=PASTE_YOUR_ACTUAL_API_KEY_HERE
```

**Example `requirements.txt`:**

```plaintext
google-adk>=1.3.0
```

**Example `_init_.py`:**

```python
from . import agent
```

**Example `agent.py`:**

```python
from google.adk import Agent
from google.adk.tools.mcp_tool import MCPToolset
from google.adk.tools.mcp_tool.mcp_toolset import StdioServerParameters

root_agent = Agent(
    model="gemini-2.5-pro-preview-06-05",
    name="root_agent",
    description="An assist to find content available in Asset Inventory.",
    instruction="Answer user questions to the best of your knowledge and provide relevant content from Asset Inventory.",
    tools=[
        MCPToolset(
            connection_params=StdioServerParameters(
                command="./asset-inventory-mcp/mcp-asset-server-macos-arm64",
                args=["stdio"],
            ),
        ),
    ],
)

if __name__ == "__main__":
    # The agent will automatically start the MCP server subprocess.
    # The agent now has access to the `list_gcp_assets` tool.
    print(root_agent.chat("Hi, can you list the GCS buckets in the 'your-gcp-project-id' project?"))

```

### How It Works

1. **Build the binary**: Compile the Go server for your specific OS and architecture as described in the "How to Build" section.
2. **Update the path**: In your Python script, change the `command` value inside `StdioServerParameters` to the correct path of the binary you just built.
3. **Run the agent**: When you execute the Python script, the `google-adk` framework automatically starts the Go program specified in the `command` path.
4. **Communicate**: The agent communicates with the server over `stdio`, discovers the `list_gcp_assets` tool, and can now use it to fulfill user requests.
