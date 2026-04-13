# Data Model: Multiplatform Desktop GUI

## Entity: GuiRunConfiguration

- Purpose: User-selected render inputs/options collected from the GUI before execution.
- Fields:
  - input_folder: string
  - output_folder: string
  - profile: enum (`fhd`, `uhd`)
  - image_effect: enum (`static`, `kenburns-low`, `kenburns-medium`, `kenburns-high`)
  - exif_overlay_enabled: bool
  - exif_font_size: int
  - fps: int
  - order_mode: enum (`name`, `exif`, `explicit`)
  - order_file: string (optional)
  - overwrite: bool
  - ffmpeg_bin: string (optional)
  - ffprobe_bin: string (optional)

## Entity: GuiRunState

- Purpose: Runtime lifecycle state shown by the status indicator.
- Fields:
  - state: enum (`idle`, `loading_files`, `processing`, `finished`, `failed`)
  - started_at: timestamp (optional)
  - finished_at: timestamp (optional)
  - last_error: string (optional)
  - active_pid: int (optional)

## Entity: GuiLogEntry

- Purpose: One line of runtime output displayed in the log box.
- Fields:
  - seq: int
  - timestamp: timestamp
  - stream: enum (`stdout`, `stderr`, `system`)
  - message: string

## Entity: GuiValidationResult

- Purpose: Preflight validation outcome before start is allowed.
- Fields:
  - ok: bool
  - messages: []string
  - supported_media_count: int

## Relationships

- GuiRunConfiguration (1) -> (1) GuiValidationResult (per start attempt)
- GuiRunConfiguration (1) -> (1) GuiRunState (active run instance)
- GuiRunState (1) -> (0..N) GuiLogEntry

## Validation Rules

- input_folder MUST exist and contain at least one supported media file before run start.
- output_folder MUST exist or be creatable; unwritable destination MUST block start with actionable message.
- state transitions MUST follow: `idle -> loading_files -> processing -> finished|failed`.
- while `state` is `loading_files` or `processing`, new run starts MUST be rejected.
- GUI output destination MUST remain folder-only; file naming MUST stay profile-auto-generated.
