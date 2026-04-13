package nvenc

import "testing"

func TestSelectEncoderAuto(t *testing.T) {
	enc, err := SelectEncoder("auto", true)
	if err != nil {
		t.Fatalf("unexpected error for auto+nvenc: %v", err)
	}
	if enc != "h264_nvenc" {
		t.Fatalf("expected h264_nvenc, got %s", enc)
	}

	enc, err = SelectEncoder("auto", false)
	if err != nil {
		t.Fatalf("unexpected error for auto+cpu: %v", err)
	}
	if enc != "libx264" {
		t.Fatalf("expected libx264, got %s", enc)
	}
}

func TestSelectEncoderRequestedModes(t *testing.T) {
	enc, err := SelectEncoder("cpu", true)
	if err != nil {
		t.Fatalf("unexpected error for cpu: %v", err)
	}
	if enc != "libx264" {
		t.Fatalf("expected libx264, got %s", enc)
	}

	enc, err = SelectEncoder("nvenc", true)
	if err != nil {
		t.Fatalf("unexpected error for nvenc available: %v", err)
	}
	if enc != "h264_nvenc" {
		t.Fatalf("expected h264_nvenc, got %s", enc)
	}

	if _, err := SelectEncoder("nvenc", false); err == nil {
		t.Fatal("expected error when nvenc requested but unavailable")
	}
	if _, err := SelectEncoder("weird", true); err == nil {
		t.Fatal("expected error for unknown encoder option")
	}
}

func TestBuildReport(t *testing.T) {
	if got := BuildReport("auto", "h264_nvenc", true); got != "encoder:auto->nvenc" {
		t.Fatalf("unexpected auto nvenc report: %s", got)
	}
	if got := BuildReport("", "libx264", false); got != "encoder:auto->cpu" {
		t.Fatalf("unexpected auto cpu report: %s", got)
	}
	if got := BuildReport("cpu", "libx264", false); got != "encoder:cpu->libx264" {
		t.Fatalf("unexpected explicit report: %s", got)
	}
}
