# pic2video Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-10

## Active Technologies
- Go 1.23+ + Standard library, `github.com/spf13/cobra` (CLI), `github.com/spf13/pflag` (flags), external FFmpeg/FFprobe binaries (001-build-photo-video-cli)
- Go 1.23+ + Standard library only (`path/filepath`, `strings`, `fmt`) — no new external dependencies (002-video-processing-status)
- N/A (filesystem I/O only) (002-video-processing-status)
- Go 1.23+ + Standard library, `github.com/spf13/cobra`, existing ffmpeg/ffprobe integration (003-ken-burns-effect)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.23+

## Code Style

Go 1.23+: Follow standard conventions

## Recent Changes
- 005-mp3-audio-fades: Added Go 1.23+ + Standard library, `github.com/spf13/cobra`, existing ffmpeg/ffprobe integration
- 003-ken-burns-effect: Added Go 1.23+ + Standard library, `github.com/spf13/cobra`, existing ffmpeg/ffprobe integration
- 002-video-processing-status: Added Go 1.23+ + Standard library only (`path/filepath`, `strings`, `fmt`) — no new external dependencies

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
