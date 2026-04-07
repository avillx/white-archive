package whitearchive_test

import (
	"context"
	"testing"
	whitearchive "white-archive"
)

// mocks

type mockStorage struct {
	files map[string][]byte
}

func newMockStorage() *mockStorage {
	return &mockStorage{files: map[string][]byte{}}
}

func (m *mockStorage) Upload(_ context.Context, name string, data []byte) error {
	m.files[name] = data
	return nil
}

func (m *mockStorage) Download(_ context.Context, name string) ([]byte, error) {
	data, ok := m.files[name]
	if !ok {
		return nil, whitearchive.ErrNotFound
	}
	return data, nil
}

func (m *mockStorage) Delete(_ context.Context, name string) error {
	delete(m.files, name)
	return nil
}

type mockFileService struct {
	files map[string][]byte
}

func newMockFileService(files map[string][]byte) *mockFileService {
	return &mockFileService{files: files}
}

func (m *mockFileService) ReadFile(path string) ([]byte, error) {
	return m.files[path], nil
}

func (m *mockFileService) SaveFile(path string, data []byte) error {
	m.files[path] = data
	return nil
}

func (m *mockFileService) Snapshot() (whitearchive.Snapshot, error) {
	snapshot := whitearchive.Snapshot{}
	for path, data := range m.files {
		snapshot[path] = whitearchive.Data{Hash: whitearchive.HashOf(data)}
	}
	return snapshot, nil
}

// helpers

func newSyncer(fs *mockFileService, storage *mockStorage) *whitearchive.Syncer {
	cipher := whitearchive.NewCipher([]byte("test-key"))
	return whitearchive.NewSyncer(fs, storage, cipher)
}

// backup tests

func TestBackup_FirstRun(t *testing.T) {
	fs := newMockFileService(map[string][]byte{
		"daily_logs/04.04.26.md": []byte("log content"),
	})
	storage := newMockStorage()
	syncer := newSyncer(fs, storage)

	if err := syncer.Backup(context.Background()); err != nil {
		t.Fatal(err)
	}

	if _, ok := storage.files["daily_logs/04.04.26.md"]; !ok {
		t.Error("file was not uploaded")
	}
	if _, ok := storage.files["index.jsonl"]; !ok {
		t.Error("index was not uploaded")
	}
}

func TestBackup_OnlyChangedFiles(t *testing.T) {
	fs := newMockFileService(map[string][]byte{
		"file_a.md": []byte("original"),
		"file_b.md": []byte("original"),
	})
	storage := newMockStorage()
	syncer := newSyncer(fs, storage)

	// first backup
	if err := syncer.Backup(context.Background()); err != nil {
		t.Fatal(err)
	}

	// change only file_a
	fs.files["file_a.md"] = []byte("changed")
	oldFileB := storage.files["file_b.md"]

	if err := syncer.Backup(context.Background()); err != nil {
		t.Fatal(err)
	}

	if string(storage.files["file_b.md"]) != string(oldFileB) {
		t.Error("file_b should not have been re-uploaded")
	}
}

// restore tests

func TestRestore_FirstRun(t *testing.T) {
	fs := newMockFileService(map[string][]byte{})
	storage := newMockStorage()

	// подготовим бекап
	backupFS := newMockFileService(map[string][]byte{
		"daily_logs/04.04.26.md": []byte("log content"),
	})
	if err := newSyncer(backupFS, storage).Backup(context.Background()); err != nil {
		t.Fatal(err)
	}

	if err := newSyncer(fs, storage).Restore(context.Background()); err != nil {
		t.Fatal(err)
	}

	if string(fs.files["daily_logs/04.04.26.md"]) != "log content" {
		t.Error("file was not restored")
	}
}

func TestRestore_OnlyMissingOrChanged(t *testing.T) {
	fs := newMockFileService(map[string][]byte{
		"file_a.md": []byte("content"),
		"file_b.md": []byte("outdated"),
	})
	storage := newMockStorage()

	backupFS := newMockFileService(map[string][]byte{
		"file_a.md": []byte("content"),
		"file_b.md": []byte("updated"),
	})
	if err := newSyncer(backupFS, storage).Backup(context.Background()); err != nil {
		t.Fatal(err)
	}

	if err := newSyncer(fs, storage).Restore(context.Background()); err != nil {
		t.Fatal(err)
	}

	if string(fs.files["file_b.md"]) != "updated" {
		t.Error("outdated file should have been restored")
	}
	if string(fs.files["file_a.md"]) != "content" {
		t.Error("unchanged file should not have been touched")
	}
}
