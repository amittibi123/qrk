package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

// HandleAdd handles the "qrk add <filename>" command execution with local git backup.
func HandleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Error: Please specify a file to add.")
		return
	}
	fileName := args[0]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer file.Close()

	repo, err := git.PlainOpen(".qrk")
	if err != nil {
		fmt.Printf("❌ Error: qrk is not initialized. Run 'qrk init' first.\n")
		return
	}
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Printf("❌ Error getting worktree: %v\n", err)
		return
	}

	const chunkSize = 1024 * 1024
	buffer := make([]byte, chunkSize)
	chunkCounter := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return
		}
		if bytesRead == 0 {
			break
		}
		chunkCounter++

		hasher := sha256.New()
		hasher.Write(buffer[:bytesRead])
		chunkHash := fmt.Sprintf("%x", hasher.Sum(nil))

		chunkFolder := ".qrk/chunks"
		_ = os.MkdirAll(chunkFolder, 0755)
		chunkPath := fmt.Sprintf("%s/%s", chunkFolder, chunkHash)

		err = os.WriteFile(chunkPath, buffer[:bytesRead], 0644)
		if err != nil {
			fmt.Printf("❌ Error saving chunk: %v\n", err)
			return
		}

		relativeChunkPath := fmt.Sprintf("chunks/%s", chunkHash)
		_, err = worktree.Add(relativeChunkPath)
		if err != nil {
			fmt.Printf("❌ Git tracking error: %v\n", err)
			return
		}

		fmt.Printf("🚀 Chunk #%d staged and tracked in git database!\n", chunkCounter)
	}

	fmt.Printf("✨ Successfully added all %d chunks to the local database!\n", chunkCounter)
}
