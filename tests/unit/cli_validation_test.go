package unit

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/loula/pic2video/internal/infra/fsio"
)

func TestCLIValidationMissingRequiredFlags(t *testing.T) {
	// With optional flags having sensible defaults, this test now validates
	// that the command fails gracefully when input directory doesn't exist
	cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for non-existent input directory")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if ee.ExitCode() != 3 {
			t.Fatalf("expected exit code 3 (ErrInputValidation), got %d", ee.ExitCode())
		}
	}
}

func TestListMP3AssetsSorted(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"ambient_b.mp3", "ambient_a.mp3", "cover.jpg", "ignored.wav"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := fsio.ListMP3Assets(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{filepath.Join(dir, "ambient_a.mp3"), filepath.Join(dir, "ambient_b.mp3")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mp3 ordering mismatch: got=%v want=%v", got, want)
	}
}

func TestListMP3AssetsNoMP3ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.jpg", "b.png", "ignored.wav"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := fsio.ListMP3Assets(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty mp3 set, got=%v", got)
	}
}

func TestListMP3AssetsCaseVarianceDeterministic(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"Ambient_a.mp3", "ambient_B.mp3", "ambient_a.MP3"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := fsio.ListMP3Assets(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(dir, "Ambient_a.mp3"),
		filepath.Join(dir, "ambient_a.MP3"),
		filepath.Join(dir, "ambient_B.mp3"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("case-variance deterministic ordering mismatch: got=%v want=%v", got, want)
	}
}

func TestCLIValidationExifFontSizeOutOfRange(t *testing.T) {
	for _, size := range []string{"35", "61"} {
		cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory", "--exif-overlay", "--exif-font-size", size)
		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected non-zero exit for out-of-range font size: %s", size)
		}
		if ee, ok := err.(*exec.ExitError); ok {
			if ee.ExitCode() != 2 {
				t.Fatalf("expected exit code 2 (ErrInvalidArguments), got %d for size %s", ee.ExitCode(), size)
			}
		}
	}
}

func TestCLIValidationExifFontSizeBoundariesAccepted(t *testing.T) {
	for _, size := range []string{"36", "60"} {
		cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory", "--exif-overlay", "--exif-font-size", size)
		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected non-zero exit for nonexistent input, size=%s", size)
		}
		if ee, ok := err.(*exec.ExitError); ok {
			if ee.ExitCode() != 3 {
				t.Fatalf("expected exit code 3 (input validation), got %d for size %s", ee.ExitCode(), size)
			}
		}
	}
}

func TestCLIValidationFPSOutOfRange(t *testing.T) {
	for _, fps := range []string{"23", "61"} {
		cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory", "--fps", fps)
		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected non-zero exit for out-of-range fps: %s", fps)
		}
		if ee, ok := err.(*exec.ExitError); ok {
			if ee.ExitCode() != 2 {
				t.Fatalf("expected exit code 2 for invalid fps, got %d", ee.ExitCode())
			}
		}
	}
}

func TestCLIValidationFPSBoundaryAccepted(t *testing.T) {
	for _, fps := range []string{"24", "60"} {
		cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory", "--fps", fps)
		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected non-zero exit for nonexistent input, fps=%s", fps)
		}
		if ee, ok := err.(*exec.ExitError); ok {
			if ee.ExitCode() != 3 {
				t.Fatalf("expected input-validation exit for nonexistent input, got %d", ee.ExitCode())
			}
		}
	}
}
