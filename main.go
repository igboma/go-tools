package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gopkg.in/yaml.v2"
)

// Config struct to match the conf.yaml structure
type Config struct {
	Version                        string `yaml:"version"`
	VersionOverride                string `yaml:"versionOverride"`
	Namespace                      string `yaml:"namespace"`
	GomplateDatasources            string `yaml:"gomplateDatasources"`
	HeoRoot                        string `yaml:"heoRoot"`
	HeoRevision                    string `yaml:"heoRevision"`
	HeoRevisionOverride            string `yaml:"heoRevisionOverride"`
	EnableArgoHookDeleteRedis      string `yaml:"enableArgoHookDeleteRedis"`
	EnableArgoHookDeleteRedisForce string `yaml:"enableArgoHookDeleteRedisForce"`
	SlackNotifyChannel             string `yaml:"slackNotifyChannel"`
	Onboarded                      string `yaml:"onboarded"`
	DeploymentSchedule             string `yaml:"deploymentSchedule"`
	DeploymentWindow               int    `yaml:"deploymentWindow"`
}

// GitRepository interface to abstract git repository operations
type GitRepository interface {
	GetFileContentFromBranch(branch, file string) (string, error)
	GetFileContentFromCommit(commitHash, file string) (string, error)
}

// GitRepositoryImpl is the actual implementation using go-git
type GitRepositoryImpl struct {
	repo *git.Repository
}

func NewGitRepository(repo *git.Repository) *GitRepositoryImpl {
	return &GitRepositoryImpl{repo: repo}
}

func (g *GitRepositoryImpl) GetFileContentFromBranch(branch, file string) (string, error) {
	ref, err := g.repo.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		return "", err
	}

	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	return getFileContentFromCommit(commit, file)
}

func (g *GitRepositoryImpl) GetFileContentFromCommit(commitHash, file string) (string, error) {
	commit, err := g.repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		return "", err
	}

	return getFileContentFromCommit(commit, file)
}

// VersionChecker struct that uses GitRepository to check versions and revisions
type VersionChecker struct {
	repo GitRepository
}

func NewVersionChecker(repo GitRepository) *VersionChecker {
	return &VersionChecker{repo: repo}
}

func (v *VersionChecker) GetVersionAndHeoRevision(branch, file string) (string, string, error) {
	content, err := v.repo.GetFileContentFromBranch(branch, file)
	if err != nil {
		return "", "", err
	}

	var config Config
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return "", "", err
	}

	return config.Version, config.HeoRevision, nil
}

func checkVersionAndHeoRevisionDiff(r *git.Repository, file string) bool {
	// Read current file contents
	currentConfig, err := getCurrentConfig(file)
	if err != nil {
		log.Fatalf("Failed to get current config: %v", err)
	}

	// Open previous config from origin/main
	gitRepo := NewGitRepository(r)
	previousConfigContent, err := gitRepo.GetFileContentFromBranch("refs/remotes/origin/main", file)
	if err != nil {
		log.Fatalf("Failed to get previous config: %v", err)
	}

	var previousConfig Config
	err = yaml.Unmarshal([]byte(previousConfigContent), &previousConfig)
	if err != nil {
		log.Fatalf("Failed to parse previous YAML: %v", err)
	}

	// Compare version and heoRevision
	if currentConfig.Version != previousConfig.Version || currentConfig.HeoRevision != previousConfig.HeoRevision {
		fmt.Println("Version or heoRevision has changed")
		return true
	}

	fmt.Println("No version or heoRevision change. Skipping creation of deployment.")
	return false
}

