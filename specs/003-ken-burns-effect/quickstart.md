# Quickstart: Ken Burns Effect Option

**Feature**: `003-ken-burns-effect`  
**Branch**: `003-ken-burns-effect`  
**Date**: 2026-04-09

## Basic Usage

Default behavior (static images):

```bash
./bin/pic2video render --input ./photos --profile uhd
```

Ken Burns low:

```bash
./bin/pic2video render --input ./photos --profile fhd --image-effect kenburns-low
```

Ken Burns medium:

```bash
./bin/pic2video render --input ./photos --profile fhd --image-effect kenburns-medium
```

Ken Burns high:

```bash
./bin/pic2video render --input ./photos --profile uhd --image-effect kenburns-high
```

All renders include a global fade-in and fade-out on the final composed video.
The fade duration is derived from `--transition-duration`.

## Expected Startup Output

```text
status=starting files=12 format=MP4
details: input=./photos output=slideshow_uhd.mp4 profile=uhd effect=kenburns-high encoder=auto overwrite=true
timing: image-duration=5.0s transition-duration=1.0s
order: mode=name order-file=-
```

## Test Commands

```bash
go test ./tests/unit/...
go test ./tests/e2e/...
go test ./...
```
