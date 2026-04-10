package cli

import "github.com/loula/pic2video/internal/app/renderjob"

func ToExitCode(err error) int {
	return renderjob.ExitCode(err)
}

func ExitCode(err error) int {
	return ToExitCode(err)
}
