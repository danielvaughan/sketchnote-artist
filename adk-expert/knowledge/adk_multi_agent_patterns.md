# Multi-Agent Systems in ADK

**Source**: [ADK Multi-Agent Systems](https://google.github.io/adk-docs/agents/multi-agents/)

## Overview
ADK enables building complex applications by combining multiple agents. This document outlines core primitives and common patterns.

## 1. ADK Primitives for Agent Composition

### 1.1. Agent Hierarchy
- **Concept**: Parent-child relationship defined in `BaseAgent`.
- **Implementation**: Pass a list of agent instances to `sub_agents` when initializing a parent agent.
- **Key Rule**: An agent instance can only be a sub-agent of one parent.

### 1.2. Workflow Agents as Orchestrators
Specialized agents that manage execution flow without performing tasks themselves.
- **SequentialAgent**: Executes sub-agents one after another. Passes the same `InvocationContext` sequentially.
- **ParallelAgent**: Runs sub-agents concurrently. Often used with a subsequent agent to aggregate results.
- **LoopAgent**: Repeats execution of sub-agents until a condition is met.

## 2. Common Multi-Agent Patterns

### Coordinator/Dispatcher Pattern
- **Structure**: Central `LlmAgent` (Coordinator) managing specialized sub-agents.
- **Goal**: Route requests to the appropriate specialist.
- **Mechanism**: LLM-Driven Delegation or Explicit Invocation.

### Sequential Pipeline Pattern
- **Structure**: `SequentialAgent` with ordered sub-agents.
- **Goal**: Multi-step process where output of one step feeds the next.
- **Mechanism**: Shared Session State (write to `output_key`, read from `context.state`).

### Parallel Fan-Out/Gather Pattern
- **Structure**: `ParallelAgent` for concurrent tasks, followed by an aggregator (often in a `SequentialAgent`).
- **Goal**: Reduce latency for independent tasks.
- **Mechanism**: Sub-agents write to distinct state keys; aggregator reads them.

### Hierarchical Task Decomposition
- **Structure**: Multi-level tree of agents.
- **Goal**: Break down complex goals into simpler steps.
- **Mechanism**: Recursive delegation down the hierarchy.

### Review/Critique Pattern (Generator-Critic)
- **Structure**: Generator agent followed by a Critic agent (within `SequentialAgent`).
- **Goal**: Improve quality via review.
- **Mechanism**: Generator writes draft; Critic reviews and provides feedback.

### Iterative Refinement Pattern
- **Structure**: Agents inside a `LoopAgent`.
- **Goal**: Progressively improve a result.
- **Mechanism**: Loop continues until quality threshold or max iterations reached.

### Human-in-the-Loop Pattern
- **Structure**: Human intervention points within a workflow.
- **Goal**: Oversight or manual tasks.
- **Mechanism**: Custom Tool to pause/request input, or delegation to a conceptual "Human Agent".
