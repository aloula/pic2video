package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/gui"
)

func TestLogStoreAppendsInOrder(t *testing.T) {
	store := gui.NewLogStore(10)
	store.Append("stdout", "line-1")
	store.Append("stderr", "line-2")
	store.Append("stdout", "line-3")
	entries := store.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Message != "line-1" || entries[1].Message != "line-2" || entries[2].Message != "line-3" {
		t.Fatalf("unexpected log order: %+v", entries)
	}
}

func TestLogStoreIsBounded(t *testing.T) {
	store := gui.NewLogStore(2)
	store.Append("stdout", "line-1")
	store.Append("stdout", "line-2")
	store.Append("stdout", "line-3")
	entries := store.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after cap, got %d", len(entries))
	}
	if !strings.Contains(store.Text(), "line-2") || !strings.Contains(store.Text(), "line-3") {
		t.Fatalf("expected most recent entries in text, got %s", store.Text())
	}
}
