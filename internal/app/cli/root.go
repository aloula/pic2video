package cli

import (
"os"

"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pic2video",
		Short:         "Create 16:9 slideshow videos from photos",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(newRenderCommand())
	return cmd
}

func Execute() error {
	return newRootCommand().Execute()
}

func envOrDefault(name, def string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return def
}
