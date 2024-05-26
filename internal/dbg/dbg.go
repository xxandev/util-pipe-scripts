package dbg

import (
	"io"
	"log"
)

var Log = &debugger{log.New(io.Discard, "[PIPE-DBG] ", 0)}

type debugger struct{ *log.Logger }

func (d *debugger) Trace(f string, v ...any) {
	d.Printf("\033[35m<trace> "+f+"\t\033[0m\n", v...)
}

func (d *debugger) Info(f string, v ...any) {
	d.Printf("\033[36m<info> "+f+"\t\033[0m\n", v...)
}

func (d *debugger) Success(f string, v ...any) {
	d.Printf("\033[32m<success> "+f+"\t\033[0m\n", v...)
}

func (d *debugger) Warn(f string, v ...any) {
	d.Printf("\033[33m<warning> "+f+"\t\033[0m\n", v...)
}

func (d *debugger) Critical(f string, v ...any) {
	d.Printf("\033[31m<critical> "+f+"\t\033[0m\n", v...)
}

func Experimental() {
	Log.Printf("\033[35m<trace> message\t\033[0m\n")
	Log.Printf("\033[34m<trace> message\t\033[0m\n")
	Log.Printf("\033[32m<success> info message\t\033[0m\n")
	Log.Printf("\033[36m<info> message\t\033[0m\n")
	Log.Printf("\033[33m<warning> message\t\033[0m\n")
	Log.Printf("\033[31m<critical> message\t\033[0m\n")
}
