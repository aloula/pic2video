package gui

import (
	"fmt"
	"sync"
	"time"
)

type RunStateMachine struct {
	mu    sync.Mutex
	state GuiRunState
}

func NewRunStateMachine() *RunStateMachine {
	return &RunStateMachine{state: GuiRunState{Status: RunStatusIdle}}
}

func (m *RunStateMachine) Snapshot() GuiRunState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

func (m *RunStateMachine) Begin(pid int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state.Status == RunStatusLoadingFiles || m.state.Status == RunStatusProcessing {
		return fmt.Errorf("a render run is already active")
	}
	m.state = GuiRunState{
		Status:    RunStatusLoadingFiles,
		StartedAt: time.Now(),
		ActivePID: pid,
	}
	return nil
}

func (m *RunStateMachine) SetProcessing() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state.Status == RunStatusLoadingFiles || m.state.Status == RunStatusProcessing {
		m.state.Status = RunStatusProcessing
	}
}

func (m *RunStateMachine) FinishSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.Status = RunStatusFinished
	m.state.FinishedAt = time.Now()
	m.state.ActivePID = 0
}

func (m *RunStateMachine) FinishFailure(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.Status = RunStatusFailed
	m.state.FinishedAt = time.Now()
	m.state.ActivePID = 0
	if err != nil {
		m.state.LastError = err.Error()
	}
}

func (m *RunStateMachine) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = GuiRunState{Status: RunStatusIdle}
}
