package gol

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

// util functions

// // Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding
// func itoa(i int, wid int) []byte {
// 	// Assemble decimal in reverse order
// 	var b [20]byte
// 	bp := len(b) - 1
// 	for i >= 10 || wid > 1 {
// 		wid--
// 		q := i / 10
// 		b[bp] = byte('0' + i - q*10)
// 		bp--
// 		i = q
// 	}
// 	// i < 10
// 	b[bp] = byte('0' + i)
// 	return b[bp:]
// }

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

func isFormatString(s string) bool {
	for idx, c := range s {
		if c == '%' && idx != len(s)-1 {
			if idx == 0 || (idx > 0 && s[idx-1] != '\\') {
				return true
			}
		}
	}
	return false
}

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

func replaceEscapePeriod(s string, flag bool) string {
	var buf []rune
	for _, c := range s {
		if flag {
			if c == '.' {
				l := len(buf)
				if l > 0 && buf[l-1] == '\\' {
					buf = buf[:l-1]
					buf = append(buf, '`')
					continue
				}
			}
		} else {
			if c == '`' {
				buf = append(buf, '.')
				continue
			}
		}
		buf = append(buf, c)
	}
	return string(buf)
}

func stringContainRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

func stringIndex(s, t string) int {
	if len(t) == 0 {
		return 0
	}
	if len(s) < len(t) {
		return -1
	}
	for i := 0; i <= len(s)-len(t); i++ {
		if string(s[i:i+len(t)]) == t {
			return i
		}
	}
	return -1
}

func stringContainStr(s, t string) bool {
	return stringIndex(s, t) >= 0
	// sr := []rune(s)
	// tr := []rune(t)
	// for i := 0; i <= len(sr)-len(tr); i++ {
	// 	if sr[i] == tr[0] {
	// 		j := 1
	// 		for ; j < len(tr); j++ {
	// 			if sr[i+j] != tr[j] {
	// 				break
	// 			}
	// 		}
	// 		if j == len(tr) {
	// 			return true
	// 		}
	// 	}
	// }
	// return false
}

func stringPrefixStr(s, t string) bool {
	return stringIndex(s, t) == 0
}

func stringSuffixStr(s, t string) bool {
	if len(t) == 0 {
		return true
	}
	if len(s) < len(t) {
		return false
	}
	return string(s[len(s)-len(t):]) == t
}

func stringSplit(s string, r rune) []string {
	var strs []string
	var runes []rune
	for i, c := range s {
		if c != r {
			runes = append(runes, c)
			if i == len(s)-1 {
				strs = append(strs, string(runes))
				break
			}
		} else {
			if runes != nil {
				strs = append(strs, string(runes))
			}
			runes = nil
		}
	}
	return strs
}

func mapi2maps(i interface{}) interface{} {
	// var body interface{}
	// _ = yaml.Unmarshal([]byte(yamlstr), &body)
	// body = mapi2maps(body)
	// jsb,  _:= json.Marshal(body);
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = mapi2maps(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = mapi2maps(v)
		}
	}
	return i
}

// func map2str(fds F) string {
// 	if fds == nil {
// 		return ""
// 	}

// 	res := ""
// 	for k, v := range fds {
// 		val := fmt.Sprintf("%v", v)
// 		if bytes.ContainsRune([]byte(val), ' ') {
// 			val = "'" + val + "'"
// 		}
// 		res += k + "=" + val + " "
// 	}
// 	return res[:len(res)-1]
// }

func colorStatusCode(statusCode int) string {
	var buff []byte

	switch {
	case statusCode < 200:
		buff = append(buff, WYellow...)
	case statusCode < 300:
		buff = append(buff, WGreen...)
	case statusCode < 400:
		buff = append(buff, WBlue...)
	case statusCode < 500:
		buff = append(buff, WRed...)
	case statusCode < 600:
		buff = append(buff, WPurple...)
	default:
		buff = append(buff, WCyan...)
	}
	itoa(&buff, statusCode, 3)
	buff = append(buff, ColorOff...)
	return string(buff)
}

func colorRequestMethod(mtd string) string {
	var buff []byte
	switch mtd {
	case "GET":
		buff = append(buff, WGreen...)
	case "POST":
		buff = append(buff, WBlue...)
	case "DELETE":
		buff = append(buff, WRed...)
	case "PUT":
		buff = append(buff, WPurple...)
	case "PATCH":
		buff = append(buff, WYellow...)
	default:
		buff = append(buff, WCyan...)
	}
	buff = append(buff, mtd...)
	for i := 5 - len(mtd); i > 0; i-- {
		buff = append(buff, ' ')
	}
	buff = append(buff, ColorOff...)
	return string(buff)
}

