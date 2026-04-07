<div align="center">

# ⚡ aigate

**OpenAI-compatible and Anthropic-compatible AI Gateway for Kiro and AWS Bedrock**

A lightweight local API proxy that lets Cursor, Cline, Continue, LangChain, and any OpenAI/Anthropic-compatible tool use Claude through Kiro — with a single binary and zero config.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/hoazgazh/aigate)](https://github.com/hoazgazh/aigate/releases/latest)

</div>

---

## What is aigate?

aigate is a lightweight OpenAI-compatible and Anthropic-compatible AI gateway for Kiro and AWS Bedrock. It lets tools like Cursor, Claude Code, Cline, Continue, and LangChain connect through a single local proxy.

Download one binary, run it, point your tools at `http://localhost:8000/v1` — done.

### Supported AI Backends

- ✅ **Kiro** — free Claude Sonnet 4.5, Haiku 4.5, Sonnet 4, DeepSeek, MiniMax, Qwen
- 🔜 AWS Bedrock
- 🔜 More coming...

### Works with Cursor, Cline, Continue, and More

Any tool that supports OpenAI or Anthropic APIs works out of the box:

Cursor • Claude Code • Codex CLI • OpenCode • Aider • Cline • Roo Code • Kilo Code • Continue • Zed • Windsurf • OpenAI SDK • Anthropic SDK • LangChain • Obsidian • any OpenAI/Anthropic-compatible tool

> **None of these tools support Kiro natively.** aigate bridges the gap — it translates their standard API calls into Kiro's internal format.

---

## Quick Start

### 1. Download

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-arm64 -o aigate && chmod +x aigate

# macOS Intel
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-amd64 -o aigate && chmod +x aigate

# Linux x86_64
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-amd64 -o aigate && chmod +x aigate

# Linux ARM64
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-arm64 -o aigate && chmod +x aigate
```

### 2. Configure

Log in to [Kiro IDE](https://kiro.dev/) or run `kiro-cli login` first. Then:

```bash
export API_KEY="my-secret-key"    # You make this up — protects your proxy
```

**That's it.** aigate auto-detects your Kiro credentials. No other config needed.

<details>
<summary>How auto-detection works</summary>

aigate checks these paths in order:

| Source | Path | Used by |
|--------|------|---------|
| kiro-cli SQLite | `~/.local/share/kiro-cli/data.sqlite3` | Linux / WSL / kiro-cli |
| amazon-q SQLite | `~/.local/share/amazon-q/data.sqlite3` | amazon-q-developer-cli |
| Kiro IDE JSON | `~/.aws/sso/cache/kiro-auth-token.json` | macOS / Kiro IDE |

To override: `export KIRO_CLI_DB_FILE=...` or `export KIRO_CREDS_FILE=...`

</details>

### 3. Start

```bash
./aigate
```

```
⚡ aigate v0.2.0
   ├─ Listening on http://0.0.0.0:8000
   ├─ OpenAI API:    /v1/chat/completions
   ├─ Anthropic API: /v1/messages
   └─ Models:        /v1/models
```

Keep this terminal open. Open a new terminal for the next step.

### 4. Use

```bash
curl http://localhost:8000/v1/chat/completions \
  -H "Authorization: Bearer my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-sonnet-4-5","messages":[{"role":"user","content":"Hello!"}]}'
```

---

## OpenAI-Compatible API

aigate exposes a fully OpenAI-compatible `/v1/chat/completions` endpoint. Any tool or SDK that works with OpenAI will work with aigate.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/chat/completions` | POST | Chat completions (streaming + non-streaming) |
| `/v1/models` | GET | List available models |
| `/health` | GET | Health check |

### Use with OpenAI Python SDK

```python
from openai import OpenAI

client = OpenAI(base_url="http://localhost:8000/v1", api_key="my-secret-key")
response = client.chat.completions.create(
    model="claude-sonnet-4-5",
    messages=[{"role": "user", "content": "Hello!"}],
    stream=True,
)
for chunk in response:
    print(chunk.choices[0].delta.content or "", end="")
```

---

## Anthropic-Compatible API

aigate also exposes a native Anthropic `/v1/messages` endpoint for tools that use the Anthropic SDK directly.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/messages` | POST | Anthropic Messages API (streaming + non-streaming) |

### Use with Anthropic Python SDK

```python
import anthropic

