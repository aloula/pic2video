package fsio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/loula/pic2video/internal/domain/media"
)

func TestListMixedAssetsIncludesMediaAndOrder(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"b.mp4", "a.jpg", "c.heic", "ignore.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	assets, err := ListMixedAssets(dir)
	if err != nil {
		t.Fatalf("expected mixed assets, got error: %v", err)
	}
	if len(assets) != 3 {
		t.Fatalf("expected 3 supported assets, got %d", len(assets))
	}
	if assets[0].MediaType != media.MediaTypeImage || filepath.Base(assets[0].Path) != "a.jpg" {
		t.Fatalf("unexpected first sorted asset: %+v", assets[0])
	}
	if assets[1].MediaType != media.MediaTypeVideo || filepath.Base(assets[1].Path) != "b.mp4" {
		t.Fatalf("unexpected second sorted asset: %+v", assets[1])
	}
	if assets[2].MediaType != media.MediaTypeImage || filepath.Base(assets[2].Path) != "c.heic" {
		t.Fatalf("unexpected third sorted asset: %+v", assets[2])
	}
	for i := range assets {
		if assets[i].OrderIndex != i {
			t.Fatalf("unexpected order index at %d: got %d", i, assets[i].OrderIndex)
		}
	}
}

func TestListMixedAssetsReturnsErrorWhenEmpty(t *testing.T) {
	dir := t.TempDir()
	if _, err := ListMixedAssets(dir); err == nil {
		t.Fatal("expected error for empty media directory")
	}
}

func TestListMP3AssetsSortsCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"B.mp3", "a.mp3", "ignore.wav"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	assets, err := ListMP3Assets(dir)
	if err != nil {
		t.Fatalf("expected mp3 list, got error: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("expected 2 mp3 assets, got %d", len(assets))
	}
	if filepath.Base(assets[0]) != "a.mp3" || filepath.Base(assets[1]) != "B.mp3" {
		t.Fatalf("unexpected mp3 sort order: %v", assets)
	}
}
