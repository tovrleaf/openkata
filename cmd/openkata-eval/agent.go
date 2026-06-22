package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// buildAgentPrompt constructs the full agent prompt from skill context and scenario.
func buildAgentPrompt(skillPath string, scenario Scenario) (string, error) {
	var parts []string

	// SKILL.md
	skillDoc, err := os.ReadFile(filepath.Join(skillPath, "SKILL.md"))
	if err != nil {
		return "", fmt.Errorf("reading SKILL.md: %w", err)
	}
	parts = append(parts, string(skillDoc))

	// references/ (excluding ACKNOWLEDGMENTS.md)
	refsDir := filepath.Join(skillPath, "references")
	if refs, err := collectDir(refsDir, func(rel string) bool {
		return rel == "ACKNOWLEDGMENTS.md"
	}); err == nil && len(refs) > 0 {
		parts = append(parts, refs...)
	}

	// scripts/
	scriptsDir := filepath.Join(skillPath, "scripts")
	if scripts, err := collectDir(scriptsDir, nil); err == nil && len(scripts) > 0 {
		parts = append(parts, scripts...)
	}

	// assets/
	assetsDir := filepath.Join(skillPath, "assets")
	if assets, err := collectDir(assetsDir, nil); err == nil && len(assets) > 0 {
		parts = append(parts, assets...)
	}

	// task.md
	taskData, err := os.ReadFile(filepath.Join(scenario.Dir, "task.md"))
	if err != nil {
		return "", fmt.Errorf("reading task.md: %w", err)
	}
	parts = append(parts, string(taskData))

	// inputs/ (as fenced code blocks)
	inputsDir := filepath.Join(scenario.Dir, "inputs")
	if inputs, err := collectInputs(inputsDir); err == nil && len(inputs) > 0 {
		parts = append(parts, inputs...)
	}

	return strings.Join(parts, "\n\n"), nil
}

// collectDir reads all files in a directory recursively, applying an exclude filter.
// Returns file contents prefixed with their relative path as a header.
func collectDir(dir string, exclude func(string) bool) ([]string, error) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, err
	}

	var files []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)
		if exclude != nil && exclude(rel) {
			return nil
		}
		files = append(files, path)
		return nil
	})

	sort.Strings(files)
	var results []string
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)
		header := fmt.Sprintf("--- %s ---", rel)
		results = append(results, header+"\n"+string(data))
	}
	return results, nil
}

// collectInputs reads files from an inputs/ directory and formats them as fenced code blocks.
func collectInputs(dir string) ([]string, error) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil, err
	}

	var files []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	sort.Strings(files)
	var results []string
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)
		block := fmt.Sprintf("```%s\n%s\n```", rel, strings.TrimRight(string(data), "\n"))
		results = append(results, block)
	}
	return results, nil
}
