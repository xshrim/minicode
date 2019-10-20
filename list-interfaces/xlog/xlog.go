package xlog

import (
	"log"
)

const (
	ERROR = iota
	WARN
	INFO
	DEBUG
	TRACE
)

var Level = INFO

func Errorf(format string, v ...interface{}) {
	if Level >= ERROR {
		log.Printf("[ERROR] "+format, v...)
	}
}

func Errorln(v ...interface{}) {
	if Level >= ERROR {
		v = append([]interface{}{"[ERROR]"}, v...)
		log.Println(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if Level >= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

func Warnln(v ...interface{}) {
	if Level >= WARN {
		v = append([]interface{}{"[WARN]"}, v...)
		log.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if Level >= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

func Infoln(v ...interface{}) {
	if Level >= INFO {
		v = append([]interface{}{"[INFO]"}, v...)
		log.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if Level >= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Debugln(v ...interface{}) {
	if Level >= DEBUG {
		v = append([]interface{}{"[DEBUG]"}, v...)
		log.Println(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if Level >= TRACE {
		log.Printf("[TRACE] "+format, v...)
	}
}

func Traceln(v ...interface{}) {
	if Level >= TRACE {
		v = append([]interface{}{"[TRACE]"}, v...)
		log.Println(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

func Fatalln(v ...interface{}) {
	v = append([]interface{}{"[FATAL]"}, v...)
	log.Fatalln(v...)
}

func Panicf(format string, v ...interface{}) {
	log.Panicf("[PANIC] "+format, v...)
}

func Panicln(v ...interface{}) {
	v = append([]interface{}{"[PANIC]"}, v...)
	log.Panicln(v...)
}
