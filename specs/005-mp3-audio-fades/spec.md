# Feature Specification: MP3 Audio and Video Fades

**Feature Branch**: `005-mp3-audio-fades`  
**Created**: 2026-04-10  
**Status**: Draft  
**Input**: User description: "If there are MP3 files in the image directory, include them in alphabetical order in the generated video. Also, include a fade-in and fade-out in the generated video."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Auto-Include MP3 Tracks (Priority: P1)

As an operator rendering a slideshow, I want MP3 files found in the input image directory to be included automatically so the exported video has background audio without extra steps.

**Why this priority**: Audio inclusion is a high-value usability improvement that reduces manual post-editing for the primary rendering workflow.

**Independent Test**: Place multiple MP3 files and image files in one input directory, render once, and confirm the output contains an audio track with MP3 segments ordered alphabetically.

**Required Test Coverage**:

- Unit Tests: Verify MP3 discovery filters to MP3 only, sorts filenames alphabetically, and preserves deterministic ordering across runs.
- E2E Test: Render from a fixture containing images and multiple MP3 files, then verify the resulting video contains an audio stream and the stream segment order matches alphabetical filename order.

**Acceptance Scenarios**:

1. **Given** an input directory with images and `a.mp3`, `b.mp3`, **When** the user renders a video, **Then** both MP3 files are included in alphabetical order.
2. **Given** an input directory with images and no MP3 files, **When** the user renders a video, **Then** rendering succeeds and output behavior remains unchanged from current silent-output behavior.

---

### User Story 2 - Smooth Intro and Outro (Priority: P2)

As a viewer, I want each generated video to start and end with visible fades so playback feels polished and less abrupt.

**Why this priority**: Fade-in and fade-out improves perceived quality and consistency for every generated video.

**Independent Test**: Render a video and inspect start/end playback to confirm visual fade-in at the beginning and visual fade-out at the end.

**Required Test Coverage**:

- Unit Tests: Verify final composition always includes both intro and outro visual fade directives.
- E2E Test: Render a short video and verify output starts from black into content and exits from content to black.

**Acceptance Scenarios**:

1. **Given** a valid image input, **When** the user renders a video, **Then** the output always includes a fade-in at start and fade-out at end.
2. **Given** user-defined transition duration, **When** rendering completes, **Then** fade timing follows configured duration rules.

---

### User Story 3 - Predictable Combined Media Behavior (Priority: P3)

As an operator, I want predictable behavior when both image transitions, global fades, and optional MP3 audio are present so I can trust output quality without trial-and-error.

**Why this priority**: Combined media behavior creates the highest regression risk and should remain deterministic for repeatable renders.

**Independent Test**: Render with images plus multiple MP3 files and validate one-pass output correctness for video fades, image timeline continuity, and successful completion.

**Required Test Coverage**:

- Unit Tests: Verify total output duration and fade timing remain valid when audio tracks are present.
- E2E Test: Render mixed inputs (images + MP3 files) and validate non-failing output with expected video and audio streams.

**Acceptance Scenarios**:

1. **Given** images with multiple MP3 tracks, **When** rendering runs, **Then** output completes successfully with both expected visual effects and ordered audio.
2. **Given** long MP3 input relative to video duration, **When** rendering runs, **Then** output duration remains bounded to configured slideshow duration.

---

### Edge Cases

- Input directory contains unsupported audio files (for example WAV) along with MP3 files.
- MP3 filenames differ only by case and must still be ordered deterministically.
- One or more MP3 files are unreadable or corrupt.
- Slideshow total duration is shorter than default fade duration; fade timing must remain valid.
- Input contains only MP3 files and no valid images.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detect MP3 files located in the same input directory used for image discovery.
- **FR-002**: System MUST include detected MP3 files in output audio in ascending alphabetical filename order.
- **FR-003**: System MUST preserve existing successful render behavior when no MP3 files are present.
- **FR-004**: System MUST ignore non-MP3 audio files for this feature scope.
- **FR-005**: System MUST fail with a clear input-validation classification if selected MP3 files cannot be processed.
- **FR-006**: System MUST apply a visual fade-in at the start of every generated video.
- **FR-007**: System MUST apply a visual fade-out at the end of every generated video.
- **FR-008**: System MUST ensure fade timing remains valid for all render durations, including short outputs.
- **FR-009**: System MUST keep final output duration bounded by slideshow timeline configuration even when included MP3 duration is longer.
- **FR-010**: System MUST define and maintain unit tests for all changed behaviors.
- **FR-011**: System MUST define and maintain end-to-end tests for primary user flows.

### Key Entities *(include if feature involves data)*

- **AudioAsset**: A discovered MP3 file eligible for inclusion, with attributes for path, filename, alphabetical rank, and validity state.
- **RenderTimeline**: The composed duration model defining image sequence timing, transition timing, and bounded final output duration.
- **FadeProfile**: Intro/outro fade behavior applied to final video output, including start point and duration.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In 100% of renders where two or more MP3 files are present, audio segment ordering in output matches alphabetical filename order.
- **SC-002**: In 100% of successful renders, generated video command assembly includes `fade=t=in:st=0:d=<x>` and `fade=t=out:st=<y>:d=<x>` directives with `x > 0` and `y >= 0`.
- **SC-003**: In a 20-run mixed-input validation suite (10 FHD + 10 UHD) using fixed fixtures and default settings, at least 19 runs complete successfully without manual intervention.
- **SC-004**: 100% of renders with no MP3 files continue to complete with no behavioral regression in output validity.

## Assumptions

- MP3 files are optional inputs; images remain mandatory for video generation.
- Alphabetical ordering uses normalized filename sorting in ascending order.
- Fade-in and fade-out apply to visual output stream; separate audio fade shaping is out of scope unless explicitly requested later.
- Existing transition behavior remains in place; this feature adds guaranteed intro/outro fades and optional MP3 inclusion.
