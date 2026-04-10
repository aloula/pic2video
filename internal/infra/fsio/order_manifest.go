package fsio

import (
"bufio"
"os"
"path/filepath"
"strings"
)

func ReadExplicitOrder(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, filepath.Clean(line))
	}
	return out, s.Err()
}