func map2json(dst []byte, fds F) []byte {
	for k, v := range fds {
		// append key
		dst = appendKey(dst, k)

		// append value
		switch val := v.(type) {
		case string:
			dst = appendStr(dst, val)
		case []string:
			dst = appendStrs(dst, val)
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
		case time.Time:
			dst = appendTime(dst, val, time.RFC3339)
		case []time.Time:
			dst = appendTimes(dst, val, time.RFC3339)
		case time.Duration:
			dst = appendDuration(dst, val, time.Millisecond)
		case []time.Duration:
			dst = appendDurations(dst, val, time.Millisecond)
		case interface{}:
			dst = appendInterface(dst, val)
		case net.IP:
			dst = appendIPAddr(dst, val)
		case net.IPNet:
			dst = appendIPPrefix(dst, val)
		case net.HardwareAddr:
			dst = appendMACAddr(dst, val)
		default:
			dst = appendObject(dst, val)
		}
	}
	return dst
}

// convert data to json-like []byte
func tojson(dst []byte, v interface{}) []byte {
	switch val := v.(type) {
	case string:
		dst = appendStr(dst, val)
	case []string:
		dst = appendStrs(dst, val)
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
	case time.Time:
		dst = appendTime(dst, val, time.RFC3339)
	case []time.Time:
		dst = appendTimes(dst, val, time.RFC3339)
	case time.Duration:
		dst = appendDuration(dst, val, time.Millisecond)
	case []time.Duration:
		dst = appendDurations(dst, val, time.Millisecond)
	case map[string]interface{}:
		for k, v := range val {
			dst = appendKey(dst, k)
			dst = tojson(dst, v)
		}
		dst = append(dst, '}')
	case []map[string]interface{}:
		dst = append(dst, '[')
		for _, m := range val {
			for k, v := range m {
				if dst[len(dst)-1] == '[' {
					dst = append(dst, '{')
				}
				dst = appendKey(dst, k)
				dst = tojson(dst, v)
			}
			dst = append(dst, '}')
			dst = append(dst, ',')
		}
		dst = dst[:len(dst)-1]
		dst = append(dst, ']')
	case F:
		for k, v := range val {
			dst = appendKey(dst, k)
			dst = tojson(dst, v)
		}
		dst = append(dst, '}')
	case []F:
		dst = append(dst, '[')
		for _, f := range val {
			for k, v := range f {
				if dst[len(dst)-1] == '[' {
					dst = append(dst, '{')
				}
				dst = appendKey(dst, k)
				dst = tojson(dst, v)
			}
			dst = append(dst, '}')
			dst = append(dst, ',')
		}
		dst = dst[:len(dst)-1]
		dst = append(dst, ']')
	case []interface{}:
		dst = append(dst, '[')
		for _, s := range val {
			dst = tojson(dst, s)
			dst = append(dst, ',')
		}
		dst = dst[:len(dst)-1]
		dst = append(dst, ']')
	case interface{}:
		dst = appendInterface(dst, val)
	case net.IP:
		dst = appendIPAddr(dst, val)
	case net.IPNet:
		dst = appendIPPrefix(dst, val)
	case net.HardwareAddr:
		dst = appendMACAddr(dst, val)
	default:
		dst = appendObject(dst, val)
	}
	return dst
}

// convert data to json-like string
func Jsonify(v interface{}) string {
	return string(tojson(nil, v))
}

