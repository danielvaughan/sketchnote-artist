# Google ADK for Go - Getting Started Guide

**Source**: [ADK Go Quickstart](https://google.github.io/adk-docs/get-started/go/)

## Overview
The Agent Development Kit (ADK) for Go allows developers to build agentic applications using the Go programming language. It provides a framework for defining agents, tools, and models.

## Prerequisites
- Go 1.24.4 or later.
- ADK Go v0.2.0 or later.
- A Google Cloud Project with Vertex AI API enabled (or Google AI Studio API key).

## Core Concepts

### 1. Project Structure
Recommended structure:
```
my_agent/
  agent.go       # Main agent code
  .env           # API keys
  go.mod         # Module definition
```

### 2. Installation
Initialize your module and install the ADK:
```bash
go mod init <your-module-path>
go get google.golang.org/adk
```

### 3. Basic Agent Implementation
A standard agent implementation involves three main steps:

#### A. Model Initialization
Initialize the generative model (e.g., Gemini).
```go
ctx := context.Background()
model, err := gemini.NewModel(ctx, "gemini-pro", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})
```

#### B. Agent Configuration
Create an agent using `llmagent.New`.
```go
myAgent, err := llmagent.New(llmagent.Config{
    Name:        "my_agent",
    Model:       model,
    Description: "Description of what the agent does.",
    Instruction: "System instruction for the agent.",
    Tools:       []tool.Tool{ /* tools here */ },
})
```

#### C. Execution (Launcher)
Use the ADK launcher to run the agent.
```go
config := &launcher.Config{
    AgentLoader: agent.NewSingleLoader(myAgent),
}
l := full.NewLauncher()
if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
    log.Fatal(err)
}
```

## Running the Agent
You can run the agent using the command line:
```bash
go run agent.go
```
The launcher provides a CLI and potentially a web interface for interacting with the agent.
