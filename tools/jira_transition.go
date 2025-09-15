package tools

import (
	"context"
	"fmt"

	"github.com/hdbrzgr/jira-mcp/v2/services"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Input types for typed tools
type TransitionIssueInput struct {
	IssueKey     string `json:"issue_key" validate:"required"`
	TransitionID string `json:"transition_id" validate:"required"`
	Comment      string `json:"comment,omitempty"`
}

func RegisterJiraTransitionTool(s *server.MCPServer) {
	jiraTransitionTool := mcp.NewTool("transition_issue",
		mcp.WithDescription("Transition an issue through its workflow using a valid transition ID. Get available transitions from jira_get_issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The issue to transition (e.g., KP-123)")),
		mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID from available transitions list")),
		mcp.WithString("comment", mcp.Description("Optional comment to add with transition")),
	)
	s.AddTool(jiraTransitionTool, mcp.NewTypedToolHandler(jiraTransitionIssueHandler))
}

func jiraTransitionIssueHandler(ctx context.Context, request mcp.CallToolRequest, input TransitionIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Use a simplified transition approach via custom request
	transitionData := map[string]interface{}{
		"transition": map[string]interface{}{
			"id": input.TransitionID,
		},
	}

	// Add comment if provided
	if input.Comment != "" {
		transitionData["update"] = map[string]interface{}{
			"comment": []map[string]interface{}{
				{
					"add": map[string]interface{}{
						"body": input.Comment,
					},
				},
			},
		}
	}

	req, err := client.NewRequest("POST", fmt.Sprintf("rest/api/2/issue/%s/transitions", input.IssueKey), transitionData)
	if err != nil {
		return nil, fmt.Errorf("failed to create transition request: %v", err)
	}

	_, err = client.Do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("transition failed: %v", err)
	}

	return mcp.NewToolResultText("Issue transition completed successfully"), nil
}
