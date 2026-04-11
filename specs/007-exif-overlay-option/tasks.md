# Tasks: EXIF Footer Overlay Option

**Input**: Design documents from `/specs/007-exif-overlay-option/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: Unit and E2E test tasks are MANDATORY for every user story and MUST be executed before implementation tasks are considered complete.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare shared files and test fixtures for EXIF overlay work.

- [X] T001 Add EXIF overlay CLI usage examples in specs/007-exif-overlay-option/quickstart.md
- [X] T002 Add EXIF overlay command examples and constraints to README.md
- [X] T003 [P] Add EXIF sample fixture notes in tests/fixtures/local-images/README.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core structures and helpers required by all user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Extend render job options/state with EXIF overlay fields in internal/app/renderjob/model.go
- [X] T005 Wire EXIF overlay options through job builder in internal/app/renderjob/builder.go
- [X] T006 [P] Add EXIF metadata model/normalizer helpers in internal/infra/fsio/exif.go
- [X] T007 [P] Add overlay text escaping helper for ffmpeg drawtext safety in internal/infra/ffmpeg/command_builder.go
- [X] T008 Add startup summary fields for EXIF overlay mode/font size in internal/app/cli/summary.go

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Enable EXIF footer overlay (Priority: P1) 🎯 MVP

**Goal**: Provide optional EXIF footer overlay with exact required field order and formatting.

**Independent Test**: Enable overlay on EXIF-rich images and verify required output string format appears in both FHD and UHD rendered video; disable overlay and verify no footer text.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T009 [P] [US1] Add unit tests for EXIF extraction fallback and Unknown normalization in tests/unit/order_mode_test.go
- [X] T010 [P] [US1] Add unit tests for overlay line formatter field order and date format in tests/unit/status_format_test.go
- [X] T011 [P] [US1] Add unit tests for ffmpeg args containing drawtext when overlay enabled and omitted when disabled in tests/unit/ffmpeg_args_test.go
- [X] T012 [US1] Add E2E render test for overlay-enabled flow and disabled baseline in tests/e2e/render_fhd_test.go
- [X] T012A [US1] Add E2E render test for overlay-enabled UHD content format in tests/e2e/render_uhd_test.go

### Implementation for User Story 1

- [X] T013 [US1] Add CLI flags --exif-overlay and --exif-font-size and pass values to build options in internal/app/cli/render_command.go
- [X] T014 [US1] Build per-image EXIF display line formatter with exact tokens in internal/app/renderjob/service.go: Camera Model - Focal Distance - Speed (1/XXXXs) - Aperture (f/X) - ISO - Captured Date (DD/MM/YYYY), using Unknown fallback for unavailable fields
- [X] T015 [US1] Extend ffmpeg arg builder to inject per-segment drawtext filters from overlay lines in internal/infra/ffmpeg/command_builder.go
- [X] T016 [US1] Integrate overlay metadata extraction and formatter invocation into render flow in internal/app/renderjob/service.go
- [X] T017 [US1] Expose overlay status values in CLI summary output in internal/app/renderjob/summary.go

**Checkpoint**: User Story 1 should be fully functional and independently testable.

---

## Phase 4: User Story 2 - Preserve placement and style across profiles (Priority: P2)

**Goal**: Keep footer placement and style consistent for FHD and UHD outputs.

**Independent Test**: Render FHD and UHD videos with overlay enabled and verify footer baseline is 10px from bottom, text is white, and background is >50% transparent.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T018 [P] [US2] Add unit tests for profile-invariant footer offset and style tokens in ffmpeg args in tests/unit/ffmpeg_args_test.go
- [X] T019 [P] [US2] Add unit tests for overlay placement rule consistency in tests/unit/profile_test.go
- [X] T020 [P] [US2] Add E2E test for FHD footer placement/style assertions in tests/e2e/render_fhd_test.go
- [X] T021 [P] [US2] Add E2E test for UHD footer placement/style assertions in tests/e2e/render_uhd_test.go

### Implementation for User Story 2

- [X] T022 [US2] Implement fixed footer y-position (10px from bottom) for FHD and UHD in internal/infra/ffmpeg/command_builder.go
- [X] T023 [US2] Enforce white font color and semi-transparent box styling (>50% transparent) in internal/infra/ffmpeg/command_builder.go
- [X] T024 [US2] Ensure overlay style/placement metadata is reflected in user-facing summary where applicable in internal/app/cli/summary.go

**Checkpoint**: User Stories 1 and 2 should both work independently.

---

## Phase 5: User Story 3 - Adjust font size for visual preference (Priority: P3)

**Goal**: Allow user-selected font size between 36 and 60 with explicit validation errors outside range.

**Independent Test**: Render with sizes 36 and 60 successfully; reject 35 and 61 before ffmpeg execution with invalid-arguments classification.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T025 [P] [US3] Add unit tests for font-size validation boundaries and error messaging in tests/unit/cli_validation_test.go
- [X] T026 [P] [US3] Add unit tests for font-size propagation into ffmpeg args in tests/unit/ffmpeg_args_test.go
- [X] T027 [P] [US3] Add E2E validation-failure test for out-of-range font size in tests/e2e/invalid_input_test.go
- [X] T028 [P] [US3] Add E2E render tests for min/max valid font sizes in tests/e2e/render_fhd_test.go

### Implementation for User Story 3

- [X] T029 [US3] Validate --exif-font-size range (36-60) in CLI argument validation in internal/app/cli/render_command.go
- [X] T030 [US3] Add font-size field to start options formatting output in internal/app/cli/summary.go
- [X] T031 [US3] Apply user-selected font size to drawtext configuration in internal/infra/ffmpeg/command_builder.go
- [X] T032 [US3] Ensure classified invalid-argument errors for font-size violations in internal/app/renderjob/errors.go

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Hardening and end-to-end validation across all user stories.

- [X] T033 [P] Refactor shared EXIF string/date formatting helpers for readability and duplication control in internal/infra/fsio/exif.go
- [X] T034 [P] Add additional mixed-input regression assertions (missing EXIF fields) in tests/e2e/render_mixed_aspect_test.go
- [X] T035 Run full unit suite for regression safety with `go test ./tests/unit/...` from repository root
- [X] T036 Run full E2E suite for regression safety with `go test ./tests/e2e/...` from repository root
- [X] T037 Run full project tests with `go test ./...` from repository root
- [X] T038 Run `make build-all` and record successful result in specs/007-exif-overlay-option/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies, starts immediately.
- **Foundational (Phase 2)**: Depends on Setup completion and blocks all user stories.
- **User Story Phases (Phase 3-5)**: Depend on Foundational completion.
- **Polish (Phase 6)**: Depends on all targeted user stories being complete.

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2 and is MVP.
- **US2 (P2)**: Can start after Phase 2; depends functionally on US1 overlay path existing.
- **US3 (P3)**: Can start after Phase 2; integrates with US1 overlay path and validation pipeline.

### Within Each User Story

- Write unit and E2E tests first and verify they fail.
- Implement code changes after tests exist.
- Re-run story tests until passing before moving to next story.

### Parallel Opportunities

- Tasks marked [P] in Setup and Foundational can run concurrently.
- In each user story, [P] unit and E2E tasks can be split among team members.
- US2 and US3 can run in parallel after US1 base overlay implementation is merged.

---

## Parallel Example: User Story 1

```bash
# Parallel test authoring for US1
Task T009: tests/unit/order_mode_test.go
Task T010: tests/unit/status_format_test.go
Task T011: tests/unit/ffmpeg_args_test.go

# Parallel implementation split after tests exist
Task T013: internal/app/cli/render_command.go
Task T015: internal/infra/ffmpeg/command_builder.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 Setup.
2. Complete Phase 2 Foundational.
3. Complete Phase 3 (US1).
4. Validate US1 independently with T009-T012.
5. Demo/deploy MVP.

### Incremental Delivery

1. Deliver US1 overlay functionality.
2. Add US2 placement/style guarantees.
3. Add US3 font-size control and validation.
4. Finish with Polish and full regressions.

### Parallel Team Strategy

1. Team completes Setup + Foundational together.
2. Developer A owns US1 core render/ffmpeg integration.
3. Developer B owns US2 placement/style tests and implementation.
4. Developer C owns US3 CLI validation and font-size propagation.

---

## Notes

- [P] tasks indicate different files and no direct dependency on incomplete tasks.
- [US1], [US2], [US3] labels map tasks to user stories for traceability.
- Keep changes incremental and independently testable per story.
