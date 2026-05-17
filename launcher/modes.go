package launcher

type LaunchMode uint8

const (
	launchInvalid LaunchMode = iota
	// LaunchDetached launches a new subprocess, then detaches it, allowing the main process to exit
	// This makes tracking the child's lifecycle a bit complicated
	// AtExit is noop in this case
	LaunchDetached
	// LaunchWait Launches a new subprocess, and waits for it to exit AtExit
	LaunchWait
	// LaunchExec does not launch the process immediately, but instead launches it AtExit by exec-ing into it.
	LaunchExec
	// LaunchNone does not actually launch anything, it's a noop
	LaunchNone
)

const (
	LaunchDetachedStr = "detached"
	LaunchWaitStr     = "wait"
	LaunchExecStr     = "exec"
	LaunchNoneStr     = "none"
)

func (m LaunchMode) String() string {
	switch m {
	case LaunchDetached:
		return LaunchDetachedStr
	case LaunchWait:
		return LaunchWaitStr
	case LaunchExec:
		return LaunchExecStr
	case LaunchNone:
		return LaunchNoneStr
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
	case LaunchNoneStr:
		return LaunchNone, true
	}
	return 0, false
}
