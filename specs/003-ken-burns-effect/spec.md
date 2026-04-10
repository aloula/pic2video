# Feature Specification: Ken Burns Effect Option

**Feature Branch**: `003-ken-burns-effect`  
**Created**: 2026-04-09  
**Status**: Draft  
**Input**: User description: "Add Ken Burns effect option, so the user can select between static images (default) and Ken Burns effects (low, medium or high). Observe the output resolution (FHD or UHD) to apply the effects accordingly. Give always preference to use smooth and high quality effects, even if the processing time will be larger"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select Motion Style (Priority: P1)

As an operator, I can choose how slideshow images move by setting one flag: static (default) or Ken Burns with low, medium, or high intensity.

**Why this priority**: Motion style control is the main user-facing value of this feature and must be explicit and deterministic.

**Independent Test**: Can be tested by invoking render with each option and verifying validation, startup status, and filter generation behavior.

**Required Test Coverage**:

- Unit Tests: motion option validation for allowed and invalid values.
- Unit Tests: effect filter generation for static and Ken Burns levels.
- E2E Test: startup status includes selected motion style option.

**Acceptance Scenarios**:

1. **Given** no motion flag is provided, **When** render starts, **Then** mode is `static`.
2. **Given** `--image-effect kenburns-medium`, **When** render starts, **Then** selected effect is shown and applied.
3. **Given** an invalid effect value, **When** render is invoked, **Then** command fails with invalid arguments error.

---

### User Story 2 - Resolution-Aware Ken Burns (Priority: P2)

As an operator, Ken Burns behavior adapts to output resolution so motion quality remains smooth and visually consistent for FHD and UHD outputs.

**Why this priority**: The request explicitly requires output-resolution-aware behavior.

**Independent Test**: Can be tested by generating filters for FHD and UHD and asserting they differ and target the correct output dimensions.

**Required Test Coverage**:

- Unit Tests: resolution-aware differences in generated motion filters.
- E2E Test: render succeeds in both FHD and UHD with Ken Burns selected.

**Acceptance Scenarios**:

1. **Given** FHD output, **When** Ken Burns is selected, **Then** generated motion filter targets FHD output and smooth movement.
2. **Given** UHD output, **When** Ken Burns is selected, **Then** generated motion filter targets UHD output and smooth movement.

---

### User Story 3 - Quality-First Motion (Priority: P3)

As an operator, the Ken Burns effect prioritizes smoothness and image quality over speed, producing cinematic motion even if render time increases.

**Why this priority**: This is an explicit quality policy from the request.

**Independent Test**: Can be tested by asserting high-quality interpolation and high frame-rate motion settings in generated filters.

**Required Test Coverage**:

- Unit Tests: generated Ken Burns filters include high-quality scaling and smooth-motion settings.

**Acceptance Scenarios**:

1. **Given** any Ken Burns mode, **When** filter graph is built, **Then** it includes high-quality scaling (`lanczos`) and smooth motion settings (`zoompan` with `fps=60`).

### Edge Cases

- What happens when effect mode is omitted? Default is `static`.
- What happens when effect mode is invalid? Command fails before rendering.
- What happens with very short image durations? Motion filter still emits a valid minimum frame count.
- What happens when effect is static? Existing framing and render behavior are preserved.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a new render option `--image-effect`.
- **FR-002**: `--image-effect` MUST support `static`, `kenburns-low`, `kenburns-medium`, and `kenburns-high`.
- **FR-003**: When `--image-effect` is omitted, system MUST default to `static`.
- **FR-004**: System MUST reject unsupported values for `--image-effect` with invalid-arguments classification.
- **FR-005**: For non-static modes, system MUST apply Ken Burns motion in the render filter pipeline.
- **FR-006**: Ken Burns behavior MUST account for selected output resolution (FHD/UHD).
- **FR-007**: For `kenburns-low`, `kenburns-medium`, and `kenburns-high`, generated motion filters MUST include high-quality scaling (`lanczos`) and smooth-motion interpolation at `fps=60`, even if processing time increases.
- **FR-008**: Startup status output MUST include the selected image effect.
- **FR-009**: System MUST define and maintain unit tests for all changed behaviors.
- **FR-010**: System MUST define and maintain end-to-end tests for primary user flows.

### Key Entities

- **ImageEffectMode**: User-selected motion style (`static`, `kenburns-low`, `kenburns-medium`, `kenburns-high`).
- **MotionFilterProfile**: Resolution-aware derived motion parameters used to build ffmpeg filter graph.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can select one of four motion modes via CLI without ambiguity.
- **SC-002**: Invalid image-effect values are rejected deterministically before render starts.
- **SC-003**: Unit tests validate static behavior, Ken Burns behavior, and resolution-aware filter differences.
- **SC-004**: E2E tests confirm startup output includes selected effect and rendering remains successful.
- **SC-005**: Unit tests verify quality directives for all Ken Burns modes by asserting `lanczos`, `zoompan`, and `fps=60` are present in generated motion filters.

## Assumptions

- Effect intensity is represented by preset low/medium/high parameter profiles.
- Resolution-aware behavior is based on selected output profile dimensions.
- Motion quality preference allows increased CPU/GPU time when effect is enabled.
- Static mode preserves legacy behavior for users not opting into motion effects.
