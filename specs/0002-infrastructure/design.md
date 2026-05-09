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

## CloudFormation Stack

Single template in `infra/`:

```text
infra/
└── template.yaml    # CloudFormation template
```

No extra tooling — just AWS CLI. State managed by AWS.

Resources:
- Lambda function (arm64, provided.al2023, 128MB, 10s timeout)
- Function URL (auth type: NONE for public access)
- IAM role with basic Lambda execution policy
- S3 bucket for deployment artifacts (Lambda zip)

## Deploy Flow

`make deploy` runs `scripts/deploy.sh`:

1. `templ generate` templates
2. `GOOS=linux GOARCH=arm64 go build -o bootstrap ./cmd/openkata-web/`
3. `zip deploy.zip bootstrap`
4. Upload zip to S3
5. `aws cloudformation deploy --template-file infra/template.yaml`

## Air Config Update

Air watches `cmd/openkata-web/` and `web/static/` and
excludes `_templ.go` files.
