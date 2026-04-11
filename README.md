# pic2video

Create YouTube-ready slideshow videos from photos using cross-fade transitions plus final fade-in/fade-out.

## Overview

`pic2video` is a Go CLI that converts ordered media (images and videos) into a single 16:9 video.
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

Render with defaults (UHD profile, current directory input, generated output under `./output`):

```bash
./bin/pic2video render
```

Render FHD:

```bash
./bin/pic2video render \
  --profile fhd \
  --image-effect kenburns-medium
```

Render UHD:

```bash
./bin/pic2video render \
  --input ./photos \
  --profile uhd
```

Force CPU encoding:

```bash
./bin/pic2video render \
  --input ./photos \
  --profile fhd \
  --encoder cpu
```

Render with auto-detected MP3 background audio from input directory:

```bash
./bin/pic2video render \
  --input ./photos-with-audio \
  --profile uhd
```

Render with EXIF footer overlay:

```bash
./bin/pic2video render \
  --input ./photos \
  --profile fhd \
  --exif-overlay \
  --exif-font-size 42
```

Render mixed photo/video media with explicit output fps:

```bash
./bin/pic2video render \
  --input ./mixed-media \
  --profile fhd \
  --fps 30
```

Status output example:

```text
status=starting files=3 format=MP4
details: input=./photos output=output/slideshow_fhd.mp4 profile=fhd effect=kenburns-medium encoder=auto overwrite=true
media: images=3 videos=0 fps=60
timing: image-duration=5.0s transition-duration=1.0s
order: mode=exif order-file=-
audio: files=2 order=alphabetical
exif-overlay: enabled=true font-size=42 footer-offset=30 box-alpha=0.40
status=success
result: profile=fhd resolution=1920x1080 encoder:auto->nvenc processed=3 files=3
media: images=3 videos=0 fps=60
exif-overlay: enabled=true font-size=42 footer-offset=30 box-alpha=0.40
output: format=MP4 elapsed=< 1s output=output/slideshow_fhd.mp4 warnings=0
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
- `--fps <int>` (optional; valid range `24` to `60`; default: profile default `60`)
- `--exif-overlay` (default: `false`; enables EXIF metadata footer overlay)
- `--exif-font-size <int>` (default: `42`; valid range: `36` to `60`, enforced when `--exif-overlay` is enabled)
- `--overwrite` (default: true — overwrite output if it exists)

Output path policy:

- The render command always writes to `./output/slideshow_fhd.mp4` or `./output/slideshow_uhd.mp4` based on the selected profile.
- The `output/` directory is created automatically before rendering.
- The CLI no longer accepts a custom output filename.

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

- The CLI scans supported image extensions (`.jpg`, `.jpeg`, `.png`, `.webp`) and supported video extensions (`.mp4`, `.mov`, `.mkv`, `.webm`).
- The CLI also scans `.mp3` files in the same input directory and includes them in ascending alphabetical filename order.
- Unsupported files are skipped during asset discovery.
- Unsupported audio types (for example `.wav`) are ignored.
- If discovered MP3 files cannot be opened, rendering fails with input-validation error (exit `3`).
- If no supported media assets remain, the command fails with input-validation error (exit `3`).
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
- Generated outputs are written under `./output`; if you want a clean rerun, remove the previous file from that directory first.
