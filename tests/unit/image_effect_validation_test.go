package unit

import (
	"os/exec"
	"testing"
)

func TestCLIValidationImageEffectInvalid(t *testing.T) {
	cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory", "--image-effect", "wrong")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for invalid image-effect")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		if ee.ExitCode() != 2 {
			t.Fatalf("expected exit code 2 (ErrInvalidArguments), got %d", ee.ExitCode())
		}
	}
}

func TestCLIValidationImageEffectAllowedValues(t *testing.T) {
	for _, effect := range []string{"static", "kenburns-low", "kenburns-medium", "kenburns-high"} {
		t.Run(effect, func(t *testing.T) {
			cmd := newUnitCLICommand(t, "render", "--input", "/nonexistent/directory", "--image-effect", effect)
			err := cmd.Run()
			if err == nil {
				t.Fatal("expected non-zero exit due to missing input directory")
			}
			if ee, ok := err.(*exec.ExitError); ok {
				if ee.ExitCode() != 3 {
					t.Fatalf("expected exit code 3 for input validation, got %d", ee.ExitCode())
				}
			}
		})
	}
}
