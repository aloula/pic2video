package gui

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/infra/ffmpeg"
	"github.com/loula/pic2video/internal/infra/fsio"
)

func shouldConfirmShortAudio(cfg GuiRunConfiguration) (bool, string, error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.AudioSource))
	if mode == "" {
		mode = "mp3"
	}
	if mode == "video" {
		return false, "", nil
	}

	assets, err := fsio.ListMixedAssets(cfg.InputFolder)
	if err != nil {
		return false, "", err
	}
	if len(assets) == 0 {
		return false, "", nil
	}

	videoDur, err := expectedVideoDurationSec(assets, cfg.ImageDuration, cfg.Transition, cfg.FFprobeBin)
	if err != nil {
		return false, "", err
	}
	if videoDur <= 0 {
		return false, "", nil
	}

	mp3Assets, err := fsio.ListMP3Assets(cfg.InputFolder)
	if err != nil || len(mp3Assets) == 0 {
		return false, "", nil
	}

	audioDur, err := totalAudioDurationSec(mp3Assets, cfg.FFprobeBin)
	if err != nil {
		return false, "", err
	}
	if audioDur+0.001 >= videoDur {
		return false, "", nil
	}

	msg := fmt.Sprintf(
		"Total audio duration (%.1fs) is shorter than expected video duration (%.1fs).\n\nDo you want to continue?",
		audioDur,
		videoDur,
	)
	return true, msg, nil
}

func expectedVideoDurationSec(assets []media.Asset, imageDur, transitionDur float64, ffprobeBin string) (float64, error) {
	if len(assets) == 0 {
		return 0, nil
	}

	total := 0.0
	for _, a := range assets {
		slotDur := imageDur
		if a.MediaType == media.MediaTypeVideo {
			p, err := ffmpeg.ProbeMedia(ffprobeBin, a.Path)
			if err != nil {
				return 0, err
			}
			if p.DurationSec > slotDur {
				slotDur = p.DurationSec
			}
		}
		total += slotDur
	}
	total -= float64(len(assets)-1) * transitionDur
	if total < 0 {
		return 0, nil
	}
	return total, nil
}

func totalAudioDurationSec(paths []string, ffprobeBin string) (float64, error) {
	total := 0.0
	for _, p := range paths {
		d, err := probeFormatDurationSec(ffprobeBin, p)
		if err != nil {
			return 0, err
		}
		total += d
	}
	return total, nil
}

func probeFormatDurationSec(ffprobeBin, path string) (float64, error) {
	if strings.TrimSpace(ffprobeBin) == "" {
		ffprobeBin = "ffprobe"
	}
	out, err := exec.Command(
		ffprobeBin,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	).CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("ffprobe duration failed for %s: %w", path, err)
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ffprobe duration for %s: %w", path, err)
	}
	if v < 0 {
		return 0, nil
	}
	return v, nil
}
