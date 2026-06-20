package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tovrleaf/openkata/internal/analytics"
)

const (
	statsDir       = ".local/stats"
	cursorFile     = ".local/stats/cursor.json"
	eventsFile     = ".local/stats/download-events.json"
	metricsFile    = ".local/stats/page-metrics.json"
	pathsFile      = ".local/stats/page-paths.json"
	lambdaFunc     = "openkata-web"
	logGroup       = "/aws/lambda/openkata-web"
)

type cursor struct {
	Downloads string `json:"downloads"`
	Metrics   string `json:"metrics"`
	Paths     string `json:"paths"`
}

type downloadEvent struct {
	Artifact  string `json:"artifact"`
	Version   string `json:"version"`
	Source    string `json:"source"`
	Client    string `json:"client"`
	Country   string `json:"country"`
	Timestamp string `json:"timestamp"`
}

type pageMetric struct {
	Date        string `json:"date"`
	Invocations int    `json:"invocations"`
}

type pagePath struct {
	Date  string `json:"date"`
	Path  string `json:"path"`
	Count int    `json:"count"`
	Bot   bool   `json:"bot"`
}

func main() {
	since := flag.String("since", "", "Start date YYYY-MM-DD (default: 30 days ago)")
	flag.Parse()

	if err := os.MkdirAll(statsDir, 0o755); err != nil {
		log.Fatalf("create stats dir: %v", err)
	}

	defaultSince := time.Now().AddDate(0, 0, -30).UTC().Format(time.RFC3339)
	if *since != "" {
		t, err := time.Parse("2006-01-02", *since)
		if err != nil {
			log.Fatalf("invalid --since: %v", err)
		}
		defaultSince = t.UTC().Format(time.RFC3339)
	}

	cur := loadCursor(defaultSince)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatalf("load AWS config: %v", err)
	}

	fetchDownloads(ctx, cfg, &cur)
	fetchMetrics(ctx, cfg, &cur)
	fetchPaths(ctx, cfg, &cur)
	saveCursor(cur)

	fmt.Println("stats-fetch complete")
}

func loadCursor(defaultSince string) cursor {
	data, err := os.ReadFile(cursorFile)
	if err != nil {
		return cursor{
			Downloads: defaultSince,
			Metrics:   defaultSince,
			Paths:     defaultSince,
		}
	}
	var c cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return cursor{Downloads: defaultSince, Metrics: defaultSince, Paths: defaultSince}
	}
	return c
}

func saveCursor(c cursor) {
	data, _ := json.MarshalIndent(c, "", "  ")
	if err := os.WriteFile(cursorFile, data, 0o644); err != nil {
		log.Printf("write cursor: %v", err)
	}
}

func fetchDownloads(ctx context.Context, cfg aws.Config, cur *cursor) {
	client := dynamodb.NewFromConfig(cfg)
	table := analytics.TableName
	filterExpr := "#ts > :since"
	input := &dynamodb.ScanInput{
		TableName:        &table,
		FilterExpression: &filterExpr,
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":since": &types.AttributeValueMemberS{Value: cur.Downloads},
		},
	}

	var newEvents []downloadEvent
	var lastTS string

	for {
		resp, err := client.Scan(ctx, input)
		if err != nil {
			log.Printf("scan downloads: %v", err)
			break
		}
		for _, item := range resp.Items {
			ev := downloadEvent{
				Artifact:  attrStr(item, "artifact"),
				Version:   attrStr(item, "version"),
				Source:    attrStr(item, "source"),
				Client:    attrStr(item, "client"),
				Country:   attrStr(item, "country"),
				Timestamp: attrStr(item, "timestamp"),
			}
			newEvents = append(newEvents, ev)
			if ev.Timestamp > lastTS {
				lastTS = ev.Timestamp
			}
		}
		if resp.LastEvaluatedKey == nil {
			break
		}
		input.ExclusiveStartKey = resp.LastEvaluatedKey
	}

	if len(newEvents) > 0 {
		existing := loadJSON[[]downloadEvent](eventsFile)
		existing = append(existing, newEvents...)
		writeJSON(eventsFile, existing)
		cur.Downloads = lastTS
		log.Printf("downloads: %d new events", len(newEvents))
	} else {
		log.Println("downloads: no new events")
	}
}

