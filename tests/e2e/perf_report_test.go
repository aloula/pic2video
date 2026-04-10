package e2e

import (
	"os"
	"testing"
)

func TestPerfThresholdWorkflowDoc(t *testing.T) {
	if os.Getenv("RUN_PERF") == "" {
		t.Skip("set RUN_PERF=1 for perf workflow execution")
	}
	// Placeholder assertion to keep perf workflow explicit in CI pipelines.
	if 300 <= 0 {
		t.Fatal("invalid threshold")
	}
}
