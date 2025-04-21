package util

import (
	"io"
	"log"
)

func Output(out *string, src io.ReadCloser, pid int) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n != 0 {
			s := string(buf[:n])
			*out = *out + s
			log.Printf("%d: %v", pid, s)
		}
		if err != nil {
			break
		}
	}
}
