package git_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/worktree"
	testtarget "github.com/kyoh86/gogh/v4/infra/git"
)

func TestWorktreeServiceAddCreateBranchUsesExistingLocalBranch(t *testing.T) {
	ctx := context.Background()
	repo := setupWorktreeTestRepository(t)
	runGit(t, repo.FullPath(), "branch", "existing")

	service := testtarget.NewWorktreeService(worktree.NewPathBuilder())
	wt, err := service.Add(ctx, *repo, "existing", worktree.AddOptions{CreateBranch: true})
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if got := currentBranch(t, wt.Path); got != "existing" {
		t.Fatalf("current branch = %q, want %q", got, "existing")
	}
}

func TestWorktreeServiceAddCreateBranchUsesRemoteBranch(t *testing.T) {
	ctx := context.Background()
	repo := setupWorktreeTestRepositoryWithRemoteBranch(t, "remote-feature")
	remoteCommit := gitOutput(t, repo.FullPath(), "rev-parse", "origin/remote-feature")

	service := testtarget.NewWorktreeService(worktree.NewPathBuilder())
	wt, err := service.Add(ctx, *repo, "remote-feature", worktree.AddOptions{CreateBranch: true})
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if got := currentBranch(t, wt.Path); got != "remote-feature" {
		t.Fatalf("current branch = %q, want %q", got, "remote-feature")
	}
	if got := gitOutput(t, wt.Path, "rev-parse", "HEAD"); got != remoteCommit {
		t.Fatalf("HEAD = %q, want origin/remote-feature %q", got, remoteCommit)
	}
}

func TestWorktreeServiceAddCreateBranchFallsBackToHeadWithoutOriginHead(t *testing.T) {
	ctx := context.Background()
	repo := setupWorktreeTestRepository(t)
	headCommit := gitOutput(t, repo.FullPath(), "rev-parse", "HEAD")

	service := testtarget.NewWorktreeService(worktree.NewPathBuilder())
	wt, err := service.Add(ctx, *repo, "new-feature", worktree.AddOptions{CreateBranch: true})
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if got := currentBranch(t, wt.Path); got != "new-feature" {
		t.Fatalf("current branch = %q, want %q", got, "new-feature")
	}
	if got := gitOutput(t, wt.Path, "rev-parse", "HEAD"); got != headCommit {
		t.Fatalf("HEAD = %q, want original HEAD %q", got, headCommit)
	}
}

func setupWorktreeTestRepository(t *testing.T) *repository.Location {
	t.Helper()

	repoPath := filepath.Join(t.TempDir(), "repo")
	runGit(t, "", "init", "--initial-branch=main", repoPath)
	runGit(t, repoPath, "config", "user.name", "Gogh Test")
	runGit(t, repoPath, "config", "user.email", "gogh-test@example.com")

	if err := os.WriteFile(filepath.Join(repoPath, "README.md"), []byte("test\n"), 0o644); err != nil {
		t.Fatalf("writing README.md: %v", err)
	}
	runGit(t, repoPath, "add", "README.md")
	runGit(t, repoPath, "commit", "-m", "initial commit")

	return repository.NewLocation(repoPath, "github.com", "kyoh86", "gogh-test")
}

func setupWorktreeTestRepositoryWithRemoteBranch(t *testing.T, branch string) *repository.Location {
	t.Helper()

	repo := setupWorktreeTestRepository(t)
	remotePath := filepath.Join(t.TempDir(), "remote.git")
	runGit(t, "", "init", "--bare", "--initial-branch=main", remotePath)
	runGit(t, repo.FullPath(), "remote", "add", "origin", remotePath)
	runGit(t, repo.FullPath(), "push", "-u", "origin", "main")
	runGit(t, repo.FullPath(), "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")

	runGit(t, repo.FullPath(), "checkout", "-b", branch)
	if err := os.WriteFile(filepath.Join(repo.FullPath(), "feature.txt"), []byte(branch+"\n"), 0o644); err != nil {
		t.Fatalf("writing feature.txt: %v", err)
	}
	runGit(t, repo.FullPath(), "add", "feature.txt")
	runGit(t, repo.FullPath(), "commit", "-m", "add remote branch")
	runGit(t, repo.FullPath(), "push", "origin", branch)
	runGit(t, repo.FullPath(), "checkout", "main")
	runGit(t, repo.FullPath(), "branch", "-D", branch)

	return repo
}

func currentBranch(t *testing.T, repoPath string) string {
	t.Helper()
	return gitOutput(t, repoPath, "branch", "--show-current")
}

func gitOutput(t *testing.T, repoPath string, args ...string) string {
	t.Helper()
	output := runGit(t, repoPath, args...)
	return strings.TrimSpace(string(output))
}

func runGit(t *testing.T, repoPath string, args ...string) []byte {
	t.Helper()

	cmdArgs := args
	if repoPath != "" {
		cmdArgs = append([]string{"-C", repoPath}, args...)
	}
	cmd := exec.Command("git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(cmdArgs, " "), err, string(output))
	}
	return output
}
