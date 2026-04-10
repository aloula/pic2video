# Tasks: Professional Photo-to-Video CLI

**Input**: Design documents from `/specs/001-build-photo-video-cli/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Unit and E2E test tasks are mandatory for every user story and must fail before implementation tasks begin.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize the CLI project baseline and essential build/docs entrypoints.

- [X] T001 Initialize module and dependencies in go.mod
- [X] T002 Create CLI entrypoint in cmd/pic2video/main.go
- [X] T003 [P] Create root command scaffold in internal/app/cli/root.go
- [X] T004 [P] Add build/test automation targets in Makefile
- [X] T005 [P] Create root usage/build/test document in README.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Implement shared domain and infrastructure components required by all stories.

**CRITICAL**: No user story work starts before this phase is complete.

- [X] T006 Define output profile model in internal/domain/profile/profile.go
- [X] T007 [P] Define media asset model in internal/domain/media/asset.go
- [X] T008 [P] Define transition segment model in internal/domain/transition/segment.go
- [X] T009 Define render job and summary models in internal/app/renderjob/model.go
- [X] T010 Implement render job validator in internal/app/renderjob/validator.go
- [X] T011 [P] Implement ordering strategies service in internal/app/pipeline/order.go
- [X] T012 [P] Implement FFprobe metadata adapter in internal/infra/ffmpeg/ffprobe.go
- [X] T013 Implement FFmpeg command argument builder in internal/infra/ffmpeg/command_builder.go
- [X] T014 [P] Implement NVENC availability detector in internal/infra/nvenc/detect.go
- [X] T015 [P] Implement encoder selection policy in internal/infra/nvenc/policy.go
- [X] T016 [P] Implement image asset filesystem listing in internal/infra/fsio/filesystem.go
- [X] T017 Define classified errors and exit class mapping in internal/app/renderjob/errors.go

**Checkpoint**: Foundation complete, user stories can proceed.

---

## Phase 3: User Story 1 - Render YouTube-ready Slideshow (Priority: P1) 🎯 MVP

**Goal**: Render FHD/UHD 16:9 slideshow videos with cross-fade transitions from input images.

**Independent Test**: CLI renders valid photo folders into playable FHD/UHD output with transitions and success status.

### Tests for User Story 1 (REQUIRED)

- [X] T018 [P] [US1] Add profile resolution mapping tests in tests/unit/profile_test.go
- [X] T019 [P] [US1] Add cross-fade timeline math tests in tests/unit/timeline_test.go
- [X] T020 [P] [US1] Add ordering mode behavior tests in tests/unit/order_mode_test.go
- [X] T021 [US1] Add FHD render happy-path E2E test in tests/e2e/render_fhd_test.go
- [X] T022 [US1] Add UHD render happy-path E2E test in tests/e2e/render_uhd_test.go
- [X] T023 [US1] Add explicit ordering E2E test in tests/e2e/render_explicit_order_test.go

### Implementation for User Story 1

- [X] T024 [US1] Implement render command flags and validation in internal/app/cli/render_command.go
- [X] T025 [P] [US1] Implement CLI options to render job builder in internal/app/renderjob/builder.go
- [X] T026 [US1] Implement cross-fade filter graph assembly in internal/app/pipeline/xfade.go
- [X] T027 [US1] Implement FFmpeg execution adapter in internal/infra/ffmpeg/executor.go
- [X] T028 [US1] Implement explicit order manifest parser in internal/infra/fsio/order_manifest.go
- [X] T029 [US1] Implement render orchestration workflow in internal/app/renderjob/service.go
- [X] T030 [US1] Wire render command execution path in cmd/pic2video/main.go

**Checkpoint**: US1 is independently functional and testable.

---

## Phase 4: User Story 2 - Preserve Professional Image Quality (Priority: P2)

**Goal**: Preserve visual quality with deterministic framing/scaling and quality warnings.

**Independent Test**: Mixed-aspect input set renders with deterministic framing and warning visibility for low-resolution sources.

### Tests for User Story 2 (REQUIRED)

- [X] T031 [P] [US2] Add framing policy unit tests in tests/unit/framing_policy_test.go
- [X] T032 [P] [US2] Add quality warning rule unit tests in tests/unit/quality_warning_test.go
- [X] T033 [US2] Add mixed-aspect quality E2E test in tests/e2e/render_mixed_aspect_test.go

### Implementation for User Story 2

- [X] T034 [US2] Implement deterministic framing/scaling policy in internal/app/pipeline/framing.go
- [X] T035 [US2] Implement quality warning evaluator in internal/app/renderjob/quality.go
- [X] T036 [US2] Integrate warnings into render summary composition in internal/app/renderjob/summary.go

**Checkpoint**: US1 and US2 are independently testable with quality guarantees.

---

## Phase 5: User Story 3 - Operate Reliably from CLI (Priority: P3)

**Goal**: Provide deterministic CLI failures, clear diagnostics, and explicit encoder reporting.

**Independent Test**: Invalid invocations fail with deterministic exits and actionable output; valid runs return concise completion summaries.

### Tests for User Story 3 (REQUIRED)

- [X] T037 [P] [US3] Add required-flag validation unit test in tests/unit/cli_validation_test.go
- [X] T038 [P] [US3] Add encoder fallback policy unit tests in tests/unit/encoder_policy_test.go
- [X] T039 [US3] Add invalid input exit behavior E2E test in tests/e2e/invalid_input_test.go
- [X] T040 [US3] Add missing FFmpeg environment E2E test in tests/e2e/missing_ffmpeg_test.go

### Implementation for User Story 3

- [X] T041 [US3] Implement deterministic exit code mapping in internal/app/cli/exitcodes.go
- [X] T042 [US3] Implement user-facing error formatter in internal/app/cli/errors.go
- [X] T043 [US3] Implement concise completion summary formatter in internal/app/cli/summary.go
- [X] T044 [US3] Implement encoder selection reporting helper in internal/infra/nvenc/report.go

**Checkpoint**: All stories are independently functional with deterministic CLI behavior.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final hardening, documentation quality, and benchmark validation across stories.

- [X] T045 [P] Document fixture requirements in tests/fixtures/README.md
- [X] T046 [P] Document benchmark fixture profile in tests/fixtures/benchmark/README.md
- [X] T047 Update quickstart usage and troubleshooting in specs/001-build-photo-video-cli/quickstart.md
- [X] T048 [P] Update root README for usage/build/run/troubleshooting changes in README.md
- [X] T049 [P] Add smoke regression test in tests/e2e/smoke_test.go
- [X] T050 Add 50-image performance E2E test in tests/e2e/perf_fhd_50_images_test.go
- [X] T051 Add performance threshold/report workflow test in tests/e2e/perf_report_test.go
- [X] T052 Add package-level app documentation in internal/app/doc.go
- [X] T053 Document explicit invalid-media validation policy (reject vs skip) in README.md and specs/001-build-photo-video-cli/quickstart.md
- [X] T054 Document measurable SC-003 quality rubric acceptance checks and mapping to tests in README.md and specs/001-build-photo-video-cli/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup (Phase 1): No dependencies, start immediately.
- Foundational (Phase 2): Depends on Setup; blocks all user stories.
- User Stories (Phase 3+): Depend on Foundational completion.
- Polish (Phase 6): Depends on completion of selected user stories.

### User Story Dependencies

- US1 (P1): Starts after Foundational; independent MVP.
- US2 (P2): Starts after Foundational; can reuse US1 pipeline outputs but remains independently testable.
- US3 (P3): Starts after Foundational; can run in parallel with US2 after shared interfaces stabilize.

### Within Each User Story

- Write unit and E2E tests first and verify failures.
- Implement domain/application logic next.
- Integrate infrastructure and CLI wiring after core logic.
- Validate story end-to-end before moving to next story.

---

## Parallel Execution Examples

### User Story 1

- Run in parallel: T018, T019, T020 (independent unit test files).
- Run in parallel: T025 and T028 (builder and manifest parser in different modules).

### User Story 2

- Run in parallel: T031 and T032 (independent quality policy unit tests).

### User Story 3

- Run in parallel: T037 and T038 (CLI validation vs encoder policy tests).
- Run in parallel: T042 and T044 (error formatting vs encoder report helper).

---

## Implementation Strategy

### MVP First

1. Complete Phase 1 and Phase 2.
2. Complete Phase 3 (US1).
3. Validate FHD/UHD rendering behavior with E2E tests.
4. Demo/deploy MVP.

### Incremental Delivery

1. Add US2 for quality fidelity and warning behavior.
2. Add US3 for operational reliability and deterministic exits.
3. Run SC-004 performance workflow and threshold validation.
4. Finish polish tasks including README/quickstart consistency.

### Team Parallelization

1. One developer completes Foundational tasks.
2. After foundation is stable:
   - Dev A focuses on US2.
   - Dev B focuses on US3.
3. Merge after each story passes unit + E2E gates.
