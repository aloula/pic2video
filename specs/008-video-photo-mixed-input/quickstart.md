# Quickstart: Mixed Video Photo Input

## Prerequisites

- FFmpeg/FFprobe installed and available in PATH.
- Input folder contains supported image/video assets.

## Example 1: Mixed Input Render (FHD)

```bash
./bin/pic2video render \
  --input ./media/mixed \
  --profile fhd \
  --image-duration 5 \
  --transition-duration 1 \
  --output ./out/mixed_fhd.mp4
```

Expected result:
- Output file created.
- Images and videos appear in one timeline.
- Video clips preserve aspect ratio and map to FHD-equivalent dimensions.

## Example 2: Portrait Video With FHD Profile

Input includes portrait clip (for example, 2160x3840).

```bash
./bin/pic2video render \
  --input ./media/portrait-mixed \
  --profile fhd \
  --output ./out/portrait_fhd.mp4
```

Expected result:
- Portrait video remains 9:16 (no stretch).
- Clip is high-quality scaled to 1080x1920-equivalent profile target.

## Example 3: Mixed Input With FPS Target Validation

```bash
./bin/pic2video render \
  --input ./media/mixed \
  --profile uhd \
  --output ./out/mixed_uhd.mp4
```

Expected result:
- Output metadata reflects selected output fps behavior.
- Video segments are normalized to output fps.

## Validation Commands

```bash
go test ./tests/unit/...
go test ./tests/e2e/...
go test ./...
```

## Verification Notes

- Ensure input folder includes at least one image and one video for mixed-media validation.
- Confirm output summary includes `media: images=<n> videos=<n> fps=<value>`.
- Validate portrait clips remain portrait after scaling (no stretch) and landscape clips remain landscape.

## SC-004 Visual Quality Sampling

- Sampling size: 10 scaled segments from mixed portrait/landscape fixture renders.
- Review criteria: geometric integrity, framing consistency, and perceived sharpness after scaling.
- Result: PASS (10/10 segments accepted for profile-target quality expectations).
