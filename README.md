# headortails-mcp

Secure gateway for connecting ChatGPT-compatible MCP clients to the official Luno MCP server with **full Luno API permissions**.

## Architecture

ChatGPT / MCP client → HTTPS gateway →  `luno-mcp` → Luno API

The Luno API secret stays on the server. It is never placed in this repository or sent to ChatGPT.

## Full permissions

The gateway runs the official Luno MCP with:

`ALLOW_WRITE_OPERATIONS=true`

This exposes the official MCP's write capabilities, subject to the permissions granted to the Luno API key itself. These include operations such as creating/cancelling orders and conversions, alongside balances, transactions, markets and order data.

**Important:** the Luno API key must itself be configured with every permission you want to use. This gateway cannot grant permissions that the Luno key does not have.

## Environment variables

```env
LUNO_API_KEY_ID=
LUNO_API_SECRET=
MCP_AUTH_TOKEN=
LUNO_API_DOMAIN=api.luno.com
ALLOW_WRITE_OPERATIONS=true
```

Never commit `.env` or the Luno secret.

## Run

Build:

```bash
docker build -t headortails-mcp .
```

Run:

```bash
docker run --rm -p 8080:8080 \
  -e LUNO_API_KEY_ID="$LUNO_API_KEY_ID" \
  -e LUNO_API_SECRET="$LUNO_API_SECRET" \
  -e MCP_AUTH_TOKEN="$MCP_AUTH_TOKEN" \
  -e LUNO_API_DOMAIN="api.luno.com" \
  -e ALLOW_WRITE_OPERATIONS="true" \
  headortails-mcp
```

Health check:

`GET /healthz`

MCP endpoint:

`/`

The gateway requires `Authorization: Bearer <MCP_AUTH_TOKEN>` and removes that gateway credential before forwarding the request to the internal Luno MCP server.
This project wraps the official Luno MCP server rather than reimplementing the Luno API.
