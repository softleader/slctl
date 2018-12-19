package verbose

import (
	"fmt"
	"github.com/softleader/slctl/pkg/environment"
	"io"
	"os"
)

// fmt.Fprintf only if environment.Settings.Verbose true
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if environment.Settings.Verbose {
		return fmt.Fprintf(w, format, a...)
	}
	return
}

// fmt.Fprintf only if environment.Settings.Verbose true
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

// fmt.Fprint only if environment.Settings.Verbose true
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	if environment.Settings.Verbose {
		return fmt.Fprint(w, a...)
	}
	return
}

// fmt.Fprint only if x true
func Print(a ...interface{}) (n int, err error) {
	return Fprint(os.Stdout, a...)
}

// fmt.Fprintln only if environment.Settings.Verbose true
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	if environment.Settings.Verbose {
		return fmt.Fprintln(w, a...)
	}
	return
}

// fmt.Fprintln only if environment.Settings.Verbose true
func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}
