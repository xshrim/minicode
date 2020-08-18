package xlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	OFF = 1 << iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Lstack                        // print stack information<package.function>
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lcolor                        // if output colorful log level or not
	Lfullcolor                    // if output colorful log or not
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type LogSaver interface {
	save([]string)
}

type LogFile struct {
	path string
}

func NewLogSaverWithFile(path string) LogSaver {
	return &LogFile{path: path}
}

func (lf *LogFile) save(data []string) {
	logfile := lf.path
	mode := os.O_RDWR | os.O_CREATE | os.O_APPEND

	writeFile(logfile, mode, data)
}

type Rotation struct {
	dir    string // log file directory
	prefix string // log file base name
	format string // log file sequence format
	size   int64  // log file max size(byte), 0 = no limit
	count  int    // log rotate threshold
	cur    int    // log current file number
}

func NewLogSaverWithRotation(dir string, size int64, count int) LogSaver {
	fname := os.Args[0]
	var prefixBytes []byte
	for i := len(fname) - 1; i >= 0; i-- {
		if fname[i] != os.PathSeparator {
			prefixBytes = append([]byte{fname[i]}, prefixBytes...)
		} else {
			break
		}
	}
	prefix := string(prefixBytes)

	if len(dir) == 0 {
		dir = "." + string(os.PathSeparator)
	}
	if !os.IsPathSeparator(dir[len(dir)-1]) {
		dir = dir + string(os.PathSeparator)
	}

	if count <= 0 {
		count = 1
	}

	if size < 0 {
		size = 0
	}

	format := ""
	switch {
	case count <= 10:
		format = "%d"
	case count <= 100:
		format = "%02d"
	case count <= 1000:
		format = "%03d"
	case count <= 10000:
		format = "%04d"
	case count <= 100000:
		format = "%05d"
	default:
		format = "%09d"
	}

	return &Rotation{dir: dir, prefix: prefix, size: size, count: count, format: format}
}

func (rt *Rotation) save(data []string) {
	logfile, mode := rt.rotate()

	writeFile(logfile, mode, data)
}

func (rt *Rotation) rotate() (string, int) {
	defaultMode := os.O_RDWR | os.O_CREATE | os.O_APPEND

	logfile := fmt.Sprintf("%s%s.%s.log", rt.dir, rt.prefix, fmt.Sprintf(rt.format, rt.cur))
	s, err := os.Stat(logfile)
	if err == nil && s != nil && !s.IsDir() {
		if rt.size > 0 && s.Size() >= rt.size {
			rt.cur = (rt.cur + 1) % rt.count
			logfile = fmt.Sprintf("%s%s.%s.log", rt.dir, rt.prefix, fmt.Sprintf(rt.format, rt.cur))
			return logfile, os.O_RDWR | os.O_CREATE | os.O_TRUNC
		}
	} else if s != nil && s.IsDir() {
		return "", defaultMode
	}

	return logfile, defaultMode
}

type Logger struct {
	mu      sync.Mutex  // ensures atomic writes; protects the following fields
	prefix  string      // prefix to write at beginning of each line
	level   int         // log level(ERROR, WARN, INFO, DEBUG, TRACE)
	hotlv   int         // hot reload log level from Environment
	flag    int         // properties
	out     io.Writer   // destination for output
	saver   LogSaver    // write log to file
	buf     []byte      // for accumulating text to write
	bufchan chan string // for logsaver to read from
	done    chan bool   // for savelog goroutine to exit
}

var std = New(os.Stderr, "", INFO, LstdFlags)

func New(out io.Writer, prefix string, level, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, level: level, flag: flag}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func toUpper(src string) string {
	var dst []rune
	for _, v := range src {
		if v >= 97 && v <= 122 {
			v -= 32
		}
		dst = append(dst, v)
	}

	return string(dst)
}

