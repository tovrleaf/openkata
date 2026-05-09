---
spec: 0002-infrastructure
---

# Design

## Static Asset Embedding

Go's `go:embed` only works with paths relative to the source
file. Static assets live in `web/static/` but the binary
source is in `cmd/openkata-web/`.

**Decision:** Add an `embed.go` file inside `web/static/`
that exports the embedded filesystem as a package. The web
server imports this package. Pattern taken from the warehouse
project's `assets/embed.go`.

```text
web/static/
├── embed.go         # package static; //go:embed all:css all:js
├── css/
│   └── style.css
└── js/
    └── htmx.min.js
```

`cmd/openkata-web/` imports
`github.com/tovrleaf/openkata/web/static` and uses
`static.FS` in production. In dev mode, it serves from the
filesystem directly for instant CSS changes without rebuild.

## Lambda Adapter

The binary detects Lambda via `AWS_LAMBDA_FUNCTION_NAME` env
var (set automatically by the runtime). When present, it
starts the Lambda handler using `aws-lambda-go-api-proxy`
which wraps the existing `http.ServeMux`. No code changes
to handlers or templates.

```go
if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
    lambda.Start(httpadapter.NewV2(mux).ProxyWithContext)
} else {
    http.ListenAndServe(addr, mux)
}
```

## CDK Stack

Single stack in `infra/`:

```text
infra/
├── go.mod           # Separate module for CDK deps
├── infra.go         # Stack definition
└── infra_test.go    # Optional: snapshot test
```

Separate `go.mod` because CDK pulls in many dependencies
that the main binary doesn't need.

Resources:
- Lambda function (arm64, provided.al2023, 128MB, 10s timeout)
- Function URL (auth type: NONE for public access)
- IAM role with basic Lambda execution policy

## Deploy Flow

`make deploy` runs `scripts/deploy.sh`:

1. `templ generate` templates
2. `GOOS=linux GOARCH=arm64 go build -o infra/lambda/bootstrap ./cmd/openkata-web/`
3. `cd infra && cdk deploy --require-approval never`

CDK references the built binary from `infra/lambda/bootstrap`
as an asset.

## Air Config Update

Air watches `cmd/openkata-web/` (including `static/`) and
excludes `_templ.go` files. No `web/` directory needed.
