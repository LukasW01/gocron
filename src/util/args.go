package util

import "os"

func SplitArgs() (flagArgs []string, execArgs []string) {
	for idx, arg := range os.Args {
		if arg == "--" {
			return os.Args[:idx], os.Args[idx+1:]
		}
	}
	return os.Args, nil
}
