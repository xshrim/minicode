package gol

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/xshrim/gol/color"
)

func replaceDoubleQuote(buf *[]byte, s string) {
	last := false
	for _, c := range []byte(s) {
		if c != '"' {
			if c == '\\' {
				last = true
			} else {
				last = false
			}
		} else {
			if !last {
				*buf = append(*buf, '\\')
			}
			last = false
		}
		*buf = append(*buf, c)
	}
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

func toLower(src string) string {
	var dst []rune
	for _, v := range src {
		if v >= 65 && v <= 90 {
			v += 32
		}
		dst = append(dst, v)
	}
	return string(dst)
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

func parseLevel(level interface{}) int {
	var lv int
	switch res := level.(type) {
	case int:
		lv = res
		if lv < 0 {
			lv = -1
		} else if lv > ALL {
			lv = OFF
		}
	case string:
		str := toUpper(res)
		switch str {
		case "OFF", "0":
			lv = OFF
		case "PANIC", "1":
			lv = PANIC
		case "FATAL", "2":
			lv = PANIC
		case "ERROR", "3":
			lv = ERROR
		case "WARN", "4":
			lv = WARN
		case "NOTIC", "5":
			lv = NOTIC
		case "INFO", "6":
			lv = INFO
		case "DEBUG", "7":
			lv = DEBUG
		case "TRACE", "8":
			lv = TRACE
		case "ALL", "9":
			lv = ALL
		default:
			lv = -1
		}
	default:
		lv = INFO
	}
	return lv
}

func colorStatusCode(statusCode int) string {
	var buff []byte

	switch {
	case statusCode < 200:
		buff = append(buff, color.WYellow...)
	case statusCode < 300:
		buff = append(buff, color.WGreen...)
	case statusCode < 400:
		buff = append(buff, color.WBlue...)
	case statusCode < 500:
		buff = append(buff, color.WRed...)
	case statusCode < 600:
		buff = append(buff, color.WPurple...)
	default:
		buff = append(buff, color.WCyan...)
	}
	itoa(&buff, statusCode, 3)
	buff = append(buff, color.ColorOff...)
	return string(buff)
}

func colorRequestMethod(mtd string) string {
	var buff []byte
	switch mtd {
	case "GET":
		buff = append(buff, color.WGreen...)
	case "POST":
		buff = append(buff, color.WBlue...)
	case "DELETE":
		buff = append(buff, color.WRed...)
	case "PUT":
		buff = append(buff, color.WPurple...)
	case "PATCH":
		buff = append(buff, color.WYellow...)
	default:
		buff = append(buff, color.WCyan...)
	}
	buff = append(buff, mtd...)
	for i := 5 - len(mtd); i > 0; i-- {
		buff = append(buff, ' ')
	}
	buff = append(buff, color.ColorOff...)
	return string(buff)
}

func writeFile(fpath string, data []byte, append ...bool) error {
	mode := os.O_RDWR | os.O_CREATE
	if len(append) > 0 && append[0] {
		mode = mode | os.O_APPEND
	} else {
		mode = mode | os.O_TRUNC
	}
	file, err := os.OpenFile(fpath, mode, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

func appendKey(dst []byte, key string) []byte {
	if len(dst) == 0 || dst[len(dst)-1] == '[' || dst[len(dst)-1] == ':' {
		dst = append(dst, '{')
	}

	// if c.buf[len(c.buf)-1] == '}' {
	// 	c.buf = c.buf[:len(c.buf)-1]
	// }

	if dst[len(dst)-1] != '{' {
		dst = append(dst, ',')
	}
	dst = appendStr(dst, key)
	return append(dst, ':')
}

func appendStrComplex(dst []byte, s string, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if b >= 0x20 && b != '\\' && b != '"' {
			i++
			continue
		}
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hexs[b>>4], hexs[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}

func appendBytesComplex(dst []byte, s []byte, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRune(s[i:])
			if r == utf8.RuneError && size == 1 {
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if b >= 0x20 && b != '\\' && b != '"' {
			i++
			continue
		}
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hexs[b>>4], hexs[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}

func appendStr(dst []byte, str string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(str); i++ {
		if !(str[i] >= 0x20 && str[i] != '\\' && str[i] != '"') {
			dst = appendStrComplex(dst, str, i)
			return append(dst, '"')
		}
	}
	dst = append(dst, str...)
	return append(dst, '"')
}

func appendStrs(dst []byte, vals []string) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendStr(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendStr(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendBytes(dst []byte, bs []byte) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(bs); i++ {
		if !(bs[i] >= 0x20 && bs[i] != '\\' && bs[i] != '"') {
			dst = appendBytesComplex(dst, bs, i)
			return append(dst, '"')
		}
	}
	dst = append(dst, bs...)
	return append(dst, '"')
}

func appendHex(dst []byte, s []byte) []byte {
	dst = append(dst, '"')
	for _, v := range s {
		dst = append(dst, hexs[v>>4], hexs[v&0x0f])
	}
	return append(dst, '"')
}

func appendJson(dst []byte, j []byte) []byte {
	return append(dst, j...)
}

func appendBool(dst []byte, b bool) []byte {
	if b {
		return append(dst, "true"...)
	} else {
		return append(dst, "false"...)
	}
}

func appendBools(dst []byte, vals []bool) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendBool(dst, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendBool(append(dst, ','), val)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt(dst []byte, val int) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func appendInts(dst []byte, vals []int) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt8(dst []byte, val int8) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func appendInts8(dst []byte, vals []int8) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt16(dst []byte, val int16) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func appendInts16(dst []byte, vals []int16) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt32(dst []byte, val int32) []byte {
	return strconv.AppendInt(dst, int64(val), 10)
}

func appendInts32(dst []byte, vals []int32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), int64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInt64(dst []byte, val int64) []byte {
	return strconv.AppendInt(dst, val, 10)
}

func appendInts64(dst []byte, vals []int64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint(dst []byte, val uint) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func appendUints(dst []byte, vals []uint) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint8(dst []byte, val uint8) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func appendUints8(dst []byte, vals []uint8) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint16(dst []byte, val uint16) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func appendUints16(dst []byte, vals []uint16) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint32(dst []byte, val uint32) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func appendUints32(dst []byte, vals []uint32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), uint64(val), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUint64(dst []byte, val uint64) []byte {
	return strconv.AppendUint(dst, uint64(val), 10)
}

func appendUints64(dst []byte, vals []uint64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendUint(dst, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = strconv.AppendUint(append(dst, ','), val, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendFloat(dst []byte, val float64, bitSize int) []byte {
	switch {
	case math.IsNaN(val):
		return append(dst, `"NaN"`...)
	case math.IsInf(val, 1):
		return append(dst, `"+Inf"`...)
	case math.IsInf(val, -1):
		return append(dst, `"-Inf"`...)
	default:
		return strconv.AppendFloat(dst, val, 'f', -1, bitSize)
	}
}

func appendFloat32(dst []byte, val float32) []byte {
	return appendFloat(dst, float64(val), 32)
}

func appendFloats32(dst []byte, vals []float32) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendFloat(dst, float64(vals[0]), 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendFloat(append(dst, ','), float64(val), 32)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendFloat64(dst []byte, val float64) []byte {
	return appendFloat(dst, val, 64)
}

func appendFloats64(dst []byte, vals []float64) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendFloat(dst, vals[0], 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendFloat(append(dst, ','), val, 64)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendTime(dst []byte, t time.Time, format string) []byte {
	switch format {
	case "":
		return appendInt64(dst, t.Unix())
	case "UNIXMS":
		return appendInt64(dst, t.UnixNano()/1000000)
	case "UNIXMICRO":
		return appendInt64(dst, t.UnixNano()/1000)
	default:
		return append(t.AppendFormat(append(dst, '"'), format), '"')
	}
}

func appendTimes(dst []byte, vals []time.Time, format string) []byte {
	switch format {
	case "":
		return appendUnixTimes(dst, vals)
	case "UNIXMS":
		return appendUnixMsTimes(dst, vals)
	}
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = append(vals[0].AppendFormat(append(dst, '"'), format), '"')
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = append(t.AppendFormat(append(dst, ',', '"'), format), '"')
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUnixTimes(dst []byte, vals []time.Time) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, vals[0].Unix(), 10)
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), t.Unix(), 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendUnixMsTimes(dst []byte, vals []time.Time) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = strconv.AppendInt(dst, vals[0].UnixNano()/1000000, 10)
	if len(vals) > 1 {
		for _, t := range vals[1:] {
			dst = strconv.AppendInt(append(dst, ','), t.UnixNano()/1000000, 10)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendDuration(dst []byte, d time.Duration, unit time.Duration) []byte {
	return appendFloat64(dst, float64(d)/float64(unit))
}

func appendDurations(dst []byte, vals []time.Duration, unit time.Duration) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendDuration(dst, vals[0], unit)
	if len(vals) > 1 {
		for _, d := range vals[1:] {
			dst = appendDuration(append(dst, ','), d, unit)
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendInterface(dst []byte, i interface{}) []byte {
	marshaled, err := json.Marshal(i)
	if err != nil {
		marshaled, err = marshal(i)
		if err != nil {
			return appendStr(dst, fmt.Sprintf("marshaling error: %v", err))
		}
	}

	return append(dst, marshaled...)
}

func appendObject(dst []byte, o interface{}) []byte {
	return appendStr(dst, fmt.Sprintf("%v", o))
}

func appendIP(dst []byte, ip net.IP) []byte {
	return appendStr(dst, ip.String())
}

func appendIPs(dst []byte, vals []net.IP) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendStr(dst, vals[0].String())
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendStr(dst, val.String())
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendIPNet(dst []byte, ipn net.IPNet) []byte {
	return appendStr(dst, ipn.String())
}

func appendIPNets(dst []byte, vals []net.IPNet) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendStr(dst, vals[0].String())
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendStr(dst, val.String())
		}
	}
	dst = append(dst, ']')
	return dst
}

func appendMac(dst []byte, mac net.HardwareAddr) []byte {
	return appendStr(dst, mac.String())
}

func appendMacs(dst []byte, vals []net.HardwareAddr) []byte {
	if len(vals) == 0 {
		return append(dst, '[', ']')
	}
	dst = append(dst, '[')
	dst = appendStr(dst, vals[0].String())
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			dst = appendStr(dst, val.String())
		}
	}
	dst = append(dst, ']')
	return dst
}

func marshal(v interface{}) ([]byte, error) {
	value := reflect.Indirect(reflect.ValueOf(v))
	typ := value.Type()
	if typ.Kind() == reflect.Struct {
		sf := make([]reflect.StructField, 0)
		for i := 0; i < typ.NumField(); i++ {
			sf = append(sf, typ.Field(i))
			if typ.Field(i).Type.Kind() == reflect.Func {
				sf[i].Tag = `json:"-"`
			}
		}
		newType := reflect.StructOf(sf)
		newValue := value.Convert(newType)
		return json.Marshal(newValue.Interface())
	} else {
		return json.Marshal(v)
	}
}

func map2json(dst []byte, fds map[string]interface{}) []byte {
	for k, v := range fds {
		// append key
		dst = appendKey(dst, k)

		// append value
		switch val := v.(type) {
		case []byte:
			dst = appendBytes(dst, val)
		case error:
			dst = appendStr(dst, val.Error())
		case bool:
			dst = appendBool(dst, val)
		case []bool:
			dst = appendBools(dst, val)
		case int:
			dst = appendInt(dst, val)
		case []int:
			dst = appendInts(dst, val)
		case int8:
			dst = appendInt8(dst, val)
		case []int8:
			dst = appendInts8(dst, val)
		case int16:
			dst = appendInt16(dst, val)
		case []int16:
			dst = appendInts16(dst, val)
		case int32:
			dst = appendInt32(dst, val)
		case []int32:
			dst = appendInts32(dst, val)
		case int64:
			dst = appendInt64(dst, val)
		case []int64:
			dst = appendInts64(dst, val)
		case uint:
			dst = appendUint(dst, val)
		case []uint:
			dst = appendUints(dst, val)
		case uint8:
			dst = appendUint8(dst, val)
		case uint16:
			dst = appendUint16(dst, val)
		case []uint16:
			dst = appendUints16(dst, val)
		case uint32:
			dst = appendUint32(dst, val)
		case []uint32:
			dst = appendUints32(dst, val)
		case uint64:
			dst = appendUint64(dst, val)
		case []uint64:
			dst = appendUints64(dst, val)
		case float32:
			dst = appendFloat32(dst, val)
		case []float32:
			dst = appendFloats32(dst, val)
		case float64:
			dst = appendFloat64(dst, val)
		case []float64:
			dst = appendFloats64(dst, val)
		case string:
			dst = appendStr(dst, val)
		case []string:
			dst = appendStrs(dst, val)
		case time.Time:
			dst = appendTime(dst, val, time.RFC3339)
		case []time.Time:
			dst = appendTimes(dst, val, time.RFC3339)
		case time.Duration:
			dst = appendDuration(dst, val, time.Millisecond)
		case []time.Duration:
			dst = appendDurations(dst, val, time.Millisecond)
		case net.IP:
			dst = appendIP(dst, val)
		case []net.IP:
			dst = appendIPs(dst, val)
		case net.IPNet:
			dst = appendIPNet(dst, val)
		case []net.IPNet:
			dst = appendIPNets(dst, val)
		case net.HardwareAddr:
			dst = appendMac(dst, val)
		case []net.HardwareAddr:
			dst = appendMacs(dst, val)
		case nil:
			dst = append(dst, []byte("null")...)
		case interface{}:
			dst = appendInterface(dst, val)
		default:
			dst = appendObject(dst, val)
		}
	}

	if len(dst) > 0 && dst[0] == '{' && dst[len(dst)-1] != '}' {
		dst = append(dst, '}')
	}
	return dst
}
