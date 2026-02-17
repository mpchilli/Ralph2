package core

import (
	"crypto/sha256"
	"fmt"
)

// LoopGuardian detects infinite loops in development tasks
type LoopGuardian struct {
	MaxRepeatedFailures int
	FailureCounts       map[string]int
}

// NewLoopGuardian creates a new loop guardian
func NewLoopGuardian() *LoopGuardian {
	return &LoopGuardian{
		MaxRepeatedFailures: 3,
		FailureCounts:       make(map[string]int),
	}
}

// RecordFailure hashes the output and returns true if a loop is detected
func (lg *LoopGuardian) RecordFailure(output string) (bool, error) {
	if output == "" {
		return false, nil
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(output)))
	
	lg.FailureCounts[hash]++
	
	count := lg.FailureCounts[hash]
	if count >= lg.MaxRepeatedFailures {
		return true, nil
	}
	
	return false, nil
}

// Reset clears the guardian state
func (lg *LoopGuardian) Reset() {
	lg.FailureCounts = make(map[string]int)
}
