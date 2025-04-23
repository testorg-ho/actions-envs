from jira import JIRA
import csv

# Function to get labels for a specific ticket
def get_ticket_labels(jira, ticket_key):
    """
    Retrieves the labels for a specific JIRA ticket.
    
    Args:
        jira: JIRA client instance
        ticket_key: The ticket ID/key (e.g., 'OPS-123')
        
    Returns:
        List of labels or empty list if no labels exist
    """
    try:
        issue = jira.issue(ticket_key)
        labels = issue.fields.labels
        return labels
    except Exception as e:
        print(f"Error retrieving labels for ticket {ticket_key}: {e}")
        return []

# Function to get labels for multiple tickets
def get_labels_for_tickets(jira, ticket_keys):
    """
    Retrieves labels for multiple tickets.
    
    Args:
        jira: JIRA client instance
        ticket_keys: List of ticket keys/IDs
        
    Returns:
        Dictionary mapping ticket keys to their labels
    """
    result = {}
    for key in ticket_keys:
        result[key] = get_ticket_labels(jira, key)
    return result

def export_issues_to_csv(jira, jql_query, output_file):
    """
    Exports JIRA issues to a CSV file with ticket number and description.
    No headers are included in the CSV.
    
    Args:
        jira: JIRA client instance
        jql_query: JQL query string to find issues
        output_file: Path to the output CSV file
    """
    try:
        # Fetch issues using the provided JQL query
        issues = jira.search_issues(jql_query, maxResults=False)
        
        # Write to CSV without headers
        with open(output_file, 'w', newline='', encoding='utf-8') as csvfile:
            csv_writer = csv.writer(csvfile)
            
            # Iterate through each issue and write ticket number and description
            for issue in issues:
                ticket_key = issue.key
                # Some issues might not have a description
                description = issue.fields.description if hasattr(issue.fields, 'description') and issue.fields.description else ""
                
                # Write the row to CSV
                csv_writer.writerow([ticket_key, description])
        
        print(f"Successfully exported {len(issues)} issues to {output_file}")
        return True
    
    except Exception as e:
        print(f"Error exporting issues to CSV: {e}")
        return False

def get_field_metadata(jira, field_name):
    """
    Retrieves metadata about a specific field, including possible values for dropdowns.
    
    Args:
        jira: JIRA client instance
        field_name: The name of the field to get metadata for (e.g., 'primaryTeam')
        
    Returns:
        Dictionary containing field metadata, including possible values for dropdown fields
    """
    try:
        # Get all available fields
        fields = jira.fields()
        
        # Find the field id for the given field name
        field_id = None
        for field in fields:
            if field['name'].lower() == field_name.lower():
                field_id = field['id']
                break
        
        if not field_id:
            print(f"Field '{field_name}' not found")
            return None
        
        # For custom fields, get metadata including allowedValues
        meta = jira.editmeta()
        
        # Some fields might be present in specific issue types
        # You might need to provide an example issue key to get complete metadata
        if not meta:
            # Try to get metadata from a sample issue
            sample_issues = jira.search_issues('project is not EMPTY', maxResults=1)
            if sample_issues:
                meta = jira.editmeta(sample_issues[0].key)
        
        # Extract the allowed values for the field
        for field_meta in meta.get('fields', {}).values():
            if field_meta.get('name', '').lower() == field_name.lower():
                return field_meta
                
        # If still not found, try to get metadata directly
        field_info = jira.field(field_id)
        return field_info
    
    except Exception as e:
        print(f"Error retrieving metadata for field {field_name}: {e}")
        return None

def get_field_allowed_values(jira, field_name):
    """
    Gets all possible values for a dropdown field like 'primaryTeam'
    
    Args:
        jira: JIRA client instance
        field_name: The name of the field to get allowed values for
        
    Returns:
        List of allowed values or None if not found/applicable
    """
    try:
        metadata = get_field_metadata(jira, field_name)
        
        if not metadata:
            return None
            
        # Extract allowed values
        allowed_values = metadata.get('allowedValues', [])
        
        # For dropdown fields, allowed values typically have a 'value' or 'name' property
        values = []
        for value in allowed_values:
            if isinstance(value, dict):
                # Extract the display value (might be under 'value' or 'name' key)
                if 'value' in value:
                    values.append(value['value'])
                elif 'name' in value:
                    values.append(value['name'])
                else:
                    values.append(value)  # Add the whole object if no standard property found
            else:
                values.append(value)
                
        return values
        
    except Exception as e:
        print(f"Error retrieving allowed values for field {field_name}: {e}")
        return None

def update_field_value(jira, ticket_key, field_name, new_value):
    """
    Updates a specific field value for a JIRA ticket.
    
    Args:
        jira: JIRA client instance
        ticket_key: The ticket ID/key (e.g., 'OPS-123')
        field_name: The name of the field to update (e.g., 'primaryTeam')
        new_value: The new value to set for the field
        
    Returns:
        True if successful, False otherwise
    """
    try:
        # Get the field ID
        fields = jira.fields()
        field_id = None
        
        for field in fields:
            if field['name'].lower() == field_name.lower():
                field_id = field['id']
                break
                
        if not field_id:
            print(f"Field '{field_name}' not found")
            return False
            
        # Prepare update data
        update_data = {field_id: new_value}
        
        # Update the issue
        jira.issue(ticket_key).update(fields=update_data)
        print(f"Successfully updated {field_name} for ticket {ticket_key}")
        return True
        
    except Exception as e:
        print(f"Error updating {field_name} for ticket {ticket_key}: {e}")
        return False

# Example usage:
# jira = JIRA(server='https://your-jira-server')
# tickets = jira.search_issues('project = OPS')
# for ticket in tickets:
#     labels = get_ticket_labels(jira, ticket.key)
#     print(f"Ticket {ticket.key}: Labels = {labels}")

"""
# Connect to JIRA
jira = JIRA(server='https://your-jira-instance.com', basic_auth=('username', 'password'))

# Define query and output file
jql_query = 'project = OPS AND primaryTeam is not EMPTY'
csv_file = '/Users/tarashrynchuk/Downloads/golang-console-project/jira_issues.csv'

# Export issues to CSV
export_issues_to_csv(jira, jql_query, csv_file)

# Example usage for the new functions:
jira = JIRA(server='https://your-jira-instance.com', basic_auth=('username', 'password'))

# Get all possible values for primaryTeam
primary_team_values = get_field_allowed_values(jira, 'primaryTeam')
print(f"Possible values for primaryTeam: {primary_team_values}")

# Update primaryTeam field for a ticket
update_field_value(jira, 'OPS-123', 'primaryTeam', 'DevOps')
"""
