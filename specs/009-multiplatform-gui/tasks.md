# Tasks: Multiplatform Desktop GUI

**Input**: Design documents from `/specs/009-multiplatform-gui/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Unit and E2E test tasks are MANDATORY for every user story and MUST be executed before implementation tasks are considered complete.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize GUI module structure and build wiring.

- [X] T001 Add desktop GUI dependency in go.mod and go.sum
- [X] T002 Create GUI entrypoint scaffold in cmd/pic2video-gui/main.go
- [X] T003 [P] Add GUI build/run targets in Makefile
- [X] T004 [P] Create GUI package scaffold in internal/app/gui/doc.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core GUI runtime primitives that block all user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T005 Define GuiRunConfiguration/GuiRunState/GuiLogEntry models in internal/app/gui/model.go
- [X] T006 Implement launch-directory defaults resolver in internal/app/gui/defaults.go
- [X] T007 Implement preflight validation service in internal/app/gui/validate.go
- [X] T008 Implement CLI-args mapping from GUI configuration in internal/app/gui/command.go
- [X] T009 Implement render process runner with output streaming in internal/app/gui/runner.go
- [X] T010 Implement state machine and active-run lock in internal/app/gui/state.go

**Checkpoint**: Foundation ready; user stories can proceed independently.

---

## Phase 3: User Story 1 - Configure Render Job Visually (Priority: P1) 🎯 MVP

**Goal**: Provide a clean GUI where users choose input/output folders, configure all user-facing options, and start rendering with folder-only output selection.

**Independent Test**: Launch GUI, verify default folders from launch directory, configure options, start run, and verify command receives selected values and output destination folder behavior.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T011 [P] [US1] Add unit tests for defaults and option-to-command mapping in tests/unit/gui_config_test.go
- [X] T012 [P] [US1] Add E2E test for GUI configuration flow in tests/e2e/gui_config_flow_test.go

### Implementation for User Story 1

- [X] T013 [P] [US1] Implement input/output folder picker form in internal/app/gui/view_form.go
- [X] T014 [US1] Implement option controls for all user-facing render flags in internal/app/gui/view_options.go
- [X] T015 [US1] Bind GUI controls to GuiRunConfiguration in internal/app/gui/controller_config.go
- [X] T016 [US1] Enforce output-folder-only policy with profile auto-filename preview in internal/app/gui/controller_output.go
- [X] T017 [US1] Integrate preflight validation into Start action in internal/app/gui/controller_start.go

**Checkpoint**: User Story 1 is independently functional and testable (MVP).

---

## Phase 4: User Story 2 - Track Run Progress and Outcome (Priority: P2)

**Goal**: Show clear lifecycle status (`idle`, `loading files`, `processing`, `finished`, `failed`) and prevent concurrent starts.

**Independent Test**: Start a run from GUI, observe lifecycle transitions, verify failed runs show `failed`, and verify second start is blocked while active.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T018 [P] [US2] Add unit tests for status transitions and run-lock behavior in tests/unit/gui_status_test.go
- [X] T019 [P] [US2] Add E2E test for status lifecycle (success and failure) in tests/e2e/gui_status_flow_test.go

### Implementation for User Story 2

- [X] T020 [P] [US2] Implement status indicator widget in internal/app/gui/view_status.go
- [X] T021 [US2] Map runner lifecycle events to status transitions in internal/app/gui/controller_run.go
- [X] T022 [US2] Enforce no-concurrent-start rule in internal/app/gui/state.go
- [X] T023 [US2] Surface actionable validation/runtime errors in internal/app/gui/view_errors.go

**Checkpoint**: User Stories 1 and 2 each work independently.

---

## Phase 5: User Story 3 - View Execution Log in GUI (Priority: P3)

**Goal**: Display chronological runtime output in a compact log box for observability and debugging.

**Independent Test**: Run render from GUI and confirm stdout/stderr lines appear in-order in the log panel and include failure output when applicable.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T024 [P] [US3] Add unit tests for log buffering and ordering in tests/unit/gui_log_test.go
- [X] T025 [P] [US3] Add E2E test for log panel runtime output flow in tests/e2e/gui_log_flow_test.go

### Implementation for User Story 3

- [X] T026 [P] [US3] Implement bounded log store and append behavior in internal/app/gui/log_store.go
- [X] T027 [US3] Implement log panel UI component in internal/app/gui/view_log.go
- [X] T028 [US3] Wire runner stdout/stderr streams into log store in internal/app/gui/runner.go
- [X] T029 [US3] Bind log updates to UI refresh and auto-scroll behavior in internal/app/gui/controller_log.go

**Checkpoint**: All user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final hardening, documentation, and verification across stories.

- [X] T030 [P] Update GUI usage/build/run documentation in README.md
- [X] T031 [P] Align quickstart verification steps with implemented GUI behavior in specs/009-multiplatform-gui/quickstart.md
- [X] T032 Execute full unit and E2E suite and record outcomes in specs/009-multiplatform-gui/tasks.md
- [X] T033 Execute make build-all and record cross-platform build result in specs/009-multiplatform-gui/tasks.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: Starts immediately.
- **Phase 2 (Foundational)**: Depends on Phase 1 and blocks all user stories.
- **Phases 3-5 (User Stories)**: Depend on Phase 2; then may run in parallel by team capacity.
- **Phase 6 (Polish)**: Depends on completion of all selected user stories.

### User Story Dependencies

- **US1 (P1)**: Starts after Foundational phase; independent MVP slice.
- **US2 (P2)**: Starts after Foundational phase; integrates with runner/state primitives but remains independently testable.
- **US3 (P3)**: Starts after Foundational phase; depends on runner output stream, remains independently testable.

### Within Each User Story

- Tests first (must fail before implementation).
- View/model primitives before controller wiring.
- Controller wiring before end-to-end validation.

## Parallel Opportunities

- Setup: T003 and T004 parallel after T001/T002 starts.
- Foundational: T007 and T008 can run in parallel after T005/T006; T009 and T010 can proceed in parallel once command/state interfaces are stable.
- US1: T011 and T012 parallel; T013 and T014 parallel.
- US2: T018 and T019 parallel; T020 parallel with early T021 scaffolding.
- US3: T024 and T025 parallel; T026 and T027 parallel.
- Polish: T030 and T031 parallel.

## Parallel Example: User Story 1

```bash
# Run US1 tests in parallel
Task: T011 tests/unit/gui_config_test.go
Task: T012 tests/e2e/gui_config_flow_test.go

# Build US1 UI components in parallel
Task: T013 internal/app/gui/view_form.go
Task: T014 internal/app/gui/view_options.go
```

## Implementation Strategy

### MVP First (US1)

1. Complete Setup (Phase 1).
2. Complete Foundational (Phase 2).
3. Deliver US1 (Phase 3) with tests passing.
4. Validate MVP manually with quickstart Scenario 1 and 2.

### Incremental Delivery

1. Add US2 status lifecycle after MVP.
2. Add US3 log panel after status lifecycle.
3. Finish with documentation and cross-platform build verification.

### Parallel Team Strategy

1. Team aligns on Foundational interfaces (`model.go`, `command.go`, `runner.go`, `state.go`).
2. Dev A owns US1 form/config tasks.
3. Dev B owns US2 status/state tasks.
4. Dev C owns US3 log panel/log stream tasks.

## Notes

- [P] tasks are safe for parallel work because they target different files with no hard dependency on incomplete tasks.
- Each task includes a concrete file path for direct execution.
- Keep commits small by completing one task or one cohesive task cluster per commit.

## Implementation Results

- T032 verification: `go test ./... -count=1` passed on 2026-04-11.
- T033 build verification: `make build-all` passed on 2026-04-11.