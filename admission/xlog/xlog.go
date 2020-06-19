package xlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	ERROR = iota
	WARN
	INFO
	DEBUG
	TRACE
)

var Level = INFO // 日志级别

var Prefix = "" // 日志前缀

var Logpath = "" // 日志目录

var Logsize = int64(2 << 25) // 单日志文件大小上限

var Lognum = 20 // 日志文件轮换阈值

var Multilog = false // 是否同时输出到文件和标准输出

var xlog *log.Logger

var mutex sync.Mutex

var clogfile = ""

// var xlog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

func init() {
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	xlog = log.New(os.Stdout, Prefix, log.Ldate|log.Ltime|log.Lshortfile)

}

func set() {
	if Prefix != "" {
		xlog.SetPrefix(Prefix)
	}

	if Logpath != "" {
		dpath, err := filepath.Abs(Logpath)
		if err != nil {
			log.Fatalln("Find log path error: ", err)
		}

		if s, err := os.Stat(dpath); err != nil || !s.IsDir() {
			log.Fatalln("Log path is not exist or not a directory")
		}

		fpath := rotate(dpath)

		file, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Open log file error: ", err)
		}

		if Multilog {
			xlog.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			xlog.SetOutput(file)
		}
	}
}

func rotate(dir string) string {
	prefix := "log"

	format := ""

	//format := fmt.Sprintf("log.\%0%vd.log", math.Ceil(math.Log10(float64(Lognum))))
	switch {
	case Lognum <= 10:
		format = "%d"
	case Lognum <= 100:
		format = "%02d"
	case Lognum <= 1000:
		format = "%03d"
	case Lognum <= 10000:
		format = "%04d"
	case Lognum <= 100000:
		format = "%05d"
	default:
		format = "%09d"
	}

	seq := fmt.Sprintf(format, 0)

	ext := "log"

	if clogfile == "" {
		clogfile = fmt.Sprintf("%s.%s.%s", prefix, seq, ext)
	} else {
		num, err := strconv.Atoi(strings.Split(clogfile, ".")[1])
		if err != nil {
			log.Fatalln("Get log file sequence error: ", err)
		}

		s, err := os.Stat(clogfile)
		if err != nil {
			log.Fatalln("Get log file state error: ", err)
		}

		if s.Size() >= Logsize {
			seq = fmt.Sprintf(format, num+1%Lognum)
			clogfile = fmt.Sprintf("%s.%s.%s", prefix, seq, ext)
		}
	}

	return filepath.Join(dir, clogfile)
}

func run(kind string, v ...interface{}) string {
	if v == nil || len(v) < 1 {
		return ""
	}
	val := ""

	if kind != "PRINT" && kind != "SPRINT" {
		mutex.Lock()
		set()
		mutex.Unlock()
	}

	if len(v) > 1 {
		if format, ok := v[0].(string); ok || strings.Contains(format, "%") {
			str := ""
			if kind != "PRINT" && kind != "SPRINT" {
				str = fmt.Sprintf("["+kind+"] "+format, v[1:]...)
			} else {
				str = fmt.Sprintf(format, v[1:]...)
			}
			if !strings.Contains(str, "%!(EXTRA") {
				val = str
			}
		}
	}

	if val == "" {
		if kind != "PRINT" && kind != "SPRINT" {
			v = append([]interface{}{"[" + kind + "]"}, v...)
			val = fmt.Sprint(v...)
		} else {
			val = fmt.Sprint(v...)
		}
	}

	if kind == "PRINT" {
		fmt.Println(val)
	} else if kind != "SPRINT" {
		_ = xlog.Output(3, val)
	}
	return val
}

func Print(v ...interface{}) {
	run("PRINT", v...)
}

func Sprint(v ...interface{}) string {
	return run("SPRINT", v...)
}

func Error(v ...interface{}) {
	if Level >= ERROR {
		run("ERROR", v...)
	}
}

func Warn(v ...interface{}) {
	if Level >= WARN {
		run("WARN", v...)
	}
}

func Info(v ...interface{}) {
	if Level >= INFO {
		run("INFO", v...)
	}
}

func Debug(v ...interface{}) {
	if Level >= DEBUG {
		run("DEBUG", v...)
	}
}

func Trace(v ...interface{}) {
	if Level >= TRACE {
		run("TRACE", v...)
	}
}

func Fatal(v ...interface{}) {
	run("FATAL", v...)
	os.Exit(1)
}

func Panic(v ...interface{}) {
	panic(run("PANIC", v...))
}
