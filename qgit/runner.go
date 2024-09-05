package qgit

import (
	"fmt"
	"os"
)

// Runner is used to validate the Git Repo. It clones or opens a Git repository and checks out a reference.
func Runner(url, directory, ref, token string) {

	options := QgitOptions{
		Url:   url,
		Path:  directory,
		Token: token,
	}
	var repo GitRepository = &GitRepo{Option: &options}
	// Handle error from NewQGit
	qGit, err := NewQGit(&options, repo)
	if err != nil {
		fmt.Println("Error initializing Git repository:", err)
		os.Exit(1)
	}

	// Perform the checkout operation
	if err := qGit.Checkout(ref); err != nil {
		fmt.Println("Error checking out reference:", err)
		os.Exit(1)
	}

	// Get and print the HEAD reference
	refInfo, err := qGit.Head()
	if err != nil {
		fmt.Println("Error getting HEAD reference:", err)
		os.Exit(1)
	}
	fmt.Println("HEAD Commit:", refInfo.Hash)

	res, err2 := qGit.GetChangedFilesByPRNumberFilesEndingWithYAML(37)
	if err2 != nil {
		fmt.Printf("err2 %v \n", err2)
	}
	fmt.Printf("returned %v \n", res)
}
