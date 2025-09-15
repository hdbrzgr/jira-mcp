package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
)

// Input types for typed tools
type AddCommentInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
	Comment  string `json:"comment" validate:"required"`
}

type GetCommentsInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
}

func RegisterJiraCommentTools(s *server.MCPServer) {
	jiraAddCommentTool := mcp.NewTool("add_comment",
		mcp.WithDescription("Add a comment to a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("comment", mcp.Required(), mcp.Description("The comment text to add to the issue")),
	)
	s.AddTool(jiraAddCommentTool, mcp.NewTypedToolHandler(jiraAddCommentHandler))

	jiraGetCommentsTool := mcp.NewTool("get_comments",
		mcp.WithDescription("Retrieve all comments from a Jira issue"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)
	s.AddTool(jiraGetCommentsTool, mcp.NewTypedToolHandler(jiraGetCommentsHandler))
}

func jiraAddCommentHandler(ctx context.Context, request mcp.CallToolRequest, input AddCommentInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	comment := &jira.Comment{
		Body: input.Comment,
	}

	createdComment, _, err := client.Issue.AddCommentWithContext(ctx, input.IssueKey, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %v", err)
	}

	result := fmt.Sprintf("Comment added successfully!\nID: %s\nAuthor: %s\nCreated: %s",
		createdComment.ID,
		createdComment.Author.DisplayName,
		createdComment.Created)

	return mcp.NewToolResultText(result), nil
}

func jiraGetCommentsHandler(ctx context.Context, request mcp.CallToolRequest, input GetCommentsInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Get comments for the issue - use custom request since GetCommentsWithContext may not exist
	req, err := client.NewRequest("GET", fmt.Sprintf("rest/api/2/issue/%s/comment", input.IssueKey), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	var commentsResponse struct {
		Comments []jira.Comment `json:"comments"`
	}
	_, err = client.Do(req, &commentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %v", err)
	}

	if len(commentsResponse.Comments) == 0 {
		return mcp.NewToolResultText("No comments found for this issue."), nil
	}

	var result strings.Builder
	for _, comment := range commentsResponse.Comments {
		authorName := "Unknown"
		if comment.Author.DisplayName != "" {
			authorName = comment.Author.DisplayName
		}

		result.WriteString(fmt.Sprintf("ID: %s\nAuthor: %s\nCreated: %s\nUpdated: %s\nBody: %s\n\n",
			comment.ID,
			authorName,
			comment.Created,
			comment.Updated,
			comment.Body))
	}

	return mcp.NewToolResultText(result.String()), nil
}
