package cmd

import (
	"fmt"
	"os"
)

// HandleInit handles the "qrk init" command execution.
// It creates the hidden .qrk directory to initialize the repository.
func HandleInit() {
	fmt.Println("🚀 Initializing qrk repository...")

	// Define the path for the hidden metadata directory
	repoDir := ".qrk"

	// Check if the repository is already initialized
	if _, err := os.Stat(repoDir); !os.IsNotExist(err) {
		fmt.Println("❌ Error: Repository is already initialized in this directory.")
		return
	}

	// Create the hidden directory with read/write/execute permissions for the owner
	err := os.Mkdir(repoDir, 0755)
	if err != nil {
		fmt.Printf("❌ Error creating metadata directory: %v\n", err)
		return
	}

	fmt.Println("✨ Successfully initialized empty qrk repository in .qrk/")
}
