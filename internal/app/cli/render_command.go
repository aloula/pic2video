package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/app/renderjob"
	"github.com/loula/pic2video/internal/infra/fsio"
	"github.com/loula/pic2video/internal/infra/nvenc"
	"github.com/spf13/cobra"
)

func newRenderCommand() *cobra.Command {
	var input, output, profileName, imageEffect, orderMode, orderFile, encoder string
	var imageDur, transDur float64
	var overwrite bool
	var ffmpegBin, ffprobeBin string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render slideshow from image folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Apply defaults for optional flags
			if input == "" {
				input = "."
			}
			if profileName == "" {
				profileName = "uhd"
			}
			if output == "" {
				// Generate output filename based on profile
				if profileName == "fhd" {
					output = "slideshow_fhd.mp4"
				} else {
					output = "slideshow_uhd.mp4"
				}
			}
			if imageEffect == "" {
				imageEffect = "static"
			}

			// Validate required orderMode values
			if orderMode != "name" && orderMode != "time" && orderMode != "exif" && orderMode != "explicit" {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInvalidArguments, Msg: "--order must be one of name|time|exif|explicit"}
			}
			if imageEffect != "static" && imageEffect != "kenburns-low" && imageEffect != "kenburns-medium" && imageEffect != "kenburns-high" {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInvalidArguments, Msg: "--image-effect must be one of static|kenburns-low|kenburns-medium|kenburns-high"}
			}
			if orderMode == "explicit" && orderFile == "" {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInvalidArguments, Msg: "--order-file required for --order explicit"}
			}
			assets, err := fsio.ListImageAssets(input)
			if err != nil {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInputValidation, Msg: "failed to read input assets", Err: err}
			}
			audioAssets, err := fsio.ListMP3Assets(input)
			if err != nil {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInputValidation, Msg: "failed to read input audio assets", Err: err}
			}
			for _, p := range audioAssets {
				f, openErr := os.Open(p)
				if openErr != nil {
					return &renderjob.ClassifiedError{Class: renderjob.ErrInputValidation, Msg: "failed to open audio asset", Err: openErr}
				}
				_ = f.Close()
			}
			explicit := []string(nil)
			if orderMode == "explicit" {
				explicit, err = fsio.ReadExplicitOrder(orderFile)
				if err != nil {
					return &renderjob.ClassifiedError{Class: renderjob.ErrInputValidation, Msg: "failed to parse explicit order file", Err: err}
				}
			}
			assets = pipeline.ApplyOrderExt(orderMode, assets, explicit, ffprobeBin)
			audioOrder := "-"
			if len(audioAssets) > 0 {
				audioOrder = "alphabetical"
			}
			fmt.Fprintln(cmd.OutOrStdout(), FormatAnnouncement(StartOptions{
				Input:              input,
				Output:             output,
				Profile:            profileName,
				ImageEffect:        imageEffect,
				ImageDuration:      imageDur,
				TransitionDuration: transDur,
				Order:              orderMode,
				OrderFile:          orderFile,
				AudioFiles:         len(audioAssets),
				AudioOrder:         audioOrder,
				Encoder:            encoder,
				Overwrite:          overwrite,
				Files:              len(assets),
			}))
			job, err := renderjob.BuildJob(renderjob.BuildOptions{
				OutputPath:      output,
				AudioAssets:     audioAssets,
				ProfileName:     profileName,
				ImageEffect:     imageEffect,
				ImageDuration:   imageDur,
				Transition:      transDur,
				Overwrite:       overwrite,
				OrderMode:       orderMode,
				OrderFile:       orderFile,
				RequestedEncode: encoder,
				FFmpegBin:       ffmpegBin,
				FFprobeBin:      ffprobeBin,
			}, assets)
			if err != nil {
				return err
			}
			service := &renderjob.Service{}
			summary, err := service.Run(context.Background(), job)
			if err != nil {
				return err
			}
			has := nvenc.Available(ffmpegBin)
			fmt.Fprintln(cmd.OutOrStdout(), FormatSummary(summary.ProfileName, summary.EffectiveResolution, encoder, summary.EffectiveEncoder, summary.OutputPath, summary.ElapsedSeconds, summary.ProcessedAssets, summary.Warnings, has))
			return nil
		},
	}

	cmd.Flags().StringVar(&input, "input", "", "Input directory containing images (default: current directory)")
	cmd.Flags().StringVar(&output, "output", "", "Output video path (default: slideshow_fhd.mp4 or slideshow_uhd.mp4 based on profile)")
	cmd.Flags().StringVar(&profileName, "profile", "", "Output profile: fhd|uhd (default: uhd)")
	cmd.Flags().StringVar(&imageEffect, "image-effect", "static", "Image effect: static|kenburns-low|kenburns-medium|kenburns-high (default: static)")
	cmd.Flags().Float64Var(&imageDur, "image-duration", 5, "Per-image duration in seconds (default: 5)")
	cmd.Flags().Float64Var(&transDur, "transition-duration", 1, "Cross-fade transition duration in seconds (default: 1)")
	cmd.Flags().StringVar(&orderMode, "order", "name", "Ordering mode: name|time|exif|explicit (default: name)")
	cmd.Flags().StringVar(&orderFile, "order-file", "", "Path to explicit order manifest file")
	cmd.Flags().StringVar(&encoder, "encoder", "auto", "Encoder preference: auto|nvenc|cpu (default: auto)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", true, "Overwrite output file if it exists (default: true)")
	cmd.Flags().StringVar(&ffmpegBin, "ffmpeg-bin", envOrDefault("P2V_FFMPEG_BIN", "ffmpeg"), "Path to ffmpeg binary")
	cmd.Flags().StringVar(&ffprobeBin, "ffprobe-bin", envOrDefault("P2V_FFPROBE_BIN", "ffprobe"), "Path to ffprobe binary")
	_ = cmd.Flags().MarkHidden("ffmpeg-bin")
	_ = cmd.Flags().MarkHidden("ffprobe-bin")
	return cmd
}
