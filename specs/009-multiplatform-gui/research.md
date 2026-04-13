# Research: Multiplatform Desktop GUI

## Decision 1: GUI Toolkit for Windows/Linux/macOS

- Decision: Use `fyne.io/fyne/v2` as the desktop UI toolkit for this feature.
- Rationale: Fyne is Go-native, cross-platform, and sufficient for a clean single-window form + status + log UX without introducing a browser stack.
- Alternatives considered:
  - Wails (webview): rejected for v1 due to higher setup complexity and larger frontend surface area.
  - GTK bindings: rejected due to packaging/runtime friction across target OSes.

## Decision 2: Execution Integration Strategy

- Decision: GUI launches rendering by invoking the existing CLI entrypoint as a child process and streaming stdout/stderr into the GUI log panel.
- Rationale: Reuses proven render workflow and validation behavior with minimal refactor risk; naturally satisfies log-box requirement.
- Alternatives considered:
  - In-process direct service invocation: rejected for v1 because current service has no structured progress event stream and would require additional API changes.

## Decision 3: Status Lifecycle Mapping

- Decision: Standardize GUI run status to `idle`, `loading files`, `processing`, `finished`, `failed`.
- Rationale: Exactly matches spec language and can be mapped from lifecycle points without deep pipeline instrumentation.
- Alternatives considered:
  - Add granular percentage progress: rejected for v1 due to uncertain reliable metric across mixed media.

## Decision 4: Render Option Exposure Model

- Decision: Expose all current user-facing render options in the GUI, with core controls visible by default and additional options grouped in a compact section.
- Rationale: Satisfies clarified requirement while keeping interface clean.
- Alternatives considered:
  - Only core options: rejected by clarification outcome.
  - Flat list of all controls: rejected for readability/usability concerns.

## Decision 5: Input/Output Path Defaults and Validation

- Decision: Initialize input and output folder fields from process working directory at app startup; block run start if input contains no supported media or output path is invalid/unwritable.
- Rationale: Matches clarified requirement and prevents avoidable long-running failures.
- Alternatives considered:
  - Allow start and fail later: rejected by clarification outcome.

## Decision 6: Output File Naming in GUI

- Decision: GUI output selection is destination folder only; output filename remains profile-based auto name from existing render behavior.
- Rationale: Matches clarification and preserves current contract (`slideshow_fhd.mp4` / `slideshow_uhd.mp4`) while giving users control of destination.
- Alternatives considered:
  - Fully user-editable filename: rejected per clarification.
  - Optional override in v1: deferred to future enhancement.
