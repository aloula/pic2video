# CLI Contract: pic2video

## Command
`pic2video render`

## Purpose
Create a YouTube-ready slideshow video from photos with cross-fade transitions.

## Inputs
- `--input <path>`: Required. Directory containing source images.
- `--output <path>`: Required. Output video file path.
- `--profile <fhd|uhd>`: Required. Output profile (`fhd` = 1920x1080, `uhd` = 3840x2160).
- `--image-duration <seconds>`: Optional. Default defined by app config.
- `--transition-duration <seconds>`: Optional. Cross-fade duration.
- `--order <name|time|explicit>`: Optional. Input ordering strategy.
- `--encoder <auto|nvenc|cpu>`: Optional. Default `auto` (prefer NVENC).
- `--overwrite`: Optional. Allow replacing existing output file.

## Output (stdout)
On success, print a concise summary object in human-readable format with at least:
- profile
- effective resolution
- effective encoder
- processed/skipped image counts
- elapsed time
- output path

## Errors (stderr + exit codes)
- `2`: Invalid arguments (missing required flags, invalid durations, unknown profile)
- `3`: Input validation errors (path missing, unreadable files, insufficient valid images)
- `4`: Environment errors (FFmpeg/FFprobe unavailable, unsupported encoder request)
- `5`: Render execution failures (FFmpeg process failed)

Errors MUST include actionable text.

## Behavioral Guarantees
- Output video MUST be 16:9 and match selected profile dimensions.
- Transition type in v1 is fixed to cross-fade between adjacent images.
- Encoder strategy for `auto` MUST prefer NVENC when available and fall back to CPU otherwise.
- Command behavior MUST be deterministic for identical inputs/options.

## Compatibility Notes
- v1 scope excludes audio track generation/mixing.
- Additional transitions and custom resolutions are future-version concerns.
