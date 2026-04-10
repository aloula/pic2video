package unit

import (
	"testing"

	"github.com/loula/pic2video/internal/infra/nvenc"
)

func TestSelectEncoderAutoFallback(t *testing.T) {
	enc, err := nvenc.SelectEncoder("auto", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc != "libx264" {
		t.Fatalf("expected libx264, got %s", enc)
	}
}

func TestSelectEncoderNvencUnavailable(t *testing.T) {
	if _, err := nvenc.SelectEncoder("nvenc", false); err == nil {
		t.Fatal("expected error when nvenc unavailable")
	}
}
