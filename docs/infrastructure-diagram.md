# Infrastructure Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                           USERS                                  │
└──────────┬──────────────────────────────────────┬────────────────┘
           │                                      │
           │ browser                              │ MCP client
           │ openkata.dev                         │ (Claude, Cursor, etc.)
           ▼                                      ▼
┌─────────────────────┐              ┌─────────────────────────┐
│     CloudFront      │              │   Lambda Function URL   │
│  (CDN + geo headers)│              │   (direct, no CDN)      │
│                     │              │                         │
│  Adds:              │              │                         │
│  CloudFront-Viewer- │              │                         │
│  Country header     │              │                         │
└─────────┬───────────┘              └────────────┬────────────┘
          │                                       │
          │ HTTPS                                  │ HTTPS
          ▼                                       ▼
┌─────────────────────┐              ┌─────────────────────────┐
│  Lambda: openkata-web│              │  Lambda: openkata-mcp   │
│  (Function URL)     │              │                         │
│                     │              │  Tools:                 │
│  Serves:            │              │  - list-skills/rules    │
│  - HTML pages       │              │  - install-skill/rule   │
│  - /X/archive       │              │  - skill/rule-versions  │
│  - static assets    │              │                         │
└───┬─────────┬───────┘              └───┬─────────┬───────────┘
    │         │                          │         │
    ▼         ▼                          ▼         ▼
┌────────┐ ┌──────────────────┐    ┌────────┐ ┌──────────────────┐
│   S3   │ │    DynamoDB      │    │   S3   │ │    DynamoDB      │
│        │ │                  │    │ (same  │ │  (same tables)   │
│openkata│ │openkata-downloads│    │ bucket)│ │                  │
│-arti-  │ │ (counters)       │    │        │ │                  │
│facts   │ │                  │    │        │ │                  │
│        │ │openkata-download-│    │        │ │                  │
│        │ │events (analytics)│    │        │ │                  │
└────────┘ └──────────────────┘    └────────┘ └──────────────────┘
```

## Flow

- **Web user** → `openkata.dev` → CloudFront → Lambda
  `openkata-web` → reads S3, writes DynamoDB
- **MCP client** → Lambda Function URL → Lambda
  `openkata-mcp` → reads S3, writes DynamoDB
- Both Lambdas share the same S3 bucket and DynamoDB tables

## Services

| Service | Resource | Purpose |
|---------|----------|---------|
| S3 | `openkata-artifacts` | Skill/rule/profile files, versions.json |
| DynamoDB | `openkata-downloads` | Per-artifact download counters |
| DynamoDB | `openkata-download-events` | Per-download event log (analytics) |
| Lambda | `openkata-web` | Website serving |
| Lambda | `openkata-mcp` | MCP server for agent clients |
| CloudFront | distribution | CDN, HTTPS, geo headers for web |
| Route 53 | `openkata.dev` | DNS |
| ACM | certificate | TLS for openkata.dev |

## Notes

- CloudFront only fronts the web server (provides
  `CloudFront-Viewer-Country` header for analytics)
- MCP has no CDN — direct Lambda Function URL
  (no country detection for MCP downloads)
- Both Lambdas run in `eu-north-1`, arm64, 128MB
- Region: `eu-north-1` (Stockholm)
