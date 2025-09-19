package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	// Create the .env file in the root directory with OPENAI_API_KEY key and value
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(option.WithAPIKey(apiKey))

	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ask me for weather: ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput) // remove trailing newline

	get_weather_tool := openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "get_weather",
		Description: openai.String("Get weather at the given location"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"location": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"location"},
		},
	})

	get_humidity_tool := openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "get_humidity",
		Description: openai.String("Get humidity for a city"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"location": map[string]string{"type": "string"},
			},
			"required": []string{"location"},
		},
	})

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userInput),
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			get_weather_tool,
			get_humidity_tool,
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	}

	// Make initial chat completion request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}

	// If there is a was a function call, continue the conversation
	params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		// Extract the location from the function call arguments
		var args map[string]interface{}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			panic(err)
		}
		if toolCall.Function.Name == "get_weather" {
			location := args["location"].(string)

			// Simulate getting weather data
			weatherData := getWeather(location)
			params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
		}
		if toolCall.Function.Name == "get_humidity" {
			location := args["location"].(string)
			humidityData := getHumidity(location)
			params.Messages = append(params.Messages, openai.ToolMessage(humidityData, toolCall.ID))
		}
	}

	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	switch strings.ToLower(location) {
	case "patiala":
		return "10°C"
	case "delhi":
		return "14°C"
	case "nainital":
		return "5°C"
	case "dehradun":
		return "12°C"
	default:
		return "Unknown"
	}
}

func getHumidity(location string) string {
	// In a real implementation, this function would call a weather API
	switch strings.ToLower(location) {
	case "patiala":
		return "40%"
	case "delhi":
		return "55%"
	case "nainital":
		return "70%"
	case "dehradun":
		return "60%"
	default:
		return "Unknown"
	}
}
