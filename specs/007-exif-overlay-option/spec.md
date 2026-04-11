# Feature Specification: EXIF Footer Overlay Option

**Feature Branch**: `007-create-feature-branch`  
**Created**: 2026-04-11  
**Status**: Draft  
**Input**: User description: "Add an EXIF overlay option. The overlay must be in footer (10 pixels above in both resolution FHD and UHD). The EXIF information must be in this format: Camera Model - Focal Distance - Speed (1/XXXXs) - Aperture (f/X) - ISO - Captured Date (DD/MM/YYYY). The font must be White, size adjustable between 36 and 60, over 50% transparent background"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enable EXIF footer overlay (Priority: P1)

As a video creator, I want an option to show EXIF metadata on the output video so viewers can see camera capture details while the slideshow plays.

**Why this priority**: The feature's core value is the ability to include metadata directly in the video output.

**Independent Test**: Can be fully tested by enabling the option on images with EXIF data and verifying the footer text appears in the requested field order and date format.

**Required Test Coverage**:

- Unit Tests: Validate EXIF field extraction and final display string formatting follows the required order and formatting.
- E2E Test: Render FHD and UHD videos with overlay enabled and verify overlay appears on output frames with required content pattern.

**Acceptance Scenarios**:

1. **Given** an image with complete EXIF metadata, **When** rendering with EXIF overlay enabled, **Then** each frame includes metadata in this order: Camera Model - Focal Distance - Speed (1/XXXXs) - Aperture (f/X) - ISO - Captured Date (DD/MM/YYYY).
2. **Given** EXIF overlay is disabled, **When** rendering a video, **Then** no EXIF footer text is shown.

---

### User Story 2 - Preserve placement and style across profiles (Priority: P2)

As a creator producing multiple resolutions, I want the overlay to keep consistent placement and readability in both FHD and UHD outputs.

**Why this priority**: Placement and legibility consistency ensure the feature is usable for all supported output profiles.

**Independent Test**: Can be fully tested by rendering one FHD and one UHD output with the overlay and checking position and visual styling constraints.

**Required Test Coverage**:

- Unit Tests: Validate profile-based placement rule uses the same 10-pixel bottom offset for both FHD and UHD outputs.
- E2E Test: Render both profiles and confirm overlay is placed in the footer with baseline 10 pixels above the bottom edge.

**Acceptance Scenarios**:

1. **Given** a render request for FHD or UHD, **When** EXIF overlay is enabled, **Then** the overlay is displayed in the footer at 10 pixels above the bottom edge.
2. **Given** EXIF overlay text is displayed, **When** viewing the output, **Then** text color is white with a background opacity greater than 50% transparency.

---

### User Story 3 - Adjust font size for visual preference (Priority: P3)

As a creator, I want to control overlay font size so the EXIF footer can be tuned for visibility and style.

**Why this priority**: Size control improves usability without changing the required metadata content.

**Independent Test**: Can be fully tested by setting minimum and maximum supported font values and verifying they are applied in output.

**Required Test Coverage**:

- Unit Tests: Validate accepted font size range is 36 through 60 inclusive; reject out-of-range values.
- E2E Test: Render with font size 36 and 60 and confirm output visually reflects selected size while preserving format and placement rules.

**Acceptance Scenarios**:

1. **Given** font size input within 36-60, **When** rendering with overlay enabled, **Then** output uses the requested font size.
2. **Given** font size input outside 36-60, **When** validating render options, **Then** the command fails with a clear validation error.

### Edge Cases

- Source image has partial EXIF metadata: missing values are displayed as `Unknown` while preserving field order and separators.
- Captured date exists but is not parseable: date output is displayed as `Unknown` rather than failing the full render.
- Camera model text is too long for frame width: overlay remains within visible frame bounds and does not overlap outside the video area.
- Mixed input set where only some images include EXIF: rendering succeeds and each image frame uses its own available metadata values.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an option to enable EXIF metadata overlay on rendered videos.
- **FR-002**: System MUST place the overlay in the footer, with text baseline positioned 10 pixels above the bottom edge for both FHD and UHD outputs.
- **FR-003**: System MUST format overlay content exactly as: `Camera Model - Focal Distance - Speed (1/XXXXs) - Aperture (f/X) - ISO - Captured Date (DD/MM/YYYY)`.
- **FR-004**: System MUST render overlay text in white.
- **FR-005**: System MUST render overlay background with alpha <= 0.5 (at least 50% transparent).
- **FR-006**: Users MUST be able to set overlay font size between 36 and 60 inclusive.
- **FR-007**: System MUST reject font size values below 36 or above 60 with a validation error before rendering starts.
- **FR-008**: System MUST produce a value for each formatted EXIF field; unavailable field values MUST be shown as `Unknown`.
- **FR-009**: System MUST format captured date as `DD/MM/YYYY` when valid source date exists.
- **FR-010**: System MUST define and maintain unit tests for all changed behaviors.
- **FR-011**: System MUST define and maintain end-to-end tests for primary user flows.

### Key Entities *(include if feature involves data)*

- **ExifOverlayOptions**: User-configurable overlay settings including enabled/disabled state and requested font size.
- **ExifFieldSet**: Normalized metadata values for one source image (camera model, focal distance, speed, aperture, ISO, captured date).
- **OverlayRenderLine**: Final display string assembled in required order with separators and formatted date.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of renders with overlay enabled show footer text at 10-pixel bottom offset in both FHD and UHD outputs.
- **SC-002**: 100% of sampled output frames from enabled renders match the required metadata field order and separator format.
- **SC-003**: 100% of accepted font size inputs in the 36-60 range are applied in output renders.
- **SC-004**: 100% of out-of-range font size requests are rejected before rendering with an actionable validation message.

## Assumptions

- EXIF overlay is optional and disabled by default unless explicitly requested.
- The feature applies only to FHD and UHD output profiles currently supported by the tool.
- Missing EXIF values are common in mixed photo sets and should not block successful video generation.
- User-provided font size is an integer value.
