package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const source = "github.com/tovrleaf/openkata"

var (
	bucket   string
	table    string
	s3Client *s3.Client
	dbClient *dynamodb.Client
	versions *versionsFile
)

type artifactInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	Tags        string `json:"tags,omitempty"`
}

type versionsFile struct {
	Skills map[string]artifactInfo `json:"skills"`
	Rules  map[string]artifactInfo `json:"rules"`
}

func main() {
	bucket = os.Getenv("OPENKATA_BUCKET")
	if bucket == "" {
		bucket = "openkata-artifacts"
	}
	table = os.Getenv("OPENKATA_TABLE")
	if table == "" {
		table = "openkata-downloads"
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "aws config: %v\n", err)
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
	dbClient = dynamodb.NewFromConfig(cfg)

	// Load versions.json from S3
	versions, err = loadVersions(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load versions: %v\n", err)
		versions = &versionsFile{
			Skills: make(map[string]artifactInfo),
			Rules:  make(map[string]artifactInfo),
		}
	}

	s := server.NewMCPServer("openkata", "0.3.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(listSkillsTool(), listSkillsHandler)
	s.AddTool(listRulesTool(), listRulesHandler)
	s.AddTool(installSkillTool(), installSkillHandler)
	s.AddTool(installRuleTool(), installRuleHandler)
	s.AddTool(skillVersionsTool(), skillVersionsHandler)
	s.AddTool(ruleVersionsTool(), ruleVersionsHandler)

	httpServer := server.NewStreamableHTTPServer(s,
		server.WithStateLess(true),
	)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(httpadapter.NewV2(httpServer).ProxyWithContext)
	} else {
		addr := os.Getenv("OPENKATA_ADDR")
		if addr == "" {
			addr = ":8081"
		}
		fmt.Fprintf(os.Stderr, "listening on %s\n", addr)
		if err := httpServer.Start(addr); err != nil {
			fmt.Fprintf(os.Stderr, "server: %v\n", err)
			os.Exit(1)
		}
	}
}

// --- Tools ---

func listSkillsTool() mcp.Tool {
	return mcp.NewTool("list_skills",
		mcp.WithDescription("List available OpenKata skills with descriptions, versions, tags, and download counts"),
		mcp.WithString("tag", mcp.Description("Filter by tag (comma-separated tags in skill metadata)")),
	)
}

func listRulesTool() mcp.Tool {
	return mcp.NewTool("list_rules",
		mcp.WithDescription("List available OpenKata rules with descriptions, versions, tags, and download counts"),
		mcp.WithString("tag", mcp.Description("Filter by tag (comma-separated tags in rule metadata)")),
	)
}

func installSkillTool() mcp.Tool {
	return mcp.NewTool("install_skill",
		mcp.WithDescription("Install an OpenKata skill. Returns all files and a .manifest.json. Write files to .agents/skills/<name>/ in your project."),
		mcp.WithString("skill", mcp.Required(), mcp.Description("Name of the skill to install")),
		mcp.WithString("version", mcp.Description("Version to install (default: latest)")),
	)
}

func installRuleTool() mcp.Tool {
	return mcp.NewTool("install_rule",
		mcp.WithDescription("Install an OpenKata rule. Returns all files and a .manifest.json. Write files to .agents/rules/<name>/ in your project."),
		mcp.WithString("rule", mcp.Required(), mcp.Description("Name of the rule to install")),
		mcp.WithString("version", mcp.Description("Version to install (default: latest)")),
	)
}

func skillVersionsTool() mcp.Tool {
	return mcp.NewTool("skill_versions",
		mcp.WithDescription("List all available versions of a skill"),
		mcp.WithString("skill", mcp.Required(), mcp.Description("Name of the skill")),
	)
}

func ruleVersionsTool() mcp.Tool {
	return mcp.NewTool("rule_versions",
		mcp.WithDescription("List all available versions of a rule"),
		mcp.WithString("rule", mcp.Required(), mcp.Description("Name of the rule")),
	)
}

// --- Handlers ---

func listSkillsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tagFilter, _ := req.RequireString("tag")
	counts := loadCounts(ctx, "skills")
	type entry struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Tags        string `json:"tags,omitempty"`
		Downloads   int    `json:"downloads"`
	}
	var result []entry
	for name, info := range versions.Skills {
		if tagFilter != "" && !hasTag(info.Tags, tagFilter) {
			continue
		}
		result = append(result, entry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["skills/"+name],
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	out, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(out)), nil
}

func listRulesHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tagFilter, _ := req.RequireString("tag")
	counts := loadCounts(ctx, "rules")
	type entry struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Tags        string `json:"tags,omitempty"`
		Downloads   int    `json:"downloads"`
	}
	var result []entry
	for name, info := range versions.Rules {
		if tagFilter != "" && !hasTag(info.Tags, tagFilter) {
			continue
		}
		result = append(result, entry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["rules/"+name],
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	out, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(out)), nil
}

func installSkillHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return installArtifact(ctx, req, "skills", "skill")
}

func installRuleHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return installArtifact(ctx, req, "rules", "rule")
}

