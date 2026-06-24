package cmd

import (
	"os"
	"testing"
)

// TestHandleInit verifies that HandleInit successfully creates the .qrk directory
// and properly handles cases where the directory already exists.
func TestHandleInit(t *testing.T) {
	// 1. Setup: Ensure we start in a clean state by removing any existing .qrk directory
	repoDir := ".qrk"
	os.RemoveAll(repoDir)

	// 2. Test initial creation
	HandleInit()

	// Verify that the directory was actually created
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		t.Errorf("Expected directory '%s' to be created, but it does not exist", repoDir)
	}

	// 3. Test duplicate initialization (should fail gracefully without crashing)
	// We run it a second time to ensure our safeguard logic works
	HandleInit()

	// 4. Teardown: Clean up the created directory after the test finishes
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("Failed to clean up test directory: %v", err)
	}
}