func fetchMetrics(ctx context.Context, cfg aws.Config, cur *cursor) {
	client := cloudwatch.NewFromConfig(cfg)

	start, err := time.Parse(time.RFC3339, cur.Metrics)
	if err != nil {
		log.Printf("parse metrics cursor: %v", err)
		return
	}
	end := time.Now().UTC()

	resp, err := client.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/Lambda"),
		MetricName: aws.String("Invocations"),
		Dimensions: []cloudwatchtypes.Dimension{
			{Name: aws.String("FunctionName"), Value: aws.String(lambdaFunc)},
		},
		StartTime:  &start,
		EndTime:    &end,
		Period:     aws.Int32(86400),
		Statistics: []cloudwatchtypes.Statistic{cloudwatchtypes.StatisticSum},
	})
	if err != nil {
		log.Printf("fetch metrics: %v", err)
		return
	}

	// Build map of new datapoints
	newData := make(map[string]int)
	for _, dp := range resp.Datapoints {
		date := dp.Timestamp.Format("2006-01-02")
		newData[date] = int(*dp.Sum)
	}

	// Merge with existing
	existing := loadJSON[[]pageMetric](metricsFile)
	byDate := make(map[string]int)
	for _, m := range existing {
		byDate[m.Date] = m.Invocations
	}
	for date, count := range newData {
		byDate[date] = count
	}

	var merged []pageMetric
	for date, count := range byDate {
		merged = append(merged, pageMetric{Date: date, Invocations: count})
	}
	sortByDate(merged)
	writeJSON(metricsFile, merged)
	cur.Metrics = end.Format(time.RFC3339)
	log.Printf("metrics: %d datapoints", len(resp.Datapoints))
}

func fetchPaths(ctx context.Context, cfg aws.Config, cur *cursor) {
	client := cloudwatchlogs.NewFromConfig(cfg)

	start, err := time.Parse(time.RFC3339, cur.Paths)
	if err != nil {
		log.Printf("parse paths cursor: %v", err)
		return
	}
	end := time.Now().UTC()

	query := `filter @message like /"req":"page"/
| parse @message /"path":"(?<path>[^"]+)"/
| parse @message /"ua":"(?<ua>[^"]+)"/
| stats count() as cnt by path, ua
| sort cnt desc
| limit 500`

	startEpoch := start.Unix()
	endEpoch := end.Unix()
	lg := logGroup

	startResp, err := client.StartQuery(ctx, &cloudwatchlogs.StartQueryInput{
		LogGroupName: &lg,
		StartTime:    &startEpoch,
		EndTime:      &endEpoch,
		QueryString:  &query,
	})
	if err != nil {
		log.Printf("start logs query: %v", err)
		writeJSON(pathsFile, []pagePath{})
		return
	}

	// Poll for results
	type pathKey struct {
		path string
		bot  bool
	}
	agg := make(map[pathKey]int)
	for {
		time.Sleep(time.Second)
		resp, err := client.GetQueryResults(ctx, &cloudwatchlogs.GetQueryResultsInput{
			QueryId: startResp.QueryId,
		})
		if err != nil {
			log.Printf("get query results: %v", err)
			break
		}
		if resp.Status == "Complete" || resp.Status == "Failed" || resp.Status == "Cancelled" {
			if resp.Status != "Complete" {
				log.Printf("logs query status: %s", resp.Status)
				break
			}
			for _, row := range resp.Results {
				var path, ua string
				var count int
				for _, field := range row {
					if field.Field != nil && field.Value != nil {
						switch *field.Field {
						case "path":
							path = *field.Value
						case "ua":
							ua = *field.Value
						case "cnt":
							fmt.Sscanf(*field.Value, "%d", &count)
						}
					}
				}
				if path != "" {
					agg[pathKey{path: path, bot: isBot(ua)}] += count
				}
			}
			break
		}
	}

	today := time.Now().UTC().Format("2006-01-02")
	var results []pagePath
	for k, count := range agg {
		results = append(results, pagePath{Date: today, Path: k.path, Count: count, Bot: k.bot})
	}

	writeJSON(pathsFile, results)
	cur.Paths = end.Format(time.RFC3339)
	log.Printf("paths: %d entries", len(results))
}

func isBot(ua string) bool {
	ua = strings.ToLower(ua)
	for _, pattern := range []string{
		"bot", "spider", "crawl", "slurp", "mediapartners",
		"lighthouse", "pagespeed", "pingdom", "uptimerobot",
		"headlesschrome", "phantomjs", "python-requests",
		"go-http-client", "wget", "curl",
	} {
		if strings.Contains(ua, pattern) {
			return true
		}
	}
	return false
}

func attrStr(item map[string]types.AttributeValue, key string) string {
	if v, ok := item[key].(*types.AttributeValueMemberS); ok {
		return v.Value
	}
	return ""
}

func loadJSON[T any](path string) T {
	var zero T
	data, err := os.ReadFile(path)
	if err != nil {
		return zero
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return zero
	}
	return v
}

func writeJSON(path string, v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Printf("write %s: %v", path, err)
	}
}

func sortByDate(metrics []pageMetric) {
	for i := 0; i < len(metrics); i++ {
		for j := i + 1; j < len(metrics); j++ {
			if metrics[i].Date > metrics[j].Date {
				metrics[i], metrics[j] = metrics[j], metrics[i]
			}
		}
	}
}
