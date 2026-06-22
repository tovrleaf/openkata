package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CommandExecutor abstracts os/exec for testing.
type CommandExecutor interface {
	Run(name string, args ...string) ([]byte, error)
}

// realExecutor shells out to the real docker CLI.
type realExecutor struct{}

func (e *realExecutor) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

// DockerSandbox manages container lifecycle for scenario execution.
type DockerSandbox struct {
	Config   SandboxConfig
	Executor CommandExecutor
}

// NewDockerSandbox creates a sandbox runner from config.
func NewDockerSandbox(cfg SandboxConfig) *DockerSandbox {
	return &DockerSandbox{Config: cfg, Executor: &realExecutor{}}
}

// containerName generates a unique container name.
func containerName(scenario string) string {
	ts := time.Now().UnixMilli()
	safe := strings.ReplaceAll(scenario, "/", "-")
	return fmt.Sprintf("openkata-eval-%d-%s", ts, safe)
}

// createArgs builds the docker create command arguments.
func (d *DockerSandbox) createArgs(name, network string) []string {
	home, _ := os.UserHomeDir()

	args := []string{
		"create",
		"--name", name,
		"--network", network,
		"-v", home + "/.aws:/root/.aws:ro",
		"-v", home + "/.kiro:/root/.kiro:ro",
		"-w", "/workspace",
	}
	if kiroPath, err := exec.LookPath("kiro-cli"); err == nil {
		args = append(args, "-v", kiroPath+":/usr/local/bin/kiro-cli:ro")
	}
	args = append(args, d.Config.Image, "sleep", "infinity")
	return args
}

// copyScenarioInputs copies task.md and inputs/ into the container.
func (d *DockerSandbox) copyScenarioInputs(name string, scenario Scenario) error {
	// Copy task.md
	taskPath := filepath.Join(scenario.Dir, "task.md")
	if _, err := os.Stat(taskPath); err == nil {
		if _, err := d.Executor.Run("docker", "cp", taskPath, name+":/workspace/task.md"); err != nil {
			return fmt.Errorf("copy task.md: %w", err)
		}
	}

	// Copy inputs/
	inputsDir := filepath.Join(scenario.Dir, "inputs")
	if info, err := os.Stat(inputsDir); err == nil && info.IsDir() {
		if _, err := d.Executor.Run("docker", "cp", inputsDir+"/.", name+":/workspace/"); err != nil {
			return fmt.Errorf("copy inputs: %w", err)
		}
	}
	return nil
}

