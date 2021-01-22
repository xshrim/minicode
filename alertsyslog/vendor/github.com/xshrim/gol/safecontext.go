package gol

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// safe context using mutex lock
type SafeContext struct {
	mu      sync.Mutex
	loggers []*Logger
	buf     []byte
}

// create thread-safe context with the fields and loggers
func NewSafeContext(fd F, loggers ...*Logger) *SafeContext {
	ctx := &SafeContext{loggers: nil, buf: map2json(nil, fd)}
	for _, l := range loggers {
		ctx.loggers = append(ctx.loggers, l)
	}
	return ctx
}

// ensure all logs are written to file for each logger
func (c *SafeContext) Flush() {
	for _, l := range c.loggers {
		if l != nil && l.done != nil {
			l.done <- true
		}
	}
}

// set loggers
func (c *SafeContext) Loggers(l ...*Logger) *SafeContext {
	if l != nil {
		c.loggers = append(c.loggers, l...)
	} else {
		c.loggers = nil
	}
	return c
}

// get fields
func (c *SafeContext) GetField() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.buf != nil && len(c.buf) > 1 && c.buf[len(c.buf)-1] != '}' {
		return append(c.buf, '}')
	}
	return c.buf
}

// set fields
func (c *SafeContext) Field(fd F) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	if fd == nil {
		c.buf = nil
		return c
	}

	c.buf = map2json(c.buf, fd)
	return c
}

// output error log
func (c *SafeContext) Error(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(ERROR) {
			l.Output(ERROR, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format error log
func (c *SafeContext) Errorf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(ERROR) {
			l.Output(ERROR, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output warn log
func (c *SafeContext) Warn(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(WARN) {
			l.Output(WARN, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format warn log
func (c *SafeContext) Warnf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(WARN) {
			l.Output(WARN, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output notice log
func (c *SafeContext) Notic(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(NOTIC) {
			l.Output(NOTIC, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format notice log
func (c *SafeContext) Noticf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(NOTIC) {
			l.Output(NOTIC, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output info log
func (c *SafeContext) Info(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(INFO) {
			l.Output(INFO, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format info log
func (c *SafeContext) Infof(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(INFO) {
			l.Output(INFO, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output debug log
func (c *SafeContext) Debug(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(DEBUG) {
			l.Output(DEBUG, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format debug log
func (c *SafeContext) Debugf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(DEBUG) {
			l.Output(DEBUG, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output trace log
func (c *SafeContext) Trace(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(TRACE) {
			l.Output(TRACE, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format trace log
func (c *SafeContext) Tracef(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(TRACE) {
			l.Output(TRACE, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// exit with code 1 after output fatal log
func (c *SafeContext) Fatal(v ...interface{}) {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(FATAL) {
			l.Output(FATAL, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	os.Exit(1)
}

// exit with code 1 after output format fatal log
func (c *SafeContext) Fatalf(format string, v ...interface{}) {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(FATAL) {
			l.Output(FATAL, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	os.Exit(1)
}

// panic after output panic log
func (c *SafeContext) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(PANIC) {
			l.Output(PANIC, 2, c.buf, s, true)
		}
	}
	panic(s)
}

// panic after output format panic log
func (c *SafeContext) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(PANIC) {
			l.Output(PANIC, 2, c.buf, s, true)
		}
	}
	panic(s)
}

// output log without newline
func (c *SafeContext) Log(level interface{}, v ...interface{}) Printer {
	s := fmt.Sprint(v...)
	lv := parseLevel(level)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(lv) {
			l.Output(lv, 2, c.buf, s, false)
		}
	}
	switch lv {
	case FATAL:
		os.Exit(1)
	case PANIC:
		panic(s)
	}
	return c
}

// output format log without newline
func (c *SafeContext) Logf(level interface{}, format string, v ...interface{}) Printer {
	s := fmt.Sprintf(format, v...)
	lv := parseLevel(level)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(lv) {
			l.Output(lv, 2, c.buf, s, false)
		}
	}
	switch lv {
	case FATAL:
		os.Exit(1)
	case PANIC:
		panic(s)
	}
	return c
}

// Str adds the field key with val as a string
func (c *SafeContext) Str(key, val string) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendStr(appendKey(c.buf, key), val)

	return c
}

// Strs adds the field key with vals as a []string
func (c *SafeContext) Strs(key string, val []string) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendStrs(appendKey(c.buf, key), val)

	return c
}

// Bytes adds the field key with val as a string
func (c *SafeContext) Bytes(key string, val []byte) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendBytes(appendKey(c.buf, key), val)

	return c
}

// Hex adds the field key with val as a hex string
func (c *SafeContext) Hex(key string, val []byte) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendHex(appendKey(c.buf, key), val)

	return c
}

// Rjson adds already encoded JSON to the log line under key
func (c *SafeContext) Json(key string, val []byte) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendJson(appendKey(c.buf, key), val)

	return c
}

// Er adds the field key with serialized err
func (c *SafeContext) Er(key string, err error) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		c.buf = appendStr(appendKey(c.buf, key), err.Error())
	}

	return c
}

// Bool adds the field key with val as a bool
func (c *SafeContext) Bool(key string, val bool) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendBool(appendKey(c.buf, key), val)

	return c
}

// Bools adds the field key with val as a []bool
func (c *SafeContext) Bools(key string, val []bool) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendBools(appendKey(c.buf, key), val)

	return c
}

// Int adds the field key with i as a int
func (c *SafeContext) Int(key string, val int) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInt(appendKey(c.buf, key), val)

	return c
}

// Ints adds the field key with i as a []int
func (c *SafeContext) Ints(key string, val []int) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInts(appendKey(c.buf, key), val)

	return c
}

// Int8 adds the field key with i as a int8
func (c *SafeContext) Int8(key string, val int8) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInt8(appendKey(c.buf, key), val)

	return c
}

// Ints8 adds the field key with i as a []int8
func (c *SafeContext) Ints8(key string, val []int8) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInts8(appendKey(c.buf, key), val)

	return c
}

// Int16 adds the field key with i as a int16
func (c *SafeContext) Int16(key string, val int16) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInt16(appendKey(c.buf, key), val)

	return c
}

// Ints16 adds the field key with i as a []int16
func (c *SafeContext) Ints16(key string, val []int16) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInts16(appendKey(c.buf, key), val)

	return c
}

// Int32 adds the field key with i as a int32
func (c *SafeContext) Int32(key string, val int32) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInt32(appendKey(c.buf, key), val)

	return c
}

// Ints32 adds the field key with i as a []int32
func (c *SafeContext) Ints32(key string, val []int32) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInts32(appendKey(c.buf, key), val)

	return c
}

// Int64 adds the field key with i as a int64
func (c *SafeContext) Int64(key string, val int64) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInt64(appendKey(c.buf, key), val)

	return c
}

// Ints64 adds the field key with i as a []int64
func (c *SafeContext) Ints64(key string, val []int64) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInts64(appendKey(c.buf, key), val)

	return c
}

// Uint adds the field key with i as a uint
func (c *SafeContext) Uint(key string, val uint) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUint(appendKey(c.buf, key), val)

	return c
}

// Uints adds the field key with i as a []uint
func (c *SafeContext) Uints(key string, val []uint) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUints(appendKey(c.buf, key), val)

	return c
}

// Uint8 adds the field key with i as a uint8
func (c *SafeContext) Uint8(key string, val uint8) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUint8(appendKey(c.buf, key), val)

	return c
}

// Uints8 adds the field key with i as a []uint8
func (c *SafeContext) Uints8(key string, val []uint8) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUints8(appendKey(c.buf, key), val)

	return c
}

// Uint16 adds the field key with i as a uint16
func (c *SafeContext) Uint16(key string, val uint16) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUint16(appendKey(c.buf, key), val)

	return c
}

// Uints16 adds the field key with i as a []uint16
func (c *SafeContext) Uints16(key string, val []uint16) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUints16(appendKey(c.buf, key), val)

	return c
}

// Uint32 adds the field key with i as a uint32
func (c *SafeContext) Uint32(key string, val uint32) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUint32(appendKey(c.buf, key), val)

	return c
}

// Uints32 adds the field key with i as a []uint32
func (c *SafeContext) Uints32(key string, val []uint32) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUints32(appendKey(c.buf, key), val)

	return c
}

// Uint64 adds the field key with i as a uint64
func (c *SafeContext) Uint64(key string, val uint64) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUint64(appendKey(c.buf, key), val)

	return c
}

// Uints64 adds the field key with i as a []uint64
func (c *SafeContext) Uints64(key string, val []uint64) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendUints64(appendKey(c.buf, key), val)

	return c
}

// Float32 adds the field key with f as a float32
func (c *SafeContext) Float32(key string, val float32) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendFloat32(appendKey(c.buf, key), val)

	return c
}

