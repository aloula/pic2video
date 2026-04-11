# Quickstart: Video Processing Status

**Feature**: `002-video-processing-status`
**Branch**: `002-video-processing-status`
**Date**: 2026-04-09

---

## What Changed

The `render` command now emits two status lines:

1. **Before encoding starts** — a pre-render announcement showing discovered file count and output format
2. **After completion** — the existing summary line, extended with `files=` and `format=` fields and a human-readable `elapsed=` value

---

## Running a Render (updated output)

```bash
pic2video render \
  --input ./photos \
  --output ./slideshow.mp4 \
  --profile fhd

# stdout:
# status=starting files=12 format=MP4
# details: input=./photos output=./slideshow.mp4 profile=fhd effect=static encoder=auto overwrite=true
# timing: image-duration=5.0s transition-duration=1.0s
# order: mode=name order-file=-
# status=success
# result: profile=fhd resolution=1920x1080 encoder=cpu processed=12 files=12
# output: format=MP4 elapsed=45.3s output=./slideshow.mp4 warnings=0
```

---

## New Fields in Completion Summary

| Field | Meaning | Example |
|-------|---------|---------|
| `files=N` | Number of input images processed | `files=12` |
| `format=EXT` | Output video container derived from file extension | `format=MP4` |
| `elapsed=...` | Total processing time, human-readable | `elapsed=1m 30s` |

---

## Elapsed Time Formats

| Duration | Format | Example |
|----------|--------|---------|
| Sub-second | `< 1s` | `< 1s` |
| Under 60 seconds | `%.1fs` | `45.3s` |
| 60 seconds or more | `Xm Ys` | `1m 30s` |

---

## Output Format Labels

Derived from the `--output` file extension:

| Extension | `format=` value |
|-----------|-----------------|
| `.mp4` | `MP4` |
| `.mov` | `MOV` |
| `.mkv` | `MKV` |
| Any other | Raw uppercased extension, e.g. `AVI` |
| No extension | `UNKNOWN` |

---

## Build

No build changes. This feature modifies only `internal/app/cli/summary.go` and `internal/app/cli/render_command.go`. Build steps are unchanged:

```bash
make build
# or
go build -o bin/pic2video ./cmd/pic2video
```

---

## Run Tests

```bash
# Unit tests (includes new status formatter tests)
go test ./tests/unit/...

# End-to-end tests (includes updated smoke test assertions)
go test ./tests/e2e/...

# All tests
go test ./...
```

---

## Backward Compatibility Notes

- `processed=N` is retained in the completion line — existing scripts relying on it continue to work
- `elapsed=` key name is unchanged; only the value format changes from raw decimal to human-readable
- `files=N` is a new field added after `processed=N`; key-value parsers that ignore unknown keys are unaffected
