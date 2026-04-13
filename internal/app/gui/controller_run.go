package gui

import "errors"

func ApplyRunnerStatus(machine *RunStateMachine, view *StatusView, status RunStatus) {
	switch status {
	case RunStatusLoadingFiles:
		view.SetStatus(RunStatusLoadingFiles)
		view.StartAnimation()
	case RunStatusProcessing:
		machine.SetProcessing()
		view.SetStatus(RunStatusProcessing)
		view.StartAnimation()
	case RunStatusFinished:
		machine.FinishSuccess()
		view.SetStatus(RunStatusFinished)
		view.StopAnimation()
	case RunStatusFailed:
		machine.FinishFailure(errors.New("run failed"))
		view.SetStatus(RunStatusFailed)
		view.StopAnimation()
	default:
		view.SetStatus(status)
		view.StopAnimation()
	}
}