// Main logic
func main() {
	// Check if the repository path is passed as a command-line argument
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go-tools <repo-path>")
	}

	repoPath := os.Args[1]

	// Retrieve environment variables (set in GitHub Actions)
	changedFiles := os.Getenv("CHANGED_FILES")
	action := os.Getenv("GITHUB_EVENT_ACTION")
	prMerged := os.Getenv("GITHUB_EVENT_PR_MERGED")
	outputFile := os.Getenv("GITHUB_OUTPUT")

	// Split the changed files into an array
	files := strings.Split(changedFiles, "\n")

	// Ensure only one file was changed
	if len(files) > 1 {
		log.Fatalf("More than one file was changed")
	}

	file := files[0]
	fmt.Printf("CHANGED FILES: %s\n", file)
	fmt.Printf("FILE: %s\n", file)

	// Ensure the file is a conf.yaml file
	if !strings.Contains(file, "components/") || !strings.HasSuffix(file, "conf.yaml") {
		log.Fatalf("The file is not a conf.yaml file")
	}

	// Extract COMPONENT and ENVIRONMENT
	component := strings.Split(file, "/")[1]
	environment := strings.Split(file, "/")[2]

	// Open the repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Failed to open repository at %s: %v", repoPath, err)
	}

	gitRepo := NewGitRepository(r)
	checker := NewVersionChecker(gitRepo)

	var version, heoRevision string
	if action == "closed" && prMerged == "true" {
		fmt.Println("PR is merged...")
		version, heoRevision, err = checker.GetVersionAndHeoRevision("refs/heads/main", file)
		if err != nil {
			log.Fatalf("Failed to get version and heoRevision: %v", err)
		}
	} else {
		fmt.Println("PR is NOT merged...")

		// Check if version or heoRevision has changed
		if !checkVersionAndHeoRevisionDiff(r, file) {
			// Exit if there is no version or heoRevision change
			os.Exit(0)
		}

		// Fetch current config for the remaining fields
		currentConfig, err := getCurrentConfig(file)
		if err != nil {
			log.Fatalf("Failed to get current config: %v", err)
		}

		// Fetch previous config for comparison
		previousConfigContent, err := gitRepo.GetFileContentFromBranch("refs/remotes/origin/main", file)
		if err != nil {
			log.Fatalf("Failed to get previous config: %v", err)
		}

		var previousConfig Config
		err = yaml.Unmarshal([]byte(previousConfigContent), &previousConfig)
		if err != nil {
			log.Fatalf("Failed to parse previous YAML: %v", err)
		}

		version = currentConfig.Version
		heoRevision = currentConfig.HeoRevision

		// Compare non-version and non-heoRevision fields
		jsonCurrentOtherFields := removeVersionAndHeoRevision(currentConfig)
		jsonPreviousOtherFields := removeVersionAndHeoRevision(&previousConfig)

		fmt.Printf("+version: %s\n", version)
		fmt.Printf("jsonCurrentOtherFields: %s\n", jsonCurrentOtherFields)
		fmt.Printf("jsonPreviousOtherFields: %s\n", jsonPreviousOtherFields)

	}

	// Determine if it is a release version
	isRelease := "true"
	if strings.Contains(version, "-") {
		isRelease = "false"
	}

	// Write outputs to the GITHUB_OUTPUT file
	writeOutput(outputFile, "COMPONENT", component)
	writeOutput(outputFile, "ENVIRONMENT", environment)
	writeOutput(outputFile, "VERSION", version)
	writeOutput(outputFile, "IS_RELEASE", isRelease)
	writeOutput(outputFile, "HEO_REVISION", heoRevision)

	fmt.Printf("COMPONENT=%s\n", component)
	fmt.Printf("ENVIRONMENT=%s\n", environment)
	fmt.Printf("VERSION=%s\n", version)
	fmt.Printf("IS_RELEASE=%s\n", isRelease)
	fmt.Printf("HEO_REVISION=%s\n", heoRevision)
}

// Utility Functions
func getFileContentFromCommit(commit *object.Commit, file string) (string, error) {
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	content, err := entry.Contents()
	if err != nil {
		return "", err
	}

	return content, nil
}

func getCurrentConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func writeOutput(outputFile, key, value string) {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open GITHUB_OUTPUT: %v", err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}
}

func removeVersionAndHeoRevision(config *Config) string {
	config.Version = ""
	config.HeoRevision = ""
	jsonData, _ := json.Marshal(config)
	return string(jsonData)
}
