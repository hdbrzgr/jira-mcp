package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/jira-mcp/services"
	"github.com/nguyenvanduocit/jira-mcp/util"
)

// Input types for typed tools
type GetIssueInput struct {
	IssueKey string `json:"issue_key" validate:"required"`
	Fields   string `json:"fields,omitempty"`
	Expand   string `json:"expand,omitempty"`
}

type CreateIssueInput struct {
	ProjectKey  string `json:"project_key" validate:"required"`
	Summary     string `json:"summary" validate:"required"`
	Description string `json:"description" validate:"required"`
	IssueType   string `json:"issue_type" validate:"required"`
	Assignee    string `json:"assignee,omitempty"`
}

type CreateChildIssueInput struct {
	ParentIssueKey string `json:"parent_issue_key" validate:"required"`
	Summary        string `json:"summary" validate:"required"`
	Description    string `json:"description" validate:"required"`
	IssueType      string `json:"issue_type,omitempty"`
	Assignee       string `json:"assignee,omitempty"`
}

type UpdateIssueInput struct {
	IssueKey    string `json:"issue_key" validate:"required"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
}

type ListIssueTypesInput struct {
	ProjectKey string `json:"project_key" validate:"required"`
}

func RegisterJiraIssueTool(s *server.MCPServer) {
	jiraGetIssueTool := mcp.NewTool("get_issue",
		mcp.WithDescription("Retrieve detailed information about a specific Jira issue including its status, assignee, description, subtasks, and available transitions"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
		mcp.WithString("fields", mcp.Description("Comma-separated list of fields to retrieve (e.g., 'summary,status,assignee'). If not specified, all fields are returned.")),
		mcp.WithString("expand", mcp.Description("Comma-separated list of fields to expand for additional details (e.g., 'transitions,changelog,subtasks'). Default: 'transitions,changelog'")),
	)
	s.AddTool(jiraGetIssueTool, mcp.NewTypedToolHandler(jiraGetIssueHandler))

	jiraCreateIssueTool := mcp.NewTool("create_issue",
		mcp.WithDescription("Create a new Jira issue with specified details. Returns the created issue's key, ID, and URL"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier where the issue will be created (e.g., KP, PROJ)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the issue")),
		mcp.WithString("issue_type", mcp.Required(), mcp.Description("Type of issue to create (common types: Bug, Task, Subtask, Story, Epic)")),
		mcp.WithString("assignee", mcp.Description("Username or email of the person to assign the issue to (optional)")),
	)
	s.AddTool(jiraCreateIssueTool, mcp.NewTypedToolHandler(jiraCreateIssueHandler))

	jiraCreateChildIssueTool := mcp.NewTool("create_child_issue",
		mcp.WithDescription("Create a child issue (sub-task) linked to a parent issue in Jira. Returns the created issue's key, ID, and URL"),
		mcp.WithString("parent_issue_key", mcp.Required(), mcp.Description("The parent issue key to which this child issue will be linked (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Brief title or headline of the child issue")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Detailed explanation of the child issue")),
		mcp.WithString("issue_type", mcp.Description("Type of child issue to create (defaults to 'Subtask' if not specified)")),
		mcp.WithString("assignee", mcp.Description("Username or email of the person to assign the issue to (optional)")),
	)
	s.AddTool(jiraCreateChildIssueTool, mcp.NewTypedToolHandler(jiraCreateChildIssueHandler))

	jiraUpdateIssueTool := mcp.NewTool("update_issue",
		mcp.WithDescription("Modify an existing Jira issue's details. Supports partial updates - only specified fields will be changed"),
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the issue to update (e.g., KP-2)")),
		mcp.WithString("summary", mcp.Description("New title for the issue (optional)")),
		mcp.WithString("description", mcp.Description("New description for the issue (optional)")),
		mcp.WithString("assignee", mcp.Description("Username or email of the person to assign the issue to (optional)")),
	)
	s.AddTool(jiraUpdateIssueTool, mcp.NewTypedToolHandler(jiraUpdateIssueHandler))

	jiraListIssueTypesTool := mcp.NewTool("list_issue_types",
		mcp.WithDescription("List all available issue types in a Jira project with their IDs, names, descriptions, and other attributes"),
		mcp.WithString("project_key", mcp.Required(), mcp.Description("Project identifier to list issue types for (e.g., KP, PROJ)")),
	)
	s.AddTool(jiraListIssueTypesTool, mcp.NewTypedToolHandler(jiraListIssueTypesHandler))
}

func jiraGetIssueHandler(ctx context.Context, request mcp.CallToolRequest, input GetIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Parse expand parameter with default values
	expand := "transitions,changelog,subtasks"
	if input.Expand != "" {
		expand = input.Expand
	}

	issue, _, err := client.Issue.GetWithContext(ctx, input.IssueKey, &jira.GetQueryOptions{
		Expand: expand,
		Fields: input.Fields,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %v", err)
	}

	// Use the new util function to format the issue
	formattedIssue := util.FormatJiraIssue(issue)

	return mcp.NewToolResultText(formattedIssue), nil
}

func jiraCreateIssueHandler(ctx context.Context, request mcp.CallToolRequest, input CreateIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     input.Summary,
			Description: input.Description,
			Project: jira.Project{
				Key: input.ProjectKey,
			},
			Type: jira.IssueType{
				Name: input.IssueType,
			},
		},
	}

	// Add assignee if provided
	if input.Assignee != "" {
		issue.Fields.Assignee = &jira.User{
			Name: input.Assignee,
		}
	}

	createdIssue, _, err := client.Issue.CreateWithContext(ctx, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %v", err)
	}

	result := fmt.Sprintf("Issue created successfully!\nKey: %s\nID: %s\nURL: %s", createdIssue.Key, createdIssue.ID, createdIssue.Self)
	return mcp.NewToolResultText(result), nil
}

func jiraCreateChildIssueHandler(ctx context.Context, request mcp.CallToolRequest, input CreateChildIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	// Get the parent issue to retrieve its project
	parentIssue, _, err := client.Issue.GetWithContext(ctx, input.ParentIssueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent issue: %v", err)
	}

	// Default issue type is Sub-task if not specified
	issueType := "Subtask"
	if input.IssueType != "" {
		issueType = input.IssueType
	}

	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     input.Summary,
			Description: input.Description,
			Project: jira.Project{
				Key: parentIssue.Fields.Project.Key,
			},
			Type: jira.IssueType{
				Name: issueType,
			},
			Parent: &jira.Parent{
				Key: input.ParentIssueKey,
			},
		},
	}

	// Add assignee if provided
	if input.Assignee != "" {
		issue.Fields.Assignee = &jira.User{
			Name: input.Assignee,
		}
	}

	createdIssue, _, err := client.Issue.CreateWithContext(ctx, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to create child issue: %v", err)
	}

	result := fmt.Sprintf("Child issue created successfully!\nKey: %s\nID: %s\nURL: %s\nParent: %s",
		createdIssue.Key, createdIssue.ID, createdIssue.Self, input.ParentIssueKey)

	if issueType == "Bug" {
		result += "\n\nA bug should be linked to a Story or Task. Next step should be to create relationship between the bug and the story or task."
	}
	return mcp.NewToolResultText(result), nil
}

func jiraUpdateIssueHandler(ctx context.Context, request mcp.CallToolRequest, input UpdateIssueInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	issue := &jira.Issue{
		Key:    input.IssueKey,
		Fields: &jira.IssueFields{},
	}

	if input.Summary != "" {
		issue.Fields.Summary = input.Summary
	}

	if input.Description != "" {
		issue.Fields.Description = input.Description
	}

	if input.Assignee != "" {
		issue.Fields.Assignee = &jira.User{
			Name: input.Assignee,
		}
	}

	_, _, err := client.Issue.UpdateWithContext(ctx, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue: %v", err)
	}

	return mcp.NewToolResultText("Issue updated successfully!"), nil
}

func jiraListIssueTypesHandler(ctx context.Context, request mcp.CallToolRequest, input ListIssueTypesInput) (*mcp.CallToolResult, error) {
	client := services.JiraClient()

	req, err := client.NewRequest("GET", "rest/api/2/issuetype", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	var issueTypes []jira.IssueType
	_, err = client.Do(req, &issueTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue types: %v", err)
	}

	if len(issueTypes) == 0 {
		return mcp.NewToolResultText("No issue types found for this project."), nil
	}

	var result strings.Builder
	result.WriteString("Available Issue Types:\n\n")

	for _, issueType := range issueTypes {
		subtaskType := ""
		if issueType.Subtask {
			subtaskType = " (Subtask Type)"
		}

		result.WriteString(fmt.Sprintf("ID: %s\nName: %s%s\n", issueType.ID, issueType.Name, subtaskType))
		if issueType.Description != "" {
			result.WriteString(fmt.Sprintf("Description: %s\n", issueType.Description))
		}
		if issueType.IconURL != "" {
			result.WriteString(fmt.Sprintf("Icon URL: %s\n", issueType.IconURL))
		}
		result.WriteString("\n")
	}

	return mcp.NewToolResultText(result.String()), nil
}
