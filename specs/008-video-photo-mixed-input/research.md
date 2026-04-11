# Research: Mixed Video Photo Input

## Decision 1: Supported Video Inputs For v1

- Decision: Accept common container formats already handled well by ffmpeg in current environment: `.mp4`, `.mov`, `.mkv`, `.webm`.
- Rationale: These cover the vast majority of user camera/export workflows and avoid introducing broad codec/container edge cases in first delivery.
- Alternatives considered:
  - Accept all ffmpeg-readable formats: rejected for v1 due to unclear validation UX and wider failure surface.
  - MP4-only: rejected as too restrictive for user-provided media.

## Decision 2: Aspect-Ratio-Preserving Profile Scaling

- Decision: Scale each video clip to fit within profile bounds while preserving aspect ratio; pad to profile frame when needed. For portrait 9:16 with FHD profile, mapped clip should be 1080x1920-equivalent then framed without distortion.
- Rationale: Preserves composition and avoids stretching while ensuring output remains profile-conformant.
- Alternatives considered:
  - Stretch to exact profile dimensions: rejected because it distorts content.
  - Crop-to-fill profile frame: rejected for v1 because important visual content may be lost.

## Decision 3: High-Quality Resampling Path

- Decision: Use high-quality ffmpeg scaling directives (`flags=lanczos`) for video transforms.
- Rationale: Existing project already favors high-quality scaling for image processing; applying same quality intent to video keeps visual consistency.
- Alternatives considered:
  - Bilinear default scaling: rejected due to visible softness/artifacts in downscale cases.

## Decision 4: Frame-Rate Normalization

- Decision: Normalize all video segments to selected output fps at filter-graph level; images and transitions remain aligned to same output fps.
- Rationale: Avoids judder/mismatch across mixed clips and keeps output metadata consistent with user profile/runtime choices.
- Alternatives considered:
  - Preserve original clip fps per segment: rejected due to mixed-fps playback inconsistency.
  - Normalize only when fps differs by threshold: rejected as unnecessary complexity for v1.

## Decision 5: Short-Clip Transition Behavior

- Decision: If clip duration is shorter than configured transition duration, clamp effective transition for that boundary to avoid negative/invalid segment windows.
- Rationale: Keeps timeline stable and prevents ffmpeg graph errors in edge cases.
- Alternatives considered:
  - Hard-fail render: rejected; punitive for valid but short user clips.
  - Disable transitions globally when one clip is short: rejected because it changes entire output unexpectedly.

## Decision 6: Error Handling For Invalid Video Inputs

- Decision: Skip unreadable/unsupported video files only when at least one valid media asset remains, while emitting warnings; fail with classified input validation error when no renderable assets remain.
- Rationale: Maximizes successful runs while keeping user feedback actionable.
- Alternatives considered:
  - Fail on first invalid video: rejected for poor UX in large mixed folders.
  - Silently skip invalid files: rejected due to low observability and confusion.
