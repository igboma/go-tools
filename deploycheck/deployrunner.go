package deploycheck

import (
	"fmt"
	"gitpkg/qgit"
	"log"
	"os"
	"strings"
)

func DeployChecker(directory string, url string, prNumber int) {
	action := os.Getenv("GITHUB_EVENT_ACTION")
	prMerged := os.Getenv("GITHUB_EVENT_PR_MERGED")
	outputFile := os.Getenv("GITHUB_OUTPUT")
	token := os.Getenv("GITHUB_TOKEN")

	fmt.Printf("GITHUB OUTPUT %s \n", outputFile)

	fmt.Printf("GITHUB directory %v \n", directory)
	fmt.Printf("url %v \n", url)
	fmt.Printf("prNumber %v \n", prNumber)

	options := qgit.QgitOptions{
		Url:   url,
		Path:  directory,
		Token: token,
	}
	var repo qgit.GitRepository = &qgit.GitRepo{Option: &options}
	// Handle error from NewQGit
	qGit, err := qgit.NewQGit(&options, repo)
	if err != nil {
		log.Fatalf("NewQGit err %v", err)
	}

	files, err := qGit.GetChangedFilesByPRNumberFilesEndingWithYAML(prNumber)

	fmt.Printf("GetChangedFilesByPRNumberFilesEndingWithYAML FILES: %s\n", files)

	if err != nil {
		log.Fatalf("GetChangedFilesByPRNumberFilesEndingWithYAML err %v", err)
	}

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

	// Use OutputWriter to handle the output to GitHub
	outputWriter := qgit.NewFileOutputWriter(outputFile)

	var version, heoRevision string
	if action == "closed" && prMerged == "true" {
		fmt.Println("PR is merged...")
		configData, err := GetConf(qGit.Repo, file, "refs/heads/main")
		fmt.Printf("configData: %v\n", configData)
		if err != nil {
			log.Fatalf("Failed to get version and heoRevision: %v", err)
		}
		version = configData.Version
		heoRevision = configData.HeoRevision
	} else {
		fmt.Println("PR is NOT merged...")

		source, destination, err := GetSourceAndDestimationConf(qGit.Repo, file, prNumber, "main")

		if err != nil {
			log.Fatalf("Error checking version and heoRevision: %v", err)
		}

		fmt.Printf("\nsource: %v\n", source.Version)
		fmt.Printf("destination: %v\n", destination.Version)

		needDeployment = source.Version != destination.Version || source.HeoRevision != destination.HeoRevision

		version = source.Version
		heoRevision = source.HeoRevision

		// Compare non-version and non-heoRevision fields
		jsonCurrentOtherFields := RemoveVersionAndHeoRevision(source)
		jsonPreviousOtherFields := RemoveVersionAndHeoRevision(destination)

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
