package main

import (
	"flag"
	"fmt"
	"gitpkg/deploycheck"
	"os"
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

	// go run main.go --workspace "/Users/rza/workspace/helprepo/cd-pipeline" \
	//       --pr-number "38" \
	//       --git-url "https://github.com/igboma/cd-pipeline"

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

	action := os.Getenv("GITHUB_EVENT_ACTION")
	prMerged := os.Getenv("GITHUB_EVENT_PR_MERGED")
	outputFile := os.Getenv("GITHUB_OUTPUT")
	token := os.Getenv("GITHUB_TOKEN")

	// type DeployCheckerOption struct {
	// 	PrNumber   int
	// 	Token      string
	// 	Url        string
	// 	Path       string
	// 	OutputFile string
	// 	Action     string
	// 	PrMerged   string
	// }

	opt := deploycheck.DeployCheckerOption{
		Token:      token,
		PrMerged:   prMerged,
		Action:     action,
		OutputFile: outputFile,
		PrNumber:   prNumber,
		Url:        gitURL,
		Path:       workspace,
	}

	checker := deploycheck.NewDeployChecker(opt)
	checker.Run()
	//deploycheck.DeployCheckRunner(workspace, gitURL, prNumber)
}

// // Main logic
// func deployChecker(directory string, url string, prNumber int) {
// 	// Check if the repository path is passed as a command-line argument
// 	// if len(os.Args) < 2 {
// 	// 	log.Fatalf("Usage: go-tools <repo-path>")
// 	// }

// 	//repoPath := os.Args[1]

// 	// Retrieve environment variables (set in GitHub Actions)
// 	//changedFiles := os.Getenv("CHANGED_FILES")
// 	action := os.Getenv("GITHUB_EVENT_ACTION")
// 	prMerged := os.Getenv("GITHUB_EVENT_PR_MERGED")
// 	outputFile := os.Getenv("GITHUB_OUTPUT")
// 	token := os.Getenv("GITHUB_TOKEN")

// 	fmt.Printf("GITHUB OUTPUT %s \n", outputFile)
// 	fmt.Printf("token %v \n", token)

// 	fmt.Printf("GITHUB directory %v \n", directory)
// 	fmt.Printf("url %v \n", url)
// 	fmt.Printf("prNumber %v \n", prNumber)

// 	options := qgit.QgitOptions{
// 		Url:    url,
// 		Path:   directory,
// 		IsBare: false,
// 		Token:  token,
// 	}
// 	var repo qgit.GitRepository = &qgit.GitRepo{Option: &options}
// 	// Handle error from NewQGit
// 	qGit, err := qgit.NewQGit(&options, repo)
// 	if err != nil {
// 		log.Fatalf("NewQGit err %v", err)
// 	}

// 	files, err := qGit.GetChangedFilesByPRNumberFilesEndingWithYAML(prNumber)

// 	fmt.Printf("GetChangedFilesByPRNumberFilesEndingWithYAML FILES: %s\n", files)

// 	if err != nil {
// 		log.Fatalf("GetChangedFilesByPRNumberFilesEndingWithYAML err %v", err)
// 	}

// 	// Ensure only one file was changed
// 	if len(files) > 1 {
// 		log.Fatalf("More than one file was changed")
// 	}

// 	file := files[0]
// 	fmt.Printf("CHANGED FILES: %s\n", file)
// 	fmt.Printf("FILE: %s\n", file)

// 	// Ensure the file is a conf.yaml file
// 	if !strings.Contains(file, "components/") || !strings.HasSuffix(file, "conf.yaml") {
// 		log.Fatalf("The file is not a conf.yaml file")
// 	}

// 	// Extract COMPONENT and ENVIRONMENT
// 	component := strings.Split(file, "/")[1]
// 	environment := strings.Split(file, "/")[2]
// 	needDeployment := false

// 	// Use OutputWriter to handle the output to GitHub
// 	outputWriter := qgit.NewFileOutputWriter(outputFile)

// 	var version, heoRevision string
// 	if action == "closed" && prMerged == "true" {
// 		fmt.Println("PR is merged...")
// 		configData, err := deploycheck.GetConf(qGit.Repo, file, "refs/heads/main")
// 		fmt.Printf("configData: %v\n", configData)
// 		if err != nil {
// 			log.Fatalf("Failed to get version and heoRevision: %v", err)
// 		}
// 		version = configData.Version
// 		heoRevision = configData.HeoRevision
// 	} else {
// 		fmt.Println("PR is NOT merged...")

// 		source, destination, err := deploycheck.GetSourceAndDestimationConf(qGit.Repo, file, prNumber, "main")

// 		if err != nil {
// 			log.Fatalf("Error checking version and heoRevision: %v", err)
// 		}

// 		fmt.Printf("\nsource: %v\n", source.Version)
// 		fmt.Printf("destination: %v\n", destination.Version)

// 		needDeployment = source.Version != destination.Version || source.HeoRevision != destination.HeoRevision

// 		version = source.Version
// 		heoRevision = source.HeoRevision

// 		// Compare non-version and non-heoRevision fields
// 		jsonCurrentOtherFields := deploycheck.RemoveVersionAndHeoRevision(source)
// 		jsonPreviousOtherFields := deploycheck.RemoveVersionAndHeoRevision(destination)

// 		fmt.Printf("+version: %s\n", version)
// 		fmt.Printf("jsonCurrentOtherFields: %s\n", jsonCurrentOtherFields)
// 		fmt.Printf("jsonPreviousOtherFields: %s\n", jsonPreviousOtherFields)

// 	}

// 	// Determine if it is a release version
// 	isRelease := "true"
// 	if strings.Contains(version, "-") {
// 		isRelease = "false"
// 	}

// 	// Write outputs to the GITHUB_OUTPUT file
// 	outputWriter.WriteOutput("COMPONENT", component)
// 	outputWriter.WriteOutput("ENVIRONMENT", environment)
// 	outputWriter.WriteOutput("VERSION", version)
// 	outputWriter.WriteOutput("IS_RELEASE", isRelease)
// 	outputWriter.WriteOutput("HEO_REVISION", heoRevision)
// 	outputWriter.WriteOutput("DEPLOYMENT_NEEDED", fmt.Sprintf("%t", needDeployment))

// 	fmt.Printf("COMPONENT=%s\n", component)
// 	fmt.Printf("ENVIRONMENT=%s\n", environment)
// 	fmt.Printf("VERSION=%s\n", version)
// 	fmt.Printf("IS_RELEASE=%s\n", isRelease)
// 	fmt.Printf("HEO_REVISION=%s\n", heoRevision)
// 	fmt.Printf("DEPLOYMENT_NEEDED=%s\n", fmt.Sprintf("%t", needDeployment))
// }
