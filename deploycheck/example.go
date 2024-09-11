package deploycheck

import (
	"flag"
	"fmt"
	"gitpkg/qgit"
	"gitpkg/utilities"
	"os"
)

func Runner() {
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

	opt := DeployCheckerOption{
		Token:      token,
		PrMerged:   prMerged,
		Action:     action,
		OutputFile: outputFile,
		PrNumber:   prNumber,
		Url:        gitURL,
		Path:       workspace,
	}
	//var repo qgit.Repository = &qgit.QGitRepo{}

	client, _ := qgit.NewClient(
		qgit.WithRepoPath(opt.Path),
		qgit.WithRepoUrl(opt.Url),
	)

	outputWriter := utilities.NewFileOutputWriter(opt.OutputFile)

	checker := &DeployChecker{gitClient: client,
		outputWriter: outputWriter,
		option:       opt}

	//checker := NewDeployChecker(opt)
	checker.Run()
}
