---
description: "Task list for 002-video-processing-status"
---

# Tasks: Video Processing Status

**Input**: Design documents from `/specs/002-video-processing-status/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/cli-contract.md ✅, quickstart.md ✅

**Tech stack**: Go 1.23+, standard library only (`path/filepath`, `strings`, `fmt`)
**Changed files**: `internal/app/cli/summary.go`, `internal/app/cli/render_command.go`
**New test file**: `tests/unit/status_format_test.go`
**Extended test file**: `tests/e2e/smoke_test.go`

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

---

## Phase 1: Setup

**Purpose**: No new dependencies or project initialization required. This phase confirms the baseline compiles and tests are green before changes begin.

- [x] T001 Verify `go test ./...` passes clean before any changes (baseline green)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The two pure formatter helper functions that both US1 and US2 depend on. Must be complete before any user story implementation.

**⚠️ CRITICAL**: `FormatOutputFormat` and `FormatElapsed` are shared by both user stories. Complete this phase before US1 or US2 implementation.

- [x] T002 [P] Write failing unit tests for `FormatOutputFormat` (known extension, unrecognized extension, no extension → `UNKNOWN`) in `tests/unit/status_format_test.go`
- [x] T003 [P] Write failing unit tests for `FormatElapsed` (sub-second → `< 1s`, under-60s → `%.1fs`, exactly-60s → `1m 0s`, over-60s → `Xm Ys`) in `tests/unit/status_format_test.go`
- [x] T004 Implement `FormatOutputFormat(outputPath string) string` in `internal/app/cli/summary.go` — derive label from `filepath.Ext`, uppercase, fallback `UNKNOWN` (unblocks T002 going green)
- [x] T005 Implement `FormatElapsed(seconds float64) string` in `internal/app/cli/summary.go` — three-branch threshold logic (unblocks T003 going green)
- [x] T006 Run `go test ./tests/unit/...` and confirm T002+T003 tests pass green

**Checkpoint**: `FormatOutputFormat` and `FormatElapsed` are implemented, tested green, and ready for use in both user stories.

---

## Phase 3: User Story 1 - Render Completion Report (Priority: P1) 🎯 MVP

**Goal**: After a successful render, stdout contains a labeled completion line with `files=N`, `format=EXT`, and a human-readable `elapsed=` value. Failing renders emit no completion report.

**Independent Test**: Run a render to completion and assert stdout contains `files=`, `format=`, and a non-raw-decimal elapsed value. Can be validated with any fixture image set.

### Tests for User Story 1

> **Write these tests FIRST — ensure they FAIL before implementation**

- [x] T007 [P] [US1] Write failing unit test: `FormatSummary` output contains `files=12` field in `tests/unit/status_format_test.go`
- [x] T008 [P] [US1] Write failing unit test: `FormatSummary` output contains `format=MP4` field derived from output path in `tests/unit/status_format_test.go`
- [x] T009 [P] [US1] Write failing unit test: `FormatSummary` output `elapsed=` value is human-readable (not raw decimal) in `tests/unit/status_format_test.go`
- [x] T010 [P] [US1] Write failing E2E assertion: `smoke_test.go` asserts stdout contains `files=` and `format=` fields after successful render in `tests/e2e/smoke_test.go`
- [x] T026 [P] [US1] Write failing E2E assertion: when render fails, stdout MUST NOT contain `status=success` or completion fields (`files=`, `format=`, `elapsed=`) in `tests/e2e/smoke_test.go`

### Implementation for User Story 1

- [x] T011 [US1] Extend `FormatSummary` signature in `internal/app/cli/summary.go` to accept `outputPath string` parameter and emit `files=<processed> format=<FormatOutputFormat(outputPath)>` fields before `elapsed=`
- [x] T012 [US1] Replace raw `elapsed=%.3fs` in `FormatSummary` with `elapsed=<FormatElapsed(elapsed)>` in `internal/app/cli/summary.go`
- [x] T013 [US1] Update `FormatSummary` call in `internal/app/cli/render_command.go` to pass `summary.OutputPath` as the new argument
- [x] T014 [US1] Run `go test ./tests/unit/...` and confirm T007–T009 tests pass green
- [x] T015 [US1] Run `go test ./tests/e2e/...` and confirm T010 and T026 E2E assertions pass
- [x] T027 [US1] Run focused failure-path E2E scenario and confirm T026 passes

**Checkpoint**: User Story 1 is fully functional and independently testable. Completion line shows `files=`, `format=`, and human-readable `elapsed=`.

---

## Phase 4: User Story 2 - Pre-Render Status Announcement (Priority: P2)

**Goal**: Before FFmpeg encoding begins, stdout shows `status=starting files=<N> format=<EXT>`, giving the operator an immediate confirmation of input count and output format.

**Independent Test**: Run a render with any fixture set; assert a `status=starting` line appears on stdout before encoding output. Story is independently testable without US1 being complete, though in practice US1 completes first.

### Tests for User Story 2

> **Write these tests FIRST — ensure they FAIL before implementation**

- [x] T016 [P] [US2] Write failing unit test: `FormatAnnouncement(files int, outputPath string)` returns a string containing `status=starting`, `files=N`, and `format=EXT` in `tests/unit/status_format_test.go`
- [x] T017 [P] [US2] Write failing E2E assertion: stdout from a render run contains a `status=starting` line before the `status=success` line in `tests/e2e/smoke_test.go`
- [x] T028 [P] [US2] Write failing E2E timing assertion: `status=starting` appears on stdout within 1 second of command invocation in `tests/e2e/smoke_test.go`

### Implementation for User Story 2

- [x] T018 [US2] Implement `FormatAnnouncement(files int, outputPath string) string` in `internal/app/cli/summary.go` — returns `status=starting files=<N> format=<FormatOutputFormat(outputPath)>`
- [x] T019 [US2] Add pre-render announcement print in `internal/app/cli/render_command.go`: after `pipeline.ApplyOrder()` resolves assets and before `renderjob.BuildJob()` is called, emit `fmt.Fprintln(cmd.OutOrStdout(), FormatAnnouncement(len(assets), output))`
- [x] T020 [US2] Run `go test ./tests/unit/...` and confirm T016 test passes green
- [x] T021 [US2] Run `go test ./tests/e2e/...` and confirm T017 and T028 E2E assertions pass
- [x] T029 [US2] Run focused timing E2E scenario and confirm T028 passes

**Checkpoint**: User Story 2 is fully functional. Both `status=starting` and `status=success` lines appear in stdout; each contains `files=` and `format=` with correct values.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Full test suite validation, README update, backward compat verification.

- [x] T022 Run full `go test ./...` and confirm all unit and E2E tests pass green with no regressions
- [x] T023 Update `README.md` to document the new `status=starting` pre-render line, the new `files=` and `format=` fields in the completion line, and the updated `elapsed=` human-readable format
- [x] T024 [P] Validate backward compat: confirm `processed=` field is still present in completion line output (existing field retained per contract)
- [x] T025 [P] Run quickstart.md validation — verify the example command output in `specs/002-video-processing-status/quickstart.md` matches actual render output

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1** (Setup): No dependencies — start immediately
- **Phase 2** (Foundational): Depends on Phase 1 — **blocks US1 and US2**
- **Phase 3** (US1): Depends on Phase 2 completion
- **Phase 4** (US2): Depends on Phase 2 completion; can run in parallel with Phase 3
- **Phase 5** (Polish): Depends on Phase 3 + Phase 4 completion

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2 — no dependency on US2
- **US2 (P2)**: Can start after Phase 2 — no dependency on US1; `FormatAnnouncement` reuses `FormatOutputFormat` from Phase 2

### Within Each User Story

1. Write unit tests → run (should fail)
2. Implement production code (formatter functions / print statement)
3. Run tests (should go green)
4. Add E2E assertion → run (should fail)
5. Verify E2E passes (no new implementation needed — E2E exercises existing CLI path)

---

## Parallel Opportunities

### Phase 2 can parallelize T002 and T003

```
T002 (write FormatOutputFormat tests)  ──┐
                                          ├── T004 (implement FormatOutputFormat)
