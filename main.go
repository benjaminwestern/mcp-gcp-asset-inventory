package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(serverName, serverVersion)

	registerTools(s)

	log.Printf("Starting %s v%s", serverName, serverVersion)

	if err := server.ServeStdio(s); err != nil {
		log.Printf("Server exited: %v", err)
	}
}
