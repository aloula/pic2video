package unit

import (
	"testing"

	"github.com/loula/pic2video/internal/app/gui"
)

func TestRunStateMachineTransitions(t *testing.T) {
	m := gui.NewRunStateMachine()
	if got := m.Snapshot().Status; got != gui.RunStatusIdle {
		t.Fatalf("expected idle initial status, got %s", got)
	}
	if err := m.Begin(42); err != nil {
		t.Fatalf("expected begin success, got %v", err)
	}
	m.SetProcessing()
	if got := m.Snapshot().Status; got != gui.RunStatusProcessing {
		t.Fatalf("expected processing status, got %s", got)
	}
	m.FinishSuccess()
	if got := m.Snapshot().Status; got != gui.RunStatusFinished {
		t.Fatalf("expected finished status, got %s", got)
	}
}

func TestRunStateMachineBlocksConcurrentStart(t *testing.T) {
	m := gui.NewRunStateMachine()
	if err := m.Begin(1); err != nil {
		t.Fatalf("expected first begin success, got %v", err)
	}
	if err := m.Begin(2); err == nil {
		t.Fatal("expected second begin to fail while active run is in progress")
	}
}
