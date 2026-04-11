# Fixture Media Requirements

- Use small placeholder images for unit/e2e runtime speed.
- Include at least one mixed image+video fixture set (`.jpg` + `.mp4`) for mixed-input discovery and timeline tests.
- Use deterministic naming (`img-001.jpg`, etc.) for ordering tests.
- Keep one mixed-aspect fixture set for framing and quality-warning validation.
- Include portrait and landscape video samples to validate aspect-ratio-preserving scaling behavior.
- Include clips with different source frame rates for output fps normalization tests.
- For audio-aware e2e tests, include deterministic MP3 names (`ambient_a.mp3`, `ambient_b.mp3`).
- Add one unsupported audio file (for example `ignored.wav`) in mixed media scenarios to validate ignore behavior.
