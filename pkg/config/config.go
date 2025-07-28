package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type StackedConfig struct {
	Branches map[string]BranchInfo `json:"branches"`
}

type BranchInfo struct {
	Parent string `json:"parent"`
}

func getConfigPath() (string, error) {
	gitDir, err := findGitDir()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	configDir := filepath.Join(gitDir, "gh-stacked")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("error creating config directory: %w", err)
	}

	return filepath.Join(configDir, "branches.json"), nil
}

func findGitDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		
		// Check if .git is a file (worktree case)
		if info, err := os.Stat(gitDir); err == nil {
			if info.IsDir() {
				return gitDir, nil
			}
			// .git is a file, read the gitdir path
			data, err := os.ReadFile(gitDir)
			if err != nil {
				return "", fmt.Errorf("error reading .git file: %w", err)
			}
			// Parse "gitdir: /path/to/git/dir"
			gitdirLine := string(data)
			if len(gitdirLine) > 8 && gitdirLine[:8] == "gitdir: " {
				return filepath.Clean(gitdirLine[8:]), nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func loadConfig() (*StackedConfig, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &StackedConfig{Branches: make(map[string]BranchInfo)}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config StackedConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if config.Branches == nil {
		config.Branches = make(map[string]BranchInfo)
	}

	return &config, nil
}

func saveConfig(config *StackedConfig) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func SetBranchParent(branch, parent string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.Branches[branch] = BranchInfo{Parent: parent}
	return saveConfig(config)
}

func GetBranchParent(branch string) (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", err
	}

	branchInfo, exists := config.Branches[branch]
	if !exists {
		return "", fmt.Errorf("no parent branch found for %s", branch)
	}

	return branchInfo.Parent, nil
}

func RemoveBranch(branch string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	delete(config.Branches, branch)
	return saveConfig(config)
}