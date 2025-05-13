// Send the add label operation without checking first
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