func parseLevel(level interface{}) (int, string) {
	var lv int
	var str string
	switch res := level.(type) {
	case int:
		lv = res
		switch lv {
		case OFF:
			str = "[OFF]"
		case ERROR:
			str = "[ERROR]"
		case WARN:
			str = "[WARN]"
		case INFO:
			str = "[INFO]"
		case DEBUG:
			str = "[DEBUG]"
		case TRACE:
			str = "[TRACE]"
		default:
			lv = -1
			str = "[FATAL]"
		}
	case string:
		str = "[" + toUpper(res) + "]"
		switch str {
		case "[OFF]", "[1]":
			lv = OFF
		case "[ERROR]", "[2]":
			lv = ERROR
		case "[WARN]", "[3]":
			lv = WARN
		case "[INFO]", "[4]":
			lv = INFO
		case "[DEBUG]", "[5]":
			lv = DEBUG
		case "[TRACE]", "[6]":
			lv = TRACE
		default:
			lv = -1
		}
	default:
		lv = INFO
		str = "[INFO]"
	}

	return lv, str
}

func writeFile(logfile string, mode int, data []string) {
	if logfile == "" {
		fmt.Println("Can not get log file")
		return
	}
	file, err := os.OpenFile(logfile, mode, 0666)
	if err != nil {
		fmt.Println("Can not open log file: ", err)
		return
	}
	for _, data := range data {
		if _, err = file.WriteString(data); err != nil {
			fmt.Println("Write log file failed: ", err)
			return
		}
	}
	file.Close()
}

func HotReload() {
	std.HotReload()
}

// ensure all logs are written to file
func Flush() {
	std.Flush()
}

func Writer() io.Writer {
	return std.Writer()
}

func SetWriter(w io.Writer) {
	std.SetWriter(w)
}

func UnSetWriter() {
	std.UnSetWriter()
}

func Saver() LogSaver {
	return std.Saver()
}

func SetSaver(ls LogSaver) {
	std.SetSaver(ls)
}

func Prefix() string {
	return std.Prefix()
}

func SetPrefix(p string) {
	std.SetPrefix(p)
}

func Level() int {
	return std.Level()
}

func SetLevel(v int) {
	std.SetLevel(v)
}

func Flag() int {
	return std.Flag()
}

func SetFlag(f int) {
	std.SetFlag(f)
}

func Println(v ...interface{}) {
	fmt.Println(formatData(v...))
}

func Print(v ...interface{}) {
	fmt.Print(formatData(v...))
}

func Sprint(v ...interface{}) string {
	return fmt.Sprint(formatData(v...))
}

func Err(v ...interface{}) error {
	return errors.New(formatData(v...))
}

func Error(v ...interface{}) {
	if std.lvcheck(ERROR) {
		_, _ = std.Output(2, true, "[ERROR]", formatData(v...))
	}
}

func Warn(v ...interface{}) {
	if std.lvcheck(WARN) {
		_, _ = std.Output(2, true, "[WARN]", formatData(v...))
	}
}

func Info(v ...interface{}) {
	if std.lvcheck(INFO) {
		_, _ = std.Output(2, true, "[INFO]", formatData(v...))
	}
}

func Debug(v ...interface{}) {
	if std.lvcheck(DEBUG) {
		_, _ = std.Output(2, true, "[DEBUG]", formatData(v...))
	}
}

func Trace(v ...interface{}) {
	if std.lvcheck(TRACE) {
		_, _ = std.Output(2, true, "[TRACE]", formatData(v...))
	}
}

func Fatal(v ...interface{}) {
	_, _ = std.Output(2, true, "[FATAL]", formatData(v...))
	os.Exit(1)
}

func Panic(v ...interface{}) {
	msg, _ := std.Output(2, true, "[PANIC]", formatData(v...))
	panic(msg)
}

// output log without newline
func Log(level interface{}, v ...interface{}) {
	var msg string
	lv, str := parseLevel(level)
	if std.lvcheck(lv) {
		msg, _ = std.Output(2, false, str, formatData(v...))
	}
	if lv < 0 {
		if str == "[FATAL]" {
			os.Exit(1)
		} else {
			panic(msg)
		}
	}
}

func (l *Logger) HotReload() {
	if l.hotlv == 0 { // hot reload can only be set once
		l.hotlv = l.level
		go l.notisfy()
	}
}

