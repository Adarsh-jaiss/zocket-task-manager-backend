package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/adarsh-jaiss/zocket/types"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (g *GeminiClient) AnalyzeTask(task types.Task) (*types.AITaskBreakdownResponse, error) {
	prompt := fmt.Sprintf(`Analyze the following task and break it down into smaller, manageable subtasks:

Task Title: %s
Description: %s
Priority: %s

Please provide:
1. A detailed analysis of the task
2. A list of suggested subtasks with descriptions
3. Estimated complexity for each subtask (High/Medium/Low)
4. Recommended order of completion
5. Any potential dependencies between subtasks

Format the response as a JSON object with the following structure:
{
    "analysis": "overall analysis text",
    "suggestions": [
        {
            "suggestion_text": "detailed breakdown and recommendation",
            "sub_tasks": [
                {
                    "title": "subtask title",
                    "description": "subtask description",
                    "priority": "High/Medium/Low"
                }
            ]
        }
    ]
}`, task.Title, task.Description, task.Priority)

	ctx := context.Background()
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	fmt.Printf("%+v\n", resp.Candidates[0].Content)

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	// Parse the response
	var aiResp types.AITaskBreakdownResponse

	// Debug logging
	fmt.Printf("Response type: %T\n", resp.Candidates[0].Content.Parts[0])
	fmt.Printf("Response value: %#v\n", resp.Candidates[0].Content.Parts[0])

	// Get the response text
	text := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Clean up the response - remove markdown code blocks
	if len(text) > 8 && text[:8] == "```json\n" {
		text = text[8:]
	}
	if len(text) > 4 && text[len(text)-4:] == "\n```" {
		text = text[:len(text)-4]
	}

	if err := json.Unmarshal([]byte(text), &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	aiResp.TaskID = task.TaskID
	return &aiResp, nil
}

func (g *GeminiClient) Close() {
	if g.client != nil {
		g.client.Close()
	}
}
