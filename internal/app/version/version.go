package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	Version   = "dev"
	BuildDate = "unknown"

	loadOnce sync.Once
	loaded   metadata
)

type metadata struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
}

func loadMetadata() {
	loaded = metadata{Version: Version, BuildDate: BuildDate}
	for _, path := range versionConfigPaths() {
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var m metadata
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		if strings.TrimSpace(m.Version) != "" {
			loaded.Version = strings.TrimSpace(m.Version)
		}
		if strings.TrimSpace(m.BuildDate) != "" {
			loaded.BuildDate = strings.TrimSpace(m.BuildDate)
		}
		return
	}
}

func getMetadata() metadata {
	loadOnce.Do(loadMetadata)
	return loaded
}

func versionConfigPaths() []string {
	paths := make([]string, 0, 16)
	seen := make(map[string]struct{})
	appendUnique := func(p string) {
		if strings.TrimSpace(p) == "" {
			return
		}
		clean := filepath.Clean(p)
		if _, ok := seen[clean]; ok {
			return
		}
		seen[clean] = struct{}{}
		paths = append(paths, clean)
	}
	if env := strings.TrimSpace(os.Getenv("P2V_VERSION_FILE")); env != "" {
		appendUnique(env)
	}
	if exe, err := os.Executable(); err == nil {
		appendUnique(filepath.Join(filepath.Dir(exe), "version.json"))
	}
	// For local go run and IDE launches, cwd may be nested (for example
	// cmd/pic2video-gui). Search current and parent directories.
	for _, prefix := range []string{".", "..", filepath.Join("..", ".."), filepath.Join("..", "..", "..")} {
		appendUnique(filepath.Join(prefix, "version.json"))
		appendUnique(filepath.Join(prefix, "bin", "gui", "version.json"))
		appendUnique(filepath.Join(prefix, "bin", "cli", "version.json"))
	}
	return paths
}

func Short() string {
	m := getMetadata()
	if m.Version == "" {
		return "dev"
	}
	return m.Version
}

func Info() string {
	m := getMetadata()
	return fmt.Sprintf("version=%s built=%s", Short(), m.BuildDate)
}
