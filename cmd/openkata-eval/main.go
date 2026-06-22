package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: openkata-eval [flags] <skill-path|scenario-path>\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  openkata-eval skills/commit-conventions\n")
	fmt.Fprintf(os.Stderr, "  openkata-eval skills/commit-conventions/evals/scenario-0\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	var flags flagValues
	flag.StringVar(&flags.backend, "backend", "", "LLM backend: kiro or http")
	flag.StringVar(&flags.model, "model", "", "model for agent calls")
	flag.StringVar(&flags.judgeModel, "judge-model", "", "model for judge calls")
	flag.IntVar(&flags.threshold, "threshold", 0, "pass threshold percentage (default 95)")
	flag.StringVar(&flags.output, "output", "", "JSON output file or directory path")
	flag.StringVar(&flags.modelLabel, "model-label", "", "human-readable label for the model")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}

	targetPath := flag.Arg(0)
	cfg := loadConfig(&flags)

	skillPath, mode, scenarioName := parsePath(targetPath)

	// Validate skill has SKILL.md
	if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
		fmt.Fprintf(os.Stderr, "error: SKILL.md not found in %s\n", skillPath)
		os.Exit(1)
	}

	var scenarios []Scenario
	if mode == modeScenario {
		all, err := discoverScenarios(skillPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		for _, s := range all {
			if s.Name == scenarioName {
				scenarios = append(scenarios, s)
				break
			}
		}
		if len(scenarios) == 0 {
			fmt.Fprintf(os.Stderr, "error: scenario %s not found\n", scenarioName)
			os.Exit(1)
		}
	} else {
		var err error
		scenarios, err = discoverScenarios(skillPath)
		if err != nil || len(scenarios) == 0 {
			fmt.Fprintf(os.Stderr, "error: no eval scenarios found in %s/evals/\n", skillPath)
			os.Exit(1)
		}
	}

	// Check kiro-cli availability
	if cfg.Backend == "kiro" {
		if _, err := exec.LookPath("kiro-cli"); err != nil {
			fmt.Fprintf(os.Stderr, "error: kiro-cli not found in PATH\n")
			fmt.Fprintf(os.Stderr, "Install: https://docs.kiro.dev/install\n")
			os.Exit(1)
		}
	}

	// Check if any scenarios need Docker
	needsDocker := false
	for _, s := range scenarios {
		if s.Sandbox {
			needsDocker = true
			break
		}
	}

	var sandboxFn SandboxRunner
	if needsDocker {
		if err := CheckDocker(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		sandboxFn = NewSandboxRunner(cfg.Sandbox, skillPath)
	}

	// Create completers
	timeout := time.Duration(cfg.DirectTimeout) * time.Second
	agent := newCompleter(cfg, cfg.Model, timeout)
	judge := newCompleter(cfg, cfg.JudgeModel, timeout)

	// For kiro backend, agent calls use --agent eval-agent (direct mode)
	if kc, ok := agent.(*KiroCompleter); ok {
		kc.AgentMode = true
		kc.SkillPath = skillPath
	}

	// Print header
	fmt.Println(skillPath)

	// Run scenarios
	runner := &Runner{
		Completer:   agent,
		JudgeModel:  judge,
		Concurrency: cfg.Concurrency,
		SkillPath:   skillPath,
		Sandbox:     sandboxFn,
	}

	results, err := runner.Run(scenarios)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Compute overall
	evalResult := computeOverall(results, cfg.Threshold)
	printOverall(evalResult)

	// Write JSON output
	direct := false
	if cfg.Output == "" && mode == modeSkill {
		cfg.Output = autoResolvePath(skillPath, cfg.Model, time.Now())
		direct = true
	}
	if cfg.Output != "" {
		skill := filepath.Base(skillPath)
		if err := writeJSONOutput(cfg.Output, skill, skillPath, cfg, evalResult, direct); err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		}
	}

	os.Exit(shouldExit(mode, evalResult))
}

// newCompleter creates the appropriate Completer based on config backend.
func newCompleter(cfg Config, model string, timeout time.Duration) Completer {
	switch cfg.Backend {
	case "http":
		return NewHTTPCompleter(cfg.HTTP.BaseURL, cfg.HTTP.APIKey, model, timeout)
	default:
		return NewKiroCompleter(model, timeout)
	}
}
