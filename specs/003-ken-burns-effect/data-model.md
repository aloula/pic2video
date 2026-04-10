# Data Model: Ken Burns Effect Option

**Phase**: 1 - Design & Contracts  
**Branch**: `003-ken-burns-effect`  
**Date**: 2026-04-09

## Entities

### ImageEffectMode

Represents user-selected image motion behavior.

| Field | Type | Description |
|------|------|-------------|
| mode | string | One of `static`, `kenburns-low`, `kenburns-medium`, `kenburns-high` |

Validation rules:
- Must match one of the supported values.
- Defaults to `static` when omitted.

### MotionFilterProfile

Derived filter parameters from selected effect mode and output resolution.

| Field | Type | Description |
|------|------|-------------|
| output_width | int | Output width from selected profile |
| output_height | int | Output height from selected profile |
| quality_priority | string | Always `high` for non-static modes |
| zoom_step | float | Per-frame zoom increment |
| zoom_max | float | Maximum zoom level cap |
| pan_amplitude_x | int | Horizontal motion amplitude |
| pan_amplitude_y | int | Vertical motion amplitude |
| fps | int | Motion interpolation target FPS |

Validation rules:
- `fps` > 0
- `zoom_step` > 0 for non-static modes
- `zoom_max` >= 1.0
- output dimensions must match selected profile

### StartupStatusDetails

Operator-visible summary emitted before processing starts.

| Field | Type | Description |
|------|------|-------------|
| effect | string | Selected image effect mode |
| profile | string | Output profile (`fhd`/`uhd`) |
| output | string | Output path |

## State Transitions

1. CLI parses image effect flag.
2. BuildJob stores `ImageEffectMode`.
3. Service builds ffmpeg args using effect-aware graph.
4. Startup status displays selected effect.
