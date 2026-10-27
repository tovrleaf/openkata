---
title: Infrastructure diagram
description: The AWS infrastructure behind openkata.dev and the MCP server.
order: 2
---

# Infrastructure Diagram

```
┌───────────────────────────────────────────────────────────────────┐
│                                                                   │
│                               USERS                               │
│                                                                   │
└───────────────┬───────────────────────────────────┬───────────────┘
                │                                   │
                │ browser                           │ MCP client
                │ openkata.dev                      │ openkata.dev/mcp
                ▼                                   ▼
┌───────────────────────────────────────────────────────────────────┐
│                             CloudFront                            │
│     CDN and HTTPS. Adds the CloudFront-Viewer-Country header.     │
│                                                                   │
│  /mcp*  ->  openkata-mcp       everything else  ->  openkata-web  │
└───────────────┬───────────────────────────────────┬───────────────┘
                │                                   │
                │ HTTPS                             │ HTTPS
                ▼                                   ▼
┌───────────────────────────────┐   ┌───────────────────────────────┐
│  Lambda: openkata-web         │   │  Lambda: openkata-mcp         │
│  (Function URL)               │   │  (Function URL)               │
│                               │   │                               │
│  Serves:                      │   │  Tools:                       │
│  - HTML pages                 │   │  - list_skills/rules/profiles │
│  - /X/archive                 │   │  - install_skill/rule/profile │
│  - static assets              │   │  - skill/rule/profile_versions│
└───┬────────────────┬──────────┘   └───┬────────────────┬──────────┘
    │                │                  │                │
    ▼                ▼                  ▼                ▼
┌────────┐ ┌──────────────────┐     ┌────────┐ ┌──────────────────┐
│   S3   │ │    DynamoDB      │     │   S3   │ │    DynamoDB      │
│        │ │                  │     │ (same  │ │ (same tables)    │
│openkata│ │openkata-downloads│     │ bucket)│ │                  │
│-arti-  │ │ (counters)       │     │        │ │                  │
│facts   │ │                  │     │        │ │                  │
│        │ │openkata-download-│     │        │ │                  │
│        │ │events (analytics)│     │        │ │                  │
└────────┘ └──────────────────┘     └────────┘ └──────────────────┘
```

## Flow

- **Web user** → `openkata.dev` → CloudFront → Lambda `openkata-web`
  → reads S3, writes DynamoDB
- **MCP client** → `openkata.dev/mcp` → CloudFront → Lambda
  `openkata-mcp` → reads S3, writes DynamoDB
- Both Lambdas share the same S3 bucket and DynamoDB tables

## Services

| Service | Resource | Purpose |
|---------|----------|---------|
| S3 | `openkata-artifacts` | Skill/rule/profile files, versions.json |
| DynamoDB | `openkata-downloads` | Per-artifact download counters |
| DynamoDB | `openkata-download-events` | Per-download event log (analytics) |
| Lambda | `openkata-web` | Website serving |
| Lambda | `openkata-mcp` | MCP server for agent clients (nine tools) |
| CloudFront | distribution | HTTPS and the viewer-country header; routes `/mcp*` to the MCP Lambda and everything else to the web Lambda |
| Route 53 | `openkata.dev` | DNS |
| ACM | certificate | TLS for openkata.dev |

## Notes

- CloudFront fronts both Lambdas. The distribution has the web
  Lambda's Function URL as its default origin and the MCP Lambda's
  Function URL as a second origin for `/mcp*`
  (`infra/add-mcp-to-cloudfront.sh`)
- CloudFront adds the `CloudFront-Viewer-Country` header. Only the
  web server reads it, for archive downloads. The MCP server
  records an empty country for every install
- Both Lambdas run in `eu-north-1` (Stockholm), arm64, 128MB
