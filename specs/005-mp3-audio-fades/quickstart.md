# Quickstart: MP3 Audio and Video Fades

**Feature**: `005-mp3-audio-fades`  
**Branch**: `005-mp3-audio-fades`  
**Date**: 2026-04-10

## Basic Usage

Render slideshow as usual (no MP3 present):

```bash
./bin/pic2video render --input ./photos --profile fhd
```

Render slideshow with MP3 files in the same directory:

```bash
./bin/pic2video render --input ./photos-with-audio --profile uhd
```

Expected behavior for mixed directory:

- MP3 tracks are auto-included
- MP3 order is alphabetical by filename
- Unsupported audio files (for example WAV) are ignored
- Output remains bounded to slideshow duration
- Video includes fade-in and fade-out

## Expected Startup Output (mixed media)

```text
status=starting files=3 format=MP4
details: input=./photos-with-audio output=slideshow_uhd.mp4 profile=uhd effect=static encoder=auto overwrite=true
timing: image-duration=5.0s transition-duration=1.0s
order: mode=name order-file=-
audio: files=2 order=alphabetical
```

## Example Input Layout

```text
photos-with-audio/
├── 001.jpg
├── 002.jpg
├── 010.jpg
├── ambient_a.mp3
└── ambient_b.mp3
```

For this input, output audio segment order is:

1. `ambient_a.mp3`
2. `ambient_b.mp3`

## Test Commands

```bash
go test ./tests/unit/...
go test ./tests/e2e/...
go test ./...
```