// Floats32 adds the field key with f as a []float32
func (c *SafeContext) Floats32(key string, val []float32) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendFloats32(appendKey(c.buf, key), val)

	return c
}

// Float64 adds the field key with f as a float64
func (c *SafeContext) Float64(key string, val float64) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendFloat64(appendKey(c.buf, key), val)

	return c
}

// Floats64 adds the field key with f as a []float64
func (c *SafeContext) Floats64(key string, val []float64) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendFloats64(appendKey(c.buf, key), val)

	return c
}

// Ts adds the current local time as UNIX timestamp
func (c *SafeContext) Ts() *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendTime(appendKey(c.buf, "time"), time.Now(), time.RFC3339)

	return c
}

// Time adds the field key with time.Time formated as string
func (c *SafeContext) Time(key string, t time.Time) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendTime(appendKey(c.buf, key), t, time.RFC3339)

	return c
}

// Times adds the field key with []time.Time formated as string
func (c *SafeContext) Times(key string, t []time.Time) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendTimes(appendKey(c.buf, key), t, time.RFC3339)

	return c
}

// Dur adds the field key with time.Duration
func (c *SafeContext) Dur(key string, d time.Duration) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendDuration(appendKey(c.buf, key), d, time.Millisecond)

	return c
}

// Durs adds the field key with []time.Duration
func (c *SafeContext) Durs(key string, d []time.Duration) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendDurations(appendKey(c.buf, key), d, time.Millisecond)

	return c
}

// TimeDiff adds the field key with positive duration between time t and start
func (c *SafeContext) TimeDiff(key string, t time.Time, start time.Time) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}
	// append value
	c.buf = appendDuration(appendKey(c.buf, key), d, time.Millisecond)

	return c
}

// If adds the field key with i marshaled using reflection
func (c *SafeContext) If(key string, i interface{}) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendInterface(appendKey(c.buf, key), i)

	return c
}

// Obj adds the field key with fmt.Sprintf raw string
// for json string, use If function
func (c *SafeContext) Obj(key string, o interface{}) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendObject(appendKey(c.buf, key), o)

	return c
}

// IPAddr adds IPv4 or IPv6 Address
func (c *SafeContext) Ip(key string, ip net.IP) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendIPAddr(appendKey(c.buf, key), ip)

	return c
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask)
func (c *SafeContext) Ipp(key string, pfx net.IPNet) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendIPPrefix(appendKey(c.buf, key), pfx)

	return c
}

// MACAddr adds MAC address
func (c *SafeContext) Mac(key string, ha net.HardwareAddr) *SafeContext {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = appendMACAddr(appendKey(c.buf, key), ha)

	return c
}
