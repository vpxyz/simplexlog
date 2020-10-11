SIMPLEXLOG leveled logging library
======================

Simplexlog simple wrapper for the go standard log package, that adds log level.

The standard log package is enough (imho) for many applications, but sometimes
is usefull to have some log level.
Simplexlog is a wrapper that adds some log level, the ability to use different
io.Writer for different level.

The log levels are: TRACE, DEBUG, CRITICAL, ERROR, WARNING, NOTICE, INFO
or ALL to switch on all levels

Simplexlog is concurrent safe.

Installation
------------

The package is go gettable:  go get -u github.com/vpxyz/simplexlog

Example
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
	// If you want color labels, just put colors escape sequence around label. For e.g. "\x1b[20;32m"+sl.LevelInfo+"\x1b[0m"
	l := sl.New(
        sl.SetDebug(sl.Config{Out: os.Stdout, Label: sl.LevelInfo + " ==> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
		sl.SetTrace(sl.Config{Out: os.Stdout, Label: sl.LevelTrace + " ===> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
		sl.SetInfo(sl.Config{Out: os.Stdout, Label: sl.LevelInfo + " =>", Flags: sl.DefaultLogFlags}),
		sl.SetNotice(sl.Config{Out: os.Stdout, Label: fmt.Sprintf("%-10s", "["+sl.LevelNotice+"]:"), Flags: sl.DefaultLogFlags}),
		sl.SetWarning(sl.Config{Out: os.Stdout, Label: sl.LevelWarning + ", ARGH! ", Flags: sl.DefaultLogFlags}),
		sl.SetError(sl.Config{Out: os.Stderr, Label: sl.LevelError + " ", Flags: sl.DefaultLogFlags}),
		sl.SetCritical(sl.Config{Out: os.Stderr, Label: sl.LevelCritical + ",GULP! ==> ", Flags: sl.DefaultLogFlags | log.Lshortfile}),
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

    // change level
    l.SwitchTo(sl.Warning)
    
    l.Info("This is hidden")

	// if you need, you can pass around an standard log.Logger, bypassing the LogLevel setting
	l.CriticalLogger().Print("test")

    // change the log level using log level name (case insensitive)
   	l.SwitchTo("error")

	l.Infof("Info log")

	l.Noticef("Notice log")

	l.Warningf("Warning log")

	l.Debugf("Debug log")

	l.Errorf("Error log")

	l.Criticalf("Critical log")

}
```
