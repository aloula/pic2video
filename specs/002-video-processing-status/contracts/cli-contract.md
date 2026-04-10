# CLI Output Contract: Video Processing Status

**Feature**: `002-video-processing-status`
**Date**: 2026-04-09
**File**: Extends existing CLI contract from `001-build-photo-video-cli`

---

## Stdout Contract

### 1. Pre-Render Announcement Line

Emitted once to stdout **before** FFmpeg encoding begins, after input assets are discovered and ordered.

**Format**:
```
status=starting files=<N> format=<FORMAT>
details: input=<DIR> output=<PATH> profile=<PROFILE> effect=<EFFECT> encoder=<ENCODER> overwrite=<BOOL>
timing: image-duration=<SECONDS>s transition-duration=<SECONDS>s
order: mode=<ORDER_MODE> order-file=<ORDER_FILE|->
```

**Fields**:
| Field | Type | Description |
|-------|------|-------------|
| `status` | string | Always `starting` |
| `files` | int | Number of input image files after ordering |
| `format` | string | Uppercased container label derived from `--output` extension (e.g., `MP4`, `MOV`). `UNKNOWN` if no extension. |
| `details` | section | Selected runtime options for the current run (`input`, `output`, `profile`, `effect`, `encoder`, `overwrite`) |
| `timing` | section | Timing options (`image-duration`, `transition-duration`) |
| `order` | section | Ordering options (`mode`, `order-file`; `-` when not applicable) |

**Example**:
```
status=starting files=20 format=MP4
details: input=./photos output=slideshow_uhd.mp4 profile=uhd effect=kenburns-medium encoder=auto overwrite=true
timing: image-duration=5.0s transition-duration=1.0s
order: mode=exif order-file=-
```

**Conditions**:
- Emitted only when `files > 0` (zero-file input fails validation before this point)
- Emitted only when all input validation passes (flags, paths, ordering)
- Stream: stdout

---

### 2. Completion Report Line (extended)

Emitted once to stdout after **successful** render completion. Extends the existing completion line format.

**Format**:
```
status=success
result: profile=<PROFILE> resolution=<WxH> <encoder_report> processed=<N> files=<N>
output: format=<FORMAT> elapsed=<ELAPSED> output=<PATH> warnings=<N>
```

**Fields** (⚡ = new in this feature):
| Field | Type | Description | Change |
|-------|------|-------------|--------|
| `status` | string | Always `success` | unchanged |
| `profile` | string | Profile name (`fhd`, `uhd`) | unchanged |
| `resolution` | string | Effective output resolution (e.g., `1920x1080`) | unchanged |
| _(encoder_report)_ | string | NVENC availability and selection inline fragment | unchanged |
| `processed` | int | Asset count (retained for backward compat) | unchanged |
| `files` | int | ⚡ Asset count (semantically labeled as "input files") | **new** |
| `format` | string | ⚡ Output container label (e.g., `MP4`, `MOV`, `UNKNOWN`) | **new** |
| `elapsed` | string | ⚡ Human-readable processing time (see below) | **changed format** |
| `output` | string | Absolute path to output video file | unchanged |
| `warnings` | int | Count of accumulated quality warnings | unchanged |

**Elapsed format rules**:
| Duration | Format | Examples |
|----------|--------|----------|
| `< 1.0s` | `< 1s` | `< 1s` |
| `1.0s – 59.9s` | `%.1fs` | `3.4s`, `45.3s` |
| `>= 60s` | `Xm Ys` | `1m 0s`, `1m 30s`, `10m 5s` |

**Example** (FHD render, 12 files, 45.3s):
```
status=success
result: profile=fhd resolution=1920x1080 encoder=cpu processed=12 files=12
output: format=MP4 elapsed=45.3s output=/tmp/slideshow.mp4 warnings=0
```

**Example** (long render, 90s):
```
status=success
result: profile=uhd resolution=3840x2160 encoder=nvenc processed=30 files=30
output: format=MOV elapsed=1m 30s output=/tmp/slideshow.mov warnings=2
```

**Conditions**:
- Emitted **only** when render exits with code `0` (success)
- NOT emitted on failure; failure uses stderr + non-zero exit code (unchanged behavior)
- Stream: stdout

---

## Stderr Contract (unchanged)

No changes to stderr output. Error messages and classification codes remain as defined in the `001-build-photo-video-cli` contract.

---

## Exit Codes (unchanged)

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Invalid arguments |
| `2` | Input validation error |
| `3` | Environment error (FFmpeg missing) |
| `4` | Execution error (FFmpeg failed) |

---

## Backward Compatibility

- The `processed=` field is **retained** in the completion report so existing scripts relying on it continue working.
- The `elapsed=` field **key is unchanged**; only the value format changes from raw decimal (e.g., `45.321s`) to human-readable (e.g., `45.3s`). Scripts that parse the key-value pair will receive a different value format — this is an intentional improvement covered by the spec.
- All other existing fields are unchanged.
