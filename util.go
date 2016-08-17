package gof

import "bytes"

func PortName(name [16]byte) string {
	b := name[:]
	if n := bytes.IndexByte(b, 0); n != -1 {
		b = b[:n]
	}
	return string(b)
}
