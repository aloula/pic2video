# Quickstart: Professional Photo-to-Video CLI

## Prerequisites
- Go 1.23+
- FFmpeg and FFprobe available in PATH
- Optional: NVIDIA GPU + NVENC-capable FFmpeg build for accelerated encoding

## 1. Build
```bash
go build -o bin/pic2video ./cmd/pic2video
```

## 2. Render FHD video
```bash
./bin/pic2video render \
  --input ./examples/photoset-a \
  --output ./out/slideshow-fhd.mp4 \
  --profile fhd \
  --image-duration 4 \
  --transition-duration 1
```

## 3. Render UHD video
```bash
./bin/pic2video render \
  --input ./examples/photoset-a \
  --output ./out/slideshow-uhd.mp4 \
  --profile uhd \
  --image-duration 5 \
  --transition-duration 1
```

## 4. Force CPU encoding fallback
```bash
./bin/pic2video render \
  --input ./examples/photoset-a \
  --output ./out/slideshow-cpu.mp4 \
  --profile fhd \
  --encoder cpu
```

## 5. Render with explicit ordering
```bash
./bin/pic2video render \
  --input ./examples/photoset-a \
  --output ./out/slideshow-ordered.mp4 \
  --profile fhd \
  --order explicit \
  --order-file ./examples/orders/sequence.txt
```

## 6. Run Unit tests
```bash
go test ./... -run Test -count=1
```

## 7. Run E2E tests
```bash
go test ./tests/e2e -count=1
```

## 8. Run performance gate (SC-004)
```bash
RUN_PERF=1 go test ./tests/e2e -run TestPerf -count=1
```

## Expected behavior
- Command exits `0` on successful render and prints summary with profile, encoder, elapsed time, and output path.
- Command exits non-zero with actionable message when validation fails.
- If NVENC is available and allowed, encoder selection reports NVIDIA path; otherwise CPU fallback is used.

## Expected summary example
```text
status=success profile=fhd resolution=1920x1080 encoder:auto->nvenc processed=50 elapsed=18.231s output=./out/slideshow-fhd.mp4 warnings=0
```

## Troubleshooting
- Exit code `2`: validate required flags and value ranges.
- Exit code `3`: inspect input directory, media formats, and explicit order file.
- Exit code `4`: verify FFmpeg/FFprobe installation and NVENC availability.
- Exit code `5`: inspect FFmpeg runtime errors and output path permissions.

## Validation policy for invalid media (FR-012)
- Unsupported file types are skipped during image discovery.
- If fewer than 2 valid images remain, render fails with exit code `3`.
- Invalid/missing explicit-order manifest entries fail with exit code `3`.

## Quality rubric acceptance mapping (SC-003)
- Geometric integrity: deterministic scale/pad framing policy (unit: `tests/unit/framing_policy_test.go`).
- Framing consistency: one framing strategy per run (e2e: `tests/e2e/render_mixed_aspect_test.go`).
- Sharpness preservation: warnings for sub-profile sources (unit: `tests/unit/quality_warning_test.go`).
- Transition smoothness: deterministic cross-fade timing offsets (unit: `tests/unit/timeline_test.go`).
