package qgit_2

import (
	"fmt"
	"os"
)

// Runner demonstrates how to use the qgit package to interact with a Git repository.
// It clones or opens a Git repository, checks out a reference, and performs some basic operations.
func Runner(url, directory, ref, token string) {
	// Step 1: Set up repository options
	options := QRepoOptions{
		Url:   url,       // Git repository URL
		Path:  directory, // Local directory path to clone or open the repository
		Token: token,     // Authentication token (e.g., GitHub personal access token)
	}

	// Step 2: Initialize the repository using qgit
	var repo Repository = &QGitRepo{}
	qGit := NewQGit(&options, repo)

	// Step 3: Perform the checkout operation to switch to the given reference (branch, tag, or commit)
	if err := qGit.Checkout(ref); err != nil {
		fmt.Println("Error checking out reference:", err)
		os.Exit(1) // Exit if there's an error during checkout
	}
	fmt.Printf("Checked out reference: %s\n", ref)

	// Step 4: Retrieve the current HEAD reference and print its commit hash
	refInfo, err := qGit.Head()
	if err != nil {
		fmt.Println("Error getting HEAD reference:", err)
		os.Exit(1) // Exit if there's an error retrieving the HEAD reference
	}
	fmt.Printf("Current HEAD Commit: %s\n", refInfo.Hash)

	// Step 5: Retrieve the list of files ending with "conf.yaml" from a pull request (PR) number
	prNumber := 37 // Example PR number
	files, err := qGit.GetConfFileChangedByPRNumber(prNumber)
	if err != nil {
		fmt.Printf("Error retrieving changed files: %v\n", err)
		os.Exit(1) // Exit if there's an error retrieving the files
	}
	fmt.Printf("Changed files in PR #%d ending with 'conf.yaml': %v\n", prNumber, files)
}

// func main() {
// 	// Example usage of the Runner function
// 	// Replace the following values with actual repository details
// 	repoUrl := "https://github.com/example/repo.git"
// 	directory := "/path/to/local/repo"
// 	reference := "refs/heads/main" // Could be a branch, tag, or commit hash
// 	token := "your-github-token"   // Authentication token for accessing the repository

// 	// Call Runner to demonstrate how to use qgit
// 	Runner(repoUrl, directory, reference, token)
// }