T003 (write FormatElapsed tests)       ──┘
                                          └── T005 (implement FormatElapsed)
                                                    └── T006 (run unit tests green)
```

### Phase 3 tests can parallelize (T007, T008, T009, T010, T026)

```
T007 [US1 unit: files field]      ──┐
T008 [US1 unit: format field]     ──┤── T011 (extend FormatSummary)
T009 [US1 unit: elapsed format]   ──┤      └── T012 (replace elapsed format)
T010 [US1 E2E success assertion]  ──┤              └── T013 (update call site)
T026 [US1 E2E failure assertion]  ──┘
```

### Phase 4 tests can parallelize (T016, T017, T028)

```
T016 [US2 unit: FormatAnnouncement] ──┐
T017 [US2 E2E order assertion]      ──┤── T018 (implement FormatAnnouncement)
T028 [US2 E2E timing assertion]     ──┘      └── T019 (add print to render_command.go)
```

### Phase 5 can parallelize T024 and T025

```
T022 (full go test ./...) ──> T023 (README) ──> T024 [P] (backward compat check)
                                               └> T025 [P] (quickstart validation)
```

---

## Implementation Strategy

**MVP scope**: Phase 1 + Phase 2 + Phase 3 (US1 only) delivers the core completion report with `files=`, `format=`, and human-readable elapsed. This is the highest-value increment.

**Full completion**: Add Phase 4 (US2 pre-render announcement) + Phase 5 (polish/docs) for the complete feature.

**Total tasks**: 29
- Phase 1: 1 task
- Phase 2 (Foundational): 5 tasks
- Phase 3 (US1): 11 tasks
- Phase 4 (US2): 8 tasks
- Phase 5 (Polish): 4 tasks

**Task count per user story**:
- US1: 11 tasks (5 test, 6 implementation)
- US2: 8 tasks (3 test, 5 implementation)
