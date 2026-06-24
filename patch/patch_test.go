package patch

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestApplyPatchStream(t *testing.T) {
	// 1. Create temporary files on your Mac disk just for this test
	inMemoryStream := strings.NewReader("Hello, this is the updated chunk data!") // Incoming data stream

	existingFileContent := []byte("Hello, this is the old chunk data!") // Outdated file content
	existingPath := "test_existing.txt"
	outputPath := "test_output.txt"

	// Write the dummy old file to disk so the function can read it
	_ = os.WriteFile(existingPath, existingFileContent, 0644)

	// Clean up these files from your Mac after the test finishes
	defer os.Remove(existingPath)
	defer os.Remove(outputPath)

	type args struct {
		incomingStream   io.Reader
		existingFilePath string
		outputFilePath   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// 2. Here we add our actual test case!
		{
			name: "Successful patch application with delta",
			args: args{
				incomingStream:   inMemoryStream,
				existingFilePath: existingPath,
				outputFilePath:   outputPath,
			},
			wantErr: false, // We expect no errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ApplyPatchStream(tt.args.incomingStream, tt.args.existingFilePath, tt.args.outputFilePath); (err != nil) != tt.wantErr {
				t.Errorf("ApplyPatchStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 3. Optional validation: Read the generated output file and make sure it applied the patch correctly
			patchedContent, err := os.ReadFile(outputPath)
			if err != nil {
				t.Errorf("Failed to read output file: %v", err)
			}

			expectedContent := "Hello, this is the updated chunk data!"
			if string(patchedContent) != expectedContent {
				t.Errorf("Patch was not applied correctly. Got = %s, Want = %s", string(patchedContent), expectedContent)
			}
		})
	}
}
