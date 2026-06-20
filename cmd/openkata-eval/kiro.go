package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ansiPattern matches ANSI escape sequences.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// KiroCompleter shells out to kiro-cli chat.
type KiroCompleter struct {
	Model     string
	Timeout   time.Duration
	AgentMode bool
	SkillPath string
}

// NewKiroCompleter creates a completer that uses kiro-cli.
func NewKiroCompleter(model string, timeout time.Duration) *KiroCompleter {
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	return &KiroCompleter{Model: model, Timeout: timeout}
}

// Complete sends a prompt to kiro-cli and returns the cleaned response.
func (k *KiroCompleter) Complete(system, user string) (string, error) {
	prompt := system + "\n\n" + user

	ctx, cancel := context.WithTimeout(context.Background(), k.Timeout)
	defer cancel()

	args := []string{"chat", "--no-interactive"}
	if k.Model != "" {
		args = append(args, "--model", k.Model)
	}

	var cmd *exec.Cmd
	var cleanupDir string

	if k.AgentMode && k.SkillPath != "" {
		// Create temp dir with .kiro/agents/eval-agent.json
		tmpDir, err := os.MkdirTemp("", "openkata-eval-agent-*")
		if err != nil {
			return "", fmt.Errorf("creating temp dir: %w", err)
		}
		cleanupDir = tmpDir

		agentDir := filepath.Join(tmpDir, ".kiro", "agents")
		if err := os.MkdirAll(agentDir, 0755); err != nil {
			os.RemoveAll(tmpDir)
			return "", fmt.Errorf("creating agent dir: %w", err)
		}

		// Write eval-agent.json (no tools for direct mode)
		agentCfg := map[string]interface{}{
			"name":      "eval-agent",
			"resources": []string{"file://skill/SKILL.md", "file://skill/references/**", "file://skill/scripts/**", "file://skill/assets/**"},
		}
		data, _ := json.Marshal(agentCfg)
		if err := os.WriteFile(filepath.Join(agentDir, "eval-agent.json"), data, 0644); err != nil {
			os.RemoveAll(tmpDir)
			return "", fmt.Errorf("writing agent config: %w", err)
		}

		// Symlink skill into temp dir
		skillLink := filepath.Join(tmpDir, "skill")
		absSkill, _ := filepath.Abs(k.SkillPath)
		os.Symlink(absSkill, skillLink)

		args = append([]string{"chat", "--agent", "eval-agent", "--no-interactive"}, args[2:]...)
		if k.Model != "" {
			args = []string{"chat", "--agent", "eval-agent", "--no-interactive", "--model", k.Model}
		} else {
			args = []string{"chat", "--agent", "eval-agent", "--no-interactive"}
		}
		args = append(args, prompt)

		cmd = exec.CommandContext(ctx, "kiro-cli", args...)
		cmd.Dir = tmpDir
	} else {
		args = append(args, prompt)
		cmd = exec.CommandContext(ctx, "kiro-cli", args...)
	}

	out, err := cmd.Output()
	if cleanupDir != "" {
		os.RemoveAll(cleanupDir)
	}
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("kiro-cli timed out after %s", k.Timeout)
		}
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("kiro-cli failed: %w\nstderr: %s", err, exitErr.Stderr)
		}
		return "", fmt.Errorf("kiro-cli failed: %w", err)
	}

	if len(out) == 0 {
		return "", fmt.Errorf("kiro-cli returned empty response")
	}

	return cleanKiroOutput(string(out)), nil
}

// cleanKiroOutput strips ANSI codes, header lines, and markdown fences from kiro-cli output.
func cleanKiroOutput(raw string) string {
	// Strip ANSI escape sequences
	cleaned := stripANSI(raw)

	// Remove header lines starting with "> " and markdown code fences
	lines := strings.Split(cleaned, "\n")
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "> ") {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			continue
		}
		result = append(result, line)
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}
