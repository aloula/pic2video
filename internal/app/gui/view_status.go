package gui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusView struct {
	StatusLabel *widget.Label
	ErrorLabel  *widget.Label
	Progress    *widget.ProgressBarInfinite
	Container   *fyne.Container

	mu          sync.Mutex
	baseStatus  RunStatus
	startedAt   time.Time
	lastElapsed time.Duration
	ticker      *time.Ticker
	stopCh      chan struct{}
	timing      bool
}

func NewStatusView() *StatusView {
	statusLabel := widget.NewLabel(string(RunStatusIdle))
	statusLabel.Alignment = fyne.TextAlignCenter
	errLabel := widget.NewLabel("")
	errLabel.Wrapping = fyne.TextWrapWord
	progress := widget.NewProgressBarInfinite()
	progress.Stop()
	progress.Hide()
	bar := container.NewStack(progress, container.NewCenter(statusLabel))
	c := container.NewVBox(bar, errLabel)
	v := &StatusView{
		StatusLabel: statusLabel,
		ErrorLabel:  errLabel,
		Progress:    progress,
		Container:   c,
		baseStatus:  RunStatusIdle,
	}
	v.updateStatusText(0)
	return v
}

func (v *StatusView) SetStatus(s RunStatus) {
	v.mu.Lock()
	v.baseStatus = s
	elapsed := v.currentElapsedLocked()
	v.mu.Unlock()
	v.updateStatusText(elapsed)
}

func (v *StatusView) SetError(msg string) {
	v.ErrorLabel.SetText(msg)
}

func (v *StatusView) StartAnimation() {
	v.mu.Lock()
	if v.timing {
		v.mu.Unlock()
		v.Progress.Show()
		v.Progress.Start()
		return
	}
	v.timing = true
	v.startedAt = time.Now()
	v.lastElapsed = 0
	v.stopCh = make(chan struct{})
	v.ticker = time.NewTicker(time.Second)
	elapsed := v.currentElapsedLocked()
	stopCh := v.stopCh
	ticker := v.ticker
	v.mu.Unlock()

	v.updateStatusText(elapsed)
	v.Progress.Show()
	v.Progress.Start()

	go func() {
		for {
			select {
			case <-ticker.C:
				v.mu.Lock()
				elapsed := v.currentElapsedLocked()
				v.mu.Unlock()
				v.updateStatusText(elapsed)
			case <-stopCh:
				return
			}
		}
	}()
}

func (v *StatusView) StopAnimation() {
	v.mu.Lock()
	if v.timing {
		v.lastElapsed = time.Since(v.startedAt)
		v.timing = false
		if v.ticker != nil {
			v.ticker.Stop()
			v.ticker = nil
		}
		if v.stopCh != nil {
			close(v.stopCh)
			v.stopCh = nil
		}
	}
	elapsed := v.currentElapsedLocked()
	v.mu.Unlock()

	v.Progress.Stop()
	v.Progress.Hide()
	v.updateStatusText(elapsed)
}

func (v *StatusView) currentElapsedLocked() time.Duration {
	if v.timing {
		return time.Since(v.startedAt)
	}
	if v.baseStatus == RunStatusFinished || v.baseStatus == RunStatusFailed {
		return v.lastElapsed
	}
	return 0
}

func (v *StatusView) updateStatusText(elapsed time.Duration) {
	v.mu.Lock()
	status := v.baseStatus
	v.mu.Unlock()

	if elapsed <= 0 {
		v.StatusLabel.SetText(string(status))
		return
	}
	v.StatusLabel.SetText(fmt.Sprintf("%s %s", status, formatElapsed(elapsed)))
}

func formatElapsed(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Seconds())
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
