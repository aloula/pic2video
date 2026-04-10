# Implementation Plan: Professional Photo-to-Video CLI

**Branch**: `[001-build-photo-video-cli]` | **Date**: 2026-04-09 | **Spec**: [/home/loula/src/pic2video/specs/001-build-photo-video-cli/spec.md](/home/loula/src/pic2video/specs/001-build-photo-video-cli/spec.md)
**Input**: Feature specification from `/home/loula/src/pic2video/specs/001-build-photo-video-cli/spec.md`

## Summary

Build a small, high-quality Go CLI that converts ordered photos into a YouTube-ready 16:9
video (FHD or UHD) with cross-fade transitions. The architecture keeps Go focused on
validation, ordering, orchestration, and reporting while delegating encode/filter execution
to FFmpeg/FFprobe for mature media processing. Encoding policy prefers NVENC when available
and falls back to CPU for deterministic portability.

## Technical Context

**Language/Version**: Go 1.23+  
**Primary Dependencies**: Standard library, `github.com/spf13/cobra` (CLI), `github.com/spf13/pflag` (flags), external FFmpeg/FFprobe binaries  
**Storage**: N/A (filesystem input/output only)  
**Testing**: `go test` for unit + E2E suites  
**Target Platform**: Linux hosts (first-class), optional NVIDIA GPU acceleration  
**Project Type**: CLI application  
**Performance Goals**: FHD render for 50 images in <= 5 minutes on baseline host (SC-004)  
**Constraints**: 16:9 output only (1920x1080 or 3840x2160), deterministic failures/exit codes, high-quality framing and transitions  
**Scale/Scope**: Single-user local/CI execution, batches of 2-500 photos, no distributed rendering in v1

## Constitution Check (Pre-Design Gate)

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- PASS: Scope remains intentionally small and focused on slideshow generation.
- PASS: Complexity is controlled by using FFmpeg as rendering backend and Go as orchestrator.
- PASS: Clean code boundaries are explicit across `internal/domain`, `internal/app`, and `internal/infra`.
- PASS: Unit test strategy is explicit for profile parsing, ordering, timeline, validation, and encoder policy.
- PASS: E2E test strategy is explicit for happy paths and failure-class exits.
- PASS: Documentation impact is addressed, with root README required and present.
- PASS: Red-green-refactor sequencing is reflected by tests-first task organization.

## Phase 0 Research

All technical clarifications are resolved and recorded in:

- `/home/loula/src/pic2video/specs/001-build-photo-video-cli/research.md`

No unresolved `NEEDS CLARIFICATION` items remain.

## Phase 1 Design Artifacts

Artifacts are defined and available at:

- Data model: `/home/loula/src/pic2video/specs/001-build-photo-video-cli/data-model.md`
- Contract: `/home/loula/src/pic2video/specs/001-build-photo-video-cli/contracts/cli-contract.md`
- Quickstart: `/home/loula/src/pic2video/specs/001-build-photo-video-cli/quickstart.md`
- Root README (constitution requirement): `/home/loula/src/pic2video/README.md`

## Project Structure

### Documentation (this feature)

```text
specs/001-build-photo-video-cli/
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
cmd/
└── pic2video/
    └── main.go

internal/
├── app/
│   ├── cli/
│   ├── pipeline/
│   └── renderjob/
├── domain/
│   ├── media/
│   ├── profile/
│   └── transition/
└── infra/
    ├── ffmpeg/
    ├── fsio/
    └── nvenc/

tests/
├── e2e/
├── fixtures/
└── unit/
```

**Structure Decision**: Single-project CLI layout is retained to satisfy simplicity, clean
boundaries, and testability goals while keeping integration with FFmpeg isolated in infra adapters.

## Constitution Check (Post-Design Recheck)

- PASS: Design stays within v1 scope (image slideshow with cross-fades only).
- PASS: Dependencies are limited and justified (`cobra`/`pflag` + FFmpeg/FFprobe).
- PASS: Unit and E2E quality gates are represented in spec/tasks and validated by test suites.
- PASS: README requirement is satisfied with build/run/test/troubleshooting guidance.
- PASS: No constitutional violations requiring waiver were introduced.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
