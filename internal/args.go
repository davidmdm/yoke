package internal

import "slices"

func CutArgs(args []string) ([]string, []string) {
	idx := slices.Index(args, "--")
	if idx == -1 {
		return args, nil
	}
	return args[:idx], args[idx+1:]
}
