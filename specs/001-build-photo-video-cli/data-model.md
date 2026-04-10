# Data Model: Professional Photo-to-Video CLI

## RenderProfile
- Purpose: Encapsulates output target and quality-related encoding settings.
- Fields:
  - `name` (enum): `fhd` | `uhd`
  - `width` (int): `1920` or `3840`
  - `height` (int): `1080` or `2160`
  - `aspect_ratio` (string): always `16:9`
  - `video_codec_policy` (enum): `prefer_nvenc_then_cpu`
  - `quality_preset` (string): profile-specific preset identifier
- Validation:
  - Width/height must match profile enum.
  - Aspect ratio must resolve to 16:9 exactly.

## MediaAsset
- Purpose: Represents one input photo participating in render timeline.
- Fields:
  - `path` (string)
  - `order_index` (int)
  - `width` (int)
  - `height` (int)
  - `format` (string)
  - `is_valid` (bool)
  - `validation_warnings` ([]string)
- Validation:
  - File must exist and be readable.
  - Format must be in allowed image types.
  - Width/height must be positive.

## TransitionSegment
- Purpose: Defines transition between two consecutive assets.
- Fields:
  - `from_asset_index` (int)
  - `to_asset_index` (int)
  - `transition_type` (enum): `crossfade`
  - `duration_seconds` (float)
  - `offset_seconds` (float)
- Validation:
  - Segment indices must refer to adjacent assets.
  - Duration must be > 0 and less than per-image hold duration.

## RenderJob
- Purpose: Aggregate user request and computed pipeline input.
- Fields:
  - `input_assets` ([]MediaAsset)
  - `output_path` (string)
  - `profile` (RenderProfile)
  - `image_duration_seconds` (float)
  - `transition_duration_seconds` (float)
  - `overwrite` (bool)
  - `requested_encoder` (string, optional)
  - `effective_encoder` (enum): `h264_nvenc` | `libx264`
  - `warnings` ([]string)
- Validation:
  - Must contain at least 2 valid assets.
  - Output path must be writable (or overwrite allowed when file exists).
  - Transition duration must satisfy timeline constraints.

## RenderSummary
- Purpose: Final immutable result for operator feedback.
- Fields:
  - `job_id` (string)
  - `started_at` (timestamp)
  - `finished_at` (timestamp)
  - `elapsed_seconds` (float)
  - `processed_assets` (int)
  - `skipped_assets` (int)
  - `profile_name` (string)
  - `effective_resolution` (string)
  - `effective_encoder` (string)
  - `output_path` (string)
  - `status` (enum): `success` | `failed`
  - `error_message` (string, optional)

## Relationships
- A `RenderJob` has one `RenderProfile`.
- A `RenderJob` has many `MediaAsset` entries.
- A `RenderJob` has many `TransitionSegment` entries derived from adjacent `MediaAsset` entries.
- A `RenderJob` produces one `RenderSummary`.

## State Transitions
- RenderJob status lifecycle:
  - `created` -> `validated` -> `rendering` -> `completed`
  - `created` -> `validation_failed`
  - `rendering` -> `failed`
- RenderSummary status reflects terminal states: `success` or `failed`.
