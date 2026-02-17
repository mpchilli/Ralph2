//go:build windows

package sandbox

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

// Windows Job Object constants
const (
	JobObjectExtendedLimitInformation = 9

	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE          = 0x2000
	JOB_OBJECT_LIMIT_JOB_MEMORY                 = 0x0200
	JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION = 0x0400
	JOB_OBJECT_LIMIT_ACTIVE_PROCESS             = 0x0008

	// Process access rights
	PROCESS_SET_QUOTA = 0x0100
	PROCESS_TERMINATE = 0x0001
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectW         = kernel32.NewProc("CreateJobObjectW")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
	procSetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
)

// JOBOBJECT_BASIC_LIMIT_INFORMATION
type basicLimitInformation struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

// IO_COUNTERS
type ioCounters struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

// JOBOBJECT_EXTENDED_LIMIT_INFORMATION
type extendedLimitInformation struct {
	BasicLimitInformation basicLimitInformation
	IoInfo                ioCounters
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

type JobObject struct {
	Handle syscall.Handle
}

// CreateJob creates a new Windows Job Object
func CreateJob() (*JobObject, error) {
	handle, _, err := procCreateJobObjectW.Call(0, 0)
	if handle == 0 {
		return nil, err
	}
	return &JobObject{Handle: syscall.Handle(handle)}, nil
}

// SetLimits sets memory and behavior limits for the job
func (job *JobObject) SetLimits(memoryLimitMB uint64) error {
	info := extendedLimitInformation{}
	info.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	if memoryLimitMB > 0 {
		info.BasicLimitInformation.LimitFlags |= JOB_OBJECT_LIMIT_JOB_MEMORY
		info.JobMemoryLimit = uintptr(memoryLimitMB * 1024 * 1024)
	}

	ret, _, err := procSetInformationJobObject.Call(
		uintptr(job.Handle),
		uintptr(JobObjectExtendedLimitInformation),
		uintptr(unsafe.Pointer(&info)),
		uintptr(unsafe.Sizeof(info)),
	)

	if ret == 0 {
		return fmt.Errorf("SetInformationJobObject failed: %w", err)
	}

	return nil
}

// Assign adds a process to the job object
func (job *JobObject) Assign(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}

	// We need a handle to the process. Since os.Process.Handle is unexported,
	// we open the process by PID.
	h, err := syscall.OpenProcess(PROCESS_SET_QUOTA|PROCESS_TERMINATE, false, uint32(cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("OpenProcess failed: %w", err)
	}
	defer syscall.CloseHandle(h)

	ret, _, err := procAssignProcessToJobObject.Call(
		uintptr(job.Handle),
		uintptr(h),
	)
	if ret == 0 {
		return fmt.Errorf("AssignProcessToJobObject failed: %w", err)
	}
	return nil
}

// Close closes the job object handle
func (job *JobObject) Close() error {
	return syscall.CloseHandle(job.Handle)
}
