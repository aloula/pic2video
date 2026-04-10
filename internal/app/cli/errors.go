package cli

import "fmt"

func FormatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: %v", err)
}
