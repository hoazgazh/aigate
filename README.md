<div align="center">

# ⚡ aigate

**Universal AI Gateway — Use Claude, GPT & more through your existing subscriptions**

One binary. Zero config. Works everywhere.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

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

```bash
# Download (macOS Apple Silicon)
curl -fsSL https://github.com/hoazgazh/aigate/releases/latest/download/aigate-darwin-arm64 -o aigate
chmod +x aigate

# Run
export API_KEY="my-secret-key"
export KIRO_CREDS_FILE="~/.aws/sso/cache/kiro-auth-token.json"
./aigate
```

API available at `http://localhost:8000`

### Build from source

```bash
git clone https://github.com/hoazgazh/aigate.git
cd aigate
make build
API_KEY=my-secret-key ./bin/aigate
```

### Docker

```bash
docker run -p 8000:8000 -e API_KEY=my-secret-key ghcr.io/hoazgazh/aigate
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Health check |
| `GET /v1/models` | List models |
| `POST /v1/chat/completions` | OpenAI Chat API |
| `POST /v1/messages` | Anthropic Messages API |

## Usage

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

## License

MIT
