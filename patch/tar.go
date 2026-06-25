package patch

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
)

// CreateTarFromDirectory takes a source directory path and writes its contents
// into a single uncompressed .tar file using streaming.
func CreateTarFromFiles(filePaths []string, outputTarPath string) error {
	out, err := os.Create(outputTarPath)
	if err != nil {
		return err
	}
	defer out.Close()

	tw := tar.NewWriter(out)
	defer tw.Close()

	for _, p := range filePaths {
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("missing file %s: %w", p, err)
		}
		info, err := f.Stat()
		if err != nil {
			f.Close()
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			f.Close()
			return err
		}
		hdr.Name = p
		if err := tw.WriteHeader(hdr); err != nil {
			f.Close()
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}
