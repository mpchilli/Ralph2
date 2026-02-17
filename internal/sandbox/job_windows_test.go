//go:build windows

package sandbox

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestJobObject(t *testing.T) {
	job, err := CreateJob()
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}
	defer job.Close()

	err = job.SetLimits(512) // 512MB
	if err != nil {
		t.Fatalf("Failed to set limits: %v", err)
	}

	// Start a long-lived command
	cmd := exec.Command("powershell", "-Command", "Start-Sleep -Seconds 10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	err = job.Assign(cmd)
	if err != nil {
		t.Fatalf("Failed to assign process to job: %v", err)
	}

	t.Logf("Process %d assigned to job", cmd.Process.Pid)
}

func TestKillOnClose(t *testing.T) {
	job, err := CreateJob()
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	err = job.SetLimits(512)
	if err != nil {
		t.Fatalf("Failed to set limits: %v", err)
	}

	// Start a command that stays alive
	cmd := exec.Command("cmd", "/c", "pause")
	// Use stdin/stdout pipe to ensure it's interactive enough to stay open
	_ , _ = cmd.StdinPipe()
	
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}
	pid := cmd.Process.Pid

	err = job.Assign(cmd)
	if err != nil {
		t.Fatalf("Failed to assign process: %v", err)
	}

	// Closing the job should kill the process because of JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	err = job.Close()
	if err != nil {
		t.Fatalf("Failed to close job: %v", err)
	}

	// Give Windows a moment to sweep
	time.Sleep(500 * time.Millisecond)

	// Check if process still exists
	existsCmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-Process -Id %d -ErrorAction SilentlyContinue", pid))
	out, _ := existsCmd.CombinedOutput()
	
	if len(out) > 0 {
		t.Errorf("Process %d still exists after job close", pid)
	} else {
		t.Logf("Process %d successfully terminated by job object", pid)
	}
}
