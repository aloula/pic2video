# Tasks: Mixed Video Photo Input

**Input**: Design documents from `/specs/008-video-photo-mixed-input/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md, quickstart.md

**Tests**: Unit and E2E test tasks are MANDATORY for every user story and MUST be executed before implementation tasks are considered complete.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare shared fixtures and docs for mixed image/video workflows.

- [X] T001 Add mixed image+video fixture guidance in tests/fixtures/README.md
- [X] T002 Add mixed-input usage examples and constraints in README.md
- [X] T003 [P] Add quickstart command validation notes in specs/008-video-photo-mixed-input/quickstart.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core media modeling and probe infrastructure required before user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Extend media asset model for mixed-media metadata in internal/domain/media/asset.go
- [X] T005 Add shared mixed-media input listing helpers in internal/infra/fsio/filesystem.go
- [X] T006 [P] Extend ffprobe metadata extraction for video width/height/duration/fps in internal/infra/ffmpeg/ffprobe.go
- [X] T007 [P] Add render job options for output fps and mixed-media stats in internal/app/renderjob/model.go
- [X] T008 Wire mixed-media build options through builder validation in internal/app/renderjob/builder.go
- [X] T009 Add shared summary fields for mixed-media and fps output in internal/app/cli/summary.go

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Render videos with photos in one timeline (Priority: P1) 🎯 MVP

**Goal**: Allow `render` to process image and video assets together using existing ordering modes.

**Independent Test**: Render a folder with JPG and MP4 assets and verify output includes both media types in expected order.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T010 [P] [US1] Add unit tests for mixed-media discovery and filtering in tests/unit/order_mode_test.go
- [X] T011 [P] [US1] Add unit tests for mixed-media timeline segment creation in tests/unit/timeline_test.go
- [X] T012 [US1] Add E2E test for mixed image+video render success path in tests/e2e/render_mixed_aspect_test.go

### Implementation for User Story 1

- [X] T013 [US1] Implement mixed-media listing and merge behavior in internal/infra/fsio/filesystem.go
- [X] T014 [US1] Apply ordering modes consistently to mixed assets in internal/app/pipeline/order.go
- [X] T015 [US1] Integrate mixed assets into render command input flow in internal/app/cli/render_command.go
- [X] T016 [US1] Build unified mixed-media timeline handling in internal/app/renderjob/service.go
- [X] T017 [US1] Update render summary reporting for mixed-media counts in internal/app/renderjob/summary.go

**Checkpoint**: User Story 1 should be fully functional and independently testable.

---

## Phase 4: User Story 2 - Preserve video aspect ratio at profile quality (Priority: P2)

**Goal**: Scale videos to profile-equivalent target using high-quality resampling without distortion.

**Independent Test**: Render portrait and landscape videos under FHD/UHD and verify aspect ratio is preserved while profile-target framing is applied.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T018 [P] [US2] Add unit tests for profile-target dimension transform policy in tests/unit/framing_policy_test.go
- [X] T019 [P] [US2] Add unit tests for ffmpeg video scale/pad/lanczos arguments in tests/unit/ffmpeg_args_test.go
- [X] T020 [US2] Add E2E test for mixed portrait/landscape clip scaling in tests/e2e/render_mixed_aspect_test.go

### Implementation for User Story 2

- [X] T021 [US2] Implement aspect-ratio-preserving video transform policy in internal/app/pipeline/framing.go
- [X] T022 [US2] Add high-quality scale+pad filter graph generation for video segments in internal/infra/ffmpeg/command_builder.go
- [X] T023 [US2] Integrate video transform policy into timeline segment preparation in internal/app/renderjob/service.go
- [X] T024 [US2] Surface transform policy outcomes in startup/summary output in internal/app/cli/summary.go

**Checkpoint**: User Stories 1 and 2 should both work independently.

---

## Phase 5: User Story 3 - Match output frame rate for video segments (Priority: P3)

**Goal**: Normalize mixed-media output to user-selected fps with stable timing.

**Independent Test**: Render mixed assets with a chosen fps and verify output metadata and segment behavior match selected fps.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T025 [P] [US3] Add unit tests for fps flag validation and defaults in tests/unit/cli_validation_test.go
- [X] T026 [P] [US3] Add unit tests for fps propagation into ffmpeg args in tests/unit/ffmpeg_args_test.go
- [X] T027 [P] [US3] Add unit tests for short-clip transition clamping with fps normalization in tests/unit/timeline_test.go
- [X] T028 [US3] Add E2E test for mixed-media render with selected fps in tests/e2e/render_fhd_test.go

### Implementation for User Story 3

- [X] T029 [US3] Add output fps selection flag and validation in internal/app/cli/render_command.go
- [X] T030 [US3] Propagate selected fps through build options and render job model in internal/app/renderjob/builder.go
- [X] T031 [US3] Apply fps normalization to video filter pipeline and output args in internal/infra/ffmpeg/command_builder.go
- [X] T032 [US3] Clamp transition boundaries for short video clips during timeline assembly in internal/app/renderjob/service.go
- [X] T033 [US3] Emit fps and normalization warnings in summary output in internal/app/renderjob/summary.go

**Checkpoint**: All user stories should now be independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Hardening, documentation alignment, and full regression validation.

- [X] T034 [P] Add mixed-media contract and behavior notes to docs in specs/008-video-photo-mixed-input/contracts/cli-contract.md
- [X] T035 [P] Refine quickstart scenarios with expected verification commands in specs/008-video-photo-mixed-input/quickstart.md
- [X] T036 Run full unit regression suite with `go test ./tests/unit/...` from repository root
- [X] T037 Run full E2E regression suite with `go test ./tests/e2e/...` from repository root
- [X] T038 Run full project tests with `go test ./...` from repository root
- [X] T039 Run `make build-all` and record successful result in README.md
- [X] T040 Execute SC-004 visual quality acceptance sampling (minimum 10 scaled segments) and record pass/fail evidence in specs/008-video-photo-mixed-input/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies, starts immediately.
- **Foundational (Phase 2)**: Depends on Setup completion and blocks all user stories.
- **User Story Phases (Phase 3-5)**: Depend on Foundational completion.
- **Polish (Phase 6)**: Depends on all user stories being complete.

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2 and is MVP.
- **US2 (P2)**: Can start after Phase 2; depends functionally on US1 mixed timeline path.
- **US3 (P3)**: Can start after Phase 2; integrates with US1 timeline and US2 transform pipeline.

### Within Each User Story

- Write unit and E2E tests first and verify they fail.
- Implement code changes after tests exist.
- Re-run story tests until passing before moving to next story.

### Parallel Opportunities

- Tasks marked [P] in Setup and Foundational can run concurrently.
- In each user story, [P] unit test tasks can run concurrently.
- US2 and US3 can run in parallel after US1 baseline mixed timeline implementation is merged.

---

## Parallel Example: User Story 2

```bash
# Parallel test authoring for US2
Task T018: tests/unit/framing_policy_test.go
Task T019: tests/unit/ffmpeg_args_test.go

# Parallel implementation split for US2
Task T021: internal/app/pipeline/framing.go
Task T022: internal/infra/ffmpeg/command_builder.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 Setup.
2. Complete Phase 2 Foundational.
3. Complete Phase 3 (US1).
4. Validate US1 independently with T010-T012.
5. Demo/deploy MVP.

### Incremental Delivery

1. Deliver US1 mixed-media timeline functionality.
2. Add US2 aspect-ratio-safe profile scaling.
3. Add US3 fps selection and normalization.
4. Finish with Polish and full regressions.

### Parallel Team Strategy

1. Team completes Setup + Foundational together.
2. Developer A owns US1 mixed discovery and timeline integration.
3. Developer B owns US2 transform policy and ffmpeg scaling.
4. Developer C owns US3 fps selection and normalization path.

---

## Notes

- [P] tasks indicate different files and no direct dependency on incomplete tasks.
- [US1], [US2], [US3] labels map tasks to user stories for traceability.
- Keep changes incremental and independently testable per story.
