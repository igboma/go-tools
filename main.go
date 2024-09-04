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
	CheckDiff(branch, file string) (bool, error)
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

// CheckDiff checks if there is any difference between the current branch and the specified branch for the given file.
func (g *GitRepositoryImpl) CheckDiff(branch, file string) (bool, error) {
	// Get the current worktree
	worktree, err := g.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %v", err)
	}

	// Get the status of the worktree
	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree status: %v", err)
	}

	// Check if the file has any changes in the working directory
	if status.File(file).Worktree != git.Unmodified {
		return true, nil
	}

	// Get the current branch
	headRef, err := g.repo.Head()
	if err != nil {
		return false, fmt.Errorf("failed to get HEAD reference: %v", err)
	}

	// Get the branch to compare with
	branchRef, err := g.repo.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		return false, fmt.Errorf("failed to get reference for branch %s: %v", branch, err)
	}

	// Compare commits between the current branch and the specified branch
	commitsIter, err := g.repo.Log(&git.LogOptions{From: branchRef.Hash()})
	if err != nil {
		return false, fmt.Errorf("failed to retrieve commit log: %v", err)
	}

	var currentCommit *object.Commit
	for {
		commit, err := commitsIter.Next()
		if err != nil {
			break
		}

		if commit.Hash == headRef.Hash() {
			currentCommit = commit
			break
		}
	}

	if currentCommit == nil {
		return false, nil
	}

	// Get the contents of the file in both branches and compare
	branchContent, err := g.GetFileContentFromBranch(branch, file)
	if err != nil {
		return false, fmt.Errorf("failed to get file content from branch: %v", err)
	}

	currentContent, err := g.GetFileContentFromCommit(currentCommit.Hash.String(), file)
	if err != nil {
		return false, fmt.Errorf("failed to get current file content: %v", err)
	}

	// Compare the content
	return branchContent != currentContent, nil
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

func (v *VersionChecker) CheckVersionDiff(branch, file string) (bool, error) {
	changed, err := v.repo.CheckDiff(branch, file)
	if err != nil {
		return false, err
	}
	return changed, nil
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

		changed, err := checker.CheckVersionDiff("refs/remotes/origin/main", file)
		if err != nil {
			log.Fatalf("Failed to check version diff: %v", err)
		}

		if !changed {
			fmt.Println("No version or heoRevision change. Skipping creation of deployment.")
			os.Exit(0)
		}

		currentConfig, err := getCurrentConfig(file)
		if err != nil {
			log.Fatalf("Failed to get current config: %v", err)
		}

		previousConfigContent, err := gitRepo.GetFileContentFromBranch("refs/remotes/origin/main", file)
		if err != nil {
			log.Fatalf("Failed to get previous config: %v", err)
		}

		var previousConfig Config
		err = yaml.Unmarshal([]byte(previousConfigContent), &previousConfig)
		if err != nil {
			log.Fatalf("Failed to parse previous YAML: %v", err)
		}

		// Compare the other fields (excluding version and heoRevision)
		jsonCurrentOtherFields := removeVersionAndHeoRevision(currentConfig)
		jsonPreviousOtherFields := removeVersionAndHeoRevision(&previousConfig)

		fmt.Printf("+version: %s\n", currentConfig.Version)
		fmt.Printf("jsonCurrentOtherFields: %s\n", jsonCurrentOtherFields)
		fmt.Printf("jsonPreviousOtherFields: %s\n", jsonPreviousOtherFields)

		version = currentConfig.Version
		heoRevision = currentConfig.HeoRevision
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
