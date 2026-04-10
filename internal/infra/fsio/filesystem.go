package fsio

import (
"fmt"
"os"
"path/filepath"
"sort"
"strings"

"github.com/loula/pic2video/internal/domain/media"
)

var supportedExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

func ListImageAssets(inputDir string) ([]media.Asset, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	assets := []media.Asset{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if !supportedExt[ext] {
			continue
		}
		assets = append(assets, media.Asset{Path: filepath.Join(inputDir, e.Name()), Format: strings.TrimPrefix(ext, "."), IsValid: true})
	}
	sort.SliceStable(assets, func(i, j int) bool { return strings.ToLower(assets[i].Path) < strings.ToLower(assets[j].Path) })
	if len(assets) == 0 {
		return nil, fmt.Errorf("no supported images found")
	}
	for i := range assets {
		assets[i].OrderIndex = i
	}
	return assets, nil
}
