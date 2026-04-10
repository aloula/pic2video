package ffmpeg

import (
"encoding/json"
"fmt"
"os/exec"
"strconv"

"github.com/loula/pic2video/internal/domain/media"
)

type probeResult struct {
	Streams []struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Codec  string `json:"codec_name"`
	} `json:"streams"`
}

func ProbeImage(ffprobeBin, path string) (media.Asset, error) {
	if ffprobeBin == "" {
		ffprobeBin = "ffprobe"
	}
	cmd := exec.Command(ffprobeBin, "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,codec_name", "-of", "json", path)
	out, err := cmd.Output()
	if err != nil {
		return media.Asset{}, err
	}
	var pr probeResult
	if err := json.Unmarshal(out, &pr); err != nil {
		return media.Asset{}, err
	}
	if len(pr.Streams) == 0 {
		return media.Asset{}, fmt.Errorf("no stream metadata for %s", path)
	}
	s := pr.Streams[0]
	return media.Asset{Path: path, Width: s.Width, Height: s.Height, Format: s.Codec, IsValid: s.Width > 0 && s.Height > 0, ValidationWarnings: []string{strconv.Itoa(s.Width), strconv.Itoa(s.Height)}[:0]}, nil
}
