# Quickstart: EXIF Footer Overlay Option

**Feature**: `007-exif-overlay-option`  
**Branch**: `007-create-feature-branch`  
**Date**: 2026-04-11

## Basic Usage

Render without EXIF overlay (baseline behavior):

```bash
./bin/pic2video render --input ./photos --profile fhd
```

Render with EXIF footer overlay enabled:

```bash
./bin/pic2video render \
  --input ./photos \
  --profile fhd \
  --exif-overlay \
  --exif-font-size 42
```

Render UHD with larger overlay text:

```bash
./bin/pic2video render \
  --input ./photos \
  --profile uhd \
  --exif-overlay \
  --exif-font-size 60
```

## Expected Behavior

- Overlay text appears in footer, 10 pixels above bottom edge in both FHD and UHD.
- Overlay text is white.
- Overlay background is greater than 50% transparent.
- EXIF text follows required field order and date format.
- Missing EXIF fields are shown as `Unknown`.

## Validation Example

Out-of-range font size fails fast before rendering:

```bash
./bin/pic2video render --input ./photos --profile fhd --exif-overlay --exif-font-size 20
```

Expected result:

- Command exits non-zero with invalid-arguments classification.
- Error indicates supported range is 36-60.

## Test Commands

```bash
go test ./tests/unit/...
go test ./tests/e2e/...
go test ./...
```

## Build Validation

- `make build-all` executed successfully on 2026-04-11, producing linux, darwin, and windows binaries.
