package ngamux

import (
	"log"
	"runtime"
)

func Recovery() {
	if r := recover(); r != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		log.Printf("panic running process: %v\n%s\n", r, buf)
	}

}
