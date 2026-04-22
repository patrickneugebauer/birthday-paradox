package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type BufferedFile struct {
	f   *os.File
	w   *bufio.Writer
	enc *json.Encoder
}

func NewBufferedFile(path string) (*BufferedFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", path, err)
	}
	w := bufio.NewWriter(f)
	return &BufferedFile{f: f, w: w, enc: json.NewEncoder(w)}, nil
}

func (bf *BufferedFile) WriteString(s string) error {
	_, err := bf.w.WriteString(s)
	return err
}

func (bf *BufferedFile) Encode(v any) error {
	return bf.enc.Encode(v)
}

func (bf *BufferedFile) Flush() error {
	return bf.w.Flush()
}

// Close flushes the writer then closes the file. Call explicitly before
// os.Rename to surface I/O errors; defers act as safety nets on error paths.
func (bf *BufferedFile) Close() error {
	if err := bf.w.Flush(); err != nil {
		return err
	}
	return bf.f.Close()
}

// OpenBufferedFiles opens multiple files and returns them as a slice.
// On error, cleans up any files already opened.
func OpenBufferedFiles(paths ...string) ([]*BufferedFile, error) {
	var writers []*BufferedFile
	for _, path := range paths {
		w, err := NewBufferedFile(path)
		if err != nil {
			for _, opened := range writers {
				opened.Close()
			}
			return nil, err
		}
		writers = append(writers, w)
	}
	return writers, nil
}

// CloseBufferedFiles closes all writers and returns combined errors if any.
func CloseBufferedFiles(writers ...*BufferedFile) error {
	var errs []error
	for _, w := range writers {
		if err := w.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("close: %v", errs)
	}
	return nil
}
