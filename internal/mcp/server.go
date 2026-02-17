package mcp

import (
	"fmt"
	"ralph2/internal/service"

	mcp_sdk "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type MCPServer struct {
	svc *service.OrchestratorService
	mcp *mcp_sdk.Server
}

func NewMCPServer() *MCPServer {
	// Initialize the transport and server
	transport := stdio.NewStdioServerTransport()
	s := &MCPServer{
		svc: service.NewOrchestratorService(),
		mcp: mcp_sdk.NewServer(transport),
	}

	// Register tools
	s.registerTools()

	return s
}

func (s *MCPServer) registerTools() {
	s.mcp.RegisterTool("run_task", "Start an autonomous coding task in Ralph2", func(args struct {
		Prompt     string `json:"prompt" jsonschema:"description=The user request or coding task"`
		Complexity string `json:"complexity,omitempty" jsonschema:"description=Complexity level (fast, streamlined, full),default=streamlined"`
	}) (*mcp_sdk.ToolResponse, error) {
		
		complexity := args.Complexity
		if complexity == "" {
			complexity = "streamlined"
		}

		err := s.svc.Run(args.Prompt, complexity)
		if err != nil {
			return mcp_sdk.NewToolResponse(mcp_sdk.NewTextContent(fmt.Sprintf("Task failed: %v", err))), nil
		}

		return mcp_sdk.NewToolResponse(mcp_sdk.NewTextContent(fmt.Sprintf("Task completed successfully. Check main/task branches for results."))), nil
	})

	s.mcp.RegisterTool("get_status", "Get the current state of the Ralph2 FSM", func(args struct{}) (*mcp_sdk.ToolResponse, error) {
		status := s.svc.SM.GetState()
		complexity := s.svc.SM.GetComplexity()
		
		msg := fmt.Sprintf("Current State: %s\nComplexity: %s", status, complexity)
		return mcp_sdk.NewToolResponse(mcp_sdk.NewTextContent(msg)), nil
	})
}

func (s *MCPServer) Serve() error {
	return s.mcp.Serve()
}
