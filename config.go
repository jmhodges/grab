// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Home              string   `toml:"home"`
	SSHPreferredHosts []string `toml:"ssh_preferred_hosts"`
}

func loadConfig() (config, error) {
	// On macOS, check $HOME/.config/grab/config.toml first before falling
	// back to the native Library/Application Support directory.
	if runtime.GOOS == "darwin" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return config{}, err
		}
		// Duplicating the appending of grab to the config dir because this is
		// the only platform we have a fallback on.
		xdgConfig := filepath.Join(homeDir, ".config", "grab", "config.toml")
		if _, err := os.Stat(xdgConfig); err == nil {
			return loadConfigFrom(filepath.Join(homeDir, ".config"))
		}
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return config{}, err
	}
	return loadConfigFrom(configDir)
}

func loadConfigFrom(configDir string) (config, error) {
	var cfg config

	configPath := filepath.Join(configDir, "grab", "config.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return cfg, err
		}
		// Missing config file is fine, proceed with defaults.
	} else {
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return cfg, err
		}
	}

	// Env vars override config file values.
	if envHome := strings.TrimSpace(os.Getenv("GRAB_HOME")); envHome != "" {
		cfg.Home = envHome
	}

	// Default home to ~/src if still unset.
	if cfg.Home == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return cfg, err
		}
		cfg.Home = filepath.Join(homeDir, "src")
	}

	return cfg, nil
}
