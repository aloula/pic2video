package unit

import (
	"os/exec"
	"testing"
)

func TestCLIValidationMissingRequiredFlags(t *testing.T) {
	cmd := newUnitCLICommand(t, "render")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing required flags")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if ee.ExitCode() != 2 {
			t.Fatalf("expected exit code 2, got %d", ee.ExitCode())
		}
	}
}
