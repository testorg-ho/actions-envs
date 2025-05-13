package main

import (
    "context"
    "fmt"
    "log"
    "os"

    v2 "github.com/ctreminiom/go-atlassian/jira/v2"
    "github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

func main() {
    // Get configuration from environment variables
    jiraURL := getEnvOrDefault("JIRA_URL", "https://your-domain.atlassian.net")
    jiraUser := getEnvOrDefault("JIRA_USER", "your-email@example.com")
    jiraToken := getEnvOrDefault("JIRA_TOKEN", "your-api-token")
    
    // Get issue key and label from arguments
    if len(os.Args) < 3 {
        fmt.Println("Usage: program <issue-key> <label>")
        fmt.Println("Example: program PROJECT-123 high-priority")
        os.Exit(1)
    }
    issueKey := os.Args[1]
    newLabel := os.Args[2]

    // Create a new Jira client
    client, err := v2.New(nil, jiraURL)
    if err != nil {
        log.Fatalf("Error creating Jira client: %v", err)
    }

    // Authenticate with Jira
    client.Auth.SetBasicAuth(jiraUser, jiraToken)

    // Set up context
    ctx := context.Background()

    // Create the update payload
    payload := &models.IssueUpdateScheme{
        Update: &models.IssueFieldOperationsScheme{
            Labels: []*models.IssueFieldOperationScheme{
                {
                    Add: newLabel,
                },
            },
        },
    }

    // Update the issue
    _, err = client.Issue.Update(ctx, issueKey, payload)
    if err != nil {
        log.Fatalf("Error updating issue label: %v", err)
    }

    fmt.Printf("Successfully added label '%s' to %s\n", newLabel, issueKey)
}

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
