package core

import "testing"

func TestLoopGuardian(t *testing.T) {
	lg := NewLoopGuardian()

	output := "Error: unexpected ';' on line 42"

	// First time
	isLoop, _ := lg.RecordFailure(output)
	if isLoop {
		t.Error("Detected loop on first failure")
	}

	// Second time
	isLoop, _ = lg.RecordFailure(output)
	if isLoop {
		t.Error("Detected loop on second failure")
	}

	// Third time
	isLoop, _ = lg.RecordFailure(output)
	if !isLoop {
		t.Error("Failed to detect loop on third identical failure")
	}

	// Different error
	isLoop, _ = lg.RecordFailure("New error")
	if isLoop {
		t.Error("Detected loop on a new, non-repeated error")
	}

	// Reset
	lg.Reset()
	isLoop, _ = lg.RecordFailure(output)
	if isLoop {
		t.Error("Detected loop after reset")
	}
}
