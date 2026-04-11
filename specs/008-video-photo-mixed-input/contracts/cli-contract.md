# CLI Contract: Mixed Video Photo Input

## Scope

Extends existing `pic2video render` behavior to support mixed image + video folders.

## Command

- Command: `pic2video render`
- Existing flags remain valid.
- Mixed-media support is automatic when supported video files are present in `--input`.
- Output frame-rate selection is provided via `--fps`.

## FPS Selection Contract

- Flag: `--fps <value>`
- Valid values: integer range 24 to 60
- Default behavior: when omitted, renderer uses profile default fps
- Validation: out-of-range or non-numeric values MUST fail with input validation classification before ffmpeg execution

## Behavioral Contract

1. Input Discovery
- The command MUST discover supported image and video assets from `--input`.
- Unsupported/unreadable media MUST produce warnings or classified input errors per policy.

2. Timeline Construction
- The command MUST produce one unified timeline containing images and videos.
- Existing ordering modes MUST apply consistently to mixed media.

3. Video Transform
- Video segments MUST preserve aspect ratio.
- Video segments MUST be scaled to profile-equivalent quality target.
- Video scaling MUST use high-quality resampling intent.

4. Frame Rate
- Video segments in output MUST align to selected output frame rate.

5. Errors
- If no valid renderable assets are found, command MUST fail with input-validation classification.

## Observable Output

Startup and summary output MUST continue to report:
- selected profile
- selected frame rate target
- processed file counts
- warnings count

## Compatibility

- Existing image-only runs MUST remain behaviorally unchanged.
- Existing EXIF overlay behavior for images MUST remain available in mixed runs (applied where relevant).
