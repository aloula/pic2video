package gui

import (
	"fmt"
	"strings"
	"time"
)

type LogStore struct {
	maxEntries int
	entries    []GuiLogEntry
	nextSeq    int
}

func NewLogStore(maxEntries int) *LogStore {
	if maxEntries <= 0 {
		maxEntries = 500
	}
	return &LogStore{maxEntries: maxEntries, entries: make([]GuiLogEntry, 0, maxEntries)}
}

func (s *LogStore) Append(stream, message string) {
	s.nextSeq++
	s.entries = append(s.entries, GuiLogEntry{
		Seq:       s.nextSeq,
		Timestamp: time.Now(),
		Stream:    stream,
		Message:   message,
	})
	if len(s.entries) > s.maxEntries {
		s.entries = s.entries[len(s.entries)-s.maxEntries:]
	}
}

func (s *LogStore) Entries() []GuiLogEntry {
	out := make([]GuiLogEntry, len(s.entries))
	copy(out, s.entries)
	return out
}

func (s *LogStore) Text() string {
	parts := make([]string, 0, len(s.entries))
	for _, e := range s.entries {
		ts := e.Timestamp.Format("15:04:05")
		parts = append(parts, fmt.Sprintf("[%s] [%s] %s", ts, e.Stream, e.Message))
	}
	return strings.Join(parts, "\n")
}
