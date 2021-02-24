package gol

import (
	"fmt"
	"net"
	"os"
	"time"
)

type Context struct {
	loggers []*Logger
	buf     []byte
}

// create context with the fields and loggers
func NewContext(fd map[string]interface{}, loggers ...*Logger) *Context {
	ctx := &Context{loggers: nil, buf: map2json(nil, fd)}
	ctx.loggers = append(ctx.loggers, loggers...)
	return ctx
}

// ensure all logs are written to file for each logger
func (c *Context) Flush() {
	for _, l := range c.loggers {
		if l != nil && l.done != nil {
			l.done <- true
		}
	}
}

// set loggers
func (c *Context) Loggers(l ...*Logger) *Context {
	if l != nil {
		c.loggers = append(c.loggers, l...)
	} else {
		c.loggers = nil
	}
	return c
}

// get fields
func (c *Context) GetField() []byte {
	return c.buf
}

// set fields
func (c *Context) Field(fd map[string]interface{}) *Context {
	if fd == nil {
		c.buf = nil
		return c
	}

	c.buf = map2json(c.buf, fd)
	return c
}

// output error log
func (c *Context) Error(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(ERROR) {
			_ = l.Output(ERROR, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format error log
func (c *Context) Errorf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(ERROR) {
			_ = l.Output(ERROR, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output warn log
func (c *Context) Warn(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(WARN) {
			_ = l.Output(WARN, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format warn log
func (c *Context) Warnf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(WARN) {
			_ = l.Output(WARN, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output notice log
func (c *Context) Notic(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(NOTIC) {
			_ = l.Output(NOTIC, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format notice log
func (c *Context) Noticf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(NOTIC) {
			_ = l.Output(NOTIC, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output info log
func (c *Context) Info(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(INFO) {
			_ = l.Output(INFO, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format info log
func (c *Context) Infof(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(INFO) {
			_ = l.Output(INFO, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output debug log
func (c *Context) Debug(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(DEBUG) {
			_ = l.Output(DEBUG, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format debug log
func (c *Context) Debugf(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(DEBUG) {
			_ = l.Output(DEBUG, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// output trace log
func (c *Context) Trace(v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(TRACE) {
			_ = l.Output(TRACE, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	return c
}

// output format trace log
func (c *Context) Tracef(format string, v ...interface{}) Printer {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(TRACE) {
			_ = l.Output(TRACE, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	return c
}

// exit with code 1 after output fatal log
func (c *Context) Fatal(v ...interface{}) {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(FATAL) {
			_ = l.Output(FATAL, 2, c.buf, fmt.Sprint(v...), true)
		}
	}
	os.Exit(1)
}

// exit with code 1 after output format fatal log
func (c *Context) Fatalf(format string, v ...interface{}) {
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(FATAL) {
			_ = l.Output(FATAL, 2, c.buf, fmt.Sprintf(format, v...), true)
		}
	}
	os.Exit(1)
}

// panic after output panic log
func (c *Context) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(PANIC) {
			_ = l.Output(PANIC, 2, c.buf, s, true)
		}
	}
	panic(s)
}

// panic after output format panic log
func (c *Context) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(PANIC) {
			_ = l.Output(PANIC, 2, c.buf, s, true)
		}
	}
	panic(s)
}

// outut log without newline
func (c *Context) Log(level interface{}, v ...interface{}) Printer {
	s := fmt.Sprint(v...)
	lv := parseLevel(level)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(lv) {
			_ = l.Output(lv, 2, c.buf, s, false)
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
func (c *Context) Logf(level interface{}, format string, v ...interface{}) Printer {
	s := fmt.Sprintf(format, v...)
	lv := parseLevel(level)
	for _, l := range c.loggers {
		if l != nil && l.lvcheck(lv) {
			_ = l.Output(lv, 2, c.buf, s, false)
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
func (c *Context) Str(key, val string) *Context {
	c.buf = appendStr(appendKey(c.buf, key), val)

	return c
}

// Strs adds the field key with vals as a []string
func (c *Context) Strs(key string, val []string) *Context {
	c.buf = appendStrs(appendKey(c.buf, key), val)

	return c
}

// Bytes adds the field key with val as a string
func (c *Context) Bytes(key string, val []byte) *Context {
	c.buf = appendBytes(appendKey(c.buf, key), val)

	return c
}

// Hex adds the field key with val as a hex string
func (c *Context) Hex(key string, val []byte) *Context {
	c.buf = appendHex(appendKey(c.buf, key), val)

	return c
}

// Rjson adds already encoded JSON to the log line under key
func (c *Context) Json(key string, val []byte) *Context {
	c.buf = appendJson(appendKey(c.buf, key), val)

	return c
}

// Er adds the field key with serialized err
func (c *Context) Er(key string, err error) *Context {
	if err != nil {
		c.buf = appendStr(appendKey(c.buf, key), err.Error())
	}

	return c
}

// Bool adds the field key with val as a bool
func (c *Context) Bool(key string, val bool) *Context {
	c.buf = appendBool(appendKey(c.buf, key), val)

	return c
}

// Bools adds the field key with val as a []bool
func (c *Context) Bools(key string, val []bool) *Context {
	c.buf = appendBools(appendKey(c.buf, key), val)

	return c
}

// Int adds the field key with i as a int
func (c *Context) Int(key string, val int) *Context {
	c.buf = appendInt(appendKey(c.buf, key), val)

	return c
}

// Ints adds the field key with i as a []int
func (c *Context) Ints(key string, val []int) *Context {
	c.buf = appendInts(appendKey(c.buf, key), val)

	return c
}

// Int8 adds the field key with i as a int8
func (c *Context) Int8(key string, val int8) *Context {
	c.buf = appendInt8(appendKey(c.buf, key), val)

	return c
}

// Ints8 adds the field key with i as a []int8
func (c *Context) Ints8(key string, val []int8) *Context {
	c.buf = appendInts8(appendKey(c.buf, key), val)

	return c
}

// Int16 adds the field key with i as a int16
func (c *Context) Int16(key string, val int16) *Context {
	c.buf = appendInt16(appendKey(c.buf, key), val)

	return c
}

// Ints16 adds the field key with i as a []int16
func (c *Context) Ints16(key string, val []int16) *Context {
	c.buf = appendInts16(appendKey(c.buf, key), val)

	return c
}

// Int32 adds the field key with i as a int32
func (c *Context) Int32(key string, val int32) *Context {
	c.buf = appendInt32(appendKey(c.buf, key), val)

	return c
}

// Ints32 adds the field key with i as a []int32
func (c *Context) Ints32(key string, val []int32) *Context {
	c.buf = appendInts32(appendKey(c.buf, key), val)

	return c
}

// Int64 adds the field key with i as a int64
func (c *Context) Int64(key string, val int64) *Context {
	c.buf = appendInt64(appendKey(c.buf, key), val)

	return c
}

// Ints64 adds the field key with i as a []int64
func (c *Context) Ints64(key string, val []int64) *Context {
	c.buf = appendInts64(appendKey(c.buf, key), val)

	return c
}

// Uint adds the field key with i as a uint
func (c *Context) Uint(key string, val uint) *Context {
	c.buf = appendUint(appendKey(c.buf, key), val)

	return c
}

// Uints adds the field key with i as a []uint
func (c *Context) Uints(key string, val []uint) *Context {
	c.buf = appendUints(appendKey(c.buf, key), val)

	return c
}

// Uint8 adds the field key with i as a uint8
func (c *Context) Uint8(key string, val uint8) *Context {
	c.buf = appendUint8(appendKey(c.buf, key), val)

	return c
}

// Uints8 adds the field key with i as a []uint8
func (c *Context) Uints8(key string, val []uint8) *Context {
	c.buf = appendUints8(appendKey(c.buf, key), val)

	return c
}

// Uint16 adds the field key with i as a uint16
func (c *Context) Uint16(key string, val uint16) *Context {
	c.buf = appendUint16(appendKey(c.buf, key), val)

	return c
}

// Uints16 adds the field key with i as a []uint16
func (c *Context) Uints16(key string, val []uint16) *Context {
	c.buf = appendUints16(appendKey(c.buf, key), val)

	return c
}

// Uint32 adds the field key with i as a uint32
func (c *Context) Uint32(key string, val uint32) *Context {
	c.buf = appendUint32(appendKey(c.buf, key), val)

	return c
}

// Uints32 adds the field key with i as a []uint32
func (c *Context) Uints32(key string, val []uint32) *Context {
	c.buf = appendUints32(appendKey(c.buf, key), val)

	return c
}

// Uint64 adds the field key with i as a uint64
func (c *Context) Uint64(key string, val uint64) *Context {
	c.buf = appendUint64(appendKey(c.buf, key), val)

	return c
}

// Uints64 adds the field key with i as a []uint64
func (c *Context) Uints64(key string, val []uint64) *Context {
	c.buf = appendUints64(appendKey(c.buf, key), val)

	return c
}

// Float32 adds the field key with f as a float32
func (c *Context) Float32(key string, val float32) *Context {
	c.buf = appendFloat32(appendKey(c.buf, key), val)

	return c
}

// Floats32 adds the field key with f as a []float32
func (c *Context) Floats32(key string, val []float32) *Context {
	c.buf = appendFloats32(appendKey(c.buf, key), val)

	return c
}

// Float64 adds the field key with f as a float64
func (c *Context) Float64(key string, val float64) *Context {
	c.buf = appendFloat64(appendKey(c.buf, key), val)

	return c
}

// Floats64 adds the field key with f as a []float64
func (c *Context) Floats64(key string, val []float64) *Context {
	c.buf = appendFloats64(appendKey(c.buf, key), val)

	return c
}

// Ts adds the current local time as UNIX timestamp
func (c *Context) Ts() *Context {
	c.buf = appendTime(appendKey(c.buf, "time"), time.Now(), time.RFC3339)

	return c
}

// Time adds the field key with time.Time formated as string
func (c *Context) Time(key string, t time.Time) *Context {
	c.buf = appendTime(appendKey(c.buf, key), t, time.RFC3339)

	return c
}

// Times adds the field key with []time.Time formated as string
func (c *Context) Times(key string, t []time.Time) *Context {
	c.buf = appendTimes(appendKey(c.buf, key), t, time.RFC3339)

	return c
}

// Dur adds the field key with time.Duration
func (c *Context) Dur(key string, d time.Duration) *Context {
	c.buf = appendDuration(appendKey(c.buf, key), d, time.Millisecond)

	return c
}

// Durs adds the field key with []time.Duration
func (c *Context) Durs(key string, d []time.Duration) *Context {
	c.buf = appendDurations(appendKey(c.buf, key), d, time.Millisecond)

	return c
}

// TimeDiff adds the field key with positive duration between time t and start
func (c *Context) TimeDiff(key string, t time.Time, start time.Time) *Context {
	var d time.Duration
	if t.After(start) {
		d = t.Sub(start)
	}

	c.buf = appendDuration(appendKey(c.buf, key), d, time.Millisecond)

	return c
}

// If adds the field key with i marshaled using reflection
func (c *Context) If(key string, i interface{}) *Context {
	c.buf = appendInterface(appendKey(c.buf, key), i)

	return c
}

// Obj adds the field key with fmt.Sprintf raw string
// for json string, use If function
func (c *Context) Obj(key string, o interface{}) *Context {
	c.buf = appendObject(appendKey(c.buf, key), o)

	return c
}

// IP adds IPv4 or IPv6 Address
func (c *Context) IP(key string, ip net.IP) *Context {
	c.buf = appendIP(appendKey(c.buf, key), ip)

	return c
}

// IPNet adds IPv4 or IPv6 Prefix (address and mask)
func (c *Context) IPNet(key string, ipn net.IPNet) *Context {
	c.buf = appendIPNet(appendKey(c.buf, key), ipn)

	return c
}

// Mac adds MAC address
func (c *Context) Mac(key string, mac net.HardwareAddr) *Context {
	c.buf = appendMac(appendKey(c.buf, key), mac)

	return c
}
