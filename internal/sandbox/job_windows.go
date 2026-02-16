//go:build windows

package sandbox

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectW    = kernel32.NewProc("CreateJobObjectW")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
	procSetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
)

type JobObject struct {
	Handle syscall.Handle
}

func CreateJob() (*JobObject, error) {
	handle, _, err := procCreateJobObjectW.Call(0, 0)
	if handle == 0 {
		return nil, err
	}
	return &JobObject{Handle: syscall.Handle(handle)}, nil
}

func (job *JobObject) Assign(cmd *exec.Cmd) error {
	// Must start process suspended or assign immediately after start.
	// For simplicity in this skeleton, we assume the process is started suspended
	// or we accept a race condition for MVP (Note: Real security needs 'CREATE_SUSPENDED')
	
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	
	ret, _, err := procAssignProcessToJobObject.Call(
		uintptr(job.Handle),
		uintptr(cmd.Process.Handle),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func (job *JobObject) Close() error {
	return syscall.CloseHandle(job.Handle)
}
