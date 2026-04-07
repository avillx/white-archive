package whitearchive

import (
	"context"
	"errors"
	"maps"
	"slices"
)

type FileProvider interface {
	ReadFile(path string) ([]byte, error)
	SaveFile(path string, data []byte) error
	Snapshot() (Snapshot, error)
}

type Storage interface {
	Download(ctx context.Context, name string) ([]byte, error)
	Upload(ctx context.Context, name string, data []byte) error
}

type Syncer struct {
	fileService FileProvider
	storage     Storage
	cipher      *Cipher
}

func NewSyncer(fileService FileProvider, storage Storage, cipher *Cipher) *Syncer {
	return &Syncer{fileService: fileService, storage: storage, cipher: cipher}
}

// backup
func (s *Syncer) Backup(ctx context.Context) error {
	// download & decode remote snapshot
	remoteSnapshot, err := s.downloadSnapshot(ctx)

	if err != nil {
		return err
	}

	// scan local files
	localSnapshot, err := s.fileService.Snapshot()
	if err != nil {
		return err
	}

	// upload changed files
	changes := diffs(remoteSnapshot, localSnapshot)
	filePaths := slices.Collect(maps.Keys(changes))

	if err := s.uploadFiles(ctx, filePaths); err != nil {
		return err
	}

	// upload updated snapshot
	return s.uploadSnapshot(ctx, localSnapshot)
}

// restore
func (s *Syncer) Restore(ctx context.Context) error {
	remoteSnapshot, err := s.downloadSnapshot(ctx)
	if err != nil {
		return err
	}

	localSnapshot, err := s.fileService.Snapshot()
	if err != nil {
		return err
	}

	for path := range diffs(localSnapshot, remoteSnapshot) {
		if err := s.downloadFile(ctx, path); err != nil {
			return err
		}
	}

	return nil
}

// helpers

func (s *Syncer) downloadFile(ctx context.Context, relPath string) error {
	data, err := s.storage.Download(ctx, relPath)
	if err != nil {
		return err
	}
	decrypted, err := s.cipher.Decrypt(data)
	if err != nil {
		return err
	}
	return s.fileService.SaveFile(relPath, decrypted)
}

func (s *Syncer) uploadFiles(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if err := s.uploadFile(ctx, path); err != nil {
			return err
		}
	}
	return nil
}

func (s *Syncer) uploadFile(ctx context.Context, relPath string) error {
	data, err := s.fileService.ReadFile(relPath)
	if err != nil {
		return err
	}

	encrypted, err := s.cipher.Encrypt(data)
	if err != nil {
		return err
	}

	return s.storage.Upload(ctx, relPath, encrypted)
}

func (s *Syncer) downloadSnapshot(ctx context.Context) (Snapshot, error) {
	data, err := s.storage.Download(ctx, "index.jsonl")

	// on first download it is empty
	if errors.Is(err, ErrNotFound) {
		return Snapshot{}, nil
	}
	if err != nil {
		return nil, err
	}

	decrypted, err := s.cipher.Decrypt(data)
	if err != nil {
		return nil, err
	}
	return UnmarshalSnapshot(decrypted)
}

func (s *Syncer) uploadSnapshot(ctx context.Context, snapshot Snapshot) error {
	data, err := MarshalSnapshot(snapshot)
	if err != nil {
		return err
	}

	encrypted, err := s.cipher.Encrypt(data)
	if err != nil {
		return err
	}
	return s.storage.Upload(ctx, "index.jsonl", encrypted)
}
