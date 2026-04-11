package pipeline

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/infra/fsio"
)

// ApplyOrder applies ordering to assets (legacy, for backward compatibility)
func ApplyOrder(mode string, assets []media.Asset, explicit []string) []media.Asset {
	return ApplyOrderExt(mode, assets, explicit, "")
}

// ApplyOrderExt applies ordering to assets with optional ffprobe support for EXIF
func ApplyOrderExt(mode string, assets []media.Asset, explicit []string, ffprobeBin string) []media.Asset {
	ordered := make([]media.Asset, len(assets))
	copy(ordered, assets)
	switch mode {
	case "exif":
		// Sort by EXIF creation date using ffprobe
		exifData := make(map[string]time.Time)
		for _, asset := range ordered {
			exif, err := fsio.ExtractExif(asset.Path, ffprobeBin)
			if err == nil && !exif.CreateDate.IsZero() {
				exifData[asset.Path] = exif.CreateDate
			} else {
				// Fallback to file modification time
				if fileInfo, err := osStat(asset.Path); err == nil {
					exifData[asset.Path] = fileInfo.ModTime()
				}
			}
		}
		sort.SliceStable(ordered, func(i, j int) bool {
			ti := exifData[ordered[i].Path]
			tj := exifData[ordered[j].Path]
			if ti.Equal(tj) {
				// Fallback to name comparison if same time
				return strings.ToLower(ordered[i].Path) < strings.ToLower(ordered[j].Path)
			}
			return ti.Before(tj)
		})
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
