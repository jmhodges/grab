# Port from golang.org/x/tools/go/vcs to Masterminds/vcs

## Phase 1: Direct repo URL support with Masterminds/vcs

### 1. Update dependencies
- `go get github.com/Masterminds/vcs`
- `go mod tidy` to drop `golang.org/x/tools/go/vcs`

### 2. Add repo root derivation helper
Write a function that extracts the repo root from a URL. For well-known hosts (github.com, gitlab.com, bitbucket.org), the repo root is always `host/owner/repo` (first 3 path segments).

### 3. Rewrite main.go

Replace the old flow:
1. Strip scheme to make an import path
2. `vcs.RepoRootForImportPath(importPath, false)` → `RepoRoot`
3. If git, rewrite repo URL to SSH via `rewriteToSSH(repoRoot.Repo, cfg.SSHPreferredHosts)`
4. `repoRoot.VCS.Create(localPath, repoURL)` → clone

With the new flow:
1. Parse the URL to get the remote URL and derive repo root for the local path
2. Apply SSH rewriting: `rewriteToSSH(remoteURL, cfg.SSHPreferredHosts)` — this already only rewrites `https`/`http` URLs, so it's safe to call unconditionally without a VCS type check
3. `vcs.NewRepo(remoteURL, localPath)` → auto-detects VCS type
4. `repo.Get()` → clone

### 4. API mapping reference

| Old (`x/tools/go/vcs`)                         | New (`Masterminds/vcs`)              |
|-------------------------------------------------|--------------------------------------|
| `vcs.RepoRootForImportPath(importPath, false)`  | `vcs.NewRepo(remoteURL, localPath)`  |
| `repoRoot.Repo`                                 | Derive from input URL                |
| `repoRoot.Root`                                 | Derive from input URL (host/owner/repo) |
| `repoRoot.VCS.Cmd == "git"` (for SSH check)      | Not needed — `rewriteToSSH` is safe to call unconditionally |
| `repoRoot.VCS.Create(root, repo)`               | `repo.Get()`                         |

## Phase 2: `<meta name="go-import">` resolution for vanity import paths

Add support for vanity import paths (e.g. `gopkg.in/yaml.v2`, `go.uber.org/zap`) by:

1. Fetching the import path URL with `?go-get=1` query parameter over HTTP
2. Parsing `<meta name="go-import" content="root vcs repo-url">` from the HTML response
3. Using the discovered VCS type and repo URL with Masterminds/vcs to clone
