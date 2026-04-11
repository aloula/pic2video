# Feature Specification: Mixed Video Photo Input

**Feature Branch**: `008-video-photo-mixed-input`  
**Created**: 2026-04-11  
**Status**: Draft  
**Input**: User description: "for this new created branch \"add-video-input\", lets give the option to render videos together with the pictures. The video aspect ratio must be preserved, but the resolution must match the output profile, without any loose of quality. Example: if an input video is 2160 x 3840 (9:16) and the profile select is FHD, the video must be resized in high quality to 1080 X 1920 (9:16) and the framerate must match the user selection for the video output"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Render Videos With Photos (Priority: P1)

As a creator, I want to include both image and video files in one render command so I can produce a single output video from mixed media without manual pre-processing.

**Why this priority**: Mixed-input rendering is the primary business value and unblocks real-world workflows where users capture both stills and clips.

**Independent Test**: Provide a folder containing JPG and MP4 files, run one render command, and verify the output includes all inputs in expected order and duration.

**Required Test Coverage**:

- Unit Tests: Mixed media input discovery, ordering behavior with mixed types, and timeline segment construction for image and video assets.
- E2E Test: Render a mixed folder and assert output creation success plus expected segment count/order.

**Acceptance Scenarios**:

1. **Given** a folder with supported image and video files, **When** the user runs render, **Then** the system includes both media types in the final output.
2. **Given** mixed media inputs and explicit ordering mode, **When** rendering starts, **Then** the output sequence follows the requested order without dropping valid files.

---

### User Story 2 - Preserve Video Aspect Ratio At Profile Quality (Priority: P2)

As a creator, I want input videos to keep their original aspect ratio while being scaled to the selected output profile quality target so the visual composition is not distorted.

**Why this priority**: Distorted video ruins output quality and is a common failure case in mixed-orientation workflows.

**Independent Test**: Use a 2160x3840 portrait source with FHD profile and verify transformed dimensions resolve to 1080x1920 equivalent profile scale while preserving 9:16 ratio.

**Required Test Coverage**:

- Unit Tests: Dimension transform policy for portrait and landscape videos across FHD and UHD profiles.
- E2E Test: Mixed render with portrait and landscape clips verifies no stretching/cropping distortion and high-quality scaling path is used.

**Acceptance Scenarios**:

1. **Given** a portrait input video (9:16), **When** rendering with FHD profile, **Then** the transformed video dimensions preserve 9:16 and map to 1080x1920 profile-equivalent scale.
2. **Given** a landscape input video (16:9), **When** rendering with UHD profile, **Then** the transformed video dimensions preserve 16:9 and map to 3840x2160 profile-equivalent scale.

---

### User Story 3 - Match Output Frame Rate For Video Segments (Priority: P3)

As a creator, I want all video segments to align to the selected output frame rate so playback remains smooth and consistent across the final timeline.

**Why this priority**: Frame-rate mismatches can create judder and inconsistent pacing when mixing different source clips.

**Independent Test**: Render mixed inputs with a selected output frame rate and verify all video segments in output conform to that frame rate target.

**Required Test Coverage**:

- Unit Tests: Frame-rate selection and propagation logic for video assets and mixed timelines.
- E2E Test: Render mixed inputs with non-default frame-rate choice and verify output metadata reports selected frame rate.

**Acceptance Scenarios**:

1. **Given** source videos with varying frame rates, **When** a target output frame rate is selected, **Then** all rendered video segments conform to the selected output frame rate.
2. **Given** mixed image/video rendering, **When** output frame rate is applied, **Then** transitions and segment boundaries remain visually stable.

---

### Edge Cases

- What happens when a very short video clip is shorter than the configured transition duration?
- How does the system handle unsupported or corrupted video files while valid image/video inputs are present?
- How does rendering behave when portrait and landscape videos are mixed in the same run?
- What happens when a video has variable frame rate metadata or missing frame rate metadata?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detect and accept supported video files in input directories alongside supported image files.
- **FR-002**: System MUST build a unified render timeline containing both image and video assets in the selected ordering mode.
- **FR-003**: System MUST preserve source video aspect ratio during transformation and MUST NOT stretch video frames.
- **FR-004**: System MUST scale each input video to the selected profile quality target while preserving orientation and aspect ratio (for example, portrait 9:16 source with FHD target resolves to 1080x1920 equivalent scale).
- **FR-005**: System MUST apply high-quality video resizing for profile scaling and MUST use high-quality resampling (`lanczos`) for transformed video segments.
- **FR-006**: System MUST align video segment playback to the user-selected output frame rate selected through `--fps`; when omitted, system MUST default to profile default fps.
- **FR-007**: System MUST ensure mixed image/video outputs maintain stable transitions and timing after video scaling and frame-rate normalization.
- **FR-008**: System MUST surface clear validation errors for unsupported video formats or unreadable media files.
- **FR-009**: System MUST define and maintain unit tests for all changed behaviors.
- **FR-010**: System MUST define and maintain end-to-end tests for primary user flows.

### Key Entities *(include if feature involves data)*

- **MixedMediaAsset**: Represents an input item in the render sequence, including type (image or video), source path, dimensions, orientation, duration, and source frame-rate metadata.
- **VideoTransformPolicy**: Represents profile-target transform rules for resizing and aspect-ratio preservation, including target dimensions, scaling mode, and quality intent.
- **TimelineSegment**: Represents a normalized output timeline segment with start time, duration, media type, transform state, and frame-rate alignment state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of valid mixed-media test runs complete successfully and include all valid input assets in the final output.
- **SC-002**: 100% of tested portrait and landscape video inputs preserve aspect ratio after rendering (no geometric distortion).
- **SC-003**: Output metadata reports the user-selected frame rate for mixed-media renders in at least 95% of automated validation runs.
- **SC-004**: In acceptance review, at least 90% of sampled scaled video segments are rated as visually equivalent to source quality expectations for the selected output profile.

## Assumptions

- Existing profile names (FHD/UHD) remain the way users choose quality level.
- Existing ordering modes continue to apply, and videos participate in ordering similarly to images.
- Existing transition behavior remains enabled unless explicitly disabled by user options.
- Existing external media tooling already available in the project can inspect video metadata needed for transform and frame-rate normalization.
