package main

import (
	"fmt"
	"os"

	appcli "github.com/loula/pic2video/internal/app/cli"
)

func main() {
	if err := appcli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(appcli.ExitCode(err))
	}
}
