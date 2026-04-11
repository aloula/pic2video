# Implementation Plan: EXIF Footer Overlay Option

**Branch**: `007-create-feature-branch` | **Date**: 2026-04-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/007-exif-overlay-option/spec.md`

## Summary

Add an optional EXIF footer overlay to rendered videos, with strict content format, fixed footer offset (10 px from bottom) for both FHD and UHD, white text on a background more than 50% transparent, and validated font size range 36-60. Implement by extending CLI/render-job options, normalizing per-image EXIF metadata, and appending deterministic ffmpeg overlay filters to the existing render graph.

## Technical Context

**Language/Version**: Go 1.23+  
**Primary Dependencies**: Standard library, `github.com/spf13/cobra`, existing ffmpeg/ffprobe integration  
**Storage**: N/A (filesystem I/O only)  
**Testing**: `go test ./tests/unit/...`, `go test ./tests/e2e/...`, `go test ./...`  
**Target Platform**: Linux primary; macOS and Windows compatible CLI behavior  
**Project Type**: CLI  
**Performance Goals**: Preserve 60 fps output behavior and avoid material slowdown beyond per-frame text overlay costs  
**Constraints**: Maintain backward compatibility when overlay is disabled; no heavy new dependencies; deterministic formatting for missing metadata values  
**Scale/Scope**: Focused changes in CLI validation, metadata extraction/normalization, ffmpeg filter assembly, tests, and README usage documentation

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Scope is intentionally small and focused on one clear user outcome.
- [x] Design keeps complexity minimal; overlay logic extends existing pipeline with isolated helpers.
- [x] Clean code approach is defined (small, testable formatter/validator/filter helpers; no duplicated string assembly).
- [x] Unit test strategy is explicit for validation, metadata formatting, and ffmpeg argument generation.
- [x] E2E test strategy is explicit for overlay-enabled render journey in FHD/UHD.
- [x] Documentation impact is explicit; README and quickstart updates are required due to new CLI options.
- [x] Red-green-refactor sequence is planned before implementation tasks.

## Project Structure

### Documentation (this feature)

```text
specs/007-exif-overlay-option/
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
└── render_command.go

internal/app/renderjob/
├── builder.go
├── errors.go
├── model.go
└── service.go

internal/infra/fsio/
└── exif.go

internal/infra/ffmpeg/
└── command_builder.go

tests/unit/
├── cli_validation_test.go
├── ffmpeg_args_test.go
└── image_effect_validation_test.go

tests/e2e/
├── render_fhd_test.go
└── render_uhd_test.go
```

**Structure Decision**: Reuse the existing CLI -> renderjob -> infra layering and add EXIF overlay functionality through targeted extensions in existing files, avoiding new package boundaries.

## Phase 0: Research Findings

- [research.md](research.md) defines rendering and formatting decisions for EXIF field normalization, date formatting, transparent background policy, and fallback behavior for missing metadata.
- All previously ambiguous items in technical context are resolved; no `NEEDS CLARIFICATION` markers remain.

## Phase 1: Design Artifacts

- [data-model.md](data-model.md) - Overlay options, normalized EXIF fields, and render line entity definitions.
- [contracts/cli-contract.md](contracts/cli-contract.md) - CLI flag behavior, validation, and output-rendering guarantees.
- [quickstart.md](quickstart.md) - Operator examples for enabling overlay and selecting font sizes.

## Constitution Check (Post-Design)

- [x] No unresolved clarifications remain.
- [x] Unit and E2E test expectations map to each user story.
- [x] Documentation update obligations are captured.
- [x] Scope remains incremental and independently testable.

## Complexity Tracking

No constitution violations requiring exceptions.
