package cmd

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/kr/binarydist"
)

func untarToWorktree(tarData []byte) error {
	tr := tar.NewReader(bytes.NewReader(tarData))
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if dir := filepath.Dir(hdr.Name); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
		f, err := os.Create(hdr.Name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return err
		}
		f.Close()
		fmt.Printf("✅ restored %s\n", hdr.Name)
	}
	return nil
}

func HandleHead(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Error: Please specify how many commits to go back.")
		return
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 0 {
		fmt.Println("❌ Error: argument must be a non-negative number.")
		return
	}
	fmt.Printf("going back %d commits\n", n)

	// 1. גיבוי HEAD.tar לפני שה-checkout עלול למחוק אותו
	backupPath := ".qrk_HEAD_backup.tar"
	hadBackup := false
	if _, err := os.Stat(headTarPath); err == nil {
		if err := copyFile(headTarPath, backupPath); err != nil {
			fmt.Println("Warning: failed to backup HEAD.tar:", err)
		} else {
			hadBackup = true
		}
	}

	repo, err := git.PlainOpen(".qrk")
	if err != nil {
		fmt.Println(err)
		return
	}
	head, err := repo.Head()
	if err != nil {
		fmt.Println(err)
		return
	}
	commitIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		fmt.Println(err)
		return
	}
	var targetHash plumbing.Hash
	found := false
	count := 0
	err = commitIter.ForEach(func(c *object.Commit) error {
		if count == n {
			targetHash = c.Hash
			found = true
			return storer.ErrStop
		}
		count++
		return nil
	})
	if err != nil && err != storer.ErrStop {
		fmt.Println(err)
		return
	}
	if !found {
		fmt.Printf("❌ Error: only %d commits exist, cannot go back %d.\n", count, n)
		return
	}

	wt, err := repo.Worktree()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := wt.Checkout(&git.CheckoutOptions{
		Hash:  targetHash,
		Force: true,
	}); err != nil {
		fmt.Println(err)
		return
	}

	// 2. שחזור HEAD.tar אחרי ה-checkout
	if hadBackup {
		if err := copyFile(backupPath, headTarPath); err != nil {
			fmt.Println("Warning: failed to restore HEAD.tar:", err)
		}
		os.Remove(backupPath)
	}

	fmt.Printf("✅ checked out commit %s\n", targetHash.String())
}

// getPatchBlob מחזיר את תוכן קובץ ה-patch השמור בתוך commit מסוים
func getPatchBlob(c *object.Commit) ([]byte, error) {
	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}
	var data []byte
	err = tree.Files().ForEach(func(f *object.File) error {
		if strings.HasPrefix(f.Name, "patches/") {
			r, err := f.Reader()
			if err != nil {
				return err
			}
			defer r.Close()
			b, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			data = b
			return storer.ErrStop
		}
		return nil
	})
	if err != nil && err != storer.ErrStop {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("no patch found in commit %s", c.Hash)
	}
	return data, nil
}

// HandleRestore משחזר tar היסטורי לפי כמה commits אחורה מ-HEAD
func HandleRestore(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Error: Please specify how many commits to go back.")
		return
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 0 {
		fmt.Println("❌ Error: argument must be a non-negative number.")
		return
	}

	repo, err := git.PlainOpen(".qrk")
	if err != nil {
		fmt.Println(err)
		return
	}
	head, err := repo.Head()
	if err != nil {
		fmt.Println(err)
		return
	}
	commitIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		fmt.Println(err)
		return
	}

	var commits []*object.Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if n >= len(commits) {
		fmt.Printf("❌ Error: only %d commits exist, cannot go back %d.\n", len(commits), n)
		return
	}

	// commits[0] = HEAD, commits[n] = המטרה. צריך מהישן ביותר עד המטרה.
	chain := commits[n:]
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	var currentTar []byte
	for i, c := range chain {
		patchData, err := getPatchBlob(c)
		if err != nil {
			fmt.Println(err)
			return
		}
		if i == 0 {
			currentTar = patchData // הפאץ' הראשון הוא ה-tar המלא עצמו
			continue
		}
		old := bytes.NewReader(currentTar)
		patchR := bytes.NewReader(patchData)
		var out bytes.Buffer
		if err := binarydist.Patch(old, &out, patchR); err != nil {
			fmt.Println(err)
			return
		}
		currentTar = out.Bytes()
	}

	if err := untarToWorktree(currentTar); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("✅ files restored to worktree")
}
