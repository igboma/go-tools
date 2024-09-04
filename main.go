package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
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
	fmt.Println("File:", file)

	// Ensure the file is a conf.yaml file
	if !strings.Contains(file, "components/") || !strings.HasSuffix(file, "conf.yaml") {
		log.Fatalf("The file is not a conf.yaml file")
	}

	// Extract COMPONENT and ENVIRONMENT
	component := strings.Split(file, "/")[1]
	environment := strings.Split(file, "/")[2]

	// Open the repository
	r, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	var version, heoRevision string

	if action == "closed" && prMerged == "true" {
		fmt.Println("PR is merged...")
		version, heoRevision = getVersionAndRevision(r, "refs/heads/main", file)
	} else {
		fmt.Println("PR is NOT merged...")
		diffExists := checkVersionAndHeoRevisionDiff(r, file)

		if !diffExists {
			fmt.Println("No version or heoRevision change. Skipping creation of deployment.")
			os.Exit(0)
		}

		version, heoRevision = getCurrentVersionAndHeoRevision(file)
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

func getVersionAndRevision(r *git.Repository, branch, file string) (string, string) {
	// Ensure 'branch' is a reference name or a specific commit hash
	ref, err := r.Reference(plumbing.ReferenceName(branch), true)
	if err != nil {
		log.Fatalf("Failed to get reference: %v", err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("Error getting commit: %v", err)
	}

	// Use the correct *object.Commit type from "github.com/go-git/go-git/v5/plumbing/object"
	content, err := getFileContentFromCommit(commit, file)
	if err != nil {
		log.Fatalf("Error getting file content: %v", err)
	}

	version := parseYAML(content, "version")
	heoRevision := parseYAML(content, "heoRevision")
	return version, heoRevision
}

func getFileContentFromCommit(commit *object.Commit, file string) (string, error) {
	// Function to get file content from a specific commit
	// Placeholder for now - replace with logic to read file contents from the commit object
	return "", nil
}

func checkVersionAndHeoRevisionDiff(r *git.Repository, file string) bool {
	// Simulate git diff to check for changes in version or heoRevision
	// Placeholder - you may implement go-git-based diff logic here
	return true
}

func getCurrentVersionAndHeoRevision(file string) (string, string) {
	// Use YAML parsing to extract version and heoRevision
	// Simulate extraction of version and heoRevision from the current file
	version := "1.0.0"       // Dummy value, replace with actual parsing
	heoRevision := "heo-123" // Dummy value, replace with actual parsing
	return version, heoRevision
}

func parseYAML(content, key string) string {
	// Simulate YAML parsing to extract a value by key
	// Placeholder for actual YAML parsing logic using Go libraries
	return ""
}
