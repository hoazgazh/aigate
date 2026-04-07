<div align="center">

# ⚡ aigate

**Universal AI Gateway — Use Claude, GPT & more through your existing subscriptions**

One binary. Zero config. Works everywhere.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/hoazgazh/aigate)](https://github.com/hoazgazh/aigate/releases/latest)

</div>

---

## What is this?

aigate is a lightweight proxy that exposes OpenAI and Anthropic-compatible APIs from various AI backends. Download a single binary, point your tools at it, done.

**Supported backends:**
- ✅ Kiro (free Claude Sonnet 4.5, Haiku 4.5, Sonnet 4)
- 🔜 AWS Bedrock
- 🔜 More coming...

**Works with:** Cursor • Claude Code • Cline • Roo Code • Continue • OpenAI SDK • LangChain • any OpenAI/Anthropic-compatible tool

## Quick Start

### 1. Download

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-arm64 -o aigate && chmod +x aigate

# macOS Intel
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-amd64 -o aigate && chmod +x aigate

# Linux x86_64
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-amd64 -o aigate && chmod +x aigate
```

### 2. Configure

You need a Kiro account. Log in to [Kiro IDE](https://kiro.dev/) or run `kiro-cli login` first.

```bash
export API_KEY="my-secret-key"    # You make this up — protects your proxy
```

aigate auto-detects your credentials from these locations (in order):

| Source | Path | Used by |
|--------|------|---------|
| kiro-cli SQLite | `~/.local/share/kiro-cli/data.sqlite3` | Linux / WSL / kiro-cli users |
| amazon-q SQLite | `~/.local/share/amazon-q/data.sqlite3` | amazon-q-developer-cli users |
| Kiro IDE JSON | `~/.aws/sso/cache/kiro-auth-token.json` | macOS / Kiro IDE users |

**No extra config needed** if you've already logged in with `kiro-cli login` or Kiro IDE.

To override auto-detection:
```bash
# Linux / kiro-cli
export KIRO_CLI_DB_FILE="~/.local/share/kiro-cli/data.sqlite3"

# macOS / Kiro IDE
export KIRO_CREDS_FILE="~/.aws/sso/cache/kiro-auth-token.json"
```

### 3. Start

```bash
./aigate
```

Server starts at `http://localhost:8000`. Keep this terminal open.

### 4. Use (in another terminal)

```bash
# Non-streaming (single JSON response)
curl http://localhost:8000/v1/chat/completions \
  -H "Authorization: Bearer my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-sonnet-4-5","messages":[{"role":"user","content":"Hello!"}]}'

# Streaming (real-time token-by-token)
curl http://localhost:8000/v1/chat/completions \
  -H "Authorization: Bearer my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-sonnet-4-5","messages":[{"role":"user","content":"Hello!"}],"stream":true}'
```

Or use with any OpenAI-compatible tool — set base URL to `http://localhost:8000/v1` and API key to your `API_KEY`.

## Download Options

| Platform | Architecture | Download |
|----------|-------------|----------|
| macOS | Apple Silicon (M1/M2/M3/M4) | [aigate-darwin-arm64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-arm64) |
| macOS | Intel | [aigate-darwin-amd64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-amd64) |
| Linux | x86_64 | [aigate-linux-amd64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-amd64) |
| Linux | ARM64 | [aigate-linux-arm64](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-linux-arm64) |
| Windows | x86_64 | [aigate-windows-amd64.exe](https://github.com/hoazgazh/aigate/releases/latest/download/aigate-windows-amd64.exe) |

## Build from source

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
  -v ~/.aws/sso/cache:/root/.aws/sso/cache:ro \
  -e KIRO_CREDS_FILE=/root/.aws/sso/cache/kiro-auth-token.json \
  ghcr.io/hoazgazh/aigate
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Health check |
| `GET /v1/models` | List available models |
| `POST /v1/chat/completions` | OpenAI Chat Completions API |
| `POST /v1/messages` | Anthropic Messages API |

## Configuration

| Env Variable | Required | Description |
|-------------|----------|-------------|
| `API_KEY` | ✅ | Password to protect your proxy (you make this up) |
| `KIRO_CREDS_FILE` | One of these | Path to Kiro IDE JSON credentials |
| `KIRO_CLI_DB_FILE` | | Path to kiro-cli SQLite database |
| `REFRESH_TOKEN` | | Kiro refresh token (manual) |
| `KIRO_REGION` | | AWS region (default: `us-east-1`) |
| `PORT` | | Server port (default: `8000`) |
| `HOST` | | Server host (default: `0.0.0.0`) |

## Usage Examples

### Python (OpenAI SDK)

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

### Python (Anthropic SDK)

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

### Cursor / Cline / Continue

Set in your IDE settings:
- **Base URL:** `http://localhost:8000/v1`
- **API Key:** your `API_KEY` value
- **Model:** `claude-sonnet-4-5`

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

## License

MIT
