# CLI Contract: Ken Burns Effect Option

**Feature**: `003-ken-burns-effect`  
**Date**: 2026-04-09

## Render Flag Contract

New flag on `pic2video render`:

- `--image-effect <mode>`

Supported values:
- `static` (default)
- `kenburns-low`
- `kenburns-medium`
- `kenburns-high`

Validation behavior:
- Invalid value returns invalid-arguments classification and non-zero exit code.

## Startup Status Contract

Startup details output includes selected effect:

`details: input=<DIR> output=<PATH> profile=<PROFILE> effect=<MODE> encoder=<ENCODER> overwrite=<BOOL>`

Example:

`details: input=./photos output=./slideshow_uhd.mp4 profile=uhd effect=kenburns-medium encoder=auto overwrite=true`

## Final Output Fade Contract

All generated videos apply final stream fades:

- Fade-in starts at `st=0`
- Fade-out starts at `total_duration - fade_duration`
- `fade_duration` defaults to `--transition-duration` and is clamped to at most half of total output duration

## Backward Compatibility

- Existing behavior remains unchanged when `--image-effect` is omitted (`static`).
- Existing output fields remain; this feature adds effect visibility in startup details.
