package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gitpkg/deploycheck"
	"gitpkg/qgit"

	"gopkg.in/yaml.v2"
)

func main() {
	// if len(os.Args) < 4 {
	// 	fmt.Println("Usage: <url> <directory> <ref>")
	// 	os.Exit(1)
	// }
	// url, directory, ref := os.Args[1], os.Args[2], os.Args[3]

	// token := os.Getenv("GITHUB_TOKEN")
	// fmt.Printf("token value %v \n", token)

	// qgit.Runner(url, directory, ref, token)

	//go run main.go https://github.com/igboma/cd-pipeline /Users/rza/workspace/helprepo/cd-pipeline  main
	//go run main.go https://github.com/qlik-trial/helm-environment-overrides /Users/rza/workspace/helprepo/override  main
	//https://github.com/qlik-trial/helm-environment-overrides

	var workspace string
	var prNumber int
	var gitURL string

	// Bind the flags to variables
	flag.StringVar(&workspace, "workspace", "", "The GitHub workspace")
	flag.IntVar(&prNumber, "pr-number", 0, "The Pull Request number")
	flag.StringVar(&gitURL, "git-url", "", "The Git URL of the PR")

	// Parse the command-line flags
	flag.Parse()

	// Check if required flags are passed
	if workspace == "" || gitURL == "" || prNumber == 0 {
		fmt.Println("Missing required flags: --workspace, --pr-number, and --git-url must all be provided.")
		return
	}

	// Use the flag values in your program
	fmt.Printf("Workspace: %s\n", workspace)
	fmt.Printf("PR Number: %d\n", prNumber)
	fmt.Printf("Git URL: %s\n", gitURL)

	deployChecker(workspace, gitURL, prNumber)
}

// Main logic
func deployChecker(directory string, url string, prNumber int) {
	// Check if the repository path is passed as a command-line argument
	// if len(os.Args) < 2 {
	// 	log.Fatalf("Usage: go-tools <repo-path>")
	// }

	//repoPath := os.Args[1]

	// Retrieve environment variables (set in GitHub Actions)
	//changedFiles := os.Getenv("CHANGED_FILES")
	action := os.Getenv("GITHUB_EVENT_ACTION")
	prMerged := os.Getenv("GITHUB_EVENT_PR_MERGED")
	outputFile := os.Getenv("GITHUB_OUTPUT")
	token := os.Getenv("GITHUB_TOKEN")
	fmt.Printf("GITHUB OUTOUT %s \n", outputFile)
	fmt.Printf("token %s \n", token)

	options := qgit.QgitOptions{
		Url:    url,
		Path:   directory,
		IsBare: false,
		Token:  token,
	}
	var repo qgit.GitRepository = &qgit.GitRepo{Option: &options}
	// Handle error from NewQGit
	qGit, err := qgit.NewQGit(&options, repo)
	if err != nil {
		log.Fatalf("NewQGit err %v", err)
	}

	files, err := qGit.GetChangedFilesByPRNumberFilesEndingWithYAML(prNumber)

	if err != nil {
		log.Fatalf("GetChangedFilesByPRNumberFilesEndingWithYAML err %v", err)
	}

	// Split the changed files into an array
	//files := strings.Split(changedFiles, "\n")

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

	// // Open the repository
	// _, err := git.PlainOpen(directory)
	// if err != nil {
	// 	log.Fatalf("Failed to open repository at %s: %v", directory, err)
	// }

	// Create GitRepository and ConfigLoader instances
	//gitRepo := deploycheck.NewGitRepository(r)
	configLoader := &deploycheck.FileConfigLoader{} // FileConfigLoader for loading the config

	checker := deploycheck.NewVersionChecker(qGit.Repo)

	// Use OutputWriter to handle the output to GitHub
	outputWriter := qgit.NewFileOutputWriter(outputFile)

	var version, heoRevision string
	if action == "closed" && prMerged == "true" {
		fmt.Println("PR is merged...")
		version, heoRevision, err = checker.GetVersionAndHeoRevision("refs/heads/main", file)
		if err != nil {
			log.Fatalf("Failed to get version and heoRevision: %v", err)
		}
	} else {
		fmt.Println("PR is NOT merged...")
		changed, err := deploycheck.CheckVersionAndHeoRevisionDiff(qGit.Repo, configLoader, file)
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
		previousConfigContent, err := qGit.Repo.GetFileContentFromBranch("refs/remotes/origin/main", file)
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
