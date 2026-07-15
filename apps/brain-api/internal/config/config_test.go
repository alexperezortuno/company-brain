package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "brain.yaml")
	if err := os.WriteFile(path, []byte("instance:\n  name: Test Brain\n  language: es\nmodules:\n  filesystem:\n    enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BRAIN_CONFIG_PATH", path)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Brain.Instance.Name != "Test Brain" {
		t.Fatalf("unexpected instance name: %s", cfg.Brain.Instance.Name)
	}
	if !cfg.Brain.Modules["filesystem"].Enabled {
		t.Fatal("filesystem module should be enabled")
	}
}
