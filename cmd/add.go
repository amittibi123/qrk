package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

// HandleAdd handles the "qrk add <filename>" command execution.
// It splits the file into simple sequential chunks and stages them in git.
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

	// Open the local git repository inside .qrk
	repo, err := git.PlainOpen(".qrk")
	if err != nil {
		fmt.Printf("❌ Error: qrk is not initialized. Run 'qrk init' first.\n")
		return
	}
	_, err = repo.Worktree()
	if err != nil {
		fmt.Printf("❌ Error getting worktree: %v\n", err)
		return
	}

	const chunkSize = 1024 * 1024 // 1MB
	buffer := make([]byte, chunkSize)
	chunkCounter := 0

	// Create the chunks directory if it doesn't exist
	chunkFolder := ".qrk/chunks"
	_ = os.MkdirAll(chunkFolder, 0755)

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("❌ Error reading file: %v\n", err)
			return
		}
		if bytesRead == 0 {
			break
		}
		chunkCounter++

		// Generate sequential chunk name based on original filename and counter
		chunkName := fmt.Sprintf("%s_chunk_%d", fileName, chunkCounter)
		chunkPath := fmt.Sprintf("%s/%s", chunkFolder, chunkName)

		// Write chunk data directly to disk
		err = os.WriteFile(chunkPath, buffer[:bytesRead], 0644)
		if err != nil {
			fmt.Printf("❌ Error saving chunk: %v\n", err)
			return
		}

		fmt.Printf("🚀 Chunk #%d (%s) staged successfully!\n", chunkCounter, chunkName)
	}

	f, err := os.OpenFile(".qrk/index.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, fileName)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf(" successfully added %s to index.txt\n", fileName)
	fmt.Printf("✨ Successfully added all %d chunks for '%s' to the local database!\n", chunkCounter, fileName)
}
