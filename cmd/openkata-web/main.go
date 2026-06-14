package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/tovrleaf/openkata/web/static"
)

var (
	s3Client *s3.Client
	dbClient *dynamodb.Client
	bucket   string
	table    string
)

func main() {
	bucket = os.Getenv("OPENKATA_BUCKET")
	if bucket == "" {
		bucket = "openkata-artifacts"
	}
	table = os.Getenv("OPENKATA_TABLE")
	if table == "" {
		table = "openkata-downloads"
	}

	// Init AWS clients
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "aws config: %v\n", err)
	} else {
		s3Client = s3.NewFromConfig(cfg)
		dbClient = dynamodb.NewFromConfig(cfg)
	}

	mux := http.NewServeMux()

	// Static files: embedded in prod, filesystem in dev
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		mux.Handle("/static/", http.StripPrefix("/static/",
			http.FileServerFS(static.FS)))
	} else {
		mux.Handle("/static/", http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static"))))
	}

	// Routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/catalog/", handleCatalog)
	mux.HandleFunc("/skills/", handleSkills)
	mux.HandleFunc("/rules/", handleRules)
	mux.HandleFunc("/profiles/", handleProfiles)
	mux.HandleFunc("/getting-started/", handleGettingStarted)
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		mux.HandleFunc("/design-system/", handleDesignSystem)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(httpadapter.NewV2(mux).ProxyWithContext)
		return
	}

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	fmt.Printf("Open Kata web server listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