// get value of the path key from json string
func Jsquery(jsonData string, keyPath string) interface{} {
	var val interface{}
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonData), &m); err != nil {
		return nil
	}
	val = m
	keyPath = replaceEscapePeriod(keyPath, true)
	for _, p := range stringSplit(keyPath, '.') {
		p = replaceEscapePeriod(p, false)
		if matched, _ := regexp.MatchString(`^\[\d+\]$`, p); matched {
			if data, ok := val.([]interface{}); ok {
				if len(data) == 0 {
					val = data
					continue
				}
				sp := regexp.MustCompile(`^\[(\d+)\]$`).FindStringSubmatch(p)
				val = getJsonItem(data, sp[1])
			} else {
				return nil
			}
		} else if matched, _ := regexp.MatchString(`^\[\w+\s?-?\w?\]$`, p); matched {
			if data, ok := val.([]interface{}); ok {
				if len(data) == 0 {
					val = data
					continue
				}
				sp := regexp.MustCompile(`^\[(\w+\s?-?\w?)\]$`).FindStringSubmatch(p)
				val = getJsonItem(data, sp[1])
			} else {
				return nil
			}
		} else if matched, _ := regexp.MatchString(`^\*?\w+\*?#?\[\d+\]$`, p); matched {
			if data, ok := val.(map[string]interface{}); ok {
				sp := regexp.MustCompile(`^(\*?\w+\*?#?)\[(\d+)\]$`).FindStringSubmatch(p)
				val = getJsonVal(data, sp[1])
				if data, ok := val.([]interface{}); ok {
					val = getJsonItem(data, sp[2])
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else if matched, _ := regexp.MatchString(`^\*?\w+\*?#?\[\w+\s?-?\w?\]$`, p); matched {
			if data, ok := val.(map[string]interface{}); ok {
				sp := regexp.MustCompile(`^(\*?\w+\*?#?)\[(\w+\s?-?\w?)\]$`).FindStringSubmatch(p)
				val = getJsonVal(data, sp[1])
				if data, ok := val.([]interface{}); ok {
					if len(data) == 0 {
						val = data
						continue
					}
					val = getJsonItem(data, sp[2])
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			if data, ok := val.(map[string]interface{}); ok {
				val = getJsonVal(data, p)
			} else {
				return nil
			}
		}
	}
	return val
}

func getJsonVal(data map[string]interface{}, p string) interface{} {
	var key string
	var i int
	ki := stringSplit(p, '#')
	key = ki[0]
	if len(ki) > 1 {
		i, _ = strconv.Atoi(ki[1])
	}
	var keys []string
	if key[0] == '*' && key[len(key)-1] == '*' {
		for k := range data {
			if stringContainStr(k, key) {
				keys = append(keys, k)
			}
		}
	} else if key[0] == '*' {
		for k := range data {
			if stringSuffixStr(k, key[1:]) {
				keys = append(keys, k)
			}
		}
	} else if key[len(key)-1] == '*' {
		for k := range data {
			if stringPrefixStr(k, key[:len(key)-1]) {
				keys = append(keys, k)
			}
		}
	} else {
		keys = append(keys, key)
	}
	if i >= len(keys) {
		return nil
	}
	return data[keys[i]]
}

func getJsonItem(data []interface{}, p string) interface{} {
	if p == "first" {
		return data[0]
	} else if p == "last" {
		return data[len(data)-1]
	} else if p == "odd" {
		var tdata []interface{}
		for idx, d := range data {
			if idx%2 == 0 {
				tdata = append(tdata, d)
			}
		}
		return tdata
	} else if p == "even" {
		var tdata []interface{}
		for idx, d := range data {
			if idx%2 == 1 {
				tdata = append(tdata, d)
			}
		}
		return tdata
	} else if p == "#" || p == "len" {
		return len(data)
	} else if stringContainRune(p, ' ') {
		var tdata []interface{}
		for _, i := range stringSplit(p, ' ') {
			if idx, err := strconv.Atoi(i); err == nil {
				if idx < len(data) {
					tdata = append(tdata, data[idx])
				}
			}
		}
		return tdata
	} else if matched, _ := regexp.MatchString(`^\d+$`, p); matched {
		idx, _ := strconv.Atoi(p)
		if idx >= len(data) {
			return nil
		}
		return data[idx]
	} else if matched, _ := regexp.MatchString(`^\d+-\d+$`, p); matched {
		ssp := regexp.MustCompile(`^(\d+)-(\d+)$`).FindStringSubmatch(p)
		start, _ := strconv.Atoi(ssp[1])
		end, _ := strconv.Atoi(ssp[2])
		if start >= len(data) {
			start = len(data) - 1
		}
		if end >= len(data) {
			end = len(data) - 1
		}
		if start == end {
			return data[start]
		} else {
			if start < end {
				return data[start:end]
			} else {
				return data[end:start]
			}
		}
	} else {
		return nil
	}
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

func appendKey(dst []byte, key string) []byte {
	if len(dst) == 0 {
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
		if noEscapeTable[b] {
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
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
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
		if noEscapeTable[b] {
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
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
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
		if !noEscapeTable[str[i]] {
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
		if !noEscapeTable[bs[i]] {
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
		dst = append(dst, hex[v>>4], hex[v&0x0f])
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
		return appendStr(dst, fmt.Sprintf("marshaling error: %v", err))
	} else {
		return append(dst, marshaled...)
	}
}

func appendObject(dst []byte, o interface{}) []byte {
	return appendStr(dst, fmt.Sprintf("%v", o))
}

func appendIPAddr(dst []byte, ip net.IP) []byte {
	return appendStr(dst, ip.String())
}

func appendIPPrefix(dst []byte, pfx net.IPNet) []byte {
	return appendStr(dst, pfx.String())
}

func appendMACAddr(dst []byte, ha net.HardwareAddr) []byte {
	return appendStr(dst, ha.String())
}
