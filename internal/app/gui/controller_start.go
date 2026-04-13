package gui

import (
	"context"
	"fmt"
	"strings"
)

func StartRun(ctx context.Context, runner *Runner, machine *RunStateMachine, cfg GuiRunConfiguration, onState func(RunStatus), onOutput OutputHandler) error {
	if err := machine.Begin(0); err != nil {
		return err
	}
	res := ValidatePreflight(cfg)
	if !res.OK {
		machine.FinishFailure(fmt.Errorf(strings.Join(res.Messages, "; ")))
		return fmt.Errorf(strings.Join(res.Messages, "; "))
	}
	if err := runner.Run(ctx, cfg, onState, onOutput); err != nil {
		machine.FinishFailure(err)
		return err
	}
	machine.FinishSuccess()
	return nil
}
