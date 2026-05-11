package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tovrleaf/openkata/cmd/openkata-web/templates"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templates.Home().Render(r.Context(), w)
}

func handleDesignSystem(w http.ResponseWriter, r *http.Request) {
	templates.DesignSystem().Render(r.Context(), w)
}

func handleSkills(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/skills/" {
		http.NotFound(w, r)
		return
	}

	skills := loadSkillsList(r.Context())
	templates.Skills(skills).Render(r.Context(), w)
}

func loadSkillsList(ctx context.Context) []templates.SkillEntry {
	if s3Client == nil {
		return nil
	}

	// Read versions.json from S3
	key := "versions.json"
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var versions struct {
		Skills map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil
	}

	// Read download counts
	counts := make(map[string]int)
	if dbClient != nil {
		scanResp, err := dbClient.Scan(ctx, &dynamodb.ScanInput{
			TableName: &table,
		})
		if err == nil {
			for _, item := range scanResp.Items {
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
		}
	}

	var skills []templates.SkillEntry
	for name, info := range versions.Skills {
		skills = append(skills, templates.SkillEntry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["skills/"+name],
		})
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/rules/" {
		http.NotFound(w, r)
		return
	}

	rules := loadRulesList(r.Context())
	templates.Rules(rules).Render(r.Context(), w)
}

func loadRulesList(ctx context.Context) []templates.SkillEntry {
	if s3Client == nil {
		return nil
	}

	key := "versions.json"
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var versions struct {
		Rules map[string]struct {
			Version     string `json:"version"`
			Description string `json:"description"`
			Tags        string `json:"tags"`
		} `json:"rules"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil
	}

	counts := make(map[string]int)
	if dbClient != nil {
		scanResp, err := dbClient.Scan(ctx, &dynamodb.ScanInput{
			TableName: &table,
		})
		if err == nil {
			for _, item := range scanResp.Items {
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
		}
	}

	var rules []templates.SkillEntry
	for name, info := range versions.Rules {
		rules = append(rules, templates.SkillEntry{
			Name:        name,
			Version:     info.Version,
			Description: info.Description,
			Tags:        info.Tags,
			Downloads:   counts["rules/"+name],
		})
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].Name < rules[j].Name })
	return rules
}
