// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("missing config file returns defaults", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("GRAB_HOME", "")

		cfg, err := loadConfigFrom(dir)
		if err != nil {
			t.Fatal(err)
		}
		if len(cfg.SSHPreferredHosts) != 0 {
			t.Errorf("expected no ssh hosts, got %v", cfg.SSHPreferredHosts)
		}
		// Home should default to ~/src.
		homeDir, _ := os.UserHomeDir()
		want := filepath.Join(homeDir, "src")
		if cfg.Home != want {
			t.Errorf("Home = %q, want %q", cfg.Home, want)
		}
	})

	t.Run("config file parsed", func(t *testing.T) {
		dir := t.TempDir()
		grabDir := filepath.Join(dir, "grab")
		os.MkdirAll(grabDir, 0755)
		os.WriteFile(filepath.Join(grabDir, "config.toml"), []byte("home = \"/tmp/src\"\nssh_preferred_hosts = [\"github.com\"]\n"), 0644)
		t.Setenv("GRAB_HOME", "")

		cfg, err := loadConfigFrom(dir)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Home != "/tmp/src" {
			t.Errorf("Home = %q, want /tmp/src", cfg.Home)
		}
		if len(cfg.SSHPreferredHosts) != 1 || cfg.SSHPreferredHosts[0] != "github.com" {
			t.Errorf("SSHHosts = %v, want [github.com]", cfg.SSHPreferredHosts)
		}
	})

	t.Run("env var overrides home", func(t *testing.T) {
		dir := t.TempDir()
		grabDir := filepath.Join(dir, "grab")
		os.MkdirAll(grabDir, 0755)
		os.WriteFile(filepath.Join(grabDir, "config.toml"), []byte("home = \"/tmp/src\"\nssh_preferred_hosts = [\"github.com\"]\n"), 0644)
		t.Setenv("GRAB_HOME", "/override/src")

		cfg, err := loadConfigFrom(dir)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Home != "/override/src" {
			t.Errorf("Home = %q, want /override/src", cfg.Home)
		}
		if len(cfg.SSHPreferredHosts) != 1 || cfg.SSHPreferredHosts[0] != "github.com" {
			t.Errorf("SSHHosts = %v, want [github.com]", cfg.SSHPreferredHosts)
		}
	})

	t.Run("malformed toml returns error", func(t *testing.T) {
		dir := t.TempDir()
		grabDir := filepath.Join(dir, "grab")
		os.MkdirAll(grabDir, 0755)
		os.WriteFile(filepath.Join(grabDir, "config.toml"), []byte(`[bad toml`), 0644)
		t.Setenv("GRAB_HOME", "")

		_, err := loadConfigFrom(dir)
		if err == nil {
			t.Fatal("expected error for malformed TOML, got nil")
		}
	})
}
