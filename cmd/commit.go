package cmd

import (
	"fmt"
	"os"
	"qrk/patch"
)

func HandleCommit() {
	fmt.Println("commiting...")

	sourcePath := ".qrk/chunks"
	newPath := ".qrk/chunks.tar"

	if err := patch.CreateTarFromDirectory(sourcePath, newPath); err != nil {
		panic(err)
	}

	if err := os.RemoveAll(sourcePath); err != nil {
		panic(err)
	}

}
