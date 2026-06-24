package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

// InitGitRepository initializes a brand new Git repository in the specified path.
// This is the programmatic equivalent of running 'git init' in the terminal.
func InitGitRepository(path string) error {
	fmt.Printf("🌿 Initializing native Git repository at '%s'...\n", path)

	// PlainInit creates a new git repository from scratch.
	// The second argument (isBare) should be false for a standard working directory repo.
	_, err := git.PlainInit(path, false)
	if err != nil {
		// If the repository already exists, go-git returns a specific error
		if err == git.ErrRepositoryAlreadyExists {
			fmt.Println("ℹ️ Git repository already exists in this location. Skipping Git init.")
			return nil
		}
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	fmt.Println("✅ Successfully initialized Git repository programmatically.")
	return nil
}

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

	if err := InitGitRepository(".qrk"); err != nil {
		fmt.Printf("❌ Git Init Warning: %v\n", err)
	}
	fmt.Println("✨ Successfully initialized empty qrk repository in .qrk/")
}
