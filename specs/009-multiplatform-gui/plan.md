# Implementation Plan: Multiplatform Desktop GUI

**Branch**: `010-create-feature-branch` | **Date**: 2026-04-11 | **Spec**: `/specs/009-multiplatform-gui/spec.md`
**Input**: Feature specification from `/specs/009-multiplatform-gui/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Deliver a clean, cross-platform desktop GUI (Windows/Linux/macOS) that lets users select input/output folders, configure all user-facing render options, start rendering, monitor lifecycle status (`loading files`, `processing`, `finished/failed`), and inspect runtime logs. The approach uses a Go-native desktop toolkit and reuses the existing CLI render workflow via process orchestration to minimize integration risk and preserve current rendering behavior.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.23+  
**Primary Dependencies**: Go standard library, `fyne.io/fyne/v2` (desktop UI), existing `github.com/spf13/cobra` render CLI stack, external FFmpeg/FFprobe binaries  
**Storage**: N/A (filesystem I/O only)  
**Testing**: `go test ./tests/unit/...`, `go test ./tests/e2e/...`, `go test ./...`  
**Target Platform**: Windows, Linux, macOS desktop environments
**Project Type**: Desktop application + existing CLI application  
**Performance Goals**: GUI starts in under 2 seconds on developer machines; status updates visible within 300ms of lifecycle changes; log panel appends output continuously during runs  
**Constraints**: Must preserve current render behavior; output selection remains folder-only with profile auto filename; block concurrent runs; block preflight on empty input media set  
**Scale/Scope**: Single-window GUI for local render workflows; dozens to low hundreds of media assets per run

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Scope is intentionally small and focused on one clear user outcome.
- Design keeps complexity minimal by reusing existing render command behavior instead of rebuilding render logic in GUI layer.
- Clean code approach is defined: isolate GUI state/orchestration from process execution and from domain render logic.
- Unit test strategy is explicit for config mapping, preflight validation, status transitions, and log buffering.
- E2E test strategy is explicit for primary GUI journeys (configure + start + status/log visibility).
- Documentation impact is explicit; README must gain GUI launch/build/run instructions.
- Red-green-refactor sequence is planned: tests for validation/state/logging before implementation of handlers.

Initial gate status: PASS

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
cmd/pic2video/
cmd/pic2video-gui/
internal/app/cli/
internal/app/gui/
internal/app/renderjob/
internal/domain/
internal/infra/ffmpeg/
internal/infra/fsio/
tests/unit/
tests/e2e/
specs/009-multiplatform-gui/
```

**Structure Decision**: Keep monorepo single-project layout and add a dedicated GUI entrypoint plus focused `internal/app/gui` orchestration package. Existing render internals remain authoritative; GUI composes options, validates preflight, launches run, and streams logs/status.

## Phase 0: Research

Resolved technical decisions in `/specs/009-multiplatform-gui/research.md`:

- Selected Fyne as cross-platform Go desktop toolkit.
- Chosen integration strategy: invoke existing CLI render workflow as child process.
- Defined status lifecycle mapping and log-box behavior.
- Confirmed folder-only output selection with profile-based auto filename.
- Defined preflight-block behavior for empty input and invalid destinations.

Output artifact: `/specs/009-multiplatform-gui/research.md`

## Phase 1: Design & Contracts

- Modeled GUI runtime entities and validation/state rules in data model.
- Defined GUI behavioral contract for defaults, validation, runtime status, and logs.
- Authored quickstart execution scenarios for success and failure paths.

Output artifacts:
- `/specs/009-multiplatform-gui/data-model.md`
- `/specs/009-multiplatform-gui/contracts/gui-contract.md`
- `/specs/009-multiplatform-gui/quickstart.md`

Post-design constitution check status: PASS

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