func installArtifact(ctx context.Context, req mcp.CallToolRequest, artifactType, paramName string) (*mcp.CallToolResult, error) {
	name, err := req.RequireString(paramName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	version := ""
	if v, err := req.RequireString("version"); err == nil {
		version = v
	}

	// Resolve version
	var info artifactInfo
	var found bool
	if artifactType == "skills" {
		info, found = versions.Skills[name]
	} else {
		info, found = versions.Rules[name]
	}
	if !found {
		return mcp.NewToolResultError(fmt.Sprintf("%s %q not found", paramName, name)), nil
	}
	if version == "" {
		version = info.Version
	}

	// Read all files from S3
	prefix := artifactType + "/" + name + "/" + version + "/"
	files, err := readAllFiles(ctx, prefix)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("read files: %v", err)), nil
	}
	if len(files) == 0 {
		return mcp.NewToolResultError(fmt.Sprintf("version %s not found for %s", version, name)), nil
	}

	// Generate checksums
	checksums := make(map[string]string)
	for path, content := range files {
		h := sha256.Sum256([]byte(content))
		checksums[path] = "sha256:" + hex.EncodeToString(h[:])
	}

	// Build manifest
	manifest := map[string]interface{}{
		"name":        name,
		"version":     version,
		"source":      source,
		"installedAt": time.Now().UTC().Format(time.RFC3339),
		"checksums":   checksums,
	}

	// Increment download counter
	incrementCount(ctx, artifactType+"/"+name)

	// Build response
	response := map[string]interface{}{
		"name":     name,
		"version":  version,
		"manifest": manifest,
		"files":    files,
	}
	out, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(out)), nil
}

func skillVersionsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return listArtifactVersions(ctx, req, "skills", "skill")
}

func ruleVersionsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return listArtifactVersions(ctx, req, "rules", "rule")
}

func listArtifactVersions(ctx context.Context, req mcp.CallToolRequest, artifactType, paramName string) (*mcp.CallToolResult, error) {
	name, err := req.RequireString(paramName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	prefix := artifactType + "/" + name + "/"
	versionList, err := listPrefixes(ctx, prefix)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("list versions: %v", err)), nil
	}
	if len(versionList) == 0 {
		return mcp.NewToolResultError(fmt.Sprintf("%s %q not found", paramName, name)), nil
	}

	sortVersions(versionList)
	out, _ := json.MarshalIndent(versionList, "", "  ")
	return mcp.NewToolResultText(string(out)), nil
}

// --- S3 helpers ---

func loadVersions(ctx context.Context) (*versionsFile, error) {
	key := "versions.json"
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var v versionsFile
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	if v.Skills == nil {
		v.Skills = make(map[string]artifactInfo)
	}
	if v.Rules == nil {
		v.Rules = make(map[string]artifactInfo)
	}
	return &v, nil
}

func listPrefixes(ctx context.Context, prefix string) ([]string, error) {
	delim := "/"
	resp, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    &bucket,
		Prefix:    &prefix,
		Delimiter: &delim,
	})
	if err != nil {
		return nil, err
	}
	var names []string
	for _, p := range resp.CommonPrefixes {
		name := strings.TrimPrefix(*p.Prefix, prefix)
		name = strings.TrimSuffix(name, "/")
		if name != "" {
			names = append(names, name)
		}
	}
	return names, nil
}

func readAllFiles(ctx context.Context, prefix string) (map[string]string, error) {
	resp, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil, err
	}

	files := make(map[string]string)
	for _, obj := range resp.Contents {
		key := *obj.Key
		relPath := strings.TrimPrefix(key, prefix)
		if relPath == "" || strings.HasSuffix(relPath, "/") {
			continue
		}
		if relPath == "tile.json" || relPath == "references/ACKNOWLEDGMENTS.md" {
			continue
		}

		getResp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			continue
		}
		data, err := io.ReadAll(getResp.Body)
		getResp.Body.Close()
		if err != nil {
			continue
		}
		files[relPath] = string(data)
	}
	return files, nil
}

// --- DynamoDB helpers ---

func incrementCount(ctx context.Context, artifact string) {
	dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &table,
		Key: map[string]types.AttributeValue{
			"artifact": &types.AttributeValueMemberS{Value: artifact},
		},
		UpdateExpression: strPtr("ADD downloads :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})
}

func loadCounts(ctx context.Context, artifactType string) map[string]int {
	counts := make(map[string]int)
	resp, err := dbClient.Scan(ctx, &dynamodb.ScanInput{
		TableName: &table,
	})
	if err != nil {
		return counts
	}
	for _, item := range resp.Items {
		artAttr, ok := item["artifact"].(*types.AttributeValueMemberS)
		if !ok {
			continue
		}
		dlAttr, ok := item["downloads"].(*types.AttributeValueMemberN)
		if !ok {
			continue
		}
		n, _ := strconv.Atoi(dlAttr.Value)
		counts[artAttr.Value] = n
	}
	return counts
}

// --- Helpers ---

func sortVersions(vers []string) {
	sort.Slice(vers, func(i, j int) bool {
		a := parseVersion(vers[i])
		b := parseVersion(vers[j])
		for k := 0; k < 3; k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return false
	})
}

func parseVersion(v string) [3]int {
	parts := strings.Split(v, ".")
	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		result[i], _ = strconv.Atoi(parts[i])
	}
	return result
}

func strPtr(s string) *string { return &s }

func hasTag(tags, filter string) bool {
	for _, t := range strings.Split(tags, ",") {
		if strings.TrimSpace(t) == filter {
			return true
		}
	}
	return false
}
