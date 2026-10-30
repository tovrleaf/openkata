package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/tovrleaf/openkata/internal/evalresults"
	"github.com/tovrleaf/openkata/internal/versionindex"
)

// bucketAPI is the part of the S3 client that building the index uses, so
// tests can supply a fake bucket.
type bucketAPI interface {
	GetObject(ctx context.Context, in *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, in *s3.ListObjectsV2Input, opts ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// buildIndexFromBucket builds the index from what is published in the
// bucket: for each artifact its newest version folder, with the description
// and tags from its main file and, for a skill, the tessl score and the eval
// results, all read from that version's own files. Things that cannot be
// read are left out and reported on stderr.
func buildIndexFromBucket(ctx context.Context, client bucketAPI, bucket string) *versionindex.Index {
	versions := versionindex.New()

	for _, artifactType := range []string{"skills", "rules", "profiles"} {
		names, err := listPrefixes(ctx, client, bucket, artifactType+"/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "listing %s: %v\n", artifactType, err)
			continue
		}

		for _, name := range names {
			versionList, err := listPrefixes(ctx, client, bucket, artifactType+"/"+name+"/")
			if err != nil || len(versionList) == 0 {
				continue
			}

			sortVersions(versionList)
			latest := versionList[len(versionList)-1]
			versionPrefix := artifactType + "/" + name + "/" + latest + "/"

			mdFile := "SKILL.md"
			if artifactType == "rules" {
				mdFile = "RULE.md"
			} else if artifactType == "profiles" {
				mdFile = name + ".md"
			}

			description, tags := readFrontmatterS3(ctx, client, bucket, versionPrefix+mdFile)
			info := versionindex.Artifact{
				Version:     latest,
				Description: description,
				Tags:        tags,
			}

			switch artifactType {
			case "skills":
				info.TesslScore, info.Published = readTesslScore(ctx, client, bucket, versionPrefix)
				info.Models = readEvalResults(ctx, client, bucket, versionPrefix, latest)
				versions.Skills[name] = info
			case "rules":
				versions.Rules[name] = info
			case "profiles":
				versions.Profiles[name] = info
			}
		}
	}
	return versions
}

// readTesslScore returns the score from the version's tessl plugin file and
// whether the skill is published, which it is unless the plugin is private.
// A missing file, or a score of zero, means no score.
func readTesslScore(ctx context.Context, client bucketAPI, bucket, versionPrefix string) (int, bool) {
	data, err := getObject(ctx, client, bucket, versionPrefix+".tessl-plugin/plugin.json")
	if err != nil {
		return 0, false
	}
	var plugin struct {
		Score   int  `json:"score"`
		Private bool `json:"private"`
	}
	if json.Unmarshal(data, &plugin) != nil || plugin.Score <= 0 {
		return 0, false
	}
	return plugin.Score, !plugin.Private
}

// readEvalResults returns the effectiveness by model from the eval reports
// in the version's own evals/results folder. Only the newest report of each
// model is downloaded, and the rules for which reports count are in
// internal/evalresults. version is what the reports must have been run
// against.
func readEvalResults(ctx context.Context, client bucketAPI, bucket, versionPrefix, version string) map[string]versionindex.Model {
	resultsPrefix := versionPrefix + "evals/results/"
	keys, err := listKeys(ctx, client, bucket, resultsPrefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listing eval results under %s: %v\n", resultsPrefix, err)
		return nil
	}

	paths := make([]string, 0, len(keys))
	for _, k := range keys {
		paths = append(paths, strings.TrimPrefix(k, resultsPrefix))
	}

	var files []evalresults.File
	for _, rel := range evalresults.Newest(paths) {
		data, err := getObject(ctx, client, bucket, resultsPrefix+rel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reading eval report %s%s: %v\n", resultsPrefix, rel, err)
			continue
		}
		files = append(files, evalresults.File{Path: rel, Data: data})
	}
	return evalresults.ForVersion(files, version)
}

// listPrefixes returns the names of the folders directly under prefix.
func listPrefixes(ctx context.Context, client bucketAPI, bucket, prefix string) ([]string, error) {
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    &bucket,
		Prefix:    &prefix,
		Delimiter: strPtr("/"),
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

// listKeys returns every key under prefix, following the pages of the
// listing.
func listKeys(ctx context.Context, client bucketAPI, bucket, prefix string) ([]string, error) {
	var keys []string
	var token *string
	for {
		resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            &bucket,
			Prefix:            &prefix,
			ContinuationToken: token,
		})
		if err != nil {
			return nil, err
		}
		for _, o := range resp.Contents {
			keys = append(keys, *o.Key)
		}
		if !aws.ToBool(resp.IsTruncated) || resp.NextContinuationToken == nil {
			return keys, nil
		}
		token = resp.NextContinuationToken
	}
}

func getObject(ctx context.Context, client bucketAPI, bucket, key string) ([]byte, error) {
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func readFrontmatterS3(ctx context.Context, client bucketAPI, bucket, key string) (description, tags string) {
	data, err := getObject(ctx, client, bucket, key)
	if err != nil {
		return "", ""
	}
	return parseFrontmatter(string(data))
}
