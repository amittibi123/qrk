package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"qrk/patch"

	"github.com/go-git/go-git/v5"
	"github.com/kr/binarydist"
	_ "github.com/mattn/go-sqlite3"
)

const (
	headTarPath  = ".qrk/HEAD.tar"
	workTarPath  = ".qrk/work.tar.tmp"
	patchDirPath = ".qrk/patches"
)

func HandleCommit(message string) {
	fmt.Println("commiting...")

	db, err := sql.Open("sqlite3", ".qrk/my_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var allEntries []string
	rows, err := db.Query("SELECT path FROM qrk")
	if err != nil {
		fmt.Println("Error reading from DB:", err)
		return
	}
	for rows.Next() {
		var filePath string
		if err := rows.Scan(&filePath); err != nil {
			fmt.Println("Error scanning row:", err)
			rows.Close()
			return
		}
		allEntries = append(allEntries, filePath)
	}
	rows.Close()

	if len(allEntries) == 0 {
		fmt.Println("ℹ️ Nothing to commit.")
		return
	}

	// 1. צילום worktree לתוך tar - בלי לזוז/למחוק קבצים מקוריים
	if err := patch.CreateTarFromFiles(allEntries, workTarPath); err != nil {
		fmt.Println(err)
		return
	}

	if err := os.MkdirAll(patchDirPath, 0755); err != nil {
		fmt.Println(err)
		os.Remove(workTarPath)
		return
	}
	patchFilePath := filepath.Join(patchDirPath, fmt.Sprintf("%d.patch", time.Now().UnixNano()))

	if _, err := os.Stat(headTarPath); os.IsNotExist(err) {
		// commit ראשון: אין head, הפאץ' הוא ה-tar עצמו
		if err := copyFile(workTarPath, patchFilePath); err != nil {
			fmt.Println(err)
			os.Remove(workTarPath)
			return
		}
	} else {
		oldF, err := os.Open(headTarPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		newF, err := os.Open(workTarPath)
		if err != nil {
			oldF.Close()
			fmt.Println(err)
			return
		}
		patchF, err := os.Create(patchFilePath)
		if err != nil {
			oldF.Close()
			newF.Close()
			fmt.Println(err)
			return
		}

		err = binarydist.Diff(oldF, newF, patchF)
		oldF.Close()
		newF.Close()
		patchF.Close()
		if err != nil {
			fmt.Println(err)
			os.Remove(workTarPath)
			os.Remove(patchFilePath)
			return
		}
	}

	// 2. ה-worktree tar החדש הופך ל-HEAD החדש
	if err := os.Rename(workTarPath, headTarPath); err != nil {
		fmt.Println(err)
		return
	}

	// 3. git add + commit - רק קובץ הפאץ' הבודד
	repo, err := git.PlainOpen(".qrk")
	if err != nil {
		fmt.Println(err)
		return
	}
	wt, err := repo.Worktree()
	if err != nil {
		fmt.Println(err)
		return
	}
	relPatch, _ := filepath.Rel(".qrk", patchFilePath)
	if _, err := wt.Add(relPatch); err != nil {
		fmt.Println(err)
		return
	}
	if _, err := wt.Commit(message, &git.CommitOptions{}); err != nil {
		fmt.Println(err)
		return
	}

	// הפאץ' כבר שמור בתוך git (.qrk/.git) - אין צורך להשאיר עותק על הדיסק
	if err := os.Remove(patchFilePath); err != nil {
		fmt.Println("Warning: failed to remove staged patch file:", err)
	}

	// 4. ריקון הטבלה לקראת ה-staging הבא
	if _, err := db.Exec("DELETE FROM qrk"); err != nil {
		fmt.Println("Warning: failed to clear staged entries:", err)
	}

	fmt.Println("✅ commit done")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
