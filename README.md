# pic2video

Create YouTube-ready slideshow videos from photos using cross-fade transitions.

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

Status output example:

```text
status=starting files=3 format=MP4
details: input=./photos output=./out/slideshow-cpu.mp4 profile=fhd encoder=auto overwrite=true
timing: image-duration=5.0s transition-duration=1.0s
order: mode=exif order-file=-
status=success
result: profile=fhd resolution=1920x1080 encoder:auto->nvenc processed=3 files=3
output: format=MP4 elapsed=< 1s output=./out/slideshow-cpu.mp4 warnings=0
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
- `--image-duration <seconds>` (default: `5`)
- `--transition-duration <seconds>` (default: `1`)
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
- Unsupported files are skipped during asset discovery.
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
