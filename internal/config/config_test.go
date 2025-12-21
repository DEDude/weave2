package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	cfg, err := Load("", "")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	wantVault := filepath.Join(tmp, "notes")
	if cfg.VaultPath != wantVault {
		t.Fatalf("VaultPath = %q, want %q", cfg.VaultPath, wantVault)
	}
	if cfg.Editor != "" {
		t.Fatalf("Editor = %q, want empty", cfg.Editor)
	}
}

func TestLoadWithOverrides(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	if err := os.MkdirAll(filepath.Join(tmp, "custom"), 0o755); err != nil {
		t.Fatalf("mkdir custom: %v", err)
	}

	cfg, err := Load(filepath.Join("custom", "vault"), "nvim")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	wantVault := filepath.Join(tmp, "custom", "vault")
	if cfg.VaultPath != wantVault {
		t.Fatalf("VaultPath = %q, want %q", cfg.VaultPath, wantVault)
	}
	if cfg.Editor != "nvim" {
		t.Fatalf("Editor = %q, want %q", cfg.Editor, "nvim")
	}
}

func TestLoadAbsolutePath(t *testing.T) {
	tmp := t.TempDir()

	cfg, err := Load(tmp, "")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.VaultPath != tmp {
		t.Fatalf("VaultPath = %q, want %q", cfg.VaultPath, tmp)
	}
}

func TestLoadErrorsWhenPathIsFile(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "vault")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if _, err := Load(filePath, ""); err == nil {
		t.Fatalf("Load() error = nil, want error for file path")
	}
}

func TestLoadErrorsWhenParentMissing(t *testing.T) {
	tmp := t.TempDir()
	missing := filepath.Join(tmp, "missing", "vault")

	if _, err := Load(missing, ""); err == nil {
		t.Fatalf("Load() error = nil, want error for missing parent")
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}
