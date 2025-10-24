# Chefly - AI Powered Recipe Generator

[![License](https://img.shields.io/github/license/voidquark/chefly)](LICENSE)

Chefly is a high-quality recipe generator powered by Claude AI.

<p align="center">
  <img src="assets/1.png" alt="Generate Recipe Page" width="400"/>
  <img src="assets/2.png" alt="My Recipes Page" width="400"/>
  <img src="assets/3.png" alt="Recipes Example" width="400"/>
</p>

**Features:**
  - Multi user support, each user can maintain their own recipe list
  -	Mobile first UI design
  - Share recipes via public link
  - Mark favorite recipes, filter them, or delete unwanted ones
  - Generate recipes by meat type, cuisine, dietary preferences, difficulty, and preparation time
  - Admin panel to manage registered users and set recipe generation limits per user
  - Generate shopping lists from recipes to easily bookmark ingredients
  - Audit logging of every request in json or pretty format
  - Supported `sk`, `en` language

> [!IMPORTANT]
> This is a **vibe coded** project created to explore what Claude AI can do, how to interact with such a model, and how to plan development around it for project in 1day.
It was a self-exploration project with a single purpose to build something for my own personal use.
I’m sharing it here in case someone wants to reuse it, fork it, extend it, or add new features.
Breaking changes may occur at any time and without notice.
**Long term maintenance is uncertain**. Anyone can open an issue for feature requests or bug fixes, but there is no guarantee that I will implement them, please don’t take it personally. There will be no human interaction with this project at the code level by design.

## Stack Information

**Backend:** `Go`

**Frontend:** `Vite`, `React`, `TypeScript`, `Tailwind CSS`

**Storage:** `SQLite`

## AI Models

**Models:** Recipe generation uses Claude SDK (tested: Haiku & Sonnet). Images are generated with OpenAI DALL·E 3 at standard quality to reduce cost.

**API Keys:** Requires two keys one for Claude and one for OpenAI. Recommended model for recipes: `claude-sonnet-4-20250514` (better results than `claude-3-haiku-20240307`).

## Quick Start

1. Obtain an API key from the Anthropic console for Claude AI.
2. Obtain an API key from the OpenAI console for DALL·E 3
3. Create folders for data storage (e.g. in the current directory: `mkdir chefly_data chefly_uploads`).
4. Run the container as a regular user in a user namespace:

```bash
podman run -d \
  --name chefly \
  --restart unless-stopped \
  --userns=keep-id:uid=1000,gid=1000 \
  -p 8080:8080 \
  -v ./chefly_data:/app/data:Z \
  -v ./chefly_uploads:/app/uploads:Z \
  -e JWT_SECRET="your-minimum-32-characters-signing-key" \
  -e CLAUDE_API_KEY="sk-ant..." \
  -e CLAUDE_MODEL="claude-sonnet-4-20250514" \
  -e OPENAI_API_KEY="sk-..." \
  -e OPENAI_MODEL="dall-e-3" \
  -e REGISTRATION_ENABLED=true \
  -e RECIPE_GENERATION_LIMIT="unlimited" \
  -e AUDIT_LOG_ENABLED=true \
  -e AUDIT_LOG_LEVEL="info" \
  -e AUDIT_LOG_FORMAT="json" \
  --security-opt=no-new-privileges \
  --cap-drop=ALL \
  voidquark/chefly:latest
```

5. The first registered user becomes the admin. This cannot be changed without directly modifying the database.

## Configuration

Environment variables:

| Variable | Description | Example Value
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing key (min 32 characters). You can generate one using `openssl rand -base64 32` | `9+RxeeHYEKAcpXbaVNy5YIU/Qk5Lr/uJ2J1tP16GayA=` |
| `CLAUDE_API_KEY` | Anthropic Claude API key | `sk-ant...` |
| `CLAUDE_MODEL` | Claude AI model to use | `claude-sonnet-4-20250514` |
| `OPENAI_API_KEY` | OpenAI API key for image generation | `sk-...` |
| `OPENAI_MODEL` | OpenAI model to use | `dall-e-3` |
| `REGISTRATION_ENABLED` | Enable registration (`true`/`false`) | `true` |
| `RECIPE_GENERATION_LIMIT` | Global recipe generation limit per user (`unlimited`, `0`, `5`, etc.), overridable per user by the admin in the admin panel. | `unlimited` |
| `AUDIT_LOG_ENABLED` | Enable audit logging | `true` |
| `AUDIT_LOG_LEVEL` | Audit log level (`debug`, `info`, `warn`, `error`) | `info` |
| `AUDIT_LOG_FORMAT` | Log format (`json` or `pretty`) | `json` |

> [!NOTE]
> The container uses the `chefly` user (`UID 1000`, `GID 1000`) inside.
If you experience permission issues with bind mounts, it’s likely due to this user mapping.

## Audit Logging

**Example JSON Log**:
```json
{
  "timestamp":"2025-01-15T14:30:45.123Z",
  "level":"warn",
  "event_type":"auth.login.failed",
  "message":"Login failed: user not found",
  "context":{
    "email_attempted":"unknown@example.com",
    "ip_address":"192.168.30.100",
    "request_id":"req-abc123",
    "status_code":401,
    "metadata":{"failure_reason":"user_not_found"}
  }
}
```

## License

MIT License
