package cmd

import (
	"database/sql"
	"fmt"
	"log"
)

// HandleAdd handles the "qrk add <filename>" command execution.
// It splits the file into simple sequential chunks and stages them in git.
func HandleAdd(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Error: Please specify a file to add.")
		return
	}
	fileName := args[0]

	db, err := sql.Open("sqlite3", ".qrk/my_database.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	if _, err := db.Exec("INSERT OR IGNORE INTO qrk (path) VALUES (?)", fileName); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" successfully added %s to index.txt\n", fileName)
}
