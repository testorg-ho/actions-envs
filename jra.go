package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/ctreminiom/go-atlassian/jira/v2"
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
    client, err := jira.New(nil, jiraURL)
    if err != nil {
        log.Fatalf("Error creating Jira client: %v", err)
    }

    // Authenticate with Jira
    client.Auth.SetBasicAuth(jiraUser, jiraToken)

    // Set up context
    ctx := context.Background()

    // Create operation to add the new label
    // Note: Jira will handle duplicates automatically (won't add if already exists)
    operation := &jira.UpdateOperationScheme{
        Operation: "add",
        Path:      "/fields/labels",
        Value:     newLabel,
    }

    // Update the issue
    _, err = client.Issue.DoTransition(ctx, issueKey, &jira.IssueUpdateScheme{
        Update: &jira.UpdateRequestScheme{
            Operations: []*jira.UpdateOperationScheme{operation},
        },
    })
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
