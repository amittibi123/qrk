package patch

import (
	"fmt"
	"io"
	"os"
)

// ApplyPatch compares a "new/patch" file with an "existing" file
// and writes the merged result into an output file.
func ApplyPatchStream(patchFilePath string, existingFilePath string, outputFilePath string) error {
	// 1. Open the patch file (the incoming data)
	patchFile, err := os.Open(patchFilePath)
	if err != nil {
		return fmt.Errorf("failed to open patch file: %v", err)
	}
	defer patchFile.Close()

	// 2. Open the existing outdated file
	existingFile, err := os.Open(existingFilePath)
	if err != nil {
		return fmt.Errorf("failed to open existing file: %v", err)
	}
	defer existingFile.Close()

	// 3. Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	const chunkSize = 1024 * 1024 // 1MB
	patchBuffer := make([]byte, chunkSize)
	existingBuffer := make([]byte, chunkSize)

	chunkCounter := 0

	for {
		// Read from both files
		patchBytes, errPatch := patchFile.Read(patchBuffer)
		if errPatch != nil && errPatch != io.EOF {
			return errPatch
		}

		existingBytes, errExisting := existingFile.Read(existingBuffer)
		if errExisting != nil && errExisting != io.EOF {
			return errExisting
		}

		if patchBytes == 0 && existingBytes == 0 {
			break
		}

		chunkCounter++

		// Compare chunks
		// Using bytes.Equal is much faster than converting to string!
		chunksMatch := (patchBytes == existingBytes) &&
			(string(patchBuffer[:patchBytes]) == string(existingBuffer[:existingBytes]))

		if chunksMatch {
			_, err = outputFile.Write(existingBuffer[:existingBytes])
			fmt.Printf("✅ Chunk #%d: No changes.\n", chunkCounter)
		} else {
			_, err = outputFile.Write(patchBuffer[:patchBytes])
			fmt.Printf("⚡ Chunk #%d: Delta applied!\n", chunkCounter)
		}

		if err != nil {
			return fmt.Errorf("failed writing to output at chunk %d: %v", chunkCounter, err)
		}
	}

	return nil
}
