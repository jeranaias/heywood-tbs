package ai

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

// HeywoodTools defines the function-calling tools available to Heywood.
// These allow the AI to take real actions: create tasks, send messages, look up data.
var HeywoodTools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "create_task",
			Description: "Create a task assigned to a staff member, SPC, or instructor. Use when the XO directs an action like 'have SSgt Diaz schedule remedial training' or 'flag Perez for counseling'.",
			Parameters: jsonSchema(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Short task title, e.g. 'Schedule remedial land nav for 2ndLt Perez'",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Detailed task description with context and expectations",
					},
					"assigned_to": map[string]interface{}{
						"type":        "string",
						"description": "Who to assign: a role ('spc', 'staff') or instructor name/ID",
					},
					"priority": map[string]interface{}{
						"type": "string",
						"enum": []string{"high", "medium", "low"},
					},
					"due_date": map[string]interface{}{
						"type":        "string",
						"description": "Due date in YYYY-MM-DD format",
					},
					"related_id": map[string]interface{}{
						"type":        "string",
						"description": "Related entity ID, e.g. student ID 'STU-042'",
					},
				},
				"required": []string{"title", "assigned_to", "priority"},
			}),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "send_message",
			Description: "Send an internal message to a role or person. Use when the XO says 'tell SSgt Diaz...' or 'notify the S-3...'.",
			Parameters: jsonSchema(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"to": map[string]interface{}{
						"type":        "string",
						"description": "Recipient: a role ('spc', 'staff', 'student') or person name/ID",
					},
					"subject": map[string]interface{}{
						"type":        "string",
						"description": "Message subject line",
					},
					"body": map[string]interface{}{
						"type":        "string",
						"description": "Message body with full context",
					},
					"related_id": map[string]interface{}{
						"type":        "string",
						"description": "Related entity ID (task, student, etc.)",
					},
				},
				"required": []string{"to", "subject", "body"},
			}),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "lookup_student",
			Description: "Look up detailed student data by name or ID. Use when the conversation requires specific student info not already in context.",
			Parameters: jsonSchema(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Student ID (e.g. 'STU-042') or partial name (e.g. 'Perez')",
					},
				},
				"required": []string{"query"},
			}),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "lookup_schedule",
			Description: "Look up training schedule for a specific date or date range.",
			Parameters: jsonSchema(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"date": map[string]interface{}{
						"type":        "string",
						"description": "Date in YYYY-MM-DD format",
					},
					"scope": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"day", "week"},
						"description": "Whether to return just that day or the full week",
					},
				},
				"required": []string{"date"},
			}),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "web_search",
			Description: "Search the web for current information via SearXNG. Use when asked about news, recent events, policies, USMC updates, or anything requiring up-to-date data not in the TBS database.",
			Parameters: jsonSchema(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
				},
				"required": []string{"query"},
			}),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "lookup_exam_results",
			Description: "Look up a student's detailed exam results — which questions they got right/wrong and topic areas. Use when a student asks how they did on a test, what they missed, or what to study. NEVER reveal actual test questions or correct answers — only topic areas and performance patterns.",
			Parameters: jsonSchema(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"student_id": map[string]interface{}{
						"type":        "string",
						"description": "Student ID (e.g. 'STU-001')",
					},
					"exam_number": map[string]interface{}{
						"type":        "integer",
						"description": "Exam number (1-4)",
					},
				},
				"required": []string{"student_id", "exam_number"},
			}),
		},
	},
}

// jsonSchema converts a map to json.RawMessage for the Parameters field.
func jsonSchema(schema map[string]interface{}) json.RawMessage {
	b, _ := json.Marshal(schema)
	return json.RawMessage(b)
}
