package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all resolved configuration for the eval runner.
type Config struct {
	Backend       string
	Model         string
	JudgeModel    string
	Threshold     int
	Output        string
	DirectTimeout int
	HTTP          HTTPConfig
	Sandbox       SandboxConfig
	Concurrency   int
}

// HTTPConfig holds HTTP backend settings.
type HTTPConfig struct {
	BaseURL string
	APIKey  string
}

// SandboxConfig holds Docker sandbox settings.
type SandboxConfig struct {
	Image   string
	Timeout int
	Network string
}

// Scenario holds a discovered scenario.
type Scenario struct {
	Dir         string
	Name        string
	Description string
	Sandbox     bool
	Network     string
}

// yamlConfig mirrors the .openkata-eval.yaml file structure.
type yamlConfig struct {
	Backend       string `yaml:"backend"`
	Model         string `yaml:"model"`
	JudgeModel    string `yaml:"judge_model"`
	Threshold     int    `yaml:"threshold"`
	Concurrency   int    `yaml:"concurrency"`
	DirectTimeout int    `yaml:"direct_timeout"`
	HTTP          struct {
		BaseURL string `yaml:"base_url"`
		APIKey  string `yaml:"api_key"`
	} `yaml:"http"`
	Sandbox struct {
		Image   string `yaml:"image"`
		Timeout int    `yaml:"timeout"`
		Network string `yaml:"network"`
	} `yaml:"sandbox"`
}

// scenarioJSON mirrors the scenario.json file.
type scenarioJSON struct {
	Description string `json:"description"`
	Sandbox     *bool  `json:"sandbox"`
	Network     string `json:"network"`
}

// defaultConfig returns configuration defaults.
func defaultConfig() Config {
	return Config{
		Backend:       "kiro",
		Model:         "claude-sonnet-4.6",
		JudgeModel:    "",
		Threshold:     95,
		DirectTimeout: 120,
		Concurrency:   2,
		Sandbox: SandboxConfig{
			Image:   "openkata-eval:latest",
			Timeout: 300,
			Network: "host",
		},
	}
}

// loadConfig resolves config in order: defaults → yaml → env → flags.
// Flag values are only applied if explicitly set (non-zero/non-empty).
func loadConfig(flags *flagValues) Config {
	cfg := defaultConfig()

	// Load YAML config
	if data, err := os.ReadFile(".openkata-eval.yaml"); err == nil {
		var yc yamlConfig
		if yaml.Unmarshal(data, &yc) == nil {
			if yc.Backend != "" {
				cfg.Backend = yc.Backend
			}
			if yc.Model != "" {
				cfg.Model = yc.Model
			}
			if yc.JudgeModel != "" {
				cfg.JudgeModel = yc.JudgeModel
			}
			if yc.Threshold > 0 {
				cfg.Threshold = yc.Threshold
			}
			if yc.Concurrency > 0 {
				cfg.Concurrency = yc.Concurrency
			}
			if yc.DirectTimeout > 0 {
				cfg.DirectTimeout = yc.DirectTimeout
			}
			if yc.HTTP.BaseURL != "" {
				cfg.HTTP.BaseURL = yc.HTTP.BaseURL
			}
			if yc.HTTP.APIKey != "" {
				cfg.HTTP.APIKey = yc.HTTP.APIKey
			}
			if yc.Sandbox.Image != "" {
				cfg.Sandbox.Image = yc.Sandbox.Image
			}
			if yc.Sandbox.Timeout > 0 {
				cfg.Sandbox.Timeout = yc.Sandbox.Timeout
			}
			if yc.Sandbox.Network != "" {
				cfg.Sandbox.Network = yc.Sandbox.Network
			}
		}
	}

	// Apply environment variables
	if v := os.Getenv("OPENKATA_EVAL_BACKEND"); v != "" {
		cfg.Backend = v
	}
	if v := os.Getenv("OPENKATA_EVAL_MODEL"); v != "" {
		cfg.Model = v
	}
	if v := os.Getenv("OPENKATA_EVAL_JUDGE_MODEL"); v != "" {
		cfg.JudgeModel = v
	}
	if v := os.Getenv("OPENKATA_EVAL_THRESHOLD"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Threshold = n
		}
	}
	if v := os.Getenv("OPENKATA_EVAL_HTTP_BASE_URL"); v != "" {
		cfg.HTTP.BaseURL = v
	}
	if v := os.Getenv("OPENKATA_EVAL_HTTP_API_KEY"); v != "" {
		cfg.HTTP.APIKey = v
	}
	if v := os.Getenv("OPENKATA_EVAL_DIRECT_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.DirectTimeout = n
		}
	}

	// Apply CLI flags (only if explicitly set)
	if flags.backend != "" {
		cfg.Backend = flags.backend
	}
	if flags.model != "" {
		cfg.Model = flags.model
	}
	if flags.judgeModel != "" {
		cfg.JudgeModel = flags.judgeModel
	}
	if flags.threshold > 0 {
		cfg.Threshold = flags.threshold
	}
	if flags.output != "" {
		cfg.Output = flags.output
	}

	// If judge model not set, use agent model
	if cfg.JudgeModel == "" {
		cfg.JudgeModel = cfg.Model
	}

	return cfg
}

// flagValues holds raw CLI flag values.
type flagValues struct {
	backend    string
	model      string
	judgeModel string
	threshold  int
	output     string
}

// discoverScenarios finds scenario directories under the skill's evals/ dir.
func discoverScenarios(skillPath string) ([]Scenario, error) {
	evalsDir := filepath.Join(skillPath, "evals")
	entries, err := os.ReadDir(evalsDir)
	if err != nil {
		return nil, fmt.Errorf("no evals directory found: %w", err)
	}

	var scenarios []Scenario
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "scenario-") {
			continue
		}
		dir := filepath.Join(evalsDir, e.Name())
		s := Scenario{
			Dir:     dir,
			Name:    e.Name(),
			Sandbox: true,
			Network: "",
		}
		// Read scenario.json if present
		if data, err := os.ReadFile(filepath.Join(dir, "scenario.json")); err == nil {
			var sj scenarioJSON
			if parseJSON(data, &sj) == nil {
				s.Description = sj.Description
				if sj.Sandbox != nil {
					s.Sandbox = *sj.Sandbox
				}
				if sj.Network != "" {
					s.Network = sj.Network
				}
			}
		}
		scenarios = append(scenarios, s)
	}

	sort.Slice(scenarios, func(i, j int) bool {
		return scenarios[i].Name < scenarios[j].Name
	})
	return scenarios, nil
}

// parseMode determines if we're running a full skill or single scenario.
type parseMode int

const (
	modeSkill    parseMode = iota
	modeScenario
)

// parsePath determines if the given path is a skill directory or a single scenario.
func parsePath(path string) (skillPath string, mode parseMode, scenarioName string) {
	// Check if path ends with a scenario directory pattern
	base := filepath.Base(path)
	if strings.HasPrefix(base, "scenario-") {
		// Single scenario mode: path is like skills/X/evals/scenario-0
		skillPath = filepath.Dir(filepath.Dir(path))
		return skillPath, modeScenario, base
	}
	return path, modeSkill, ""
}

// parseJSON is a helper to unmarshal JSON data.
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
