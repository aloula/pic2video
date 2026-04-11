# Implementation Plan: Mixed Video Photo Input

**Branch**: `008-video-photo-mixed-input` | **Date**: 2026-04-11 | **Spec**: `/specs/008-video-photo-mixed-input/spec.md`
**Input**: Feature specification from `/specs/008-video-photo-mixed-input/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Enable rendering folders that contain both images and videos in one timeline while preserving video aspect ratio, scaling each clip to profile-equivalent quality targets, and normalizing all video segments to the selected output frame rate. The approach extends existing media discovery, timeline building, and ffmpeg filter graph construction with explicit video-segment transform policy and frame-rate alignment.

## Technical Context

**Language/Version**: Go 1.23+  
**Primary Dependencies**: Go standard library, `github.com/spf13/cobra`, external FFmpeg/FFprobe binaries  
**Storage**: N/A (filesystem-based media processing)  
**Testing**: `go test ./tests/unit/...`, `go test ./tests/e2e/...`, `go test ./...`  
**Target Platform**: Linux/macOS/Windows CLI environments with FFmpeg available
**Project Type**: CLI media processing tool  
**Performance Goals**: Maintain smooth output playback at selected profile FPS; avoid visible scaling artifacts for portrait and landscape clips  
**Constraints**: Preserve input video aspect ratio, no stretching, profile-equivalent scaling, mixed-media ordering compatibility, deterministic behavior in tests  
**Scale/Scope**: Single-command render for folders with dozens to low hundreds of mixed media assets

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Scope is intentionally small and focused on one clear user outcome.
- Design keeps complexity minimal; only existing ffmpeg graph paths are extended, avoiding new services.
- Clean code approach is defined: keep media typing in domain/app layers, keep ffmpeg transforms isolated in infra.
- Unit test strategy is explicit for media detection, transform policy, frame-rate normalization, and ordering.
- E2E test strategy is explicit for mixed image+video rendering, aspect-ratio preservation, and FPS conformity.
- Documentation impact is explicit: README must include mixed-input usage and constraints.
- Red-green-refactor sequence is planned before implementation tasks.

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

```text
cmd/pic2video/
internal/app/cli/
internal/app/pipeline/
internal/app/renderjob/
internal/domain/media/
internal/domain/profile/
internal/infra/ffmpeg/
internal/infra/fsio/
tests/unit/
tests/e2e/
specs/008-video-photo-mixed-input/
```

**Structure Decision**: Keep the existing single CLI project structure. Implement mixed-media detection in fsio/media domain, timeline behavior in app pipeline/renderjob, and scaling/FPS normalization in ffmpeg command builder paths; validate in unit and e2e suites under existing test folders.

## Phase 0: Research

- Decide supported video extensions for v1 mixed-input discovery and failure policy.
- Define scaling policy that preserves aspect ratio and maps to profile-equivalent dimensions for portrait/landscape clips.
- Define frame-rate normalization approach for clips with fixed/variable/missing fps metadata.
- Define transition handling for short video clips.

Output artifact: `/specs/008-video-photo-mixed-input/research.md`

## Phase 1: Design & Contracts

- Model mixed media asset and normalized timeline segment states.
- Design transform policy contract for profile-target scaling and fps normalization.
- Define CLI contract updates for mixed media behavior and output expectations.
- Draft quickstart scenarios for mixed inputs (FHD/UHD, portrait/landscape, fps selection).

Output artifacts:
- `/specs/008-video-photo-mixed-input/data-model.md`
- `/specs/008-video-photo-mixed-input/contracts/cli-contract.md`
- `/specs/008-video-photo-mixed-input/quickstart.md`

Post-design constitution check status: PASS

## Complexity Tracking
| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
