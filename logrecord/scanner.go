package logrecord

import (
	"bufio"
	"io"
)

var liveLogs []string
var LogUpdateChan = make(chan struct{}, 100)


const maxLogs = 20

func AddLiveLog(line string) {
	renderLock.Lock()
	defer renderLock.Unlock()

	liveLogs = append(liveLogs, line)

	if len(liveLogs) > maxLogs {
		liveLogs = liveLogs[1:]
	}

	select {
	case LogUpdateChan <- struct{}{}:
	default:
	}
}

func ScanStream(r io.Reader, isErr bool, traceChan chan<- StackTrace) {
	scanner := bufio.NewScanner(r)

	var current []string
	capturing := false

	for scanner.Scan() {
		line := scanner.Text()


		// Print live output
		if isErr {
			AddLiveLog(line)
		} else {
			AddLiveLog(line)
		}

		// Trace start
		if IsTraceStart(line) {
			if capturing && len(current) > 0 {
				traceChan <- StackTrace{Lines: current}
			}
			capturing = true
			current = []string{line}
			continue
		}

		// Continuation
		if capturing && IsTraceContinuation(line) {
			current = append(current, line)
			continue
		}

		// End trace
		if capturing {
			traceChan <- StackTrace{Lines: current}
			capturing = false
			current = nil
		}
	}

	// Flush last trace
	if capturing && len(current) > 0 {
		traceChan <- StackTrace{Lines: current}
	}
}
