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

Render FHD:

```bash
./bin/pic2video render \
  --input ./photos \
  --output ./out/slideshow-fhd.mp4 \
  --profile fhd \
  --image-duration 4 \
  --transition-duration 1
```

Render UHD:

```bash
./bin/pic2video render \
  --input ./photos \
  --output ./out/slideshow-uhd.mp4 \
  --profile uhd
```

Force CPU encoding:

```bash
./bin/pic2video render \
  --input ./photos \
  --output ./out/slideshow-cpu.mp4 \
  --profile fhd \
  --encoder cpu
```

Explicit ordering:

```bash
./bin/pic2video render \
  --input ./photos \
  --output ./out/slideshow-ordered.mp4 \
  --profile fhd \
  --order explicit \
  --order-file ./order.txt
```

## Command flags

`pic2video render` supports:

- `--input <dir>` (required)
- `--output <file>` (required)
- `--profile <fhd|uhd>` (required)
- `--image-duration <seconds>` (default: `4`)
- `--transition-duration <seconds>` (default: `1`)
- `--order <name|time|explicit>` (default: `name`)
- `--order-file <file>` (required with `--order explicit`)
- `--encoder <auto|nvenc|cpu>` (default: `auto`)
- `--overwrite` (overwrite output if it exists)

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
