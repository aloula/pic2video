---
description: "Task list for 005-mp3-audio-fades"
---

# Tasks: MP3 Audio and Video Fades

**Input**: Design documents from `/specs/005-mp3-audio-fades/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: Unit and E2E test tasks are mandatory and must be executed before implementation tasks are considered complete.

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm branch baseline and ensure fixtures for audio-aware tests exist.

- [x] T001 Verify baseline regression health with `go test ./...` from project root `go.mod`
- [x] T002 Verify feature context points to `specs/005-mp3-audio-fades` in `.specify/feature.json`
- [x] T003 Prepare mixed media fixture documentation in `tests/fixtures/README.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build shared timeline/audio plumbing required by all user stories.

**CRITICAL**: No user story implementation starts before this phase is complete.

- [x] T004 Add optional audio assets collection to `RenderJob` in `internal/app/renderjob/model.go`
- [x] T005 Add MP3 asset propagation through build options in `internal/app/renderjob/builder.go`
- [x] T006 Add MP3 discovery helper and deterministic filename sort in `internal/infra/fsio/filesystem.go`
- [x] T007 Add audio-aware ffmpeg args entry point in `internal/infra/ffmpeg/command_builder.go`
- [x] T008 Wire service layer to pass optional audio assets to ffmpeg args builder in `internal/app/renderjob/service.go`

**Checkpoint**: Shared audio/timeline plumbing is complete and user stories can proceed.

---

## Phase 3: User Story 1 - Auto-Include MP3 Tracks (Priority: P1) MVP

**Goal**: Operator gets automatic MP3 inclusion in alphabetical order when MP3 files exist in input directory.

**Independent Test**: Render with image + MP3 input folder and verify audio stream exists with deterministic alphabetical source ordering.

### Tests for User Story 1 (Required)

- [x] T009 [P] [US1] Add unit test for MP3 discovery filtering and sorting in `tests/unit/cli_validation_test.go`
- [x] T010 [P] [US1] Add unit test for no-MP3 fallback behavior in `tests/unit/cli_validation_test.go`
- [x] T011 [P] [US1] Add E2E test for image+MP3 successful render with audio stream in `tests/e2e/render_fhd_test.go`
- [x] T012 [P] [US1] Add E2E test for image-only render unchanged behavior in `tests/e2e/smoke_test.go`
- [x] T041 [P] [US1] Add E2E test for images + mp3 + unsupported audio (e.g., wav) to assert non-MP3 audio is ignored in `tests/e2e/render_fhd_test.go`

### Implementation for User Story 1

- [x] T013 [US1] Extend input discovery flow to collect sorted MP3 assets in `internal/app/cli/render_command.go`
- [x] T014 [US1] Validate discovered MP3 readability and classify invalid audio input errors in `internal/app/cli/render_command.go`
- [x] T015 [US1] Populate ordered MP3 assets in `BuildJob` call in `internal/app/cli/render_command.go`
- [x] T016 [US1] Implement ffmpeg audio input mapping for sorted MP3 files in `internal/infra/ffmpeg/command_builder.go`
- [x] T017 [US1] Bound audio output to slideshow duration in ffmpeg args generation in `internal/infra/ffmpeg/command_builder.go`
- [x] T018 [US1] Run `go test ./tests/unit -run 'CLI|Validation'` to validate MP3 discovery behavior in `tests/unit/cli_validation_test.go`
- [x] T019 [US1] Run `go test ./tests/e2e -run 'RenderFHD|Smoke'` to validate MP3 inclusion and fallback behavior in `tests/e2e/render_fhd_test.go`

**Checkpoint**: User Story 1 is fully functional and independently testable.

---

## Phase 4: User Story 2 - Smooth Intro and Outro (Priority: P2)

**Goal**: Every generated video includes deterministic visual fade-in and fade-out timing.

**Independent Test**: Generated ffmpeg graph always includes fade-in at start and fade-out near end with valid timing.

### Tests for User Story 2 (Required)

- [x] T020 [P] [US2] Add unit test for mandatory fade-in and fade-out directives in `tests/unit/ffmpeg_args_test.go`
- [x] T021 [P] [US2] Add unit test for short-duration fade clamping behavior in `tests/unit/ffmpeg_args_test.go`
- [x] T022 [P] [US2] Add E2E test asserting fade-enabled output generation in `tests/e2e/render_fhd_test.go`

### Implementation for User Story 2

