# Implementation Plan: MP3 Audio and Video Fades

**Branch**: `005-mp3-audio-fades` | **Date**: 2026-04-10 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/005-mp3-audio-fades/spec.md`

## Summary

Extend slideshow rendering to optionally include MP3 files found in the input directory, sorted alphabetically, while preserving existing render behavior when no MP3 is present. Ensure generated videos consistently apply intro and outro fades, with robust timing for short timelines and deterministic output duration bounds.

## Technical Context

**Language/Version**: Go 1.23+  
**Primary Dependencies**: Standard library, `github.com/spf13/cobra`, existing ffmpeg/ffprobe integration  
**Storage**: N/A (filesystem I/O only)  
**Testing**: `go test ./...` (unit + e2e suites)  
**Target Platform**: Linux primary; macOS and Windows compatible CLI behavior  
**Project Type**: CLI  
**Performance Goals**: Preserve current render throughput envelope; keep startup status immediate; no material slowdown when no MP3 files exist  
**Constraints**: Backward compatibility for silent workflows; deterministic MP3 ordering; bounded output duration; no heavy new dependencies  
**Scale/Scope**: Focused changes in discovery, render command assembly, and regression tests/docs

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Scope is intentionally small and focused on one clear user outcome.
- [x] Design keeps complexity minimal; audio discovery and mapping remain isolated in existing layers.
- [x] Clean code approach is defined (small helpers, deterministic ordering, no duplication).
- [x] Unit test strategy is explicit for discovery, ordering, duration bounding, and fade timing.
- [x] E2E test strategy is explicit for image-only and image+MP3 primary journeys.
- [x] Documentation impact is explicit; README and feature quickstart/contract are updated.
- [x] Red-green-refactor sequence is planned before implementation tasks.

## Project Structure

### Documentation (this feature)

```text
specs/005-mp3-audio-fades/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/app/cli/
├── render_command.go
└── summary.go

internal/app/renderjob/
├── builder.go
├── model.go
└── service.go

internal/infra/fsio/
└── filesystem.go

internal/infra/ffmpeg/
└── command_builder.go

tests/unit/
├── ffmpeg_args_test.go
└── cli_validation_test.go

tests/e2e/
├── smoke_test.go
└── render_fhd_test.go
```

**Structure Decision**: Reuse current CLI -> renderjob -> ffmpeg pipeline; add MP3 discovery and argument wiring within existing packages without introducing new package boundaries.

## Phase 0: Research Findings

- Detect MP3 assets from the input directory in one scan flow and keep ordering deterministic via filename sort.
- Keep audio handling optional: no MP3 means existing behavior and command path remain valid.
- Use a single deterministic ffmpeg audio concatenation strategy for discovered MP3 files, then trim to slideshow timeline bounds.
- Keep global video fade behavior guaranteed on output and explicitly bounded for short durations.

## Phase 1: Design Artifacts

- [research.md](research.md) - Decisions, trade-offs, and alternatives
- [data-model.md](data-model.md) - Audio asset, timeline bounds, and fade profile entities
- [contracts/cli-contract.md](contracts/cli-contract.md) - Discovery and render-output behavior contract
- [quickstart.md](quickstart.md) - Operator workflows for image-only and image+MP3 runs

## Constitution Check (Post-Design)

- [x] No unresolved clarifications remain.
- [x] Unit and E2E test expectations are mapped to each user story.
- [x] Documentation update obligations are captured.
- [x] Scope remains incremental and independently testable.

## Complexity Tracking

No constitution violations requiring exceptions.
