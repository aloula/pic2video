# Fixture Media Requirements

- Use small placeholder images for unit/e2e runtime speed.
- Use deterministic naming (`img-001.jpg`, etc.) for ordering tests.
- Keep one mixed-aspect fixture set for framing and quality-warning validation.
- For audio-aware e2e tests, include deterministic MP3 names (`ambient_a.mp3`, `ambient_b.mp3`).
- Add one unsupported audio file (for example `ignored.wav`) in mixed media scenarios to validate ignore behavior.
