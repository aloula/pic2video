# Data Model: Video Processing Status

**Phase**: 1 — Design & Contracts
**Branch**: `002-video-processing-status`
**Date**: 2026-04-09

---

## Entities

### OutputFormat

A normalized container label derived from the output file path extension.

| Field | Type | Description |
|-------|------|-------------|
| label | string | Uppercased extension without leading dot, e.g. `MP4`, `MOV`. Falls back to `UNKNOWN` when path has no extension. |

**Derivation rules**:
1. Extract extension with `filepath.Ext(outputPath)` → e.g. `.mp4`
2. Strip leading `.` → `mp4`
3. Uppercase → `MP4`
4. If result is empty string → `UNKNOWN`

**No storage.** Pure computed value from the output path string.

---

### HumanDuration

A human-readable representation of processing time in seconds.

| Condition | Format | Example |
|-----------|--------|---------|
| `seconds < 1.0` | `< 1s` | `< 1s` |
| `1.0 <= seconds < 60.0` | `%.1fs` | `45.3s` |
| `seconds >= 60.0` | `Xm Ys` | `1m 30s` |

**No storage.** Pure computed value from `RenderSummary.ElapsedSeconds`.

---

### ProcessingAnnouncement

Emitted to stdout before encoding begins. Contains pre-render discovery facts.

| Field | Source | Example |
|-------|--------|---------|
| file_count | `len(assets)` after `ApplyOrder()` in `render_command.go` | `20` |
| output_format | `FormatOutputFormat(outputPath)` | `MP4` |

**Output line format**: `status=starting files=20 format=MP4`

---

### ProcessingReport (extends existing completion line)

Emitted to stdout after successful render. Extends the existing `FormatSummary` output.

| Field | Key | Source | Change? |
|-------|-----|--------|---------|
| status | `status` | hardcoded `success` | No change |
| profile | `profile` | `RenderSummary.ProfileName` | No change |
| resolution | `resolution` | `RenderSummary.EffectiveResolution` | No change |
| encoder report | _(inline)_ | `nvenc.BuildReport(...)` | No change |
| files processed | `files` | `RenderSummary.ProcessedAssets` | **New field** |
| output format | `format` | `FormatOutputFormat(RenderSummary.OutputPath)` | **New field** |
| elapsed | `elapsed` | `FormatElapsed(RenderSummary.ElapsedSeconds)` | **Changed** (raw decimal → human-readable) |
| output path | `output` | `RenderSummary.OutputPath` | No change |
| warning count | `warnings` | `len(RenderSummary.Warnings)` | No change |

**Output line example** (before):
```
status=success profile=fhd resolution=1920x1080 encoder=cpu processed=12 elapsed=45.321s output=/tmp/out.mp4 warnings=0
```

**Output line example** (after):
```
status=success profile=fhd resolution=1920x1080 encoder=cpu processed=12 files=12 format=MP4 elapsed=45.3s output=/tmp/out.mp4 warnings=0
```

---

## State Transitions

```
[command invoked]
      |
      v
[assets discovered + ordered]  →  emit ProcessingAnnouncement
      |
      v
[FFmpeg encoding]
      |             \
      v              v
[success]         [failure]
      |              |
emit             no report;
ProcessingReport  error exits
```

---

## Validation Rules

- `OutputFormat.label` is never empty (guaranteed by `UNKNOWN` fallback)
- `HumanDuration` is never negative (elapsed is computed as `time.Now().Sub(started)` which is non-negative in normal flow)
- `ProcessingAnnouncement` is only emitted when `len(assets) > 0` (assets=0 triggers existing validation error before announcement point)
