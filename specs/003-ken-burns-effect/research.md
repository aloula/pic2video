# Research: Ken Burns Effect Option

**Phase**: 0 - Outline & Research  
**Branch**: `003-ken-burns-effect`  
**Date**: 2026-04-09

## Decision 1: CLI option shape

**Decision**: Use one option `--image-effect` with values `static`, `kenburns-low`, `kenburns-medium`, `kenburns-high`.

**Rationale**: Explicit values are easier to validate, document, and test than multiple booleans.

**Alternatives considered**:
- Separate boolean `--ken-burns` plus strength flag: rejected due to invalid combinations and more validation complexity.

## Decision 2: Default behavior

**Decision**: Default effect is `static`.

**Rationale**: Preserves existing output behavior and avoids surprise motion for current users.

**Alternatives considered**:
- Default to low Ken Burns: rejected because it changes existing renders unexpectedly.

## Decision 3: Resolution-aware motion profile

**Decision**: Derive motion profile from output resolution (FHD/UHD) and selected intensity level.

**Rationale**: Motion scale and zoom progression must match output dimensions to look consistent.

**Alternatives considered**:
- Same constants for all resolutions: rejected because UHD would feel under-scaled or over-aggressive.

## Decision 4: Quality-first filter strategy

**Decision**: Prioritize smoothness and quality using high-quality scaling and smooth-motion settings even if slower.

**Rationale**: This directly matches the requirement to prefer quality over processing speed.

**Alternatives considered**:
- Speed-first settings: rejected due to visible quality loss and requirement mismatch.

## Decision 5: Observability

**Decision**: Include selected effect in startup status output.

**Rationale**: Operators can confirm selected motion behavior immediately before processing.

**Alternatives considered**:
- Do not expose effect in status output: rejected due to poorer operator feedback.
