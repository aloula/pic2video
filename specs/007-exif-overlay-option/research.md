# Research: EXIF Footer Overlay Option

**Phase**: 0 - Outline & Research  
**Branch**: `007-create-feature-branch`  
**Date**: 2026-04-11

## Decision 1: Overlay rendering approach

**Decision**: Use ffmpeg filtergraph text rendering with a white foreground and a semi-transparent box background to draw metadata per image segment.

**Rationale**: This fits the existing ffmpeg-centric pipeline and avoids introducing rendering libraries or image pre-processing passes.

**Alternatives considered**:
- Pre-render text into image assets before ffmpeg: rejected due to extra IO and pipeline complexity.
- Burn subtitles via external subtitle files: rejected because dynamic per-image metadata mapping is harder to keep deterministic.

## Decision 2: EXIF extraction source and fallback policy

**Decision**: Extend the existing ffprobe-based metadata extraction flow to collect required EXIF values and normalize missing values to `Unknown`.

**Rationale**: The repository already extracts EXIF creation time through ffprobe, so extending this path keeps dependencies small and behavior consistent.

**Alternatives considered**:
- Add exiftool dependency: rejected to avoid heavy external dependency and portability overhead.
- Fail rendering when one field is missing: rejected because requirements prioritize successful renders with graceful fallback.

## Decision 3: Output string format and date normalization

**Decision**: Build a single normalized overlay line in this exact sequence: Camera Model - Focal Distance - Speed (1/XXXXs) - Aperture (f/X) - ISO - Captured Date (DD/MM/YYYY), with strict date formatting and `Unknown` placeholder for unavailable values.

**Rationale**: A deterministic formatter enables straightforward unit tests and guarantees requirement compliance.

**Alternatives considered**:
- Locale-dependent date formatting: rejected because output would vary across environments.
- Omitting missing fields: rejected because it breaks required ordering.

## Decision 4: Placement and style constraints

**Decision**: Position overlay in the footer at fixed vertical offset (10 pixels above bottom) for both FHD and UHD, enforce white text, and configure box opacity to exceed 50% transparency.

**Rationale**: Directly fulfills explicit visual constraints while keeping behavior profile-invariant.

**Alternatives considered**:
- Scale offset by resolution: rejected because requirement mandates 10-pixel offset for both profiles.
- Fully opaque background: rejected due to transparency requirement.

## Decision 5: Font-size validation contract

**Decision**: Add CLI validation for integer font size range 36 through 60 inclusive when overlay is enabled.

**Rationale**: Early validation produces clear operator feedback and prevents invalid ffmpeg argument generation.

**Alternatives considered**:
- Silently clamp out-of-range values: rejected because hidden mutation is less predictable for users.
- Accept any positive size: rejected due to explicit bounded requirement.

## Decision 6: Test strategy mapping

**Decision**: Add unit coverage for formatter, range validation, and ffmpeg args assembly; add E2E coverage for overlay-enabled FHD/UHD render outputs.

**Rationale**: Satisfies constitution requirements for unit and E2E coverage on each user-facing behavior.

**Alternatives considered**:
- Unit-only verification: rejected because full integration behavior with ffmpeg filters must be validated end-to-end.
