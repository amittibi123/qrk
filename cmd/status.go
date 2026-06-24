package cmd

import (
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
	status := HandleStatus()
	fmt.Println("🟢 Staged changes ready for commit in qrk database:")
	// 5. רצים על כל הקבצים/צ'אנקים שגיט זיהה ב-Staging
	for path, fileStatus := range status {
		// אנחנו מחפשים קבצים שנוספו או השתנו ב-Staging Area
		if fileStatus.Staging == git.Added || fileStatus.Staging == git.Modified {
			fmt.Printf("  • [Staged] %s\n", path)
		}
	}

}
