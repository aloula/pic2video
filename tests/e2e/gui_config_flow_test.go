package e2e

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/gui"
)

func TestGUIConfigFlowBuildsExpectedRenderArgs(t *testing.T) {
	cfg := gui.GuiRunConfiguration{
		InputFolder:   "/tmp/in",
		OutputFolder:  "/tmp/out",
		Profile:       "uhd",
		ImageEffect:   "kenburns-medium",
		ImageDuration: 5,
		Transition:    1,
		FPS:           60,
		OrderMode:     "name",
		ExifOverlay:   true,
		ExifFontSize:  42,
		Encoder:       "auto",
		Overwrite:     true,
	}
	args := gui.BuildRenderCommandArgs(cfg)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--input /tmp/in") || !strings.Contains(joined, "--profile uhd") {
		t.Fatalf("expected input/profile in args, got: %s", joined)
	}
	if strings.Contains(joined, "--output") {
		t.Fatalf("did not expect output flag in GUI args, got: %s", joined)
	}
}
