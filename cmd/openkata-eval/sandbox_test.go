package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// mockExecutor records commands for verification.
type mockExecutor struct {
	commands [][]string
	outputs  map[string][]byte
	errors   map[string]error
}

func newMockExecutor() *mockExecutor {
	return &mockExecutor{
		outputs: make(map[string][]byte),
		errors:  make(map[string]error),
	}
}

func (m *mockExecutor) Run(name string, args ...string) ([]byte, error) {
	cmd := append([]string{name}, args...)
	m.commands = append(m.commands, cmd)
	key := strings.Join(cmd, " ")
	if err, ok := m.errors[key]; ok {
		return nil, err
	}
	if out, ok := m.outputs[key]; ok {
		return out, nil
	}
	return nil, nil
}

func TestContainerName(t *testing.T) {
	name := containerName("scenario-0")
	if !strings.HasPrefix(name, "openkata-eval-") {
		t.Errorf("containerName() = %q, want prefix openkata-eval-", name)
	}
	if !strings.HasSuffix(name, "-scenario-0") {
		t.Errorf("containerName() = %q, want suffix -scenario-0", name)
	}
}

func TestContainerNameSanitizesSlashes(t *testing.T) {
	name := containerName("path/to/scenario")
	if strings.Contains(name, "/") {
		t.Errorf("containerName() = %q, should not contain /", name)
	}
}

func TestCreateArgs(t *testing.T) {
	sandbox := &DockerSandbox{
		Config: SandboxConfig{
			Image:   "openkata-eval:latest",
			Timeout: 300,
			Network: "host",
		},
	}

	args := sandbox.createArgs("test-container", "host")

	// Verify key arguments
	found := map[string]bool{}
	for i, a := range args {
		switch a {
		case "create":
			found["create"] = true
		case "--name":
			if i+1 < len(args) && args[i+1] == "test-container" {
				found["name"] = true
			}
		case "--network":
			if i+1 < len(args) && args[i+1] == "host" {
				found["network"] = true
			}
		case "-v":
			if i+1 < len(args) && strings.Contains(args[i+1], ".aws") {
				found["aws-mount"] = true
			}
		case "-w":
			if i+1 < len(args) && args[i+1] == "/workspace" {
				found["workdir"] = true
			}
		}
	}
	if args[len(args)-2] != "sleep" || args[len(args)-1] != "infinity" {
		t.Error("expected sleep infinity as container command")
	}

	for _, key := range []string{"create", "name", "network", "aws-mount", "workdir"} {
		if !found[key] {
			t.Errorf("createArgs missing %s, got: %v", key, args)
		}
	}
}

func TestCreateArgsCustomNetwork(t *testing.T) {
	sandbox := &DockerSandbox{
		Config: SandboxConfig{
			Image:   "openkata-eval:latest",
			Network: "bridge",
		},
	}
	args := sandbox.createArgs("c1", "none")
	for i, a := range args {
		if a == "--network" && i+1 < len(args) {
			if args[i+1] != "none" {
				t.Errorf("network = %q, want %q", args[i+1], "none")
			}
			return
		}
	}
	t.Error("--network not found in args")
}

func TestAgentConfigJSON(t *testing.T) {
	data := agentConfigJSON()

	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("agentConfigJSON() produced invalid JSON: %v", err)
	}
	if cfg["name"] != "eval-agent" {
		t.Errorf("name = %v, want eval-agent", cfg["name"])
	}

	tools, ok := cfg["tools"].([]interface{})
	if !ok || len(tools) != 5 {
		t.Errorf("tools = %v, want 5 items", cfg["tools"])
	}

	resources, ok := cfg["resources"].([]interface{})
	if !ok || len(resources) != 4 {
		t.Errorf("resources = %v, want 4 items", cfg["resources"])
	}
}

func TestSandboxRunCallsDocker(t *testing.T) {
	executor := newMockExecutor()
	sandbox := &DockerSandbox{
		Config: SandboxConfig{
			Image:   "openkata-eval:latest",
			Timeout: 60,
			Network: "host",
		},
		Executor: executor,
	}

	scenario := Scenario{
		Dir:  t.TempDir(),
		Name: "scenario-0",
	}

	// Use sandboxRun which is the live code path
	sandboxRun(sandbox, t.TempDir(), scenario, "test prompt")

	if len(executor.commands) == 0 {
		t.Fatal("no docker commands were executed")
	}

	// First command should be docker create
	first := executor.commands[0]
	if first[0] != "docker" || first[1] != "create" {
		t.Errorf("first command = %v, want docker create", first[:2])
	}
}
