package gol

const (
	OFF = iota
	PANIC
	FATAL
	ERROR
	WARN
	NOTIC
	INFO
	DEBUG
	TRACE
	ALL
)

const (
	Ldate    = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                    // the time in the local time zone: 01:23:23
	Lmsec                    // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Lstack                   // print stack information<package.function>
	Lnolvl                   // not print level information
	Lfile                    // final file name element and line number: d.go:23. overrides Llongfilec/d.go:23
	Llfile                   // full file name and line number: /a/b/
	Ljson                    // output json, this flag will override Lcolor and Lfullcolor flag
	Lcolor                   // if output colorful log level or not
	Lfcolor                  // if output colorful log or not
	Lutc                     // if Ldate or Ltime is set, use UTC rather than the local time zone
	Ldefault = Ldate | Ltime // initial values for the standard logger
)

const (
	OVERRIDE = iota
	APPEND
	DELETE
)

// default logger
var std = New()

const hexs = "0123456789abcdef"
