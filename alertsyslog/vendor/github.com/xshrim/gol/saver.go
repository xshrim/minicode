package gol

import (
	"fmt"
	"os"
)

type LogSaver interface {
	save([]string)
}

// single file for log persistence
type LogFile struct {
	path string
}

// rotatable files for log persistence
type Rotation struct {
	dir    string // log file directory
	prefix string // log file base name
	format string // log file sequence format
	size   int64  // max size(byte) of each log file, 0 = no limit
	count  int    // log rotate threshold
	cur    int    // current log file number
}

// create log saver with single file
func NewLogSaverWithFile(path string) LogSaver {
	return &LogFile{path: path}
}

func (lf *LogFile) save(data []string) {
	logfile := lf.path
	mode := os.O_RDWR | os.O_CREATE | os.O_APPEND

	writeFile(logfile, mode, data)
}

// create log saver with rotatable files
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