client = anthropic.Anthropic(base_url="http://localhost:8000", api_key="my-secret-key")
response = client.messages.create(
    model="claude-sonnet-4-5",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Hello!"}],
)
print(response.content[0].text)
```

---

## Kiro Proxy — Use Kiro as an OpenAI API

aigate acts as a proxy between your tools and Kiro. It translates OpenAI/Anthropic API calls into Kiro's internal format, handles authentication, token refresh, and streaming automatically.

```
Your Tool (Cursor, Cline, SDK...)
    ↓ OpenAI/Anthropic API
  aigate (localhost:8000)
    ↓ Kiro internal API
  Kiro → Claude Sonnet 4.5, Haiku 4.5, etc.
```

---

## Use with Cursor / Cline / Continue

Set in your IDE settings:

| Setting | Value |
|---------|-------|
| Base URL | `http://localhost:8000/v1` |
| API Key | your `API_KEY` value (e.g. `my-secret-key`) |
| Model | `claude-sonnet-4-5` |

### Codex CLI

```toml
# ~/.codex/config.toml
openai_base_url = "http://localhost:8000/v1"
model = "claude-sonnet-4-5"
```
```bash
export OPENAI_API_KEY="my-secret-key"
codex
```

### OpenCode

```json
// opencode.json
{
  "provider": {
    "openai-compatible": {
      "apiKey": "my-secret-key",
      "baseURL": "http://localhost:8000/v1",
      "models": {
        "claude-sonnet-4-5": { "maxTokens": 16384 }
      }
    }
  }
}
```

### Aider

```bash
aider --openai-api-base http://localhost:8000/v1 --openai-api-key my-secret-key --model openai/claude-sonnet-4-5
```

### Claude Code

```bash
export ANTHROPIC_BASE_URL="http://localhost:8000"
export ANTHROPIC_API_KEY="my-secret-key"
claude
```

### Zed

```json
// settings.json
{
  "language_models": {
    "openai": {
      "api_url": "http://localhost:8000/v1",
      "available_models": [{"name": "claude-sonnet-4-5", "max_tokens": 16384}]
    }
  }
}
```

---

## Supported Models

| Model | Description |
|-------|-------------|
| `claude-sonnet-4-5` | Claude Sonnet 4.5 — balanced performance |
| `claude-haiku-4-5` | Claude Haiku 4.5 — fast responses |
| `claude-sonnet-4` | Claude Sonnet 4 — previous gen |
| `claude-3.7-sonnet` | Claude 3.7 Sonnet — legacy |
| `deepseek-v3.2` | DeepSeek V3.2 — open MoE model |
| `minimax-m2.1` | MiniMax M2.1 — planning & workflows |
| `qwen3-coder-next` | Qwen3 Coder — coding focused |

> Model availability depends on your Kiro subscription tier.

---

## Download Options

| Platform | Architecture | Download |
|----------|-------------|----------|
| macOS | Apple Silicon (M1/M2/M3/M4) | [aigate-darwin-arm64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-arm64) |
| macOS | Intel | [aigate-darwin-amd64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-amd64) |
| Linux | x86_64 | [aigate-linux-amd64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-amd64) |
| Linux | ARM64 | [aigate-linux-arm64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-arm64) |
| Windows | x86_64 | [aigate-windows-amd64.exe](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-windows-amd64.exe) |

---

## Build from Source

```bash
git clone https://github.com/hoazgazh/aigate.git
cd aigate
make build
API_KEY=my-secret-key ./bin/aigate
```

## Docker

```bash
docker run -p 8000:8000 \
  -e API_KEY=my-secret-key \
  -v ~/.local/share/kiro-cli:/root/.local/share/kiro-cli:ro \
  ghcr.io/hoazgazh/aigate
```

---

## Configuration

| Env Variable | Required | Description |
|-------------|----------|-------------|
| `API_KEY` | ✅ | Password to protect your proxy (you make this up) |
| `KIRO_CREDS_FILE` | Auto-detected | Path to Kiro IDE JSON credentials |
| `KIRO_CLI_DB_FILE` | Auto-detected | Path to kiro-cli SQLite database |
| `REFRESH_TOKEN` | | Kiro refresh token (manual) |
| `KIRO_REGION` | | AWS region (default: `us-east-1`) |
| `PORT` | | Server port (default: `8000`) |
| `HOST` | | Server host (default: `0.0.0.0`) |

---

## License

MIT
