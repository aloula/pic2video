# Research: MP3 Audio and Video Fades

**Phase**: 0 - Outline & Research  
**Branch**: `005-mp3-audio-fades`  
**Date**: 2026-04-10

## Decision 1: MP3 discovery scope

**Decision**: Discover only `.mp3` files from the render input directory and sort by normalized filename ascending.

**Rationale**: Matches requirement scope exactly and guarantees deterministic output ordering.

**Alternatives considered**:
- Support all audio formats in v1: rejected to keep scope small and avoid codec/permutation complexity.
- Preserve filesystem order: rejected due to non-determinism across OS/filesystems.

## Decision 2: Optional audio behavior

**Decision**: If no MP3 files are found, render using existing video-only behavior with no extra flags required.

**Rationale**: Preserves backward compatibility and avoids breaking existing scripts.

**Alternatives considered**:
- Fail when no MP3 exists: rejected as incompatible with current silent slideshow workflows.

## Decision 3: Duration bounding with audio present

**Decision**: Keep output duration bound to slideshow timeline and trim audio to final video length when audio exceeds it.

**Rationale**: Enforces predictable output length and aligns with slideshow-first behavior.

**Alternatives considered**:
- Extend video duration to fit audio: rejected because it changes slideshow timing semantics.
- Loop images to match audio length: rejected as out of scope and higher complexity.

## Decision 4: Fade timing policy

**Decision**: Apply fade-in at `st=0` and fade-out at `total_duration - fade_duration`, using transition duration as baseline and clamping to half timeline for short outputs.

**Rationale**: Produces consistent visual polish while keeping filter timings valid.

**Alternatives considered**:
- Fixed fade duration independent of transition: rejected due to inconsistent operator expectations.

## Decision 5: Error handling for invalid MP3 assets

**Decision**: Treat unreadable/corrupt selected MP3 files as input-validation failure with clear classification.

**Rationale**: Keeps failures explicit and actionable for operators.

**Alternatives considered**:
- Silently skip invalid MP3 files: rejected because it can hide media issues and produce unexpected missing audio.
