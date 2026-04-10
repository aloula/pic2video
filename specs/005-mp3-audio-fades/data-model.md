# Data Model: MP3 Audio and Video Fades

**Phase**: 1 - Design & Contracts  
**Branch**: `005-mp3-audio-fades`  
**Date**: 2026-04-10

## Entities

### AudioAsset

Represents an MP3 file eligible for inclusion in output audio.

| Field | Type | Description |
|------|------|-------------|
| path | string | Absolute or workspace-relative asset path |
| file_name | string | Base filename used for ordering |
| order_index | int | Alphabetical order position |
| readable | bool | Whether file can be opened/probed |

Validation rules:
- Extension must be `.mp3` (case-insensitive).
- `file_name` sorting is deterministic and stable.
- `readable` must be true for all selected assets before render execution.

### AudioSelection

Represents the ordered audio set attached to one render request.

| Field | Type | Description |
|------|------|-------------|
| assets | []AudioAsset | Ordered MP3 asset collection |
| has_audio | bool | True when at least one MP3 is selected |
| source_dir | string | Input directory scanned |

Validation rules:
- `has_audio` is false when `assets` is empty.
- `assets` order must follow ascending normalized filename.

### RenderTimeline

Represents slideshow timing bounds used for final output.

| Field | Type | Description |
|------|------|-------------|
| image_count | int | Number of selected image assets |
| image_duration_sec | float64 | Per-image duration |
| transition_duration_sec | float64 | Cross-fade duration |
| total_duration_sec | float64 | Computed video duration |

Validation rules:
- `total_duration_sec > 0` for valid render jobs.
- `transition_duration_sec < image_duration_sec` for valid overlap math.

### FadeProfile

Represents final visual fade directives.

| Field | Type | Description |
|------|------|-------------|
| fade_in_start_sec | float64 | Always `0` |
| fade_out_start_sec | float64 | `total_duration_sec - fade_duration_sec` |
| fade_duration_sec | float64 | Derived/clamped fade duration |

Validation rules:
- `fade_duration_sec > 0`
- `fade_out_start_sec >= 0`
- `fade_duration_sec <= total_duration_sec / 2`

## State Transitions

1. CLI receives input directory and render options.
2. Image assets and MP3 assets are discovered from the input directory.
3. MP3 assets are filtered and sorted alphabetically.
4. Render timeline and fade profile are computed.
5. ffmpeg command is assembled with video graph and optional ordered audio inputs.
6. Output is produced with deterministic duration and valid fade behavior.
