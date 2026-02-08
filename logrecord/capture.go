package logrecord

import "sync"

var renderLock sync.Mutex

type CaptureEntry struct {
	Trace    StackTrace
	Selected bool
}

var buffer []CaptureEntry
var cursor int

func AddTrace(t StackTrace) {
	buffer = append(buffer, CaptureEntry{Trace: t})
}
