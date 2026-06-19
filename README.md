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
│   ├── support-bot.yaml           # Simple support triage example
│   ├── coffee-order-bot.yaml      # Multi-path order flow example
│   ├── insurance-claim-bot.yaml   # Complex multi-branch claims flow
│   ├── it-helpdesk-bot.yaml       # IT helpdesk with troubleshooting trees
│   └── components-demo-bot.yaml # Showcases all input/action components
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

### Insurance Claim Bot (Complex Multi-Branch Flow)

```bash
./chatbot --bot examples/insurance-claim-bot.yaml
```

This example demonstrates:

- **Main menu branching** — file claim, check status, speak to agent, FAQ
- **Intent-based category parsing** — auto / home / health / other with rich example phrases for rule matching
- **Parallel sub-flows** — each incident type has its own multi-step input chain
- **Confirmation gates** — review summary before submission
- **Session variables** — `claim_category`, `claim_status`, `queue_position` via `set_var`
- **Loop-back menus** — return to `start` from FAQ, cancelled, or post-submit states

Example paths:

- `file a claim` → enter policy number → `car accident` → complete auto form → `confirm`
- `check claim status` → enter claim ID → view status
- `speak to agent` → queue simulation → connect or return to menu

### IT Helpdesk Bot (Troubleshooting Trees + Parsing)

```bash
./chatbot --bot examples/it-helpdesk-bot.yaml
```

This example demonstrates:

- **Four top-level paths** — report issue, password reset, request access, check ticket
- **Category intent parsing** — hardware / software / network / email routed via substring and word-overlap matching
- **Troubleshooting trees** — each category runs guided fixes before escalating to ticket creation
- **Yes/no parsing** — `fixed` vs `not_fixed` intents after troubleshoot steps
- **Priority assignment** — `set_var` sets `ticket_priority` (HIGH/MEDIUM/LOW) per category
- **Access request workflow** — multi-step form with manager approval routing

Example paths:

- `something is broken` → `wifi` → describe issue → try fixes → `no` → create ticket
- `forgot password` → enter email → reset link sent
- `need access` → software → Figma → justification → manager → `confirm`

### Components Demo Bot (All Input & Action Types)

```bash
./chatbot --bot examples/components-demo-bot.yaml
```

Showcases every bot component:

- **Input types** — `email`, `regex`, `choice`, `confirm`, `text`, `number`
- **Actions** — `set_var`, `append_var`, `copy_var`, `increment_var`, `delete_var`
- **Global intents** — `menu`, `help`, `exit` available on every node
- **Terminal nodes** — explicit `terminal: true` for hard exits

### With Ollama LLM

```bash
./chatbot --bot examples/support-bot.yaml --llm ollama --ollama-url http://localhost:11434 --ollama-model llama2
```

## YAML Bot Definition

The bot definition follows this schema:

```yaml
bot:
  name: SupportBot
  global_intents:
    - name: menu
      examples: ["menu", "back"]
      next: start

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
    terminal: true
```

### Node Types

1. **Intent-based nodes**: Use `intents` to route user input to different flows
2. **Input capture nodes**: Use `input` to capture and save user input to variables
3. **Terminal nodes**: Set `terminal: true`, or omit `next`/`intents`/`input` when no global intents exist

### Input Types

| Type | Purpose | Key fields |
|------|---------|------------|
| `text` | Free-form text (default) | `save_as` |
| `choice` | Pick from fixed options | `options`, `save_as` |
| `confirm` | Yes/no branching | `save_as`, `on_yes`, `on_no` |
| `email` | Email validation | `save_as`, `retry_message` |
| `number` | Numeric validation | `save_as`, `retry_message` |
| `regex` | Pattern validation | `save_as`, `pattern`, `retry_message` |

```yaml
# Choice input
input:
  type: choice
  save_as: size
  options: [small, medium, large]

# Confirm with branching
input:
  type: confirm
  save_as: confirmed
  on_yes: success
  on_no: cancelled

# Regex validation
input:
  type: regex
  save_as: policy_id
  pattern: "^POL-[0-9]{6}$"
  retry_message: "Must be POL-XXXXXX format."
```

Invalid input triggers `retry_message` and re-prompts on the same node.

### Global Intents

Define intents available on **every node** under `bot.global_intents`:

```yaml
bot:
  name: MyBot
  global_intents:
    - name: menu
      examples: ["menu", "back", "start over"]
      next: start
    - name: help
      examples: ["help"]
      next: show_help
```

Node-specific intents take priority over global intents during matching.

### Variable Interpolation

Variables can be interpolated in messages using `{{variable_name}}`:

```yaml
message: "Your order ID is {{order_id}}"
```

### Actions

| Action | Purpose | Args |
|--------|---------|------|
| `set_var` | Set a session variable | `name`, optional `value` (defaults to user input) |
| `append_var` | Append to a variable | `name`, optional `value`, optional `separator` (default `", "`) |
| `copy_var` | Copy between variables | `from`, `to` |
| `increment_var` | Increment numeric variable | `name`, optional `by` (default 1) |
| `delete_var` | Remove a variable | `name` |

```yaml
actions:
  - type: set_var
    args:
      name: status
      value: "ACTIVE"

  - type: append_var
    args:
      name: tags
      separator: " | "

  - type: copy_var
    args:
      from: rating
      to: last_rating

  - type: increment_var
    args:
      name: counter
      by: 1

  - type: delete_var
    args:
      name: temp_data
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
