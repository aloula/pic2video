package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/cli"
)

func TestFormatOutputFormat(t *testing.T) {
	cases := []struct {
		name string
		path string
		want string
	}{
		{name: "known extension", path: "/tmp/out.mp4", want: "MP4"},
		{name: "unrecognized extension", path: "/tmp/out.avi", want: "AVI"},
		{name: "missing extension", path: "/tmp/out", want: "UNKNOWN"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cli.FormatOutputFormat(tc.path)
			if got != tc.want {
				t.Fatalf("format mismatch: got=%s want=%s", got, tc.want)
			}
		})
	}
}

func TestFormatElapsed(t *testing.T) {
	cases := []struct {
		name    string
		seconds float64
		want    string
	}{
		{name: "sub-second", seconds: 0.4, want: "< 1s"},
		{name: "under-60s", seconds: 45.34, want: "45.3s"},
		{name: "exactly-60s", seconds: 60.0, want: "1m 0s"},
		{name: "over-60s", seconds: 90.9, want: "1m 30s"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cli.FormatElapsed(tc.seconds)
			if got != tc.want {
				t.Fatalf("elapsed mismatch: got=%s want=%s", got, tc.want)
			}
		})
	}
}

func TestFormatSummaryIncludesNewFields(t *testing.T) {
	got := cli.FormatSummary("fhd", "1920x1080", "auto", "cpu", "/tmp/out.mp4", 45.321, 12, nil, false)
	if !strings.Contains(got, "status=success") {
		t.Fatalf("expected success status in summary: %s", got)
	}
	if !strings.Contains(got, "result:") {
		t.Fatalf("expected result section in summary: %s", got)
	}
	if !strings.Contains(got, "output:") {
		t.Fatalf("expected output section in summary: %s", got)
	}
	if !strings.Contains(got, "processed=12") {
		t.Fatalf("expected processed field in summary: %s", got)
	}
	if !strings.Contains(got, "files=12") {
		t.Fatalf("expected files field in summary: %s", got)
	}
	if !strings.Contains(got, "format=MP4") {
		t.Fatalf("expected format field in summary: %s", got)
	}
	if !strings.Contains(got, "elapsed=45.3s") {
		t.Fatalf("expected human-readable elapsed in summary: %s", got)
	}
}

func TestFormatAnnouncement(t *testing.T) {
	got := cli.FormatAnnouncement(cli.StartOptions{
		Input:              "/tmp/photos",
		Output:             "/tmp/out.mov",
		Profile:            "uhd",
		ImageEffect:        "kenburns-medium",
		ImageDuration:      5,
		TransitionDuration: 1,
		Order:              "exif",
		OrderFile:          "",
		Encoder:            "auto",
		Overwrite:          true,
		Files:              3,
	})
	if !strings.Contains(got, "status=starting") {
		t.Fatalf("expected status in announcement: %s", got)
	}
	if !strings.Contains(got, "files=3") {
		t.Fatalf("expected files in announcement: %s", got)
	}
	if !strings.Contains(got, "format=MOV") {
		t.Fatalf("expected format in announcement: %s", got)
	}
	if !strings.Contains(got, "details:") {
		t.Fatalf("expected details section in announcement: %s", got)
	}
	if !strings.Contains(got, "effect=kenburns-medium") {
		t.Fatalf("expected image effect in announcement details: %s", got)
	}
	if !strings.Contains(got, "timing:") {
		t.Fatalf("expected timing section in announcement: %s", got)
	}
	if !strings.Contains(got, "order:") {
		t.Fatalf("expected order section in announcement: %s", got)
	}
	if !strings.Contains(got, "order-file=-") {
		t.Fatalf("expected order-file placeholder when missing: %s", got)
	}
}
