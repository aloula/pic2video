package unit

import (
	"os/exec"
	"testing"
)

func TestCLIValidationMissingRequiredFlags(t *testing.T) {
	// With optional flags having sensible defaults, this test now validates
	// that the command fails gracefully when input directory doesn't exist
	cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for non-existent input directory")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if ee.ExitCode() != 3 {
			t.Fatalf("expected exit code 3 (ErrInputValidation), got %d", ee.ExitCode())
		}
	}
}
