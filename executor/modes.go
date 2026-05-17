package executor

type LaunchMode uint8

const (
	// LaunchDetached launches a new subprocess, then detaches it, allowing the main process to exit
	// This makes tracking the child's lifecycle a bit complicated
	// AtExit is noop in this case
	LaunchDetached LaunchMode = iota

	// LaunchWait Launches a new subprocess, and waits for it to exit AtExit
	LaunchWait
	// LaunchExec does not launch the process immediately, but instead launches it AtExit by exec-ing into it.
	LaunchExec
)

const (
	LaunchDetachedStr = "detached"
	LaunchWaitStr     = "wait"
	LaunchExecStr     = "exec"
)

func (m LaunchMode) String() string {
	switch m {
	case LaunchDetached:
		return LaunchDetachedStr
	case LaunchWait:
		return LaunchWaitStr
	case LaunchExec:
		return LaunchExecStr
	}
	return ""
}

func ModeFromString(name string) (LaunchMode, bool) {
	switch name {
	case LaunchDetachedStr:
		return LaunchDetached, true
	case LaunchWaitStr:
		return LaunchWait, true
	case LaunchExecStr:
		return LaunchExec, true
	}
	return 0, false
}
