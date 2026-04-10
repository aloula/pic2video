# Research: Professional Photo-to-Video CLI

## Decision 1: Orchestration in Go, rendering in FFmpeg
- Decision: Implement command orchestration, validation, timeline assembly, and backend selection in Go; delegate encoding and transition filtering to FFmpeg.
- Rationale: FFmpeg is mature, stable, and high-performance for codec/filter pipelines, while Go provides maintainable CLI logic and robust process control.
- Alternatives considered:
  - Pure-Go video processing libraries only: rejected due to weaker maturity/performance for production-grade transitions and codec support.
  - Full custom renderer pipeline: rejected as too complex for small-scope v1 and violates simplicity principle.

## Decision 2: Stable Go modules only for CLI/runtime glue
- Decision: Use `cobra`/`pflag` for CLI UX plus Go standard library for filesystem, process execution, and validation.
- Rationale: These modules are stable, widely adopted, and minimize dependency risk while preserving clean architecture.
- Alternatives considered:
  - Heavy framework stack: rejected to keep scope and maintenance cost low.
  - Hand-rolled argument parsing only: rejected due to poorer ergonomics and discoverability for CLI users.

## Decision 3: Prefer NVENC when available, deterministic fallback otherwise
- Decision: Detect NVIDIA encoder availability (via FFmpeg encoder probing and optional host checks) and use NVENC path first; if unavailable, use CPU encoder with quality-preserving defaults.
- Rationale: Meets explicit requirement for fast processing on capable hosts while preserving cross-host portability.
- Alternatives considered:
  - CPU-only encoding: rejected because it does not satisfy performance priority when GPU is available.
  - NVENC-only strict mode: rejected because it would fail on non-NVIDIA hosts and hurt usability.

## Decision 4: Output profile model limited to FHD/UHD 16:9
- Decision: Define two explicit output profiles only: FHD (1920x1080) and UHD (3840x2160), both strict 16:9.
- Rationale: Keeps product simple and aligned with YouTube publication goals.
- Alternatives considered:
  - Arbitrary custom dimensions: rejected as out-of-scope for first increment.
  - Additional profile presets (vertical/social): rejected as future work.

## Decision 5: Quality-first scaling and framing policy
- Decision: Apply deterministic, documented framing/scaling policy for mixed aspect ratios with warnings on quality risk scenarios.
- Rationale: Professional source photos require predictable high-quality output and transparent behavior.
- Alternatives considered:
  - Implicit automatic behavior without warnings: rejected due to poor operator trust and diagnosability.
  - User-defined advanced filter graph in v1: rejected as complexity beyond small scope.

## Decision 6: Test strategy with mandatory Unit + E2E
- Decision: Require unit tests for core business logic and E2E tests invoking the compiled CLI with fixture photos for FHD/UHD and error cases.
- Rationale: Aligns with constitution and ensures both local correctness and real workflow reliability.
- Alternatives considered:
  - Unit-only strategy: rejected because integration regressions in FFmpeg invocation would go undetected.
  - E2E-only strategy: rejected because it is slower and weaker for precise logic validation.
