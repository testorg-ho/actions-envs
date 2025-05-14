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
    client, err := v2.New(nil, jiraURL)
    if err != nil {
        log.Fatalf("Error creating Jira client: %v", err)
    }

    // Authenticate with Jira
    client.Auth.SetBasicAuth(jiraUser, jiraToken)

    // Set up context
    ctx := context.Background()

    // Get the current issue to retrieve existing labels
    issue, _, err := client.Issue.Get(ctx, issueKey, nil)
    if err != nil {
        log.Fatalf("Error getting issue: %v", err)
    }

    // Get existing labels and add the new one
    existingLabels := issue.Fields.Labels
    
    // Check if label already exists to avoid unnecessary updates
    labelExists := false
    for _, label := range existingLabels {
        if label == newLabel {
            labelExists = true
            break
        }
    }
    
    if labelExists {
        fmt.Printf("Label '%s' already exists on issue %s. No update needed.\n", newLabel, issueKey)
        return
    }
    
    // Add the new label
    updatedLabels := append(existingLabels, newLabel)
    
    // Create the update payload
    payload := map[string]interface{}{
        "fields": map[string]interface{}{
            "labels": updatedLabels,
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





if err != nil {
    // Check if error has response data (depends on the client)
    if apiErr, ok := err.(*v2.ErrorResponse); ok {
        log.Printf("API Error: %s", apiErr.Response.Status)
        log.Printf("Error details: %s", apiErr.Message)
    }
    log.Fatalf("Error updating issue label: %v", err)
}

    // Create the update payload using models.IssueSchemaV1
    payload := &models.IssueSchemaV1{
        Fields: &models.IssueFieldsSchemeV1{
            Labels: updatedLabels,
        },
    }
