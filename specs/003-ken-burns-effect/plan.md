# Implementation Plan: Ken Burns Effect Option

**Branch**: `003-ken-burns-effect` | **Date**: 2026-04-09 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-ken-burns-effect/spec.md`

## Summary

Add a render flag to select static images (default) or Ken Burns motion intensity (`low`, `medium`, `high`). Integrate effect selection across CLI, render job model, and ffmpeg filter graph generation. Motion behavior must be resolution-aware for FHD/UHD and explicitly quality-first even if render time increases.

## Technical Context

**Language/Version**: Go 1.23+  
**Primary Dependencies**: Standard library, `github.com/spf13/cobra`, existing ffmpeg/ffprobe integration  
**Storage**: N/A (filesystem I/O only)  
**Testing**: `go test ./...` (unit + e2e suites)  
**Target Platform**: Linux primary; macOS and Windows compatible CLI behavior
**Project Type**: CLI  
**Performance Goals**: Smooth motion quality over speed; startup status remains immediate  
**Constraints**: Preserve static behavior by default; no breaking CLI changes; no heavy new dependencies  
**Scale/Scope**: Focused change in CLI flags + pipeline filter generation + status output/tests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Scope is intentionally small and focused on one clear user outcome.
- [x] Design keeps complexity minimal; non-trivial filter logic remains isolated in pipeline functions.
- [x] Clean code approach is defined (single-responsibility option parsing, filter generation helpers).
- [x] Unit test strategy is explicit for effect validation and filter generation.
- [x] E2E test strategy is explicit for startup output and successful render flow.
- [x] Documentation impact is explicit; README is updated for new CLI option.
- [x] Red-green-refactor sequence is planned before implementation tasks.

## Project Structure

### Documentation (this feature)

```text
specs/003-ken-burns-effect/
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

internal/app/pipeline/
├── framing.go
└── xfade.go

internal/app/renderjob/
├── model.go
├── builder.go
└── service.go

internal/infra/ffmpeg/
└── command_builder.go

tests/unit/
├── framing_policy_test.go
├── status_format_test.go
└── image_effect_validation_test.go

tests/e2e/
├── smoke_test.go
├── render_fhd_test.go
└── render_uhd_test.go
```

**Structure Decision**: Reuse existing CLI and pipeline layers; no new packages required.

## Phase 0: Research Findings

- Use a single string flag (`--image-effect`) with four explicit modes.
- Keep `static` as default for full backward compatibility.
- For Ken Burns, use quality-first ffmpeg primitives (`lanczos`, higher motion fps) to prioritize visual smoothness.
- Make motion resolution-aware by deriving effect parameters from target output dimensions.

## Phase 1: Design Artifacts

- [research.md](research.md) - Decisions and trade-offs
- [data-model.md](data-model.md) - Effect mode and motion profile entities
- [contracts/cli-contract.md](contracts/cli-contract.md) - CLI flag and startup output contract updates
- [quickstart.md](quickstart.md) - Operator usage examples

## Constitution Check (Post-Design)

- [x] No unresolved clarifications remain.
- [x] Test strategy covers changed behavior (unit and e2e).
- [x] Documentation updates are included for changed CLI usage.
- [x] Scope remains incremental and independently testable.

## Complexity Tracking

No constitution violations requiring exceptions.
