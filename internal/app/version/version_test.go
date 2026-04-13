package version

import (
	"path/filepath"
	"testing"
)

func TestVersionConfigPathsDefaultOrder(t *testing.T) {
	t.Setenv("P2V_VERSION_FILE", "")
	got := versionConfigPaths()
	if len(got) < 6 {
		t.Fatalf("expected several version search paths, got %d: %v", len(got), got)
	}
	wantContains := []string{
		filepath.Clean("version.json"),
		filepath.Clean(filepath.Join("bin", "gui", "version.json")),
		filepath.Clean(filepath.Join("bin", "cli", "version.json")),
		filepath.Clean(filepath.Join("..", "version.json")),
		filepath.Clean(filepath.Join("..", "..", "version.json")),
	}
	for _, want := range wantContains {
		if !contains(got, want) {
			t.Fatalf("expected search paths to contain %q, got: %v", want, got)
		}
	}
}

func TestVersionConfigPathsEnvPrecedence(t *testing.T) {
	t.Setenv("P2V_VERSION_FILE", "/tmp/custom-version.json")
	got := versionConfigPaths()
	if len(got) == 0 {
		t.Fatalf("expected non-empty version search paths")
	}
	if got[0] != "/tmp/custom-version.json" {
		t.Fatalf("expected env-configured path first, got: %s", got[0])
	}
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if filepath.Clean(item) == filepath.Clean(want) {
			return true
		}
	}
	return false
}
