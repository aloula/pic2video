---
description: "Task list for 003-ken-burns-effect"
---

# Tasks: Ken Burns Effect Option

**Input**: Design documents from `/specs/003-ken-burns-effect/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: Unit and E2E test tasks are mandatory and must be executed before implementation tasks are considered complete.

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm baseline and prepare artifacts for this feature.

- [x] T001 Verify baseline test health with `go test ./...` from project root `go.mod`
- [x] T002 Verify active feature pointer and branch context (`.specify/feature.json` -> `specs/003-ken-burns-effect`, branch `003-ken-burns-effect`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared motion/effect plumbing needed by all user stories.

**CRITICAL**: No user story implementation starts before this phase is complete.

- [x] T003 Add `ImageEffect` field to `RenderJob` in `internal/app/renderjob/model.go`
- [x] T004 Add `ImageEffect` to `BuildOptions` and map it in `internal/app/renderjob/builder.go`
- [x] T005 Add effect-aware xfade entry point in `internal/app/pipeline/xfade.go`
- [x] T006 Add effect-aware ffmpeg args builder function in `internal/infra/ffmpeg/command_builder.go`
- [x] T007 Wire effect-aware args builder from service layer in `internal/app/renderjob/service.go`

**Checkpoint**: Shared effect plumbing is complete and user stories can proceed.

---

## Phase 3: User Story 1 - Select Motion Style (Priority: P1) MVP

**Goal**: Operator can choose `static`, `kenburns-low`, `kenburns-medium`, `kenburns-high` via CLI with static as default.

**Independent Test**: Validate accepted/rejected flag values and startup output showing selected effect.

### Tests for User Story 1 (Required)

- [x] T008 [P] [US1] Add invalid effect value test in `tests/unit/image_effect_validation_test.go`
- [x] T009 [P] [US1] Add accepted effect values test in `tests/unit/image_effect_validation_test.go`
- [x] T010 [P] [US1] Add startup status effect assertion in `tests/unit/status_format_test.go`
- [x] T011 [P] [US1] Add e2e assertion for `effect=` in startup details in `tests/e2e/smoke_test.go`

### Implementation for User Story 1

- [x] T012 [US1] Add `--image-effect` flag (default `static`) and validation in `internal/app/cli/render_command.go`
- [x] T013 [US1] Pass selected image effect into `BuildJob` options in `internal/app/cli/render_command.go`
- [x] T014 [US1] Include `ImageEffect` in startup status payload in `internal/app/cli/render_command.go`
- [x] T015 [US1] Add effect field to startup formatting in `internal/app/cli/summary.go`
- [x] T016 [US1] Run `go test ./tests/unit -run ImageEffect` and confirm US1 tests pass in `tests/unit/image_effect_validation_test.go`
- [x] T017 [US1] Run `go test ./tests/e2e -run Smoke` and confirm effect appears in startup output in `tests/e2e/smoke_test.go`

**Checkpoint**: User Story 1 is fully functional and independently testable.

---

## Phase 4: User Story 2 - Resolution-Aware Ken Burns (Priority: P2)

**Goal**: Ken Burns adapts to FHD/UHD output dimensions with deterministic behavior.

**Independent Test**: Generated filters differ for FHD and UHD and target the selected output dimensions.

### Tests for User Story 2 (Required)

- [x] T018 [P] [US2] Add motion filter static behavior test in `tests/unit/framing_policy_test.go`
- [x] T019 [P] [US2] Add motion filter medium mode content test in `tests/unit/framing_policy_test.go`
- [x] T020 [P] [US2] Add FHD vs UHD motion filter difference test in `tests/unit/framing_policy_test.go`
- [x] T021 [P] [US2] Update timeline graph test to effect-aware builder in `tests/unit/timeline_test.go`
- [x] T036 [P] [US2] Add E2E test for FHD render success with `--image-effect kenburns-medium` in `tests/e2e/render_fhd_test.go`
- [x] T037 [P] [US2] Add E2E test for UHD render success with `--image-effect kenburns-medium` in `tests/e2e/render_uhd_test.go`

### Implementation for User Story 2

- [x] T022 [US2] Implement `BuildMotionFilter(effect, width, height, imageDur)` in `internal/app/pipeline/framing.go`
- [x] T023 [US2] Implement effect-aware graph generation in `internal/app/pipeline/xfade.go`
- [x] T024 [US2] Apply static-vs-kenburns graph logic in `internal/infra/ffmpeg/command_builder.go`
- [x] T025 [US2] Pass `ImageEffect` into ffmpeg argument generation in `internal/app/renderjob/service.go`
- [x] T026 [US2] Run `go test ./tests/unit -run 'Framing|Timeline'` and confirm resolution-aware behavior in `tests/unit/framing_policy_test.go`
- [x] T038 [US2] Run `go test ./tests/e2e -run RenderFHD` and confirm FHD Ken Burns E2E passes in `tests/e2e/render_fhd_test.go`
- [x] T039 [US2] Run `go test ./tests/e2e -run RenderUHD` and confirm UHD Ken Burns E2E passes in `tests/e2e/render_uhd_test.go`

**Checkpoint**: User Story 2 is fully functional and independently testable.

---

## Phase 5: User Story 3 - Quality-First Motion (Priority: P3)

**Goal**: Ken Burns prioritizes smooth, high-quality motion even if processing time increases.

**Independent Test**: Generated Ken Burns filters include quality/smoothness directives.

### Tests for User Story 3 (Required)

- [x] T027 [P] [US3] Add table-driven quality directives assertions for `kenburns-low`, `kenburns-medium`, and `kenburns-high` requiring `lanczos` in `tests/unit/framing_policy_test.go`
- [x] T028 [P] [US3] Add table-driven smooth-motion assertions for `kenburns-low`, `kenburns-medium`, and `kenburns-high` requiring `zoompan` and `fps=60` in `tests/unit/framing_policy_test.go`

### Implementation for User Story 3

- [x] T029 [US3] Tune Ken Burns presets for low/medium/high and quality-priority defaults in `internal/app/pipeline/framing.go`
- [x] T030 [US3] Ensure per-resolution scaling and pan amplitude preserve smoothness in `internal/app/pipeline/framing.go`
- [x] T031 [US3] Run `go test ./tests/unit -run Framing` and verify quality-first motion settings in `tests/unit/framing_policy_test.go`

**Checkpoint**: User Story 3 is fully functional and independently testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, regression, and final validation.

- [x] T032 [P] Update usage and flag docs in `README.md` for `--image-effect`
- [x] T033 [P] Align startup output examples in `specs/003-ken-burns-effect/quickstart.md`
- [x] T034 [P] Align contract details in `specs/003-ken-burns-effect/contracts/cli-contract.md`
- [x] T035 Run full regression `go test ./...` and confirm no regressions from project root `go.mod`
- [x] T040 Run `make build-all` after successful implementation and confirm multi-platform binaries are generated in `bin/`
- [x] T041 Add final global fade-in/fade-out chain for all output videos in `internal/infra/ffmpeg/command_builder.go`
- [x] T042 Add unit assertions for global fade-in/fade-out in static and Ken Burns modes in `tests/unit/ffmpeg_args_test.go`
- [x] T043 Update docs for final output fade behavior in `README.md`
- [x] T044 Run full regression `go test ./...` and confirm fade behavior introduces no regressions

---

## Dependencies & Execution Order

### Phase Dependencies

- Phase 1: no dependencies
- Phase 2: depends on Phase 1 and blocks all user stories
- Phase 3 (US1): depends on Phase 2
- Phase 4 (US2): depends on Phase 2
- Phase 5 (US3): depends on Phase 2 and benefits from US2 infrastructure
- Phase 6: depends on completion of US1, US2, and US3

### User Story Dependencies

- US1 (P1): no dependency on other stories after foundational work
- US2 (P2): no functional dependency on US1 after foundational work
- US3 (P3): builds on motion filter behavior implemented in US2

### Within Each User Story

- Tests first and failing
- Implementation second
- Story-specific tests green before moving on

---

## Parallel Opportunities

- [P] tasks in each story can be run in parallel when they touch different files.
- US1 test tasks T008-T011 can run in parallel.
- US2 test tasks T018-T021 and T036-T037 can run in parallel.
- US3 test tasks T027-T028 can run in parallel.
- Polish docs tasks T032-T034 can run in parallel.

## Parallel Example: User Story 1

- Run T008 and T009 together (CLI validation tests)
- Run T010 and T011 together (status assertions)

## Implementation Strategy

### MVP First (US1)

1. Complete Phase 1 and Phase 2
2. Complete Phase 3 (US1)
3. Validate US1 independently

### Incremental Delivery

1. Add US2 resolution-aware behavior
2. Add US3 quality-first tuning
3. Finish with Phase 6 documentation and full regression
