# Research: Video Processing Status

**Phase**: 0 — Outline & Research
**Branch**: `002-video-processing-status`
**Date**: 2026-04-09

---

## FR-007 Clarification Resolution

**Question**: Should an unrecognized container extension be silently accepted (emit raw uppercase label) or rejected with a validation error before encoding begins?

**Decision**: Accept silently; emit raw uppercased extension label (e.g., `.avi` → `AVI`). Emit `UNKNOWN` when the output path has no extension at all.

**Rationale**: The spec's own Edge Cases section already answers this: "What happens when the output file has an unrecognized extension (e.g., `.avi`)? The format field shows the raw uppercased extension (`AVI`)." Adding a rejection path would require maintaining an allow-list of container extensions, which is out of scope and conflicts with the spec's stated behavior. The existing validator already handles invalid output paths via other rules (e.g., directory not writable, overwrite guard).

**Alternatives considered**: Reject with a validation error (option B) — rejected because it requires a maintained allow-list, couples the status formatter to encoding knowledge, and contradicts the spec edge case text.

---

## Duration Formatting

**Question**: How to implement human-readable duration formatting in Go without new dependencies?

**Decision**: Use standard library `math` package and simple integer arithmetic. No external libraries needed.

**Rationale**:
- For `< 1s`: `seconds < 1.0`
- For `< 60s`: `fmt.Sprintf("%.1fs", seconds)` — one decimal place
- For `>= 60s`: `fmt.Sprintf("%dm %ds", int(seconds)/60, int(seconds)%60)`

Rule table:
| Condition | Format | Example |
|-----------|--------|---------|
| `seconds < 1.0` | `< 1s` | `< 1s` |
| `seconds < 60.0` | `%.1fs` | `45.3s` |
| `seconds >= 60.0` | `Xm Ys` | `1m 30s` |

**Alternatives considered**: `time.Duration.String()` — rejected because it produces output like `1m30.001234567s` which is not operator-friendly and uses nanosecond precision not appropriate here.

---

## Output Format Label Derivation

**Question**: How to derive a clean container label from the output file path extension?

**Decision**: Use `path/filepath.Ext(outputPath)` (standard library), strip the leading `.`, and uppercase. If the result is empty, return `UNKNOWN`.

**Rationale**: `filepath.Ext` handles edge cases (no extension returns `""`; multiple dots returns only the last extension). `strings.ToUpper(strings.TrimPrefix(ext, "."))` is a two-step, zero-dependency solution.

**Alternatives considered**: Manual string splitting — rejected because `filepath.Ext` handles OS path edge cases correctly and is already used in the codebase.

---

## Pre-Render Announcement Injection Point

**Question**: Where in the call chain should the pre-render announcement be emitted?

**Decision**: Emit from `render_command.go`, after `pipeline.ApplyOrder()` resolves the final asset list and before `service.Run()` is called. This is the earliest point where the final asset count is known and the render has not yet started.

**Rationale**: The announcement requires two facts: (1) final input file count (after ordering/filtering) and (2) output container format (from `--output` flag). Both are available in `render_command.go` after `ApplyOrder` and before `service.Run`. Pushing this into `Service.Run` would mix I/O concerns into the service layer, violating separation of concerns.

**Code location**: `internal/app/cli/render_command.go`, between the `ApplyOrder` call and the `BuildJob` call.

---

## Backward Compatibility of Completion Summary Line

**Question**: Can the `FormatSummary` output format be extended without breaking existing scripts?

**Decision**: Extend the existing key=value line by appending new fields (`files=N`, `format=EXT`) and replacing the existing `elapsed=` field value with a human-readable string. The key names and order of existing fields are preserved.

**Rationale**: The current output is: `status=success profile=fhd resolution=1920x1080 encoder=cpu processed=12 elapsed=45.321s output=/tmp/out.mp4 warnings=0`

New output adds `files=` (same value as `processed=` — retained for backward compat) and `format=`, and changes `elapsed=` from raw decimal to human-readable. Since `processed=` is retained, any script relying on it continues to work. The `elapsed=` format change is a minor behavioral change covered by the spec.

---

## Files Modified

| File | Change |
|------|--------|
| `internal/app/cli/summary.go` | Add `FormatOutputFormat()`, `FormatElapsed()`, `FormatAnnouncement()`; extend `FormatSummary()` |
| `internal/app/cli/render_command.go` | Add pre-render announcement print before `service.Run()` |
| `tests/unit/status_format_test.go` | New file — unit tests for all formatter functions |
| `tests/e2e/smoke_test.go` | Extend smoke test to assert `files=`, `format=`, human-readable elapsed |

No new dependencies. No schema changes. No new packages.
