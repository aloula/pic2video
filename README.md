# pic2video

Create YouTube-ready slideshow videos from photos using cross-fade transitions plus final fade-in/fade-out.

## Overview

`pic2video` is a Go CLI that converts ordered images into a single 16:9 video.
It supports these output profiles:

- `fhd` -> 1920x1080
- `uhd` -> 3840x2160

The renderer uses FFmpeg/FFprobe and prefers NVIDIA NVENC when available.

## Prerequisites

- Go 1.23+
- FFmpeg and FFprobe available in `PATH`
- Optional: NVIDIA GPU with NVENC-capable FFmpeg build

## Build

Build host binary:

```bash
go build -o bin/pic2video ./cmd/pic2video
```

Or with make:

```bash
make build
```

### Cross-platform builds

```bash
make build-all
```

Generated artifacts:

- `bin/pic2video-linux-amd64`
- `bin/pic2video-darwin-amd64`
- `bin/pic2video-windows-amd64.exe`

You can also build each platform separately:

```bash
make build-linux
make build-macos
make build-windows
```

## Usage

Render with defaults (UHD profile, current directory input, slideshow_uhd.mp4 output):

```bash
./bin/pic2video render
```

Render FHD with custom output location:

```bash
./bin/pic2video render \
  --profile fhd \
  --image-effect kenburns-medium \
  --output ./out/slideshow.mp4
```

Render UHD:

```bash
./bin/pic2video render \
  --input ./photos \
  --output ./out/slideshow-uhd.mp4
```

Force CPU encoding:

```bash
./bin/pic2video render \
  --input ./photos \
  --output ./out/slideshow-cpu.mp4 \
  --profile fhd \
  --encoder cpu
```

Render with auto-detected MP3 background audio from input directory:

```bash
./bin/pic2video render \
  --input ./photos-with-audio \
  --profile uhd
```

Status output example:

```text
status=starting files=3 format=MP4
details: input=./photos output=./out/slideshow-cpu.mp4 profile=fhd effect=kenburns-medium encoder=auto overwrite=true
timing: image-duration=5.0s transition-duration=1.0s
order: mode=exif order-file=-
audio: files=2 order=alphabetical
status=success
result: profile=fhd resolution=1920x1080 encoder:auto->nvenc processed=3 files=3
output: format=MP4 elapsed=< 1s output=./out/slideshow-cpu.mp4 warnings=0
```

Use Ken Burns effect (resolution-aware, high quality):

```bash
./bin/pic2video render \
  --input ./photos \
  --profile uhd \
  --image-effect kenburns-high
```

Explicit ordering:

```bash
./bin/pic2video render \
  --profile fhd \
  --order explicit \
  --order-file ./order.txt
```

EXIF date-based ordering (use photo capture time):

```bash
./bin/pic2video render \
  --order exif
```

## Command flags

`pic2video render` supports:

- `--input <dir>` (default: current directory)
- `--output <file>` (default: `slideshow_fhd.mp4` or `slideshow_uhd.mp4` based on profile)
- `--profile <fhd|uhd>` (default: `uhd`)
- `--image-effect <static|kenburns-low|kenburns-medium|kenburns-high>` (default: `static`)
  - `static`: no zoom/pan motion
  - `kenburns-low`: subtle slow motion
  - `kenburns-medium`: balanced cinematic motion
  - `kenburns-high`: stronger cinematic motion
- `--image-duration <seconds>` (default: `5`)
- `--transition-duration <seconds>` (default: `1`)
  - Also controls global output fade-in/fade-out duration (clamped to half of total video length)
- `--order <name|time|exif|explicit>` (default: `name`)
  - `name`: alphabetical filename order
  - `time`: OS file modification time
  - `exif`: EXIF capture date/time when present, with mod-time fallback
  - `explicit`: manifest file order
- `--order-file <file>` (required with `--order explicit`)
- `--encoder <auto|nvenc|cpu>` (default: `auto`)
- `--overwrite` (default: true — overwrite output if it exists)

## Render status fields

- `status=starting`: pre-render announcement with discovered input count and output container format
- `files=<N>`: number of input images processed
- `format=<EXT>`: output container label derived from output extension (for example, `MP4`, `MOV`, `AVI`, `UNKNOWN`)
- `elapsed=<value>`: human-readable processing time
- `audio: files=<N> order=<MODE>`: discovered MP3 file count and ordering mode (`alphabetical` or `-` when no MP3 is present)

Elapsed format rules:

- `< 1s` for sub-second renders
- `%.1fs` for durations under 60 seconds (for example, `45.3s`)
- `Xm Ys` for durations of 60 seconds or more (for example, `1m 30s`)

## Exit codes

- `0`: success
- `2`: invalid arguments
- `3`: input validation error
- `4`: environment error (for example, missing ffmpeg)
- `5`: render execution failure

## Validation policy (invalid media)

- The CLI scans supported image extensions only (`.jpg`, `.jpeg`, `.png`, `.webp`).
- The CLI also scans `.mp3` files in the same input directory and includes them in ascending alphabetical filename order.
- Unsupported files are skipped during asset discovery.
- Unsupported audio types (for example `.wav`) are ignored.
- If discovered MP3 files cannot be opened, rendering fails with input-validation error (exit `3`).
- If no supported images remain, the command fails with input-validation error (exit `3`).
- Rendering requires at least 2 valid image assets; otherwise the command fails with exit `3`.
- For explicit ordering mode, invalid/missing manifest entries cause input-validation failure (exit `3`).

This policy satisfies FR-012 with deterministic reject/skip behavior.

## Quality rubric acceptance (SC-003)

Quality acceptance is evaluated using the four rubric dimensions below:

- Geometric integrity: enforced by deterministic framing filter (scale + pad, no stretch).
- Framing consistency: one framing policy is applied for all frames in a render run.
- Sharpness preservation: warnings are emitted when source resolution is below target profile.
- Transition smoothness: cross-fade timeline uses deterministic offsets across segments.
- Intro/outro continuity: final output always applies global fade-in and fade-out.

Validation mapping:

- Unit: `tests/unit/framing_policy_test.go`
- Unit: `tests/unit/quality_warning_test.go`
- Unit: `tests/unit/timeline_test.go`
- E2E: `tests/e2e/render_mixed_aspect_test.go`

## Testing

Run all tests:

```bash
go test ./...
```

Run with make:

```bash
make test
```

Run performance-gated tests:

```bash
RUN_PERF=1 go test ./tests/e2e -run TestPerf -count=1
```

## Troubleshooting

- If FFmpeg is not found, ensure `ffmpeg` and `ffprobe` are installed and on `PATH`.
- If `--order explicit` fails, verify each line in the manifest points to a valid image filename/path.
- If output exists and render fails, add `--overwrite`.
