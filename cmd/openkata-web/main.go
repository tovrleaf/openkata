package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/tovrleaf/openkata/web/static"
)

func main() {
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
