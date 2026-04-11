# Data Model: EXIF Footer Overlay Option

**Phase**: 1 - Design & Contracts  
**Branch**: `007-create-feature-branch`  
**Date**: 2026-04-11

## Entities

### ExifOverlayOptions

Represents user-provided render options controlling overlay behavior.

| Field | Type | Description |
|------|------|-------------|
| enabled | bool | Whether EXIF footer overlay is applied |
| font_size | int | Requested text size (36-60 inclusive) |
| footer_offset_px | int | Bottom offset in pixels (fixed value: 10) |
| text_color | string | Overlay text color (fixed value: white) |
| background_alpha | float | Background opacity scalar (> 0.5 transparency) |

Validation rules:
- `font_size` MUST be in range 36-60 inclusive when `enabled=true`.
- `footer_offset_px` is fixed to 10 for FHD and UHD.
- `text_color` is fixed to white.
- `background_alpha` MUST correspond to greater than 50% transparency.

### ExifFieldSet

Represents normalized EXIF metadata values per source image.

| Field | Type | Description |
|------|------|-------------|
| camera_model | string | Camera model value or `Unknown` |
| focal_distance | string | Focal distance value or `Unknown` |
| speed_fraction | string | Exposure speed formatted as `1/XXXXs` or `Unknown` |
| aperture | string | Aperture formatted as `f/X` or `Unknown` |
| iso | string | ISO value or `Unknown` |
| captured_date | string | Date formatted as `DD/MM/YYYY` or `Unknown` |

Validation rules:
- Every field MUST have a non-empty value after normalization.
- Missing or unparsable metadata MUST map to `Unknown`.
- `captured_date` MUST be formatted as `DD/MM/YYYY` when source date is valid.

### OverlayRenderLine

Represents the final display string composed for one image segment.

| Field | Type | Description |
|------|------|-------------|
| image_path | string | Source image path this line belongs to |
| text | string | Fully formatted overlay line |
| start_sec | float64 | Segment start time in final timeline |
| end_sec | float64 | Segment end time in final timeline |

Validation rules:
- `text` MUST follow exact field sequence and separators.
- `start_sec < end_sec` for all segments.
- Segment timing MUST align with existing image-duration and transition rules.

## Relationships

- One `RenderJob` has one `ExifOverlayOptions` object.
- One input image maps to one `ExifFieldSet` and one `OverlayRenderLine`.
- `OverlayRenderLine` collection is ordered identically to the final image timeline order.

## State Transitions

1. CLI parses overlay flags and validates option range.
2. Input image list is ordered by selected order mode.
3. EXIF fields are extracted and normalized into `ExifFieldSet` per image.
4. Formatter assembles deterministic `OverlayRenderLine.text` for each image.
5. ffmpeg command builder injects footer text filters into render graph.
6. Render executes with overlay enabled or bypasses overlay logic when disabled.
