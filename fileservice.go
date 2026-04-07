package whitearchive

import (
	"crypto/sha256"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type FileService struct {
	*Scanner
}

func NewFileService(dir string) *FileService {
	return &FileService{
		Scanner: NewScanner(dir),
	}
}

func (s FileService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.directory, path))
}

func (s FileService) SaveFile(path string, data []byte) error {
	fullPath := filepath.Join(s.directory, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

type Scanner struct {
	directory string
}

func NewScanner(dir string) *Scanner {
	return &Scanner{
		directory: dir,
	}
}

func (s *Scanner) Snapshot() (Snapshot, error) {
	snapshot := Snapshot{}

	err := filepath.WalkDir(s.directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(s.directory, path)
		if err != nil {
			return err
		}

		snapshot[relPath] = Data{
			Hash:   HashOf(data),
			Update: time.Now(),
		}

		return nil
	})

	return snapshot, err
}

func HashOf(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
