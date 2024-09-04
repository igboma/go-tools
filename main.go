package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"qgit/deploycheck"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

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
	fmt.Printf("GITHUB OUTOUT %s \n", outputFile)

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
	needDeployment := false

	// Open the repository
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Failed to open repository at %s: %v", repoPath, err)
	}

	// Create GitRepository and ConfigLoader instances
	gitRepo := deploycheck.NewGitRepository(r)
	configLoader := &deploycheck.FileConfigLoader{} // FileConfigLoader for loading the config

	checker := deploycheck.NewVersionChecker(gitRepo)

	// Use OutputWriter to handle the output to GitHub
	outputWriter := deploycheck.NewFileOutputWriter(outputFile)

	var version, heoRevision string
	if action == "closed" && prMerged == "true" {
		fmt.Println("PR is merged...")
		version, heoRevision, err = checker.GetVersionAndHeoRevision("refs/heads/main", file)
		if err != nil {
			log.Fatalf("Failed to get version and heoRevision: %v", err)
		}
	} else {
		fmt.Println("PR is NOT merged...")

		changed, err := deploycheck.CheckVersionAndHeoRevisionDiff(gitRepo, configLoader, file)
		if err != nil {
			log.Fatalf("Error checking version and heoRevision: %v", err)
		}

		needDeployment = changed

		// Fetch current config for the remaining fields
		currentConfig, err := deploycheck.GetCurrentConfig(file)
		if err != nil {
			log.Fatalf("Failed to get current config: %v", err)
		}

		// Fetch previous config for comparison
		previousConfigContent, err := gitRepo.GetFileContentFromBranch("refs/remotes/origin/main", file)
		if err != nil {
			log.Fatalf("Failed to get previous config: %v", err)
		}

		var previousConfig deploycheck.Config
		err = yaml.Unmarshal([]byte(previousConfigContent), &previousConfig)
		if err != nil {
			log.Fatalf("Failed to parse previous YAML: %v", err)
		}

		version = currentConfig.Version
		heoRevision = currentConfig.HeoRevision

		// Compare non-version and non-heoRevision fields
		jsonCurrentOtherFields := deploycheck.RemoveVersionAndHeoRevision(currentConfig)
		jsonPreviousOtherFields := deploycheck.RemoveVersionAndHeoRevision(&previousConfig)

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
	outputWriter.WriteOutput("COMPONENT", component)
	outputWriter.WriteOutput("ENVIRONMENT", environment)
	outputWriter.WriteOutput("VERSION", version)
	outputWriter.WriteOutput("IS_RELEASE", isRelease)
	outputWriter.WriteOutput("HEO_REVISION", heoRevision)
	outputWriter.WriteOutput("DEPLOYMENT_NEEDED", fmt.Sprintf("%t", needDeployment))

	fmt.Printf("COMPONENT=%s\n", component)
	fmt.Printf("ENVIRONMENT=%s\n", environment)
	fmt.Printf("VERSION=%s\n", version)
	fmt.Printf("IS_RELEASE=%s\n", isRelease)
	fmt.Printf("HEO_REVISION=%s\n", heoRevision)
	fmt.Printf("DEPLOYMENT_NEEDED=%s\n", fmt.Sprintf("%t", needDeployment))
}
