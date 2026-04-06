// Copyright 2021 Jeffrey M Hodges.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	mvcs "github.com/Masterminds/vcs"
)

// fakeRepo implements mvcs.Repo for testing. Only Get, Remote, and LocalPath
// are meaningful; every other method is a no-op stub.
type fakeRepo struct {
	remote    string
	localPath string
	getErr    error
}

func (f *fakeRepo) Vcs() mvcs.Type                            { return mvcs.Git }
func (f *fakeRepo) Remote() string                             { return f.remote }
func (f *fakeRepo) LocalPath() string                          { return f.localPath }
func (f *fakeRepo) Get() error                                 { return f.getErr }
func (f *fakeRepo) Init() error                                { return nil }
func (f *fakeRepo) Update() error                              { return nil }
func (f *fakeRepo) UpdateVersion(string) error                 { return nil }
func (f *fakeRepo) Version() (string, error)                   { return "", nil }
func (f *fakeRepo) Current() (string, error)                   { return "", nil }
func (f *fakeRepo) Date() (time.Time, error)                   { return time.Time{}, nil }
func (f *fakeRepo) CheckLocal() bool                           { return true }
func (f *fakeRepo) Branches() ([]string, error)                { return nil, nil }
func (f *fakeRepo) Tags() ([]string, error)                    { return nil, nil }
func (f *fakeRepo) IsReference(string) bool                    { return false }
func (f *fakeRepo) IsDirty() bool                              { return false }
func (f *fakeRepo) CommitInfo(string) (*mvcs.CommitInfo, error) { return nil, nil }
func (f *fakeRepo) TagsFromCommit(string) ([]string, error)    { return nil, nil }
func (f *fakeRepo) Ping() bool                                 { return true }
func (f *fakeRepo) RunFromDir(string, ...string) ([]byte, error) { return nil, nil }
func (f *fakeRepo) CmdFromDir(cmd string, args ...string) *exec.Cmd {
	return exec.Command(cmd, args...)
}
func (f *fakeRepo) ExportDir(string) error { return nil }

func TestGrab(t *testing.T) {
	tests := []struct {
		name          string
		rawURL        string
		cfg           config
		resolveImport func(string) (goImport, error)
		getErr        error
		wantRemote    string
		wantLocal     string
		wantErr       bool
		wantErrMsg    string
	}{
		{
			name:   "standard https github URL",
			rawURL: "https://github.com/jmhodges/grab",
			cfg:    config{Home: "/tmp/src"},
			resolveImport: func(string) (goImport, error) {
				return goImport{}, errors.New("should not be called")
			},
			wantRemote: "https://github.com/jmhodges/grab",
			wantLocal:  "/tmp/src/github.com/jmhodges/grab",
		},
		{
			name:   "URL without scheme gets https prepended",
			rawURL: "github.com/jmhodges/grab",
			cfg:    config{Home: "/tmp/src"},
			resolveImport: func(string) (goImport, error) {
				return goImport{}, errors.New("should not be called")
			},
			wantRemote: "https://github.com/jmhodges/grab",
			wantLocal:  "/tmp/src/github.com/jmhodges/grab",
		},
		{
			name:   "vanity import falls back to resolveImport",
			rawURL: "https://go.uber.org/zap",
			cfg:    config{Home: "/tmp/src"},
			resolveImport: func(importPath string) (goImport, error) {
				return goImport{
					Root:    "go.uber.org/zap",
					VCS:     "git",
					RepoURL: "https://github.com/uber-go/zap",
				}, nil
			},
			wantRemote: "https://github.com/uber-go/zap",
			wantLocal:  "/tmp/src/go.uber.org/zap",
		},
		{
			name:   "SSH rewriting applied when host in SSHPreferredHosts",
			rawURL: "https://github.com/jmhodges/grab",
			cfg:    config{Home: "/tmp/src", SSHPreferredHosts: []string{"github.com"}},
			resolveImport: func(string) (goImport, error) {
				return goImport{}, errors.New("should not be called")
			},
			wantRemote: "git@github.com:jmhodges/grab.git",
			wantLocal:  "/tmp/src/github.com/jmhodges/grab",
		},
		{
			name:   "SSH rewriting applied to vanity import resolved URL",
			rawURL: "https://go.uber.org/zap",
			cfg:    config{Home: "/tmp/src", SSHPreferredHosts: []string{"github.com"}},
			resolveImport: func(importPath string) (goImport, error) {
				return goImport{
					Root:    "go.uber.org/zap",
					VCS:     "git",
					RepoURL: "https://github.com/uber-go/zap",
				}, nil
			},
			wantRemote: "git@github.com:uber-go/zap.git",
			wantLocal:  "/tmp/src/go.uber.org/zap",
		},
		{
			name:   "repo.Get error propagated",
			rawURL: "https://github.com/jmhodges/grab",
			cfg:    config{Home: "/tmp/src"},
			resolveImport: func(string) (goImport, error) {
				return goImport{}, errors.New("should not be called")
			},
			getErr:     errors.New("clone failed"),
			wantErr:    true,
			wantErrMsg: "unable to download",
		},
		{
			name:   "resolveImport error propagated when repoRoot fails",
			rawURL: "https://go.uber.org/zap",
			cfg:    config{Home: "/tmp/src"},
			resolveImport: func(string) (goImport, error) {
				return goImport{}, errors.New("network timeout")
			},
			wantErr:    true,
			wantErrMsg: "unable to determine repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRemote, gotLocal string

			deps := grabDeps{
				newRepo: func(remote, local string) (mvcs.Repo, error) {
					gotRemote = remote
					gotLocal = local
					return &fakeRepo{
						remote:    remote,
						localPath: local,
						getErr:    tt.getErr,
					}, nil
				},
				resolveImport: tt.resolveImport,
			}

			err := grab(tt.rawURL, tt.cfg, deps)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("error = %q, want it to contain %q", err, tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if gotRemote != tt.wantRemote {
				t.Errorf("newRepo remote = %q, want %q", gotRemote, tt.wantRemote)
			}
			if gotLocal != tt.wantLocal {
				t.Errorf("newRepo local = %q, want %q", gotLocal, tt.wantLocal)
			}
		})
	}
}

func TestGrab_newRepoError(t *testing.T) {
	deps := grabDeps{
		newRepo: func(remote, local string) (mvcs.Repo, error) {
			return nil, fmt.Errorf("unsupported VCS")
		},
		resolveImport: func(string) (goImport, error) {
			return goImport{}, errors.New("should not be called")
		},
	}

	err := grab("https://github.com/jmhodges/grab", config{Home: "/tmp/src"}, deps)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unable to determine VCS type") {
		t.Errorf("error = %q, want it to contain 'unable to determine VCS type'", err)
	}
}
