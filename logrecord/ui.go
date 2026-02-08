package logrecord

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var screen tcell.Screen
var traceScroll int
var logScroll int
var traceViewportHeight int

func InitUI() error {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	if err := screen.Init(); err != nil {
		return err
	}

	screen.Clear()
	return nil
}

func PromptAndCapture(traceChan <-chan StackTrace, processDone <-chan struct{}, cancel context.CancelFunc) error {

	if err := InitUI(); err != nil {
		return err
	}
	defer screen.Fini()

	eventChan := make(chan tcell.Event)

	go func() {
		for {
			ev := screen.PollEvent()
			if ev == nil {
				close(eventChan)
				return
			}
			eventChan <- ev
		}
	}()

	RenderTcell("Process Running")

	isProcessRunning := true
	processStatus := "Process Running"

	for {
		select {

		case ev := <-eventChan:
			if ev == nil {
				return nil
			}
			switch e := ev.(type) {
			case *tcell.EventKey:
				if e.Key() == tcell.KeyCtrlC {
					if isProcessRunning {
						cancel() // Kill the process
						// Wait for process to exit
						<-processDone
					}
					return nil
				}
				msg := handleKeyEvent(e)
				if msg != "" {
					RenderTcell(msg + (func() string {
						if isProcessRunning {
							return " [Running]"
						}
						return " [Done]"
					})())
				} else {
					RenderTcell(processStatus)
				}
			case *tcell.EventResize:
				screen.Sync()
				RenderTcell(processStatus)
			}

		case trace := <-traceChan:
			AddTrace(trace)
			RenderTcell(processStatus)

		case <-LogUpdateChan:
			// Drain remaining updates to avoid excessive rendering
		drainLoop:
			for {
				select {
				case <-LogUpdateChan:
				default:
					break drainLoop
				}
			}
			RenderTcell(processStatus)

		case <-processDone:
			isProcessRunning = false
			processStatus = "Process Finished (Ctrl+C to Exit)"
			RenderTcell(processStatus)
		}
	}
}

func RenderTcell(status string) {

	screen.Clear()

	w, h := screen.Size()

	drawBox(0, 0, w, h/2, "Live Logs ["+status+"]")
	drawBox(0, h/2, w, h/2, "Captured Stack Traces")

	drawLiveLogs(1, 1, w-2, h/2-2)
	drawTraces(1, h/2+1, w-2, h/2-2)

	screen.Show()
}

func drawBox(x, y, w, h int, title string) {

	style := tcell.StyleDefault

	for i := x; i < x+w; i++ {
		screen.SetContent(i, y, '─', nil, style)
		screen.SetContent(i, y+h-1, '─', nil, style)
	}

	for i := y; i < y+h; i++ {
		screen.SetContent(x, i, '│', nil, style)
		screen.SetContent(x+w-1, i, '│', nil, style)
	}

	screen.SetContent(x, y, '┌', nil, style)
	screen.SetContent(x+w-1, y, '┐', nil, style)
	screen.SetContent(x, y+h-1, '└', nil, style)
	screen.SetContent(x+w-1, y+h-1, '┘', nil, style)

	drawTextClipped(x+2, y, w, title, style)
}

func drawLiveLogs(x, y, w, h int) {
	renderLock.Lock()
	defer renderLock.Unlock()

	start := 0
	if len(liveLogs) > h {
		start = len(liveLogs) - h
	}

	row := 0

	for i := start; i < len(liveLogs) && row < h; i++ {
		drawTextClipped(x, y+row, w, liveLogs[i], tcell.StyleDefault)
		row++
	}
}

func drawTraces(x, y, w, h int) {

	traceViewportHeight = h
	row := 0
	i := traceScroll

	for i < len(buffer) && row < h {

		e := buffer[i]

		cursorMark := " "
		if i == cursor {
			cursorMark = ">"
		}

		selectMark := " "
		if e.Selected {
			selectMark = "x"
		}

		header := cursorMark + " [" + selectMark + "] " + e.Trace.Lines[0]

		drawTextClipped(x, y+row, w, header, tcell.StyleDefault)
		row++

		for _, l := range e.Trace.Lines[1:] {
			if row >= h {
				break
			}
			drawTextClipped(x+4, y+row, w-4, l, tcell.StyleDefault)
			row++
		}

		row++
		i++
	}
}

func drawTextClipped(x, y, w int, text string, style tcell.Style) {

	runes := []rune(text)

	for i := 0; i < len(runes) && i < w; i++ {
		screen.SetContent(x+i, y, runes[i], nil, style)
	}
}

func handleKeyEvent(ev *tcell.EventKey) string {

	switch ev.Key() {

	case tcell.KeyUp:
		if cursor > 0 {
			cursor--
		}
		if cursor < traceScroll {
			traceScroll--
		}

	case tcell.KeyDown:
		if cursor < len(buffer)-1 {
			cursor++
		}
		if cursor >= traceScroll+traceViewportHeight {
			traceScroll++
		}

	case tcell.KeyEnter:
		return RecordSelected()
	}

	if ev.Rune() == ' ' && len(buffer) > 0 {
		buffer[cursor].Selected = !buffer[cursor].Selected
	}
	return ""
}

func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	return string(r[:w-3]) + "..."
}

func RecordSelected() string {

	f, _ := os.Create("debugo_record.txt")
	defer f.Close()

	count := 0
	for _, e := range buffer {
		if !e.Selected {
			continue
		}

		f.WriteString(strings.Join(e.Trace.Lines, "\n"))
		f.WriteString("\n\n---\n\n")
		count++
	}

	return fmt.Sprintf("Recorded %d traces -> debugo_record.txt", count)
}
