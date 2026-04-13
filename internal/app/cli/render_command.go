package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/app/renderjob"
	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/infra/fsio"
	"github.com/loula/pic2video/internal/infra/nvenc"
	"github.com/spf13/cobra"
)

func newRenderCommand() *cobra.Command {
	var input, profileName, imageEffect, orderMode, orderFile, encoder, audioSource string
	var imageDur, transDur float64
	var outputFPS int
	var exifOverlay bool
	var exifFontSize int
	var debugExif bool
	exifFooterOffsetPx := 30
	exifBoxAlpha := 0.4
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
			output := defaultOutputPath(profileName)
			if imageEffect == "" {
				imageEffect = "static"
			}
			exifFooterOffsetPx = exifFooterOffsetForProfile(profileName)

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
			if exifOverlay && (exifFontSize < 36 || exifFontSize > 60) {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInvalidArguments, Msg: "--exif-font-size must be between 36 and 60"}
			}
			if outputFPS != 0 && (outputFPS < 24 || outputFPS > 60) {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInvalidArguments, Msg: "--fps must be between 24 and 60"}
			}
			if audioSource != "mp3" && audioSource != "video" && audioSource != "mix" {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInvalidArguments, Msg: "--audio-source must be one of mp3|video|mix"}
			}
			assets, err := fsio.ListMixedAssets(input)
			if err != nil {
				return &renderjob.ClassifiedError{Class: renderjob.ErrInputValidation, Msg: "failed to read input assets", Err: err}
			}
			if outputFPS == 0 {
				outputFPS = 60
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
			if debugExif {
				for i, a := range assets {
					exif, exifErr := fsio.ExtractExif(a.Path, ffprobeBin)
					if exifErr != nil {
						fmt.Fprintf(cmd.OutOrStdout(), "debug-exif: index=%d path=%s error=%v\n", i, a.Path, exifErr)
						continue
					}
					fmt.Fprintf(
						cmd.OutOrStdout(),
						"debug-exif: index=%d path=%s model=%q focal=%q speed=%q aperture=%q iso=%q captured=%s\n",
						i,
						a.Path,
						exif.CameraModel,
						exif.FocalDistance,
						exif.ShutterSpeed,
						exif.Aperture,
						exif.ISO,
						fsio.FormatCapturedDate(exif.CreateDate),
					)
				}
			}
			audioOrder := "-"
			if len(audioAssets) > 0 {
				audioOrder = "alphabetical"
			}
			fmt.Fprintln(cmd.OutOrStdout(), FormatAnnouncement(StartOptions{
				Input:              input,
				Output:             output,
				Profile:            profileName,
				OutputFPS:          outputFPS,
				ImageFiles:         countAssetsByType(assets, "image"),
				VideoFiles:         countAssetsByType(assets, "video"),
				ImageEffect:        imageEffect,
				ImageDuration:      imageDur,
				TransitionDuration: transDur,
				Order:              orderMode,
				OrderFile:          orderFile,
				AudioFiles:         len(audioAssets),
				AudioOrder:         audioOrder,
				AudioSource:        audioSource,
				ExifOverlay:        exifOverlay,
				ExifFontSize:       exifFontSize,
				ExifFooterOffsetPx: exifFooterOffsetPx,
				ExifBoxAlpha:       exifBoxAlpha,
				Encoder:            encoder,
				Overwrite:          overwrite,
				Files:              len(assets),
			}))
			job, err := renderjob.BuildJob(renderjob.BuildOptions{
				OutputPath:         output,
				AudioAssets:        audioAssets,
				AudioSource:        audioSource,
				OutputFPS:          outputFPS,
				ExifOverlay:        exifOverlay,
				ExifFontSize:       exifFontSize,
				ExifFooterOffsetPx: exifFooterOffsetPx,
				ExifBoxAlpha:       exifBoxAlpha,
				ProfileName:        profileName,
				ImageEffect:        imageEffect,
				ImageDuration:      imageDur,
				Transition:         transDur,
				Overwrite:          overwrite,
				OrderMode:          orderMode,
				OrderFile:          orderFile,
				RequestedEncode:    encoder,
				FFmpegBin:          ffmpegBin,
				FFprobeBin:         ffprobeBin,
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
			fmt.Fprintln(cmd.OutOrStdout(), FormatSummaryWithMedia(summary.ProfileName, summary.EffectiveResolution, summary.ExifOverlayEnabled, summary.ExifFontSize, exifFooterOffsetPx, exifBoxAlpha, encoder, summary.EffectiveEncoder, summary.OutputPath, summary.ElapsedSeconds, summary.ProcessedAssets, summary.ImageCount, summary.VideoCount, summary.OutputFPS, summary.Warnings, has))
			return nil
		},
	}

	cmd.Flags().StringVar(&input, "input", "", "Input directory containing images (default: current directory)")
	cmd.Flags().StringVar(&profileName, "profile", "", "Output profile: fhd|uhd (default: uhd)")
	cmd.Flags().StringVar(&imageEffect, "image-effect", "static", "Image effect: static|kenburns-low|kenburns-medium|kenburns-high (default: static)")
	cmd.Flags().Float64Var(&imageDur, "image-duration", 5, "Per-image duration in seconds (default: 5)")
	cmd.Flags().Float64Var(&transDur, "transition-duration", 1, "Cross-fade transition duration in seconds (default: 1)")
	cmd.Flags().IntVar(&outputFPS, "fps", 0, "Output frames per second (24-60, default: profile default)")
	cmd.Flags().StringVar(&orderMode, "order", "name", "Ordering mode: name|time|exif|explicit (default: name)")
	cmd.Flags().StringVar(&orderFile, "order-file", "", "Path to explicit order manifest file")
	cmd.Flags().StringVar(&audioSource, "audio-source", "mp3", "Audio source: mp3|video|mix (default: mp3)")
	cmd.Flags().BoolVar(&exifOverlay, "exif-overlay", false, "Enable EXIF metadata footer overlay")
	cmd.Flags().IntVar(&exifFontSize, "exif-font-size", 42, "EXIF overlay font size (36-60)")
	cmd.Flags().BoolVar(&debugExif, "debug-exif", false, "Print extracted EXIF values for each image before rendering")
	cmd.Flags().StringVar(&encoder, "encoder", "auto", "Encoder preference: auto|nvenc|cpu (default: auto)")
	cmd.Flags().BoolVar(&overwrite, "overwrite", true, "Overwrite output file if it exists (default: true)")
	cmd.Flags().StringVar(&ffmpegBin, "ffmpeg-bin", envOrDefault("P2V_FFMPEG_BIN", "ffmpeg"), "Path to ffmpeg binary")
	cmd.Flags().StringVar(&ffprobeBin, "ffprobe-bin", envOrDefault("P2V_FFPROBE_BIN", "ffprobe"), "Path to ffprobe binary")
	_ = cmd.Flags().MarkHidden("ffmpeg-bin")
	_ = cmd.Flags().MarkHidden("ffprobe-bin")
	return cmd
}

func defaultOutputPath(profileName string) string {
	profileName = strings.ToLower(strings.TrimSpace(profileName))
	if profileName == "fhd" {
		return filepath.Join("output", "slideshow_fhd.mp4")
	}
	return filepath.Join("output", "slideshow_uhd.mp4")
}

func countAssetsByType(assets []media.Asset, mediaType string) int {
	count := 0
	for _, a := range assets {
		if string(a.MediaType) == mediaType {
			count++
		}
	}
	return count
}

func exifFooterOffsetForProfile(profileName string) int {
	if strings.EqualFold(strings.TrimSpace(profileName), "uhd") {
		return 60
	}
	return 30
}
