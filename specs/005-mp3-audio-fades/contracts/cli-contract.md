# CLI Contract: MP3 Audio and Video Fades

**Feature**: `005-mp3-audio-fades`  
**Date**: 2026-04-10

## Input Discovery Contract

For `pic2video render` with `--input <dir>`:

- Images are discovered as before from supported image extensions.
- MP3 files in the same directory are discovered automatically.
- MP3 ordering is ascending alphabetical by filename.
- Non-MP3 audio files are ignored by this feature scope.

## Render Output Contract

- If one or more MP3 files are discovered, output contains an audio stream assembled in deterministic alphabetical order.
- If no MP3 files are discovered, output remains valid and preserves existing video-only behavior.
- Final output duration remains bounded by slideshow timeline rules.
- Audio composition uses deterministic MP3 concatenation and is trimmed to slideshow duration bounds.

## Startup Status Contract

- Startup output includes `audio: files=<N> order=<MODE>`.
- `<MODE>` is `alphabetical` when one or more MP3 files are discovered.
- `<MODE>` is `-` when no MP3 files are discovered.

## Fade Contract

- Final video stream always includes visual fade-in and fade-out.
- Fade-in starts at `st=0`.
- Fade-out starts at `total_duration - fade_duration`.
- `fade_duration` is derived from transition duration and clamped to at most half total duration.

## Error Contract

- Corrupt or unreadable selected MP3 assets produce input-validation classification and non-zero exit code.
- Directories with no valid images continue to fail as input validation errors under existing policy.
