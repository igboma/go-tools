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

// Define the Config struct to match the conf.yaml structure exactly
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

	fmt.Println("The file is a conf.yaml file, because its path is 'components/.../conf.yaml'")
	fmt.Printf("FILE: %s\n", file)

	// Extract COMPONENT and ENVIRONMENT
	component := strings.Split(file, "/")[1]
	environment := strings.Split(file, "/")[2]

	// Open the repository using the provided path
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Failed to open repository at %s: %v", repoPath, err)
	}

	var version, heoRevision string
	if action == "closed" && prMerged == "true" {
		fmt.Println("PR is merged...")
		version, heoRevision = getVersionAndHeoRevision(r, "refs/heads/main", file)
	} else {
		fmt.Println("PR is NOT merged...")

		if !checkVersionAndHeoRevisionDiff(r, file) {
			fmt.Println("No version or heoRevision change. Skipping creation of deployment.")
			os.Exit(0)
		}

		currentConfig, _ := getCurrentConfig(file)
		previousConfig, _ := getPreviousConfig(r, "refs/remotes/origin/main", file)

		version = currentConfig.Version
		heoRevision = currentConfig.HeoRevision

		// Compare non-version and non-heoRevision fields
		jsonCurrentOtherFields := removeVersionAndHeoRevision(currentConfig)
		jsonPreviousOtherFields := removeVersionAndHeoRevision(previousConfig)

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

// getVersionAndHeoRevision retrieves the version and heoRevision from the specified branch
func getVersionAndHeoRevision(r *git.Repository, branch, file string) (string, string) {
	// Get the reference for the branch
	ref, err := r.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		log.Fatalf("Failed to get reference for branch %s: %v", branch, err)
	}

	// Get the commit for the branch
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("Failed to get commit for branch %s: %v", branch, err)
	}

	// Get the file content from the commit
	content, err := getFileContentFromCommit(commit, file)
	if err != nil {
		log.Fatalf("Failed to get file content from branch %s: %v", branch, err)
	}

	// Parse the file content as YAML and extract the version and heoRevision
	var config Config
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML content from branch %s: %v", branch, err)
	}

	// Return the version and heoRevision
	return config.Version, config.HeoRevision
}

func getCurrentConfig(file string) (*Config, error) {
	// Read the current YAML file to get the config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", file, err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}
	return &config, nil
}

func getPreviousConfig(r *git.Repository, branch, file string) (*Config, error) {
	// Get the file content from the previous commit (from the specified branch)
	ref, err := r.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		log.Fatalf("Failed to get reference: %v", err)
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("Error getting commit: %v", err)
	}
	content, err := getFileContentFromCommit(commit, file)
	if err != nil {
		log.Fatalf("Error getting file content: %v", err)
	}
	var config Config
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}
	return &config, nil
}

func getFileContentFromCommit(commit *object.Commit, file string) (string, error) {
	// Get the tree object from the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	// Get the file entry
	entry, err := tree.File(file)
	if err != nil {
		return "", err
	}

	// Read the file content
	content, err := entry.Contents()
	if err != nil {
		return "", err
	}

	return content, nil
}

func removeVersionAndHeoRevision(config *Config) string {
	// Remove version and heoRevision fields and return JSON representation of other fields
	config.Version = ""
	config.HeoRevision = ""
	jsonData, _ := json.Marshal(config)
	return string(jsonData)
}

func writeOutput(outputFile, key, value string) {
	// Append output key=value to the GITHUB_OUTPUT file
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to open GITHUB_OUTPUT: %v", err)
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}
}

func checkVersionAndHeoRevisionDiff(r *git.Repository, file string) bool {
	// Get the current worktree (the state of the repository as currently checked out)
	worktree, err := r.Worktree()
	if err != nil {
		log.Fatalf("Failed to get worktree: %v", err)
	}

	// Get the status of the working tree to see if there are changes to the file
	status, err := worktree.Status()
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	// Check if the specific file has changes in the working directory
	if status.File(file).Worktree != git.Unmodified {
		fmt.Println("File has changes in the working tree")
	}

	// Get the commit of the main branch (origin/main)
	ref, err := r.Reference(plumbing.ReferenceName("refs/remotes/origin/main"), true)
	if err != nil {
		log.Fatalf("Failed to get reference for origin/main: %v", err)
	}

	// Get the commit object for the reference
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("Failed to get commit object for origin/main: %v", err)
	}

	// Get the file content from origin/main
	previousContent, err := getFileContentFromCommit(commit, file)
	if err != nil {
		log.Fatalf("Failed to get previous file content from origin/main: %v", err)
	}

	// Get the current file content from the worktree (the file as it is now)
	currentData, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read the current file %s: %v", file, err)
	}

	// Parse both current and previous versions as YAML
	var currentConfig, previousConfig Config
	err = yaml.Unmarshal(currentData, &currentConfig)
	if err != nil {
		log.Fatalf("Failed to parse current YAML: %v", err)
	}
	err = yaml.Unmarshal([]byte(previousContent), &previousConfig)
	if err != nil {
		log.Fatalf("Failed to parse previous YAML: %v", err)
	}

	// Compare the version and heoRevision between the two
	if currentConfig.Version != previousConfig.Version || currentConfig.HeoRevision != previousConfig.HeoRevision {
		fmt.Println("Version or heoRevision has changed")
		return true
	}

	fmt.Println("No changes in version or heoRevision")
	return false
}
