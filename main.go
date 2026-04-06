// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"
	"os"

	mvcs "github.com/Masterminds/vcs"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) <= 1 || len(os.Args) > 2 {
		log.Fatal("usage: grab REPO_URL")
	}
	rawURL := os.Args[1]

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err)
	}

	deps := grabDeps{
		newRepo:       mvcs.NewRepo,
		resolveImport: resolveGoImport,
	}
	if err := grab(rawURL, cfg, deps); err != nil {
		log.Fatal(err)
	}
}
