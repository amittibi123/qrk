package patch

import (
	"fmt"
	"io"
	"os"
)

// ApplyPatchStream compares an incoming stream with an existing file and writes
// the updated, patched version into a new output file seamlessly.
func ApplyPatchStream(incomingStream io.Reader, existingFilePath string, outputFilePath string) error {
	// Open the existing outdated file for reading
	existingFile, err := os.Open(existingFilePath)
	if err != nil {
		return fmt.Errorf("failed to open existing file: %v", err)
	}
	defer existingFile.Close()

	// Create the new output file where the patched result will be saved
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output patched file: %v", err)
	}
	defer outputFile.Close()

	const chunkSize = 1024 * 1024 // 1MB chunks
	incomingBuffer := make([]byte, chunkSize)
	existingBuffer := make([]byte, chunkSize)

	chunkCounter := 0

	for {
		incomingBytes, errIncoming := incomingStream.Read(incomingBuffer)
		if errIncoming != nil && errIncoming != io.EOF {
			return errIncoming
		}

		existingBytes, errExisting := existingFile.Read(existingBuffer)
		if errExisting != nil && errExisting != io.EOF {
			return errExisting
		}

		// If both streams are empty, the patch operation is complete
		if incomingBytes == 0 && existingBytes == 0 {
			break
		}

		chunkCounter++

		// Check if the current incoming chunk matches the existing chunk perfectly
		chunksMatch := (incomingBytes == existingBytes) &&
			string(incomingBuffer[:incomingBytes]) == string(existingBuffer[:existingBytes])

		if chunksMatch {
			// No changes in this chunk -> Write the existing data to the output file
			_, err = outputFile.Write(existingBuffer[:existingBytes])
			fmt.Printf("✅ Chunk #%d: No changes. Keeping existing data.\n", chunkCounter)
		} else {
			// Delta detected -> Write the new incoming data to the output file (Applies the patch)
			_, err = outputFile.Write(incomingBuffer[:incomingBytes])
			fmt.Printf("⚡ Chunk #%d: Delta applied! Writing new data.\n", chunkCounter)
		}

		if err != nil {
			return fmt.Errorf("failed to write to output file at chunk %d: %v", chunkCounter, err)
		}
	}

	fmt.Println("🎉 Patch applied perfectly! New file generated without loading the full asset into RAM.")
	return nil
}
