package main

import (
	"fmt"
	"os"

	// Importing our custom internal command package
	"qrk/cmd"
)

func main() {
	// os.Args[0] is always the program name. We need at least one more argument (the command).
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Extract the main tracking command (e.g., "init", "add")
	command := os.Args[1]

	switch command {
	case "init":
		cmd.HandleInit()

	case "add":
		// Pass only the remaining arguments to the add handler (the file path)
		cmd.HandleAdd(os.Args[2:])

	default:
		fmt.Printf("❌ Unknown command: '%s'\n", command)
		fmt.Println("Run 'qrk' without arguments to see available commands.")
	}
}

// printUsage displays helpful CLI instructions when the user inputs invalid commands.
func printUsage() {
	fmt.Println("Usage: qrk <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  init    Initialize a new high-performance repository")
	fmt.Println("  add     Track any large binary or file seamlessly")
}
