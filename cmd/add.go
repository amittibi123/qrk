package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// HandleAdd handles the "qrk add <filename>" command execution.
// It streams the file in chunks to prevent high memory consumption and prepares it for tracking.
func HandleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Error: Please specify a file to add.")
		fmt.Println("Usage: qrk add <filename>")
		return
	}

	fileName := args[0]

	// Open the target file for reading
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("❌ Error: Could not open file '%s': %v\n", fileName, err)
		return
	}
	// Ensure the file descriptor is closed when the function finishes
	defer file.Close()

	fmt.Printf("📦 Tracking file: %s\n", fileName)
	fmt.Println("⏳ Processing file chunks...")

	// Define a fixed chunk size (e.g., 1MB = 1024 * 1024 bytes)
	// In the future, this can be upgraded to dynamic content-defined chunking (CDC)
	const chunkSize = 1024 * 1024
	buffer := make([]byte, chunkSize)

	chunkCounter := 0

	for {
		// Read up to chunkSize bytes from the file into our buffer
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("❌ Error reading file data: %v\n", err)
			return
		}

		// EOF (End of File) means we have successfully finished reading the entire file
		if bytesRead == 0 {
			break
		}

		chunkCounter++

		// Generate a unique SHA-256 hash for the current slice of data (the chunk)
		// This hash acts as the unique identifier for deduplication
		hasher := sha256.New()
		hasher.Write(buffer[:bytesRead])
		chunkHash := fmt.Sprintf("%x", hasher.Sum(nil))

		// Log the processed chunk (printing only the first 12 characters of the hash for clean output)
		fmt.Printf("  -> Processed Chunk #%d | Size: %d bytes | Hash: %s...\n", chunkCounter, bytesRead, chunkHash[:12])
	}

	fmt.Printf("✨ Successfully processed %d chunks for '%s'.\n", chunkCounter, fileName)
}
