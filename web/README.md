# Jira MCP Test Interface

This directory contains a web-based test interface for the Jira MCP (Model Context Protocol) server.

## Files

- `index.html` - Interactive test UI for all Jira MCP tools
- `swagger.html` - Swagger UI for API documentation
- `openapi.yaml` - OpenAPI 3.0 specification for the MCP tools
- `server.go` - Simple HTTP server to serve the test interface
- `start.sh` - Convenience script to start both servers

## Quick Start

### Option 1: Using the Go server (Recommended)

```bash
# Start the test UI server
cd web
go run server.go

# Or build and run
go build -o test-server server.go
./test-server
```

The test UI will be available at:
- **Test Interface**: http://localhost:3000
- **Swagger Documentation**: http://localhost:3000/swagger.html

### Option 2: Using Python's built-in server

```bash
cd web
python3 -m http.server 3000
```

### Option 3: Using Node.js http-server

```bash
cd web
npx http-server -p 3000
```

## Usage

1. **Start your Jira MCP server** first:
   ```bash
   # From the project root
   ./jira-mcp -http_port 8080
   ```

2. **Open the test interface** at http://localhost:3000

3. **Initialize session** by clicking "Initialize Session" button

4. **Select a tool** from the grid (Issues, Search, Comments, Workflow)

5. **Fill in the parameters** and click "Execute"

## Available Tools

### Issue Management
- **get_issue** - Retrieve detailed information about a specific issue
- **create_issue** - Create a new Jira issue
- **create_child_issue** - Create a sub-task linked to a parent issue
- **update_issue** - Modify an existing issue's details
- **list_issue_types** - List all available issue types in a project

### Search
- **search_issue** - Search for issues using JQL (Jira Query Language)

### Comments
- **add_comment** - Add a comment to an issue
- **get_comments** - Retrieve all comments from an issue

### Workflow
- **transition_issue** - Move an issue through its workflow states

## Environment Variables

Make sure your Jira MCP server has the required environment variables:

```bash
# Required
export JIRA_HOST=http://your-jira-instance.com

# Authentication (choose one method)
# Method 1: Personal Access Token (recommended)
export JIRA_PAT=your-personal-access-token

# Method 2: Username/Password (for older Jira versions)
export JIRA_USERNAME=your-username
export JIRA_PASSWORD=your-password
```

## API Documentation

The OpenAPI specification is available at:
- **Interactive Docs**: http://localhost:3000/swagger.html
- **Raw Spec**: http://localhost:3000/openapi.yaml

## Features

### Test Interface
- ✅ Modern, responsive UI built with Tailwind CSS and Alpine.js
- ✅ Session management with automatic initialization
- ✅ Form validation and error handling
- ✅ Real-time results display
- ✅ Organized by tool categories
- ✅ Copy-paste friendly for issue keys and JQL queries

### API Documentation
- ✅ Complete OpenAPI 3.0 specification
- ✅ Interactive Swagger UI
- ✅ Request/response examples
- ✅ Parameter validation schemas
- ✅ Authentication documentation

## Troubleshooting

### Common Issues

1. **"Invalid session ID" error**
   - Make sure to click "Initialize Session" first
   - Check that the server URL is correct (default: http://localhost:8080/mcp)

2. **CORS errors**
   - The MCP server should handle CORS automatically
   - If issues persist, try using the same origin for both servers

3. **Connection refused**
   - Ensure the Jira MCP server is running on port 8080
   - Check that your Jira credentials are properly configured

4. **Tool execution fails**
   - Verify your Jira credentials and permissions
   - Check that the issue keys and project keys exist
   - Review the error message in the Results section

### Testing Tips

1. **Start with get_issue** to test connectivity:
   - Use an existing issue key like "TEST-1"
   - This will also show available transitions for workflow testing

2. **For JQL searches**, try simple queries first:
   - `project = YOUR_PROJECT`
   - `assignee = currentUser()`
   - `status != Done`

3. **For transitions**, get available transition IDs from get_issue first

## Development

To modify or extend the test interface:

1. Edit `index.html` for UI changes
2. Update `openapi.yaml` for API documentation
3. Modify `server.go` for server behavior

The interface uses:
- **Tailwind CSS** for styling
- **Alpine.js** for reactivity
- **Fetch API** for HTTP requests
- **JSON-RPC 2.0** for MCP protocol communication
