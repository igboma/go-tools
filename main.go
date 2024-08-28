package main

import (
	"fmt"
	"os"
	"qgit/gitpkg" // Import the qgit package from your githelper module
)

// Main is the entry point for the program. It clones or opens a Git repository and checks out a reference.
func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: <url> <directory> <ref>")
		os.Exit(1)
	}
	url, directory, ref := os.Args[1], os.Args[2], os.Args[3]

	var repo gitpkg.GitRepository = &gitpkg.GitRepo{}
	option := gitpkg.QgitOptions{
		Url:    url,
		Path:   directory,
		IsBare: false,
	}

	// Handle error from NewQGit
	qGit, err := gitpkg.NewQGit(option, repo)
	if err != nil {
		fmt.Println("Error initializing Git repository:", err)
		os.Exit(1)
	}

	// if err := gitpkg.Fetch(""); err != nil {
	// 	fmt.Println("Error fetching remote references:", err)
	// 	os.Exit(1)
	// }

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

	err1 := qGit.PR()

	fmt.Println("err PR:", err1)

}
