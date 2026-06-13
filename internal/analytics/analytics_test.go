package analytics

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockPutter struct {
	called bool
	input  *dynamodb.PutItemInput
}

func (m *mockPutter) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	m.called = true
	m.input = params
	return &dynamodb.PutItemOutput{}, nil
}

func TestRecordDownload(t *testing.T) {
	t.Run("nil client no-ops", func(t *testing.T) {
		RecordDownload(context.Background(), nil, Event{Artifact: "skills/test"})
	})

	t.Run("writes event to table", func(t *testing.T) {
		m := &mockPutter{}
		RecordDownload(context.Background(), m, Event{
			Artifact: "skills/create-adr",
			Version:  "1.2.0",
			Source:   "web",
			Client:   "browser",
			Country:  "FI",
		})
		if !m.called {
			t.Error("RecordDownload() did not call PutItem")
		}
		if *m.input.TableName != TableName {
			t.Errorf("RecordDownload() table = %q, want %q", *m.input.TableName, TableName)
		}
	})
}

func TestParseClient(t *testing.T) {
	tests := []struct {
		userAgent string
		want      string
	}{
		{"Claude-Desktop/1.0", "Claude-Desktop"},
		{"Cursor/0.45.2", "Cursor"},
		{"Kiro/1.2.0", "Kiro"},
		{"curl/8.1.0", "curl"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X)", "browser"},
		{"", "unknown"},
		{"some-random-agent", "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.userAgent, func(t *testing.T) {
			got := ParseClient(tt.userAgent)
			if got != tt.want {
				t.Errorf("ParseClient(%q) = %q, want %q", tt.userAgent, got, tt.want)
			}
		})
	}
}
