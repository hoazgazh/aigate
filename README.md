<div align="center">

# ⚡ aigate

**Free AI Gateway — Use Claude and GPT through Kiro and GitHub Copilot, no API key needed**

One binary. Multiple free AI providers. OpenAI and Anthropic compatible.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/hoazgazh/aigate)](https://github.com/hoazgazh/aigate/releases/latest)

</div>

---

## What is aigate?

aigate is a lightweight OpenAI-compatible and Anthropic-compatible AI gateway that combines **multiple free AI providers** into a single local API endpoint. Use Claude Sonnet 4.5 from Kiro and GPT-4.1 from GitHub Copilot — all through one proxy, no paid API keys required.

Download one binary, run it, point your tools at `http://localhost:8000/v1` — done.

### Free AI Models — No API Key Required

| Provider | Models | Cost | Setup |
|----------|--------|------|-------|
| **Kiro** | Claude Sonnet 4.5, Haiku 4.5, Sonnet 4, DeepSeek V3.2, MiniMax, Qwen | Free | `kiro-cli login` or Kiro IDE |
| **GitHub Copilot** | GPT-4.1, Claude 3.5 Sonnet | Free (2000 req/mo) | `./aigate --copilot-login` |
| 🔜 Gemini | Gemini 2.5 Pro, Flash | Free | Coming soon |

> **No paid API keys needed.** aigate uses your existing free-tier accounts to provide API access.

### Not Just Free Tier — Works with Pro Subscriptions Too

Many AI subscriptions (Kiro Pro, GitHub Copilot Pro/Pro+) give you access to powerful models but **don't provide an API key**. aigate unlocks API access from these subscriptions too:

| Subscription | Models you unlock | API key provided? | aigate? |
|-------------|-------------------|-------------------|---------|
| Kiro Free | Claude Sonnet 4.5, Haiku 4.5 | ❌ No API | ✅ Works |
| Kiro Pro (paid) | Claude Opus 4.5, higher limits | ❌ No API | ✅ Works |
| Copilot Free | GPT-4.1, Claude 3.5 | ❌ No API | ✅ Works |
| Copilot Pro ($10/mo) | Unlimited completions, more models | ❌ No API | ✅ Works |
| Copilot Pro+ ($39/mo) | GPT-5, o3-pro, highest limits | ❌ No API | ✅ Works |

> **If you're paying for a Pro subscription but can't use it outside the IDE — aigate fixes that.** Use your Pro models with Cursor, Cline, OpenClaw, Aider, or any OpenAI-compatible tool.

### Works with OpenClaw, Cursor, Cline, Codex CLI, and More

aigate provides the OpenAI-compatible backend that tools like [OpenClaw](https://github.com/openclaw/openclaw) and coding IDEs need. Any tool that supports OpenAI or Anthropic APIs works out of the box:

[OpenClaw](https://github.com/openclaw/openclaw) • Cursor • Claude Code • Codex CLI • OpenCode • Aider • Cline • Roo Code • Kilo Code • Continue • Zed • Windsurf • OpenAI SDK • Anthropic SDK • LangChain • Obsidian

> **None of these tools support Kiro or Copilot Free natively.** aigate bridges the gap — it translates their standard API calls into each provider's internal format.

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

### 2. Login to providers

```bash
# Kiro (free Claude Sonnet 4.5, Haiku 4.5)
# Option A: Kiro IDE — just open and log in, credentials auto-detected
# Option B: kiro-cli
kiro-cli login

# GitHub Copilot (free GPT-4.1, Claude 3.5) — optional
./aigate --copilot-login
```

### 3. Start

```bash
export API_KEY="my-secret-key"    # You make this up — protects your proxy
./aigate
```

```
⚡ aigate v0.3.0
   ├─ Listening on http://0.0.0.0:8000
   ├─ OpenAI API:    /v1/chat/completions
   ├─ Anthropic API: /v1/messages
   └─ Models:        /v1/models
```

Keep this terminal open. Open a new terminal for the next step.

### 4. Use

```bash
# Kiro — Claude Sonnet 4.5
curl http://localhost:8000/v1/chat/completions \
  -H "Authorization: Bearer my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-sonnet-4-5","messages":[{"role":"user","content":"Hello!"}]}'

# GitHub Copilot — GPT-4.1 (prefix with copilot/)
curl http://localhost:8000/v1/chat/completions \
  -H "Authorization: Bearer my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"model":"copilot/gpt-4.1","messages":[{"role":"user","content":"Hello!"}]}'
```

---

## Use with OpenClaw

[OpenClaw](https://github.com/openclaw/openclaw) is a personal AI assistant (351k+ stars) that runs on your own devices and connects to WhatsApp, Telegram, Slack, Discord, and 20+ other channels. It needs an OpenAI-compatible API endpoint — aigate provides exactly that, for free.

### Setup

1. Start aigate (see Quick Start above)
2. Configure OpenClaw to use aigate as its model provider:

```json
{
  "agent": {
    "model": "claude-sonnet-4-5"
  },
  "providers": {
    "openai": {
      "baseURL": "http://localhost:8000/v1",
      "apiKey": "my-secret-key"
    }
  }
}
```

3. Run OpenClaw:
```bash
openclaw onboard --install-daemon
```

Now your OpenClaw assistant uses Claude Sonnet 4.5 (via Kiro) or GPT-4.1 (via Copilot) — completely free. Switch models by changing the `model` field.

### Why aigate + OpenClaw?

- **Free AI backend** — no Anthropic or OpenAI API key needed
- **Multiple models** — switch between Claude and GPT without changing providers
- **Co-host on same server** — aigate uses <10MB RAM, runs alongside OpenClaw on a Raspberry Pi or $5 VPS
- **Always-on** — aigate handles token refresh automatically, OpenClaw stays connected 24/7

---

## Multi-Provider Model Routing

aigate routes requests to the right provider based on model name prefix:

| Model | Provider | Description |
|-------|----------|-------------|
| `claude-sonnet-4-5` | Kiro | Claude Sonnet 4.5 — balanced |
| `claude-haiku-4-5` | Kiro | Claude Haiku 4.5 — fast |
| `claude-sonnet-4` | Kiro | Claude Sonnet 4 |
| `deepseek-v3.2` | Kiro | DeepSeek V3.2 — open MoE |
| `copilot/gpt-4.1` | GitHub Copilot | GPT-4.1 |
| `copilot/claude-3.5-sonnet` | GitHub Copilot | Claude 3.5 Sonnet |
| `copilot/o4-mini` | GitHub Copilot | OpenAI o4-mini |

> Use `copilot/` prefix for Copilot models. Everything else routes to Kiro.

---

## OpenAI-Compatible API

aigate exposes a fully OpenAI-compatible `/v1/chat/completions` endpoint. Any tool or SDK that works with OpenAI will work with aigate.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/chat/completions` | POST | Chat completions (streaming + non-streaming) |
| `/v1/models` | GET | List available models from all providers |
| `/health` | GET | Health check |

### Use with OpenAI Python SDK

```python
from openai import OpenAI

client = OpenAI(base_url="http://localhost:8000/v1", api_key="my-secret-key")

# Use Kiro (Claude Sonnet 4.5)
response = client.chat.completions.create(
    model="claude-sonnet-4-5",
    messages=[{"role": "user", "content": "Hello!"}],
)

# Use Copilot (GPT-4.1)
response = client.chat.completions.create(
    model="copilot/gpt-4.1",
    messages=[{"role": "user", "content": "Hello!"}],
)
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
Your Tool (OpenClaw, Cursor, Cline, SDK...)
    ↓ OpenAI/Anthropic API
  aigate (localhost:8000)
    ├─ claude-*     → Kiro → Claude Sonnet 4.5, Haiku 4.5
    └─ copilot/*    → GitHub Copilot → GPT-4.1, Claude 3.5
```

---

## GitHub Copilot Proxy — Use Copilot Free as an OpenAI API

aigate can also proxy GitHub Copilot's free tier (2000 requests/month) as an OpenAI-compatible API. This gives you access to GPT-4.1 and Claude 3.5 Sonnet for free.

```bash
# One-time login (opens browser for GitHub OAuth)
./aigate --copilot-login

# Then start normally — Copilot models appear with copilot/ prefix
./aigate
```

---

## Use with Cursor / Cline / Continue

Set in your IDE settings:

| Setting | Value |
|---------|-------|
| Base URL | `http://localhost:8000/v1` |
| API Key | your `API_KEY` value (e.g. `my-secret-key`) |
| Model | `claude-sonnet-4-5` or `copilot/gpt-4.1` |

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
{
  "provider": {
    "openai-compatible": {
      "apiKey": "my-secret-key",
      "baseURL": "http://localhost:8000/v1",
      "models": {
        "claude-sonnet-4-5": { "maxTokens": 16384 },
        "copilot/gpt-4.1": { "maxTokens": 16384 }
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
{
  "language_models": {
    "openai": {
      "api_url": "http://localhost:8000/v1",
      "available_models": [
        {"name": "claude-sonnet-4-5", "max_tokens": 16384},
        {"name": "copilot/gpt-4.1", "max_tokens": 16384}
      ]
    }
  }
}
```

---

## Download Options

| Platform | Architecture | Download |
|----------|-------------|----------|
| macOS | Apple Silicon (M1/M2/M3/M4) | [aigate-darwin-arm64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-arm64) |
| macOS | Intel | [aigate-darwin-amd64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-amd64) |
| Linux | x86_64 | [aigate-linux-amd64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-amd64) |
| Linux | ARM64 | [aigate-linux-arm64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-arm64) |
| Windows | x86_64 | [aigate-windows-amd64.exe](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-windows-amd64.exe) |

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

## Resource Usage

aigate is a pure proxy — it doesn't run AI models, doesn't cache responses, and doesn't need a database. Resource usage is minimal:

| Metric | Value |
|--------|-------|
| Binary size | ~10MB |
| RAM (idle) | ~5-8MB |
| RAM (100 concurrent streams) | ~20MB |
| CPU | Near zero (JSON parse + HTTP forward) |
| Disk | None (only a ~1KB token file) |

Runs comfortably on a Raspberry Pi, a $5 VPS, or alongside other services on the same machine. Ideal for co-hosting with [OpenClaw](https://github.com/openclaw/openclaw), OpenCode, Aider, or any OpenAI-compatible application — aigate adds virtually no overhead to your server.

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

<details>
<summary>Credential auto-detection</summary>

aigate checks these paths in order:

| Source | Path | Used by |
|--------|------|---------|
| kiro-cli SQLite | `~/.local/share/kiro-cli/data.sqlite3` | Linux / WSL / kiro-cli |
| amazon-q SQLite | `~/.local/share/amazon-q/data.sqlite3` | amazon-q-developer-cli |
| Kiro IDE JSON | `~/.aws/sso/cache/kiro-auth-token.json` | macOS / Kiro IDE |
| Copilot token | `~/.local/share/aigate/github_token` | `--copilot-login` |

</details>

---

## License

MIT
