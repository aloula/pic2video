package pipeline

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/loula/pic2video/internal/domain/media"
)

func ApplyOrder(mode string, assets []media.Asset, explicit []string) []media.Asset {
	ordered := make([]media.Asset, len(assets))
	copy(ordered, assets)
	switch mode {
	case "time":
		sort.SliceStable(ordered, func(i, j int) bool {
			a, _ := osModTime(ordered[i].Path)
			b, _ := osModTime(ordered[j].Path)
			return a.Before(b)
		})
	case "explicit":
		idx := map[string]int{}
		for i, p := range explicit {
			idx[filepath.Clean(p)] = i
			idx[filepath.Base(filepath.Clean(p))] = i
		}
		sort.SliceStable(ordered, func(i, j int) bool {
			ai, aok := idx[filepath.Base(ordered[i].Path)]
			bi, bok := idx[filepath.Base(ordered[j].Path)]
			if aok && bok {
				return ai < bi
			}
			if aok != bok {
				return aok
			}
			return strings.ToLower(ordered[i].Path) < strings.ToLower(ordered[j].Path)
		})
	default:
		sort.SliceStable(ordered, func(i, j int) bool {
			return strings.ToLower(ordered[i].Path) < strings.ToLower(ordered[j].Path)
		})
	}
	for i := range ordered {
		ordered[i].OrderIndex = i
	}
	return ordered
}

func osModTime(path string) (time.Time, bool) {
	fi, err := osStat(path)
	if err != nil {
		return time.Time{}, false
	}
	return fi.ModTime(), true
}

var osStat = func(name string) (interface{ ModTime() time.Time }, error) {
	return os.Stat(name)
}
