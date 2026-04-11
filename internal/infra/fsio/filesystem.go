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

var supportedVideoExt = map[string]bool{
	".mp4": true, ".mov": true, ".mkv": true, ".webm": true,
}

var supportedAudioExt = map[string]bool{
	".mp3": true,
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
		assets = append(assets, media.Asset{Path: filepath.Join(inputDir, e.Name()), MediaType: media.MediaTypeImage, Format: strings.TrimPrefix(ext, "."), IsValid: true})
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

func ListMixedAssets(inputDir string) ([]media.Asset, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	assets := []media.Asset{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		switch {
		case supportedExt[ext]:
			assets = append(assets, media.Asset{
				Path:      filepath.Join(inputDir, name),
				MediaType: media.MediaTypeImage,
				Format:    strings.TrimPrefix(ext, "."),
				IsValid:   true,
			})
		case supportedVideoExt[ext]:
			assets = append(assets, media.Asset{
				Path:      filepath.Join(inputDir, name),
				MediaType: media.MediaTypeVideo,
				Format:    strings.TrimPrefix(ext, "."),
				IsValid:   true,
			})
		}
	}
	sort.SliceStable(assets, func(i, j int) bool { return strings.ToLower(assets[i].Path) < strings.ToLower(assets[j].Path) })
	if len(assets) == 0 {
		return nil, fmt.Errorf("no supported media files found")
	}
	for i := range assets {
		assets[i].OrderIndex = i
	}
	return assets, nil
}

func ListMP3Assets(inputDir string) ([]string, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	assets := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !supportedAudioExt[strings.ToLower(filepath.Ext(e.Name()))] {
			continue
		}
		assets = append(assets, filepath.Join(inputDir, e.Name()))
	}
	sort.SliceStable(assets, func(i, j int) bool {
		ai := strings.ToLower(filepath.Base(assets[i]))
		aj := strings.ToLower(filepath.Base(assets[j]))
		if ai == aj {
			return filepath.Base(assets[i]) < filepath.Base(assets[j])
		}
		return ai < aj
	})
	return assets, nil
}
