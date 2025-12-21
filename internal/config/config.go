package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	VaultPath string
	Editor    string
}

func Default() Config {
	return Config{
		VaultPath: "notes",
		Editor:    "",
	}
}

func Load(vaultPath, editor string) (Config, error) {
	cfg := Default()

	if vaultPath != "" {
		cfg.VaultPath = vaultPath
	}
	if editor != "" {
		cfg.Editor = editor
	}

	absVault, err := validateVaultPath(cfg.VaultPath)
	if err != nil {
		return Config{}, err
	}
	cfg.VaultPath = absVault

	return cfg, nil
}

func validateVaultPath(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("resolve vault path: %w", err)
	}

	info, err := os.Stat(abs)
	if err == nil {
		if !info.IsDir() {
			return "", fmt.Errorf("vault path exists but is not a directory: %s", abs)
		}
		return abs, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("stat vault path: %w", err)
	}

	parent := filepath.Dir(abs)
	pinfo, perr := os.Stat(parent)
	if perr != nil {
		return "", fmt.Errorf("vault parent missing or not accessible: %w", perr)
	}
	if !pinfo.IsDir() {
		return "", fmt.Errorf("vault parent is not a directory: %s", parent)
	}

	return abs, nil
}
