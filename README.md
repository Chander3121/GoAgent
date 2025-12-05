# GoAgent

An AI Agent implementation in Go using OpenAI's ChatGPT LLM with function calling capabilities. This project demonstrates how to build an intelligent agent that can understand user queries and call external functions to retrieve real-time data.

## Overview

GoAgent is a learning project that showcases building an autonomous AI agent in Go. The agent can understand natural language queries about weather and humidity, determine which functions to call, execute those functions, and provide intelligent responses to users.

**Key Features:**
- Natural language processing using ChatGPT (GPT-4o model)
- Function calling/tool use capabilities
- Multi-turn conversation support
- Mock weather and humidity data retrieval
- Interactive command-line interface

## Technology Stack

- **Language:** Go 1.x
- **AI Model:** OpenAI ChatGPT (GPT-4o)
- **OpenAI SDK:** [openai-go](https://github.com/openai/openai-go)
- **Environment Management:** godotenv

## Project Structure

```
GoAgent/
├── main.go           # Main application with agent logic
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
└── .gitignore        # Git ignore file
```

## How It Works

### Architecture

1. **User Input:** The agent accepts natural language queries from the user via command-line input
2. **LLM Processing:** Sends the query to OpenAI's ChatGPT with defined tools/functions
3. **Tool Selection:** ChatGPT determines which functions need to be called based on the query
4. **Function Execution:** The agent executes the selected functions and retrieves data
5. **Response Generation:** The agent sends the function results back to ChatGPT for a final response
6. **Output:** Returns an intelligent response to the user

### Available Tools/Functions

#### 1. `get_weather`
Retrieves the current temperature for a given location.

**Parameters:**
- `location` (string, required): The city or location name

**Mock Data:**
- Patiala: 10°C
- Delhi: 14°C
- Nainital: 5°C
- Dehradun: 12°C

#### 2. `get_humidity`
Retrieves the humidity percentage for a given location.

**Parameters:**
- `location` (string, required): The city or location name

**Mock Data:**
- Patiala: 40%
- Delhi: 55%
- Nainital: 70%
- Dehradun: 60%

## Prerequisites

- Go 1.16 or higher
- OpenAI API key (available from [platform.openai.com](https://platform.openai.com))

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Chander3121/GoAgent.git
cd GoAgent
```

2. Install dependencies:
```bash
go mod download
go mod tidy
```

3. Create a `.env` file in the root directory:
```bash
touch .env
```

4. Add your OpenAI API key to the `.env` file:
```
OPENAI_API_KEY=your_api_key_here
```

## Usage

Run the application:
```bash
go run main.go
```

The agent will prompt you to ask a weather-related question:
```
Ask me for weather: What is the weather in Delhi?
```

### Example Queries

- "What is the temperature in Patiala?"
- "How humid is it in Nainital?"
- "Tell me the weather and humidity in Delhi"
- "Is it cold in Dehradun?"
- "What's the weather like in multiple cities?"

## How Function Calling Works

### Step 1: Initial Request
The user query is sent to ChatGPT along with the available tools definition.

```go
params := openai.ChatCompletionNewParams{
    Messages: []openai.ChatCompletionMessageParamUnion{
        openai.UserMessage(userInput),
    },
    Tools: []openai.ChatCompletionToolUnionParam{
        get_weather_tool,
        get_humidity_tool,
    },
    Model: openai.ChatModelGPT4o,
}
```

### Step 2: Function Call Response
ChatGPT analyzes the query and determines which functions to call:

```go
toolCalls := completion.Choices[0].Message.ToolCalls
```

### Step 3: Execute Functions
The agent executes the required functions:

```go
for _, toolCall := range toolCalls {
    if toolCall.Function.Name == "get_weather" {
        // Extract location and retrieve weather data
    }
    if toolCall.Function.Name == "get_humidity" {
        // Extract location and retrieve humidity data
    }
}
```

### Step 4: Final Response
The function results are sent back to ChatGPT to generate a natural language response:

```go
completion, err = client.Chat.Completions.New(ctx, params)
println(completion.Choices[0].Message.Content)
```

## Code Example

Here's a simplified example of the main flow:

```go
// Load environment and API key
err := godotenv.Load()
apiKey := os.Getenv("OPENAI_API_KEY")
client := openai.NewClient(option.WithAPIKey(apiKey))

// Get user input
fmt.Print("Ask me for weather: ")
userInput, _ := reader.ReadString('\n')

// Define tools
get_weather_tool := openai.ChatCompletionFunctionTool(
    openai.FunctionDefinitionParam{
        Name:        "get_weather",
        Description: openai.String("Get weather at the given location"),
        // ... parameters definition
    })

// Make initial request with tools
params := openai.ChatCompletionNewParams{
    Messages: []openai.ChatCompletionMessageParamUnion{
        openai.UserMessage(userInput),
    },
    Tools: []openai.ChatCompletionToolUnionParam{
        get_weather_tool,
        get_humidity_tool,
    },
    Model: openai.ChatModelGPT4o,
}

completion, _ := client.Chat.Completions.New(ctx, params)
toolCalls := completion.Choices[0].Message.ToolCalls

// Execute functions and get results
// ...

// Get final response
completion, _ = client.Chat.Completions.New(ctx, params)
println(completion.Choices[0].Message.Content)
```

## Building for Production

Create an executable binary:
```bash
go build -o goagent
./goagent
```

Or build for specific OS:
```bash
# For macOS
GOOS=darwin GOARCH=amd64 go build -o goagent

# For Linux
GOOS=linux GOARCH=amd64 go build -o goagent

# For Windows
GOOS=windows GOARCH=amd64 go build -o goagent.exe
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `OPENAI_API_KEY` | Your OpenAI API key | Yes |

## Extending the Project

### Adding New Tools

To add a new tool to the agent:

1. **Define the function** in `main.go`:
```go
func getAirQuality(location string) string {
    // Your implementation
    return "data"
}
```

2. **Create the tool definition**:
```go
get_air_quality_tool := openai.ChatCompletionFunctionTool(
    openai.FunctionDefinitionParam{
        Name:        "get_air_quality",
        Description: openai.String("Get air quality for a location"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]any{
                "location": map[string]string{"type": "string"},
            },
            "required": []string{"location"},
        },
    })
```

3. **Add to tools list**:
```go
Tools: []openai.ChatCompletionToolUnionParam{
    get_weather_tool,
    get_humidity_tool,
    get_air_quality_tool,  // New tool
}
```

4. **Handle function calls**:
```go
if toolCall.Function.Name == "get_air_quality" {
    location := args["location"].(string)
    airQualityData := getAirQuality(location)
    params.Messages = append(params.Messages, 
        openai.ToolMessage(airQualityData, toolCall.ID))
}
```

### Connecting Real APIs

Replace the mock functions with actual API calls:
- OpenWeatherMap API for real weather data
- AQI API for air quality information
- Any other third-party service

## Learning Objectives

This project demonstrates:
- ✅ Integration with OpenAI's Go SDK
- ✅ Function calling/tool use in LLMs
- ✅ Multi-turn conversation patterns
- ✅ Error handling in Go
- ✅ Environment configuration management
- ✅ Building intelligent agents

## Troubleshooting

### Issue: "Error loading .env file"
**Solution:** Ensure you have a `.env` file in the root directory with your OpenAI API key.

### Issue: "API key not found"
**Solution:** Verify that `OPENAI_API_KEY` is correctly set in your `.env` file.

### Issue: "No function call"
**Solution:** This is expected behavior if ChatGPT doesn't determine that a function call is needed based on your query.

### Issue: Unknown location
**Solution:** The mock data only includes Patiala, Delhi, Nainital, and Dehradun. Query these locations or modify the mock functions to add more.

## References

- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [OpenAI Function Calling Documentation](https://platform.openai.com/docs/guides/function-calling)
- [Go Documentation](https://golang.org/doc/)

## Author

Created by [Chander Prakash](https://github.com/Chander3121) for self-learning purposes.

## Disclaimer

This is a learning project demonstrating AI agent concepts. The mock weather data is not real. For production use, integrate with real weather APIs and handle errors appropriately.
