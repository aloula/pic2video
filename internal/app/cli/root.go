package cli

import (
	"fmt"
	"os"

	appversion "github.com/loula/pic2video/internal/app/version"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pic2video",
		Short:         "Create 16:9 slideshow videos from photos",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       appversion.Short(),
	}
	cmd.SetVersionTemplate("{{.Use}} {{.Version}}\n")
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print detailed build version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), appversion.Info())
		},
	})
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
