package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/gui"
)

func TestDefaultConfigurationUsesLaunchDirectory(t *testing.T) {
	cfg := gui.DefaultConfiguration()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InputFolder != wd {
		t.Fatalf("expected input default to launch directory %q, got %q", wd, cfg.InputFolder)
	}
	wantOutput := filepath.Join(wd, "output")
	if cfg.OutputFolder != wantOutput {
		t.Fatalf("expected output default to input/output %q, got %q", wantOutput, cfg.OutputFolder)
	}
}

func TestBuildRenderCommandArgsMapsOptions(t *testing.T) {
	cfg := gui.GuiRunConfiguration{
		InputFolder:   "/tmp/in",
		OutputFolder:  "/tmp/out",
		Profile:       "fhd",
		ImageEffect:   "kenburns-low",
		ImageDuration: 4,
		Transition:    1,
		FPS:           30,
		OrderMode:     "name",
		ExifOverlay:   true,
		ExifFontSize:  48,
		Encoder:       "cpu",
		Overwrite:     true,
	}
	args := gui.BuildRenderCommandArgs(cfg)
	joined := strings.Join(args, " ")
	for _, want := range []string{"render", "--input /tmp/in", "--profile fhd", "--image-effect kenburns-low", "--image-duration 4", "--transition-duration 1", "--fps 30", "--exif-overlay", "--exif-font-size 48", "--encoder cpu"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected args to include %q, got %s", want, joined)
		}
	}
	if !strings.Contains(joined, "--audio-source mp3") {
		t.Fatalf("expected default audio source flag in GUI args, got %s", joined)
	}
	if strings.Contains(joined, "--output") {
		t.Fatalf("did not expect --output in GUI args, got %s", joined)
	}
}

func TestValidatePreflightBlocksEmptyMediaFolder(t *testing.T) {
	dir := t.TempDir()
	cfg := gui.GuiRunConfiguration{InputFolder: dir, OutputFolder: dir, Profile: "uhd"}
	res := gui.ValidatePreflight(cfg)
	if res.OK {
		t.Fatal("expected preflight failure for empty media folder")
	}
}

func TestOutputPreviewUsesProfileAutoFileName(t *testing.T) {
	cfg := gui.GuiRunConfiguration{OutputFolder: filepath.Clean("/tmp/output"), Profile: "fhd", OutputFileName: "slideshow_fhd.mp4"}
	preview := gui.OutputPreviewText(cfg)
	if !strings.Contains(preview, "slideshow_fhd.mp4") {
		t.Fatalf("expected auto filename in preview, got %s", preview)
	}
}