func (l *Logger) Flush() {
	if l.done != nil {
		l.done <- true
	}
}

func (l *Logger) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetWriter(writer io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = writer
}

func (l *Logger) UnSetWriter() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = nil
}

func (l *Logger) Saver() LogSaver {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.saver
}

func (l *Logger) SetSaver(ls LogSaver) {
	if l.saver != nil || ls == nil { // log saver can only be set once
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.saver = ls
	l.done = make(chan bool)
	l.bufchan = make(chan string, 3000)
	go l.saveLog()
}

func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) Level() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) Flag() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

func (l *Logger) SetFlag(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

func (l *Logger) Println(v ...interface{}) {
	fmt.Println(formatData(v...))
}

func (l *Logger) Print(v ...interface{}) {
	fmt.Print(formatData(v...))
}

func (l *Logger) Sprint(v ...interface{}) string {
	return fmt.Sprint(formatData(v...))
}

func (l *Logger) Err(v ...interface{}) error {
	return errors.New(formatData(v...))
}

func (l *Logger) Error(v ...interface{}) {
	if l.lvcheck(ERROR) {
		_, _ = l.Output(2, true, "[ERROR]", formatData(v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.lvcheck(WARN) {
		_, _ = l.Output(2, true, "[WARN]", formatData(v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.lvcheck(INFO) {
		_, _ = l.Output(2, true, "[INFO]", formatData(v...))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.lvcheck(DEBUG) {
		_, _ = l.Output(2, true, "[DEBUG]", formatData(v...))
	}
}

func (l *Logger) Trace(v ...interface{}) {
	if l.lvcheck(TRACE) {
		_, _ = l.Output(2, true, "[TRACE]", formatData(v...))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	_, _ = l.Output(2, true, "[FATAL]", formatData(v...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	msg, _ := l.Output(2, true, "[PANIC]", formatData(v...))
	panic(msg)
}

func (l *Logger) Log(level interface{}, v ...interface{}) {
	var msg string
	lv, str := parseLevel(level)
	if l.lvcheck(lv) {
		msg, _ = l.Output(2, false, str, formatData(v...))
	}
	if lv < 0 {
		if str == "[FATAL]" {
			os.Exit(1)
		} else {
			panic(msg)
		}
	}
}

func formatData(v ...interface{}) string {
	if v == nil || len(v) < 1 {
		return ""
	}
	val := ""

	if len(v) > 1 {
		if format, ok := v[0].(string); ok || (bytes.Contains([]byte(format), []byte("%")) && !bytes.Contains([]byte(format), []byte("\\%"))) {
			str := fmt.Sprintf(format, v[1:]...)
			if !bytes.Contains([]byte(str), []byte("%!(EXTRA")) {
				val = str
			}
		}
	}

	if val == "" {
		val = fmt.Sprint(v...)
	}

	return val
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, kind, fn, file string, line int) {
	if l.flag&Lfullcolor != 0 && l.prefix != "" {
		*buf = append(*buf, []byte("\033[0;34m")...)
	}
	*buf = append(*buf, l.prefix...)
	if l.flag&Lfullcolor != 0 && l.prefix != "" {
		*buf = append(*buf, []byte("\033[0m")...)
	}

	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			if l.flag&Lfullcolor != 0 {
				*buf = append(*buf, []byte("\033[0;33m")...)
			}
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
			if l.flag&Lfullcolor != 0 {
				*buf = append(*buf, []byte("\033[0m")...)
			}
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			if l.flag&Lfullcolor != 0 {
				*buf = append(*buf, []byte("\033[0;33m")...)
			}
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
			if l.flag&Lfullcolor != 0 {
				*buf = append(*buf, []byte("\033[0m")...)
			}
		}
	}
	if l.flag&Lstack != 0 {
		if l.flag&Lfullcolor != 0 {
			*buf = append(*buf, []byte("\033[0;35m")...)
		}
		*buf = append(*buf, '<')
		*buf = append(*buf, fn...)
		*buf = append(*buf, '>')
		*buf = append(*buf, ' ')
		if l.flag&Lfullcolor != 0 {
			*buf = append(*buf, []byte("\033[0m")...)
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lfullcolor != 0 {
			*buf = append(*buf, []byte("\033[0;36m")...)
		}
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
		if l.flag&Lfullcolor != 0 {
			*buf = append(*buf, []byte("\033[0m")...)
		}
	}

	if l.flag&Lcolor != 0 || l.flag&Lfullcolor != 0 { // colorful output feature only supports linux
		switch kind {
		case "[ERROR]":
			*buf = append(*buf, []byte("\033[1;31m")...)
		case "[WARN]":
			*buf = append(*buf, []byte("\033[1;33m")...)
		case "[INFO]":
			*buf = append(*buf, []byte("\033[1;32m")...)
		case "[DEBUG]":
			*buf = append(*buf, []byte("\033[1;35m")...)
		case "[TRACE]":
			*buf = append(*buf, []byte("\033[1;36m")...)
		case "[FATAL]":
			*buf = append(*buf, []byte("\033[1;34m")...)
		case "[PANIC]":
			*buf = append(*buf, []byte("\033[1;37m")...)
		}
	}
	*buf = append(*buf, kind...)
	*buf = append(*buf, ' ')
	if l.flag&Lcolor != 0 || l.flag&Lfullcolor != 0 {
		*buf = append(*buf, []byte("\033[0m")...)
	}
}

func (l *Logger) Output(calldepth int, feed bool, kind, s string) (string, error) {
	now := time.Now() // get this early.
	var pc uintptr
	var fn string
	var file string
	var line int
	var err error
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		pc, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
			fn = "???"
		} else {
			if l.flag&Lstack != 0 {
				fn = runtime.FuncForPC(pc).Name()
			}
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, kind, fn, file, line)

	if l.flag&Lfullcolor != 0 {
		l.buf = append(l.buf, []byte("\033[0;32m")...)
	}
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	l.buf = append(l.buf, s...)
	if l.flag&Lfullcolor != 0 {
		l.buf = append(l.buf, []byte("\033[0m")...)
	}

	if feed {
		l.buf = append(l.buf, '\n')
	}

	if l.bufchan != nil {
		l.bufchan <- string(l.buf)
	}

	if l.out != nil {
		_, err = l.out.Write(l.buf)
	}

	return string(l.buf), err
}

func (l *Logger) notisfy() {
	var lastModifyTime int64
	for {
		// lv, _ := parseLevel(os.Getenv("XLOG_LEVEL"))
		cpath := "/tmp/.xlog"
		if runtime.GOOS == "windows" {
			cpath = "C:\\.xlog"
		}
		file, err := os.Open(cpath)
		if err == nil {
			fileInfo, err := file.Stat()
			if err == nil {
				curModifyTime := fileInfo.ModTime().Unix()
				if curModifyTime > lastModifyTime {
					lastModifyTime = curModifyTime
					var line []byte
					for {
						b := make([]byte, 1)
						n, _ := file.Read(b)
						if n > 0 {
							c := b[0]
							if c == '\n' {
								break
							}
							line = append(line, c)
						}
					}
					lv, _ := parseLevel(string(line))
					if lv == -1 {
						l.hotlv = l.level
					} else {
						l.hotlv = lv
					}
				}
			}
		}
		file.Close()

		time.Sleep(3 * time.Second)
	}
}

func (l *Logger) lvcheck(lv int) bool {
	if l.hotlv >= lv || (l.hotlv == 0 && l.level >= lv) {
		return true
	}
	return false
}

func (l *Logger) saveLog() {
	for {
		select {
		case data := <-l.bufchan:
			l.saver.save([]string{data})
		case <-l.done:
			close(l.bufchan)
			break
		}
	}
	// for data := range l.bufchan {
	// 	l.saver.save([]string{data})
	// }
}

// func (l *Logger) writelog() error {
// fw := l.rotatelog()
// if fw != nil {
// 	if l.out != nil {
// 		l.out = io.MultiWriter(l.out, fw)
// 	} else {
// 		l.out = fw
// 	}
// }
// }
