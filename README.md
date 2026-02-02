# Chatbot-Go

A CLI-based chatbot framework in Go that supports deterministic conversational flows and is LLM-ready for integration with local LLMs (e.g., Ollama, llama.cpp, LM Studio).

## Core Design Principles

- **Deterministic Flow Control**: Conversation flow is defined in YAML and executed via FSM. The engine controls all state transitions.
- **LLM as Assistive Component**: LLMs can classify intent, extract entities, and generate response text, but they **never** change conversation state, choose next nodes, or execute actions.
- **Offline-First**: The system works fully without any LLM configured.
- **CLI-First**: No web UI - CLI renderer only.

## Architecture

```
CLI (cobra)
  ↓
Bot Loader + Validator
  ↓
Conversation Engine (FSM)
  ↓
Input Interpretation
   ├─ Rule Router (first priority)
   └─ LLM Router (optional, fallback)
  ↓
Action Executor
  ↓
Response Renderer (CLI)
```

## Project Structure

```
chatbot-go/
├── cmd/
│   ├── main.go              # Entry point
│   └── root.go              # Cobra CLI setup
│
├── internal/
│   ├── bot/                 # Bot definition and loading
│   │   ├── loader.go        # YAML → AST
│   │   ├── schema.go        # Basic validation
│   │   └── types.go         # Bot, Node, Intent types
│   │
│   ├── engine/              # FSM-based conversation engine
│   │   ├── engine.go        # Main conversation loop
│   │   ├── session.go       # Session management
│   │   ├── fsm.go           # State transitions
│   │   └── types.go         # Engine and Session types
│   │
│   ├── router/              # Input routing
│   │   ├── router.go        # Router interface
│   │   ├── rule_router.go   # Rule-based routing
│   │   └── llm_router.go    # LLM-based routing
│   │
│   ├── llm/                 # LLM provider abstraction
│   │   ├── provider.go      # LLM interface
│   │   ├── noop.go          # No-op provider (default)
│   │   └── ollama.go        # Ollama HTTP stub
│   │
│   ├── actions/             # Action execution
│   │   └── executor.go      # Action executor (set_var)
│   │
│   ├── render/              # Output rendering
│   │   └── cli.go          # CLI renderer
│   │
│   └── validate/            # Flow validation
│       └── flow.go         # Comprehensive validation
│
├── examples/
│   └── support-bot.yaml    # Example bot definition
│
├── go.mod
└── README.md
```

## Installation

```bash
go mod download
go build -o chatbot ./cmd/main.go
```

## Usage

### Basic Usage (No LLM)

```bash
./chatbot --bot examples/support-bot.yaml
```

### Another End-to-End Example (Coffee Order Bot)

```bash
./chatbot --bot examples/coffee-order-bot.yaml
```

Example conversation:

- Choose a path:
  - Type `order coffee` to place an order end-to-end (captures `size`, `drink`, `milk`, `pickup_time`, `customer_name`, then confirms).
  - Type `track my order` to enter an `order_number` and see a demo status.
  - Type `hours` to see store hours and return to the start menu.

### With Ollama LLM

```bash
./chatbot --bot examples/support-bot.yaml --llm ollama --ollama-url http://localhost:11434 --ollama-model llama2
```

## YAML Bot Definition

The bot definition follows this schema:

```yaml
bot:
  name: SupportBot

flows:
  start:
    message: "Hi! How can I help you?"
    intents:
      - name: order_issue
        examples:
          - "problem with order"
          - "order not delivered"
        next: ask_order_id

      - name: refund
        examples:
          - "refund"
          - "money back"
        next: refund_flow

  ask_order_id:
    message: "Please enter your order ID"
    input:
      type: text
      save_as: order_id
    next: end

  refund_flow:
    message: "Refund process started"
    next: end

  end:
    message: "Thank you!"
```

### Node Types

1. **Intent-based nodes**: Use `intents` to route user input to different flows
2. **Input capture nodes**: Use `input` to capture and save user input to variables
3. **Terminal nodes**: Nodes with no `next` and no `intents` end the conversation

### Variable Interpolation

Variables can be interpolated in messages using `{{variable_name}}`:

```yaml
message: "Your order ID is {{order_id}}"
```

### Actions

Currently supported actions:

- `set_var`: Save a value to a session variable
  ```yaml
  actions:
    - type: set_var
      args:
        name: order_id
        value: "12345"
  ```

## How It Works

1. **Engine Loop**:
   - Starts at the `start` node
   - Renders the node's message (with variable interpolation)
   - Reads user input
   - Routes input using RuleRouter first, then LLMRouter if needed
   - Executes any declared actions
   - Transitions to the next node
   - Repeats until a terminal node is reached

2. **Input Routing**:
   - **RuleRouter** (first priority): Exact match, keyword match, simple similarity
   - **LLMRouter** (optional): Only used if no rule match or multiple ambiguous matches

3. **Session Management**:
   - Tracks current node
   - Maintains variables map
   - Records conversation history

## Adding an LLM Provider

To add a new LLM provider:

1. Implement the `llm.Provider` interface in `internal/llm/`:

```go
type MyLLMProvider struct {
    // your fields
}

func (p *MyLLMProvider) ClassifyIntent(ctx context.Context, input string, intents []Intent) (string, error) {
    // your implementation
}

func (p *MyLLMProvider) ExtractEntities(ctx context.Context, input string, schema map[string]string) (map[string]string, error) {
    // your implementation
}

func (p *MyLLMProvider) GenerateText(ctx context.Context, prompt Prompt) (string, error) {
    // your implementation
}
```

2. Add a case in `cmd/root.go` to initialize your provider:

```go
case "myllm":
    llmProvider = llm.NewMyLLMProvider(...)
```

3. Use it:

```bash
./chatbot --bot examples/support-bot.yaml --llm myllm
```

## Testing

The engine is designed to be testable without CLI:

```go
bot, _ := bot.LoadFromFile("examples/support-bot.yaml")
llmProvider := llm.NewNoopProvider()
engine := engine.NewConversationEngine(bot, llmProvider)
// Test engine.Run() with mocked input/output
```

## Key Constraints

- **LLMs cannot**: Change conversation state directly, choose next nodes, execute actions, mutate session data
- **LLMs can**: Classify intent, extract entities, generate response text (optional)
- **No hard-coded flows**: All flows defined in YAML
- **No global state**: All state in Session
- **No web UI**: CLI only
- **No database**: In-memory session only

## License

MIT
