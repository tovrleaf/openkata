package analytics

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TableName = "openkata-download-events"

type Event struct {
	Artifact  string
	Version   string
	Source    string
	Client    string
	Country   string
	Timestamp string
}

type DynamoPutter interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func RecordDownload(ctx context.Context, client DynamoPutter, event Event) {
	if client == nil {
		return
	}
	if event.Timestamp == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	}
	table := TableName
	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &table,
		Item: map[string]types.AttributeValue{
			"artifact":  &types.AttributeValueMemberS{Value: event.Artifact},
			"timestamp": &types.AttributeValueMemberS{Value: event.Timestamp},
			"version":   &types.AttributeValueMemberS{Value: event.Version},
			"source":    &types.AttributeValueMemberS{Value: event.Source},
			"client":    &types.AttributeValueMemberS{Value: event.Client},
			"country":   &types.AttributeValueMemberS{Value: event.Country},
		},
	})
	if err != nil {
		log.Printf("analytics: record download: %v", err)
	}
}

func ParseClient(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "claude"):
		return "Claude-Desktop"
	case strings.Contains(ua, "cursor"):
		return "Cursor"
	case strings.Contains(ua, "kiro"):
		return "Kiro"
	case strings.Contains(ua, "curl"):
		return "curl"
	case strings.Contains(ua, "mozilla"):
		return "browser"
	default:
		return "unknown"
	}
}
