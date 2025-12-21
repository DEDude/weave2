package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func resetFlags() {
	vaultFlag = ""
	editorFlag = ""
	rootCmd.SetArgs(nil)
}

func TestRootLoadsConfigFromEnv(t *testing.T) {
	t.Setenv("WEAVE_VAULT", filepath.Join(t.TempDir(), "vault"))
	t.Setenv("WEAVE_EDITOR", "nano")
	resetFlags()

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	wantVault := os.Getenv("WEAVE_VAULT")
	if cfg.VaultPath != wantVault {
		t.Fatalf("VaultPath = %q, want %q", cfg.VaultPath, wantVault)
	}
	if cfg.Editor != "nano" {
		t.Fatalf("Editor = %q, want %q", cfg.Editor, "nano")
	}
}

func TestFlagsOverrideEnv(t *testing.T) {
	t.Setenv("WEAVE_VAULT", filepath.Join(t.TempDir(), "envvault"))
	t.Setenv("WEAVE_EDITOR", "nano")
	resetFlags()
	rootCmd.SetArgs([]string{"--vault", "flagvault", "--editor", "vim"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	wantVault, _ := filepath.Abs("flagvault")
	if cfg.VaultPath != wantVault {
		t.Fatalf("VaultPath = %q, want %q", cfg.VaultPath, wantVault)
	}
	if cfg.Editor != "vim" {
		t.Fatalf("Editor = %q, want %q", cfg.Editor, "vim")
	}
}
