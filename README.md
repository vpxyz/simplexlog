SIMPLEXLOG leveled logging library
======================

Simplexlog simple wrapper for the go standard log package, that adds logging levels.

The standard `log` package is enough (IMHO) for many applications, but sometimes
is usefull to have some logging levels.
Simplexlog adds some logging levels and the ability to use different
io.Writer for different level. All the standard logging flags can be used.

The logging levels are: TRACE, DEBUG, CRITICAL, ERROR, WARNING, NOTICE, INFO
or ALL.
Functions like Panic or Fatal are also available, their massage are redirect to CRITICAL.

Since simplexlog write to an io.Writer, you can provvide, for e.g, your custom io.Writer for writting name of function. Take a look at the examples.

Simplexlog is concurrent safe.

Installation
------------

The package is go gettable: go get -u github.com/vpxyz/simplexlog

Example
-------

basic usage:

``` go

package main

import (
    "github.com/vpxyz/simplexlog"
)

func main() {
    // by default use os.Stderr for error, critical, fatal and panic, and os.Stdout for others
    l := simplexlog.New()

    fmt.Printf("available levels: %s\n", l.LevelNames())

    // the defaul log level is Info
    l.Trace("Trace log")

    l.Info("Info log")

    l.Notice("Notice log")

    l.Warning("Warning log")

    l.Debug("Debug log")

    l.Error("Error log")

    l.Critical("Critical log")

    l.SwitchTo(sl.Warning)

    l.Info("This is hidden")

    // if you need, you can pass around an standard log.Logger, bypassing the LogLevel setting
    l.CriticalLogger().Print("test")

}

```

custom io.Writer, this approach can be useful even with the standard `log` package:


``` go

package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	sl "github.com/vpxyz/simplexlog"
)

var (
	logger *sl.Logger
)

type LogWriter struct{}

func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	fmt.Printf("%s: %s, line %d, %s\n", strings.ReplaceAll(string(p), "\n", ""), filepath.Base(file), line, fn.Name())
	return len(p), nil
}

func dumb() {
	logger.Debug("who I am ?")
}

func main() {
	// write func name for debugging
	debugLw := LogWriter{}
	logger = sl.New(sl.SetDebug(sl.Config{Out: debugLw, Label: fmt.Sprintf("%-9s", sl.LevelDebug), Flags: sl.DefaultLogFlags}))
	// logger = sl.New()
	logger.SwitchTo(sl.Debug)

	dumb()
}


```

weird  example:


``` go
package main

import (
    sl "github.com/vpxyz/simplexlog"
    "log"
    "os"
    "bytes"
)

func main() {
    // Set different tag for any level
    // If you need, you can use a different io.Writer for each level witch different flags and prefix
    // If you want color labels, simply put colors escape sequence around label. For e.g. "\x1b[20;32m"+sl.LevelInfo+"\x1b[0m"
    l := sl.New(
        sl.SetDebug(sl.Config{Out: os.Stdout, Label: sl.LevelInfo + " ==> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
        sl.SetTrace(sl.Config{Out: os.Stdout, Label: sl.LevelTrace + " ==> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
        sl.SetInfo(sl.Config{Out: os.Stdout, Label: sl.LevelInfo + " =>", Flags: sl.DefaultLogFlags}),
        sl.SetNotice(sl.Config{Out: os.Stdout, Label: fmt.Sprintf("%-10s", "["+sl.LevelNotice+"]:"), Flags: sl.DefaultLogFlags}),
        sl.SetWarning(sl.Config{Out: os.Stdout, Label: " 😒 ", Flags: sl.DefaultLogFlags}),
        sl.SetError(sl.Config{Out: os.Stderr, Label: " 🥲 " + " ", Flags: sl.DefaultLogFlags}),
        sl.SetCritical(sl.Config{Out: os.Stderr, Label: " 😡 ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
    )

    // enable all log level
    l.SwitchTo(sl.All)

    l.Tracef("Trace log %s", "!!!")

    l.Info("Info log")

    l.Notice("Notice log")

    l.Warning("Warning log")

    l.Debug("Debug log")

    l.Error("Error log")

    l.Critical("Critical log")

    // switch log level
    l.SwitchTo(sl.Warning)

    l.Info("This is hidden")

    // if you need, you can pass around an standard log.Logger, bypassing the LogLevel setting
    l.CriticalLogger().Print("test")

    // switch logging level using logging level name (case insensitive)
    l.SwitchTo("error")

    l.Infof("Info log %s", "🦇")

    l.Noticef("Notice log")

    l.Warningf("Warning log")

    l.Debugf("Debug log")

    l.Errorf("Error log")

    l.Criticalf("Critical log")

    // change the output of the Info level
    var buf bytes.Buffer
    l.SetOutput(sl.Info, &buf)

    // change label of Error level
    l.SetLabel(sl.Error, " 🐛 ")

    // change flags of Info level
    l.SetFlags(sl.Info, sl.DefaultLogFlags | log.Lshortfile)
}
```
