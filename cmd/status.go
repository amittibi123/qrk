package cmd

import (
	"database/sql"
	"fmt"

	"github.com/go-git/go-git/v5"
)

// HandleStatus checks the local staging area inside .qrk
// and prints which chunks are ready to be committed.
func HandleStatus() git.Status {
	// 1. פותחים את התיקייה המקומית (100% אופליין)
	repo, err := git.PlainOpen(".qrk")
	if err != nil {
		fmt.Println("❌ Error: qrk repository not found. Run 'qrk init' first.")
		return nil
	}

	// 2. קוראים את ה-Worktree הנוכחי
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Printf("❌ Error getting worktree: %v\n", err)
		return nil
	}

	// 3. שולפים את הסטטוס של ה-Index
	status, err := worktree.Status()
	if err != nil {
		fmt.Printf("❌ Error reading status: %v\n", err)
		return nil
	}

	// 4. בודקים אם ה-Index ריק
	if status.IsClean() {
		fmt.Println("ℹ️ Nothing staged. Your workspace is clean. Use 'qrk add <file>' to track files.")
		return nil
	}
	return status

}

func PrintStatus() {
	db, err := sql.Open("sqlite3", ".qrk/my_database.db")
	if err != nil {
		fmt.Println("❌ Error: qrk repository not found. Run 'qrk init' first.")
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT path FROM qrk")
	if err != nil {
		fmt.Println("❌ Error reading staged files:", err)
		return
	}
	defer rows.Close()

	var entries []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		entries = append(entries, p)
	}

	if len(entries) == 0 {
		fmt.Println("ℹ️ Nothing staged. Use 'qrk add <file>' to track files.")
		return
	}

	fmt.Println("🟢 Staged changes ready for commit:")
	for _, p := range entries {
		fmt.Printf("  • [Staged] %s\n", p)
	}
}
