// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Home string `toml:"home"`
}

func loadConfig() (config, error) {
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
