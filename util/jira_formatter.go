package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
)

// FormatJiraIssue converts a Jira issue struct to a formatted string representation
// It handles the andygrunwald/go-jira Issue structure
func FormatJiraIssue(issue *jira.Issue) string {
	var sb strings.Builder

	// Basic issue information
	sb.WriteString(fmt.Sprintf("Key: %s\n", issue.Key))

	if issue.ID != "" {
		sb.WriteString(fmt.Sprintf("ID: %s\n", issue.ID))
	}

	if issue.Self != "" {
		sb.WriteString(fmt.Sprintf("URL: %s\n", issue.Self))
	}

	// Fields information
	if issue.Fields != nil {
		fields := issue.Fields

		// Summary and Description
		if fields.Summary != "" {
			sb.WriteString(fmt.Sprintf("Summary: %s\n", fields.Summary))
		}

		if fields.Description != "" {
			sb.WriteString(fmt.Sprintf("Description: %s\n", fields.Description))
		}

		// Issue Type
		if fields.Type.Name != "" {
			sb.WriteString(fmt.Sprintf("Type: %s\n", fields.Type.Name))
			if fields.Type.Description != "" {
				sb.WriteString(fmt.Sprintf("Type Description: %s\n", fields.Type.Description))
			}
		}

		// Status
		if fields.Status != nil && fields.Status.Name != "" {
			sb.WriteString(fmt.Sprintf("Status: %s\n", fields.Status.Name))
			if fields.Status.Description != "" {
				sb.WriteString(fmt.Sprintf("Status Description: %s\n", fields.Status.Description))
			}
		}

		// Priority
		if fields.Priority != nil && fields.Priority.Name != "" {
			sb.WriteString(fmt.Sprintf("Priority: %s\n", fields.Priority.Name))
		} else {
			sb.WriteString("Priority: None\n")
		}

		// Resolution
		if fields.Resolution != nil && fields.Resolution.Name != "" {
			sb.WriteString(fmt.Sprintf("Resolution: %s\n", fields.Resolution.Name))
			if fields.Resolution.Description != "" {
				sb.WriteString(fmt.Sprintf("Resolution Description: %s\n", fields.Resolution.Description))
			}
		}

		// People
		if fields.Reporter != nil {
			sb.WriteString(fmt.Sprintf("Reporter: %s", fields.Reporter.DisplayName))
			if fields.Reporter.EmailAddress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Reporter.EmailAddress))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("Reporter: Unassigned\n")
		}

		if fields.Assignee != nil {
			sb.WriteString(fmt.Sprintf("Assignee: %s", fields.Assignee.DisplayName))
			if fields.Assignee.EmailAddress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Assignee.EmailAddress))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("Assignee: Unassigned\n")
		}

		if fields.Creator != nil {
			sb.WriteString(fmt.Sprintf("Creator: %s", fields.Creator.DisplayName))
			if fields.Creator.EmailAddress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Creator.EmailAddress))
			}
			sb.WriteString("\n")
		}

		// Dates - convert jira.Time to time.Time to check if zero
		if !time.Time(fields.Created).IsZero() {
			sb.WriteString(fmt.Sprintf("Created: %s\n", time.Time(fields.Created).Format("2006-01-02 15:04:05")))
		}

		if !time.Time(fields.Updated).IsZero() {
			sb.WriteString(fmt.Sprintf("Updated: %s\n", time.Time(fields.Updated).Format("2006-01-02 15:04:05")))
		}

		// Project information
		if fields.Project.Name != "" {
			sb.WriteString(fmt.Sprintf("Project: %s", fields.Project.Name))
			if fields.Project.Key != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", fields.Project.Key))
			}
			sb.WriteString("\n")
		}

		// Parent issue
		if fields.Parent != nil && fields.Parent.Key != "" {
			sb.WriteString(fmt.Sprintf("Parent: %s", fields.Parent.Key))
			sb.WriteString("\n")
		}

		// Labels
		if len(fields.Labels) > 0 {
			sb.WriteString(fmt.Sprintf("Labels: %s\n", strings.Join(fields.Labels, ", ")))
		}

		// Components
		if len(fields.Components) > 0 {
			sb.WriteString("Components:\n")
			for _, component := range fields.Components {
				sb.WriteString(fmt.Sprintf("- %s", component.Name))
				if component.Description != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", component.Description))
				}
				sb.WriteString("\n")
			}
		}

		// Fix Versions
		if len(fields.FixVersions) > 0 {
			sb.WriteString("Fix Versions:\n")
			for _, version := range fields.FixVersions {
				sb.WriteString(fmt.Sprintf("- %s", version.Name))
				if version.Description != "" {
					sb.WriteString(fmt.Sprintf(" (%s)", version.Description))
				}
				sb.WriteString("\n")
			}
		}

		// Subtasks
		if len(fields.Subtasks) > 0 {
			sb.WriteString("Subtasks:\n")
			for _, subtask := range fields.Subtasks {
				sb.WriteString(fmt.Sprintf("- %s", subtask.Key))
				if subtask.Fields.Summary != "" {
					sb.WriteString(fmt.Sprintf(": %s", subtask.Fields.Summary))
				}
				if subtask.Fields.Status != nil {
					sb.WriteString(fmt.Sprintf(" [%s]", subtask.Fields.Status.Name))
				}
				sb.WriteString("\n")
			}
		}

		// Issue Links
		if len(fields.IssueLinks) > 0 {
			sb.WriteString("Issue Links:\n")
			for _, link := range fields.IssueLinks {
				if link.OutwardIssue != nil {
					sb.WriteString(fmt.Sprintf("- %s %s", link.Type.Outward, link.OutwardIssue.Key))
					if link.OutwardIssue.Fields != nil && link.OutwardIssue.Fields.Summary != "" {
						sb.WriteString(fmt.Sprintf(": %s", link.OutwardIssue.Fields.Summary))
					}
					sb.WriteString("\n")
				}
				if link.InwardIssue != nil {
					sb.WriteString(fmt.Sprintf("- %s %s", link.Type.Inward, link.InwardIssue.Key))
					if link.InwardIssue.Fields != nil && link.InwardIssue.Fields.Summary != "" {
						sb.WriteString(fmt.Sprintf(": %s", link.InwardIssue.Fields.Summary))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	// Available Transitions
	if len(issue.Transitions) > 0 {
		sb.WriteString("\nAvailable Transitions:\n")
		for _, transition := range issue.Transitions {
			sb.WriteString(fmt.Sprintf("- %s (ID: %s)\n", transition.Name, transition.ID))
		}
	}

	return sb.String()
}

// FormatJiraIssueCompact returns a compact single-line representation of a Jira issue
// Useful for search results or lists
func FormatJiraIssueCompact(issue *jira.Issue) string {
	if issue == nil {
		return ""
	}

	var parts []string

	parts = append(parts, fmt.Sprintf("Key: %s", issue.Key))

	if issue.Fields != nil {
		fields := issue.Fields

		if fields.Summary != "" {
			parts = append(parts, fmt.Sprintf("Summary: %s", fields.Summary))
		}

		if fields.Status != nil && fields.Status.Name != "" {
			parts = append(parts, fmt.Sprintf("Status: %s", fields.Status.Name))
		}

		if fields.Assignee != nil {
			parts = append(parts, fmt.Sprintf("Assignee: %s", fields.Assignee.DisplayName))
		} else {
			parts = append(parts, "Assignee: Unassigned")
		}

		if fields.Priority != nil && fields.Priority.Name != "" {
			parts = append(parts, fmt.Sprintf("Priority: %s", fields.Priority.Name))
		}
	}

	return strings.Join(parts, " | ")
}
