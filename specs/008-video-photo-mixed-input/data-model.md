# Data Model: Mixed Video Photo Input

## Entity: MixedMediaAsset

- Purpose: Unified representation of image or video input used by discovery, ordering, and rendering.
- Fields:
  - id: string (stable per-run identifier)
  - path: string (source filesystem path)
  - media_type: enum (`image`, `video`)
  - width: int
  - height: int
  - aspect_ratio: string (normalized, for example `16:9`, `9:16`)
  - duration_seconds: float (0 for still images before timeline expansion)
  - source_fps: float (0 when not applicable/unknown)
  - order_index: int
  - validity_state: enum (`valid`, `invalid`, `skipped`)
  - validation_message: string

## Entity: VideoTransformPolicy

- Purpose: Defines how a source video maps to output profile and fps without distortion.
- Fields:
  - profile_name: enum (`fhd`, `uhd`)
  - target_width: int
  - target_height: int
  - scale_mode: enum (`fit_preserve_aspect`)
  - pad_mode: enum (`center_pad`)
  - scaling_quality: enum (`high`)
  - target_fps: float
  - transition_clamp_enabled: bool

## Entity: TimelineSegment

- Purpose: Final segment in render timeline after normalization.
- Fields:
  - segment_id: string
  - media_asset_id: string
  - media_type: enum (`image`, `video`)
  - start_seconds: float
  - end_seconds: float
  - effective_duration_seconds: float
  - transition_in_seconds: float
  - transition_out_seconds: float
  - effective_width: int
  - effective_height: int
  - effective_fps: float

## Relationships

- MixedMediaAsset (1) -> (0..N) TimelineSegment
- VideoTransformPolicy (1) -> (0..N) TimelineSegment for video media_type

## Validation Rules

- media_type MUST be determined from extension + probe metadata.
- width/height MUST be > 0 for valid video assets.
- target_fps MUST be > 0.
- effective_duration_seconds MUST be > 0 for all TimelineSegment records.
- transition values MUST be clamped so `start_seconds < end_seconds` remains true.
