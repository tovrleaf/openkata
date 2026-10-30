// Package versionindex defines versions.json, the index of released
// artifacts. The generator writes it and the web and MCP servers read
// it, so they all use these types and agree on its shape.
package versionindex

import (
	"encoding/json"
	"fmt"
)

// Index is the content of versions.json: one map of artifacts for each
// type, keyed by artifact name.
type Index struct {
	Skills   map[string]Artifact `json:"skills"`
	Rules    map[string]Artifact `json:"rules"`
	Profiles map[string]Artifact `json:"profiles"`
}

// Artifact is one entry of the index. Version is the newest released
// version. The remaining fields are only written when there is something
// to say, and only skills have them.
type Artifact struct {
	Version     string           `json:"version"`
	Description string           `json:"description"`
	Tags        string           `json:"tags,omitempty"`
	Models      map[string]Model `json:"models,omitempty"`
	TesslScore  int              `json:"tessl_score,omitempty"`
	Published   bool             `json:"published,omitempty"`
}

// Model is the effectiveness of a skill for one model, from its evals.
type Model struct {
	Label         string     `json:"label"`
	Effectiveness int        `json:"effectiveness"`
	Scenarios     []Scenario `json:"scenarios"`
}

// Scenario is the result of one eval scenario for a model.
type Scenario struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pass        bool   `json:"pass"`
	Score       int    `json:"score"`
}

// New returns an empty index whose maps are ready to use.
func New() *Index {
	return &Index{
		Skills:   map[string]Artifact{},
		Rules:    map[string]Artifact{},
		Profiles: map[string]Artifact{},
	}
}

// Parse decodes versions.json. A section that is missing from the file is
// an empty map, never nil.
func Parse(data []byte) (*Index, error) {
	idx := New()
	if err := json.Unmarshal(data, idx); err != nil {
		return nil, fmt.Errorf("parse versions.json: %w", err)
	}
	if idx.Skills == nil {
		idx.Skills = map[string]Artifact{}
	}
	if idx.Rules == nil {
		idx.Rules = map[string]Artifact{}
	}
	if idx.Profiles == nil {
		idx.Profiles = map[string]Artifact{}
	}
	return idx, nil
}

// Of returns the artifacts of a type, "skills", "rules" or "profiles". Any
// other type has none.
func (i *Index) Of(artifactType string) map[string]Artifact {
	switch artifactType {
	case "skills":
		return i.Skills
	case "rules":
		return i.Rules
	case "profiles":
		return i.Profiles
	}
	return nil
}

// Marshal encodes the index as versions.json is written: indented with two
// spaces and ending in a newline.
func (i *Index) Marshal() ([]byte, error) {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}
