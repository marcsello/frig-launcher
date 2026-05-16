package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseResolution(res string) (int32, int32, error) {
	parts := strings.SplitN(res, "x", 2)

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format")
	}

	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	if w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("zero or negative resolution")
	}

	return int32(w), int32(h), nil
}
