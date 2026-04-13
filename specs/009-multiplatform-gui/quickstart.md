# Quickstart: Multiplatform Desktop GUI

## Goal

Run a render end-to-end from the new GUI on Windows, Linux, or macOS.

## Prerequisites

- FFmpeg and FFprobe available on PATH (or configured via GUI options if exposed)
- A folder with supported media files

## Scenario 1: Basic Run With Defaults

1. Launch the GUI from a terminal in a project/media folder.
2. Verify input and output fields default to the launch directory.
3. Keep default profile/options and click Start.
4. Observe status transitions: `idle -> loading files -> processing -> finished`.
5. Confirm output file appears in selected output folder with profile-based auto filename.

## Scenario 2: Configure Options and Run

1. Select a different input folder containing supported media.
2. Select a destination output folder.
3. Set options such as profile, image effect, EXIF overlay, and FPS.
4. Start render and verify log box shows runtime output lines.
5. Confirm run completes and output reflects selected options.

## Scenario 3: Preflight Validation

1. Select an empty input folder with no supported media.
2. Attempt to start.
3. Verify GUI blocks start and shows actionable validation message.
4. Add media or choose a valid folder, then rerun successfully.

## Scenario 4: Failure Visibility

1. Configure an invalid FFmpeg binary path (if option exposed) or induce a runtime failure.
2. Start render.
3. Verify status becomes `failed` and log box contains error output.