// writeAgentConfig creates the eval-agent.json inside the container.
func (d *DockerSandbox) writeAgentConfig(name string) error {
	agentJSON := agentConfigJSON()

	// Write to temp file, then docker cp
	tmp, err := os.CreateTemp("", "eval-agent-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(agentJSON); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	// Ensure .kiro/agents/ exists in container
	d.Executor.Run("docker", "exec", name, "mkdir", "-p", "/workspace/.kiro/agents")

	// Copy agent config
	if _, err := d.Executor.Run("docker", "cp", tmp.Name(), name+":/workspace/.kiro/agents/eval-agent.json"); err != nil {
		return fmt.Errorf("copy agent config: %w", err)
	}
	return nil
}

// agentConfigJSON returns the eval-agent.json content.
func agentConfigJSON() []byte {
	cfg := map[string]interface{}{
		"name":         "eval-agent",
		"tools":        []string{"fs_read", "fs_write", "execute_bash", "grep", "glob"},
		"allowedTools": []string{"fs_read", "fs_write", "execute_bash", "grep", "glob"},
		"resources":    []string{"file://skill/SKILL.md", "file://skill/references/**", "file://skill/scripts/**", "file://skill/assets/**"},
	}
	data, _ := json.Marshal(cfg)
	return data
}

// NewSandboxRunner creates a SandboxRunner from config with skill path for file copying.
func NewSandboxRunner(cfg SandboxConfig, skillPath string) SandboxRunner {
	sandbox := NewDockerSandbox(cfg)
	return func(scenario Scenario, prompt string) (string, error) {
		return sandboxRun(sandbox, skillPath, scenario, prompt)
	}
}

// sandboxRun handles the full sandbox lifecycle including skill file copy and workspace diff.
func sandboxRun(d *DockerSandbox, skillPath string, scenario Scenario, prompt string) (string, error) {
	name := containerName(scenario.Name)
	network := scenario.Network
	if network == "" {
		network = d.Config.Network
	}

	createArgs := d.createArgs(name, network)
	if _, err := d.Executor.Run("docker", createArgs...); err != nil {
		return "", fmt.Errorf("docker create: %w", err)
	}
	defer d.Executor.Run("docker", "rm", "-f", name)

	// Start container first so exec works
	if _, err := d.Executor.Run("docker", "start", name); err != nil {
		return "", fmt.Errorf("docker start: %w", err)
	}

	// Copy skill files
	if err := copySkillToContainer(d.Executor, name, skillPath); err != nil {
		return "", err
	}

	// Copy scenario inputs
	if err := d.copyScenarioInputs(name, scenario); err != nil {
		return "", err
	}

	// Write agent config
	if err := d.writeAgentConfig(name); err != nil {
		return "", err
	}

	// Snapshot workspace before
	beforeDir, err := copyWorkspaceOut(d.Executor, name)
	if err != nil {
		return "", fmt.Errorf("snapshot before: %w", err)
	}
	defer os.RemoveAll(beforeDir)
	before, _ := snapshotDir(beforeDir)

	// Run kiro-cli
	timeout := d.Config.Timeout
	if timeout == 0 {
		timeout = 300
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	_ = ctx // timeout is handled by the timeout command wrapper

	execArgs := []string{
		"exec", name,
		"timeout", fmt.Sprintf("%d", timeout),
		"kiro-cli", "chat",
		"--agent", "eval-agent",
		"--trust-all-tools",
		"--no-interactive",
		prompt,
	}
	out, execErr := d.Executor.Run("docker", execArgs...)

	// Snapshot workspace after
	afterDir, err := copyWorkspaceOut(d.Executor, name)
	if err != nil {
		return string(out), fmt.Errorf("snapshot after: %w", err)
	}
	defer os.RemoveAll(afterDir)
	after, _ := snapshotDir(afterDir)

	// Compute diff and format result
	diff := computeDiff(before, after)
	result := formatDiffForJudge(string(out), diff)

	if execErr != nil {
		return result, fmt.Errorf("kiro-cli: %w", execErr)
	}
	return result, nil
}

// copySkillToContainer copies skill files into /workspace/skill/ in the container.
func copySkillToContainer(executor CommandExecutor, name, skillPath string) error {
	executor.Run("docker", "exec", name, "mkdir", "-p", "/workspace/skill")

	// Copy SKILL.md
	skillMD := filepath.Join(skillPath, "SKILL.md")
	if _, err := os.Stat(skillMD); err == nil {
		if _, err := executor.Run("docker", "cp", skillMD, name+":/workspace/skill/SKILL.md"); err != nil {
			return fmt.Errorf("copy SKILL.md: %w", err)
		}
	}

	// Copy subdirectories
	for _, sub := range []string{"references", "scripts", "assets"} {
		dir := filepath.Join(skillPath, sub)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			if _, err := executor.Run("docker", "cp", dir, name+":/workspace/skill/"+sub); err != nil {
				return fmt.Errorf("copy %s: %w", sub, err)
			}
		}
	}
	return nil
}

// copyWorkspaceOut copies /workspace from the container to a temp directory.
func copyWorkspaceOut(executor CommandExecutor, name string) (string, error) {
	tmp, err := os.MkdirTemp("", "eval-workspace-*")
	if err != nil {
		return "", err
	}
	if _, err := executor.Run("docker", "cp", name+":/workspace/.", tmp); err != nil {
		os.RemoveAll(tmp)
		return "", err
	}
	return tmp, nil
}

// CheckDocker verifies docker is available (lazy check).
func CheckDocker() error {
	_, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker not found in PATH")
	}
	out, err := exec.Command("docker", "info").CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker not running: %s", string(out))
	}
	return nil
}