- [x] T023 [US2] Refine global fade profile helper for valid start/end timing in `internal/infra/ffmpeg/command_builder.go`
- [x] T024 [US2] Ensure fade profile is applied in audio and non-audio render paths in `internal/infra/ffmpeg/command_builder.go`
- [x] T025 [US2] Run `go test ./tests/unit -run BuildRenderCommandArgs` to validate fade timing logic in `tests/unit/ffmpeg_args_test.go`
- [x] T026 [US2] Run `go test ./tests/e2e -run RenderFHD` to validate fade behavior in end-to-end rendering in `tests/e2e/render_fhd_test.go`

**Checkpoint**: User Story 2 is fully functional and independently testable.

---

## Phase 5: User Story 3 - Predictable Combined Media Behavior (Priority: P3)

**Goal**: Mixed image+MP3 renders remain deterministic with bounded duration and valid fades.

**Independent Test**: Rendering mixed inputs repeatedly produces stable stream layout and bounded output duration.

### Tests for User Story 3 (Required)

- [x] T027 [P] [US3] Add unit test for bounded timeline when MP3 duration exceeds slideshow duration in `tests/unit/ffmpeg_args_test.go`
- [x] T028 [P] [US3] Add unit test for deterministic MP3 order with case-variance filenames in `tests/unit/cli_validation_test.go`
- [x] T029 [P] [US3] Add E2E mixed-media stability test in `tests/e2e/render_fhd_test.go`

### Implementation for User Story 3

- [x] T030 [US3] Normalize filename ordering rules for deterministic MP3 sequencing in `internal/infra/fsio/filesystem.go`
- [x] T031 [US3] Ensure mixed-media command assembly preserves deterministic input index mapping in `internal/infra/ffmpeg/command_builder.go`
- [x] T032 [US3] Add startup reporting for discovered MP3 count and ordering mode in `internal/app/cli/summary.go`
- [x] T033 [US3] Populate startup reporting fields from render command flow in `internal/app/cli/render_command.go`
- [x] T034 [US3] Run `go test ./tests/unit -run 'BuildRenderCommandArgs|CLI|Validation'` to verify deterministic mixed-media behavior in `tests/unit/ffmpeg_args_test.go`
- [x] T035 [US3] Run `go test ./tests/e2e -run RenderFHD` to verify mixed-media bounded output behavior in `tests/e2e/render_fhd_test.go`

**Checkpoint**: User Story 3 is fully functional and independently testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final documentation alignment and full-project validation.

- [x] T036 [P] Update root usage and behavior documentation for MP3 auto-inclusion and fades in `README.md`
- [x] T037 [P] Align operator workflow examples in `specs/005-mp3-audio-fades/quickstart.md`
- [x] T038 [P] Align behavior guarantees in `specs/005-mp3-audio-fades/contracts/cli-contract.md`
- [x] T039 Run full regression suite `go test ./...` from project root `go.mod`
- [x] T040 Run `make build-all` after successful implementation and confirm binaries are generated in `bin/`

---

## Dependencies & Execution Order

### Phase Dependencies

- Phase 1: no dependencies
- Phase 2: depends on Phase 1 and blocks all user stories
- Phase 3 (US1): depends on Phase 2
- Phase 4 (US2): depends on Phase 2
- Phase 5 (US3): depends on Phase 2 and integrates US1 + US2 behavior
- Phase 6: depends on completion of US1, US2, and US3

### User Story Dependencies

- US1 (P1): no dependency on other stories after foundational work
- US2 (P2): no dependency on US1 for fade logic after foundational work
- US3 (P3): depends on both US1 audio flow and US2 fade/timeline behavior

### Within Each User Story

- Tests first and failing
- Implementation second
- Story-specific tests green before moving on

---

## Parallel Opportunities

- [P] tasks in each story can run in parallel when they touch different files.
- US1 tests T009-T012 and T041 can run in parallel.
- US2 tests T020-T022 can run in parallel.
- US3 tests T027-T029 can run in parallel.
- Polish docs tasks T036-T038 can run in parallel.

## Parallel Example: User Story 1

- Run T009 and T010 together (unit test coverage for discovery and fallback)
- Run T011 and T012 together (e2e tests for mixed-media and image-only journeys)

---

## Implementation Strategy

### MVP First (US1)

1. Complete Phase 1 and Phase 2
2. Complete Phase 3 (US1)
3. Validate US1 independently

### Incremental Delivery

1. Add US2 guaranteed fade behavior
2. Add US3 deterministic combined-media behavior
3. Finish with Phase 6 docs and full regression/build validation
