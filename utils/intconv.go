package utils

// https://dev.to/chigbeef_77/bool-int-but-stupid-in-go-3jb3

func BoolToInt(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
