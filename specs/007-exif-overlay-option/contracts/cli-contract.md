# CLI Contract: EXIF Footer Overlay Option

**Feature**: `007-exif-overlay-option`  
**Date**: 2026-04-11

## New Render Options Contract

For `pic2video render`:

- `--exif-overlay` enables EXIF footer overlay.
- `--exif-font-size <int>` sets overlay font size.
- If `--exif-overlay` is enabled, `--exif-font-size` MUST be between 36 and 60 inclusive.

## Overlay Content Contract

When overlay is enabled, rendered output MUST show text in this exact format and order:

`Camera Model - Focal Distance - Speed (1/XXXXs) - Aperture (f/X) - ISO - Captured Date (DD/MM/YYYY)`

Additional requirements:

- Missing or unavailable values are rendered as `Unknown`.
- Captured date uses `DD/MM/YYYY` when source value is valid.

## Placement and Style Contract

When overlay is enabled:

- Overlay appears in footer for FHD and UHD outputs.
- Overlay baseline is 10 pixels above the bottom edge.
- Overlay text color is white.
- Overlay background is more than 50% transparent.

## Validation and Error Contract

- Out-of-range font size (`<36` or `>60`) produces a classified invalid-arguments error before ffmpeg execution.
- Rendering MUST continue for images with partial EXIF metadata, substituting `Unknown` values.

## Backward Compatibility Contract

- If `--exif-overlay` is not provided, render behavior and output remain unchanged from current baseline.
- Existing flags and ordering behavior (`name|time|exif|explicit`) remain compatible.
