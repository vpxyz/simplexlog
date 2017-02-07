SIMPLEXLOG log library
======================

Simplexlog simple wrapper for the go standard log package, that adds log level.

The standard log package in enough (imho) for many applications, but sometimes
is usefull to have some log level.
Simplexlog is a wrapper that adds some log level, the ability to use different
io.Writer for different level.

The log levels are: TRACE, DEBUG, CRITICAL, ERROR, WARNING, NOTICE, INFO

Simplexlog is concurrent safe.

INSTALLATION
------------

The package is go gettable:  go get -u github.com/vpxyz/simplexlog

EXAMPLE
-------

``` go

package main

import (
	"github.com/vpxyz/simplexlog"
	"log"
	"os"
)

func main() {
    // by default use os.Stderr for error, critical, fatal and panic, and os.Stdout for others
    l := simplexlog.New() 

    // the defaul log level is Info
	l.TracePrint("Trace log")

	l.InfoPrint("Info log")

	l.NoticePrint("Notice log")

	l.WarningPrint("Warning log")

	l.DebugPrint("Debug log")

	l.ErrorPrint("Error log")

	l.CriticalPrint("Critical log")
    
    l.SetLevel(sl.Warning)
    
    l.InfoPrint("This is hidden")

	// if you need, you can pass around an standard log.Logger, bypassing the LogLevel setting
	l.CriticalLogger().Print("test")

}

```

More "complex" example

``` go
package main

import (
	sl "github.com/vpxyz/simplexlog"
	"log"
	"os"
)

func main() {
	// Set different tag for any level
	// If you need, you can use a different io.Writer for each level witch different flags and prefix
	l := sl.New(
		sl.SetDebug(sl.Option{Out: os.Stdout, Label: sl.LabelDebug + "==> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
		sl.SetTrace(sl.Option{Out: os.Stdout, Label: sl.LabelTrace, Flags: sl.DefaultLogFlags | log.Lshortfile}),
		sl.SetInfo(sl.Option{Out: os.Stdout, Label: sl.LabelInfo, Flags: sl.DefaultLogFlags}),
		sl.SetNotice(sl.Option{Out: os.Stdout, Label: sl.LabelNotice, Flags: sl.DefaultLogFlags}),
		sl.SetWarning(sl.Option{Out: os.Stdout, Label: sl.LabelWarning + " ==> ", Flags: sl.DefaultLogFlags}),
		sl.SetError(sl.Option{Out: os.Stderr, Label: sl.LabelError, Flags: sl.DefaultLogFlags}),
		sl.SetCritical(sl.Option{Out: os.Stderr, Label: sl.LabelCritical + ",GULP! ==> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
	)

    // print all log
	l.SetLevel(sl.All)

	l.TracePrintf("Trace log %s", "!!!")

	l.InfoPrint("Info log")

	l.NoticePrint("Notice log")

	l.WarningPrint("Warning log")

	l.DebugPrint("Debug log")

	l.ErrorPrint("Error log")

	l.CriticalPrint("Critical log")
    
    l.SetLevel(sl.Warning)
    
    l.InfoPrint("This is hidden")

	// if you need, you can pass around an standard log.Logger, bypassing the LogLevel setting
	l.CriticalLogger().Print("test")

}
```
