# Feature Specification: Video Processing Status

**Feature Branch**: `002-video-processing-status`
**Created**: 2026-04-09
**Status**: Draft
**Input**: User description: "Add video processing status, including the number of input files, the output format, the total processing time"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Render Completion Report (Priority: P1)

After a render job finishes, the operator wants to see a clear, labeled completion report stating exactly how many source images were processed, what video container format was produced, and how long the render took. The current one-line key=value output omits the output format entirely and presents elapsed time as a raw decimal number, making it hard to quickly confirm a job ran correctly.

**Why this priority**: Most critical — without a reliable completion report, an operator cannot tell at a glance whether input count, format, and duration matched expectations. All three data points are post-processing facts, so this story alone delivers the full value of the feature.

**Independent Test**: Can be fully tested by running a render command to completion and asserting stdout contains labeled fields for file count, output format, and elapsed time.

**Required Test Coverage**:

- Unit Tests:
  - Status formatter emits a labeled input file count field (e.g., `files=15`)
  - Status formatter derives and emits a container format label from the output path extension (e.g., `.mp4` → `MP4`, `.mov` → `MOV`)
  - Status formatter renders processing time as human-readable: decimal seconds for durations under 60 seconds; minutes and seconds for 60 seconds or more
  - Status formatter renders sub-second durations as `< 1s`
  - Status formatter handles unrecognized or missing extensions by emitting the raw uppercased extension or `UNKNOWN`

- E2E Test: Run a render against a fixture image set and assert that stdout contains a completion line with `files=N`, `format=EXT`, and a non-zero elapsed time in human-readable form.

**Acceptance Scenarios**:

1. **Given** a successful render of 12 JPEGs to `output.mp4` taking 45 seconds, **When** render completes, **Then** stdout includes `files=12`, `format=MP4`, and `elapsed=45.0s`
2. **Given** a successful render taking 90 seconds, **When** render completes, **Then** elapsed time is shown as `1m 30s`, not `90.000s`
3. **Given** a successful render producing a `.mov` file, **When** render completes, **Then** stdout shows `format=MOV`
4. **Given** a render that fails, **When** the error is returned, **Then** no completion report is printed; the existing error output path is used

---

### User Story 2 - Pre-Render Status Announcement (Priority: P2)

Before the potentially long encoding phase begins, the operator sees a single status line confirming how many images were discovered and what output format will be created, allowing them to catch misconfigured jobs before minutes of compute time are spent.

**Why this priority**: A pre-render announcement prevents silent waste. If the wrong folder is passed or the output extension is wrong, the operator sees it within one second rather than discovering it after a multi-minute render.

**Independent Test**: Can be tested independently by running a render command and asserting a status line containing input count and format appears before encoding completes.

**Required Test Coverage**:

- Unit Tests:
  - Pre-render announcement formatter produces a line containing input image count and output format label

- E2E Test: Run a render with a multi-image fixture set; assert a status line appears on stdout before the render finishes, containing input file count and output container format.

**Acceptance Scenarios**:

1. **Given** an input folder with 20 images and `--output out.mp4`, **When** the render command is invoked, **Then** a status line containing `20 files` and `MP4` is printed before encoding begins
2. **Given** an input folder with 0 valid images, **When** the render command is invoked, **Then** no pre-render announcement is printed; the job fails via the existing validation path

---

### Edge Cases

- What happens when the output file has an unrecognized extension (e.g., `.avi`)? The format field shows the raw uppercased extension (`AVI`).
- What happens when the output path has no extension? The format field shows `UNKNOWN`.
- What happens when elapsed time is under 1 second? Show `< 1s`.
- What happens when elapsed time is exactly 60 seconds? Show `1m 0s`.
- What happens when warnings were accumulated during the job? Warning count must remain visible in the completion report.
- What happens when rendering fails mid-job? No completion report is emitted; the error path handles output.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST emit a completion report on stdout after each successful render
- **FR-002**: The completion report MUST include the number of input images processed, labeled as a distinct field (e.g., `files=N`)
- **FR-003**: The completion report MUST include the output container format derived from the output file path extension, labeled as a distinct field (e.g., `format=MP4`)
- **FR-004**: The completion report MUST display total processing time in human-readable form: decimal seconds for durations under 60 seconds; minutes and seconds for durations of 60 seconds or more; `< 1s` for sub-second durations
- **FR-005**: System MUST emit a pre-render announcement line before encoding begins, stating the discovered input image count and the intended output container format
- **FR-006**: The completion report MUST NOT be emitted when rendering fails; failure output follows the existing error reporting path
- **FR-007**: Output container format MUST be derived from the output file path extension; unrecognized extensions MUST be accepted and emitted as raw uppercased labels (for example, `AVI`), and missing extensions MUST emit `UNKNOWN`; the system MUST NOT reject the render solely because of an unrecognized extension
- **FR-008**: System MUST define unit tests for all status-formatting behaviors, including edge cases (sub-second, >= 60s, missing/unknown extension)
- **FR-009**: System MUST define end-to-end tests verifying both the pre-render announcement and completion report are present with correctly labeled fields

### Key Entities

- **ProcessingStatus**: The set of information emitted during and after a render job; includes input file count, output container format label, elapsed time, and job outcome
- **OutputFormat**: A normalized container label derived from the output file path extension (e.g., `.mp4` → `MP4`); falls back to `UNKNOWN` when no extension is present

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can confirm input file count, output container format, and processing time from a single glance at terminal output, without parsing file paths or decoding raw decimal numbers
- **SC-002**: The pre-render announcement appears on stdout within 1 second of command invocation under normal local execution conditions, before any FFmpeg encoding output
- **SC-003**: All status-formatting behaviors are covered by unit tests with zero regression in the existing unit and end-to-end test suites
- **SC-004**: End-to-end tests assert specific labeled fields (`files=`, `format=`, elapsed time) in status output, providing automated regression coverage for the feature

## Assumptions

- "Output format" means the video container type (e.g., MP4, MOV, MKV) derived from the file extension of the `--output` path; codec details are out of scope for this feature
- Human-readable time formatting uses minutes and seconds for durations >= 60s and decimal seconds (one decimal place) for shorter durations
- Status output (both pre-render and completion) goes to stdout, consistent with the existing summary line; stderr remains reserved for errors
- The existing completion summary line is extended in-place with the new labeled fields; its key=value structure is preserved for backward compatibility with any existing scripts that parse it
- Live per-frame or per-second progress reporting during encoding is out of scope for this feature
