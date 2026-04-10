package nvenc

import "fmt"

func SelectEncoder(requested string, hasNVENC bool) (string, error) {
	switch requested {
	case "", "auto":
		if hasNVENC {
			return "h264_nvenc", nil
		}
		return "libx264", nil
	case "nvenc":
		if !hasNVENC {
			return "", fmt.Errorf("nvenc requested but not available")
		}
		return "h264_nvenc", nil
	case "cpu":
		return "libx264", nil
	default:
		return "", fmt.Errorf("unknown encoder option: %s", requested)
	}
}
