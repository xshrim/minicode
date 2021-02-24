package tk

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	rd "crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type V = struct{}
type I = interface{}
type M = map[string]interface{}

const hexs = "0123456789abcdef"

// var noEscapeTable = [256]bool{}

// func init() {
// 	for i := 0; i <= 0x7e; i++ {
// 		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
// 	}
// }

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

func toUpper(str string) string {
	var dst []rune
	for _, v := range str {
		if v >= 97 && v <= 122 {
			v -= 32
		}
		dst = append(dst, v)
	}
	return string(dst)
}

func toLower(str string) string {
	var dst []rune
	for _, v := range str {
		if v >= 65 && v <= 90 {
			v += 32
		}
		dst = append(dst, v)
	}
	return string(dst)
}

func toCapitalize(str string) string {
	var dst []rune
	for idx, v := range str {
		if idx == 0 {
			if v >= 97 && v <= 122 {
				v -= 32
			}
		} else {
			if v >= 65 && v <= 90 {
				v += 32
			}
		}
		dst = append(dst, v)
	}
	return string(dst)
}

func replaceEscapePeriod(str string, flag bool) string {
	var buf []rune
	for _, c := range str {
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

func stringEscapeSep(str string, sep rune) string {
	var buf []rune
	for _, c := range str {
		if c == sep {
			buf = append(buf, '\\')
		}
		buf = append(buf, c)
	}
	return string(buf)
}

func stringRepeat(str string, times int) string {
	out := ""
	for i := 0; i < times; i++ {
		out += str
	}

	return out
}

func stringContainRune(str string, r rune) bool {
	for _, c := range str {
		if c == r {
			return true
		}
	}
	return false
}

func stringIndex(str, sub string) int {
	if len(sub) == 0 {
		return 0
	}
	if len(str) < len(sub) {
		return -1
	}
	for i := 0; i <= len(str)-len(sub); i++ {
		if string(str[i:i+len(sub)]) == sub {
			return i
		}
	}
	return -1
}

func stringContainStr(str, sub string) bool {
	return stringIndex(str, sub) >= 0
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

func stringPrefixStr(str, sub string) bool {
	return stringIndex(str, sub) == 0
}

func stringSuffixStr(str, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(str) < len(sub) {
		return false
	}
	return string(str[len(str)-len(sub):]) == sub
}

func stringSplit(str string, r rune) []string {
	var strs []string
	var runes []rune
	for i, c := range str {
		if c != r {
			runes = append(runes, c)
			if i == len(str)-1 {
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

// convert data to json-like []byte
func tojson(dst []byte, v interface{}) []byte {
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
	case map[string]interface{}:
		dst = append(dst, '{')
		for k, v := range val {
			dst = appendKey(dst, k)
			dst = tojson(dst, v)
		}
		dst = append(dst, '}')
	case []map[string]interface{}:
		dst = append(dst, '[')
		for _, m := range val {
			if dst[len(dst)-1] == '[' || dst[len(dst)-1] == ',' {
				dst = append(dst, '{')
			}
			for k, v := range m {
				dst = appendKey(dst, k)
				dst = tojson(dst, v)
			}
			dst = append(dst, '}')
			dst = append(dst, ',')
		}
		if len(val) > 0 {
			dst = dst[:len(dst)-1]
		}
		dst = append(dst, ']')
	case []interface{}:
		dst = append(dst, '[')
		for _, s := range val {
			dst = tojson(dst, s)
			dst = append(dst, ',')
		}
		if len(val) > 0 {
			dst = dst[:len(dst)-1]
		}
		dst = append(dst, ']')
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
	return dst
}

// marshal json with skipping func fields
func Marshal(v interface{}) ([]byte, error) {
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

// convert data to map[string]interface{}
func Imapify(data interface{}) map[string]interface{} {
	var err error
	m := make(map[string]interface{})

	switch v := data.(type) {
	case string:
		err = json.Unmarshal([]byte(v), &m)
	case []byte:
		err = json.Unmarshal(v, &m)
	default:
		if d, er := json.Marshal(v); er != nil {
			return nil
		} else {
			err = json.Unmarshal(d, &m)
		}
	}

	if err != nil {
		return nil
	}

	return m
}

// modify map[string]interface{}
func Imapset(data map[string]interface{}, keyPath string, val interface{}) error {
	ok := false
	keys := stringSplit(replaceEscapePeriod(keyPath, true), '.')

	for idx, key := range keys {
		kidx := -1
		key = replaceEscapePeriod(key, false)
		if key == "" {
			continue
		}
		if matched, _ := regexp.MatchString(`^\w+\[\d+\]$`, key); matched {
			res := regexp.MustCompile(`^(\w+)\[(\d+)\]$`).FindStringSubmatch(key)
			key = res[1]
			var err error
			if kidx, err = strconv.Atoi(res[2]); err != nil {
				return err
			}
		}

		if idx == len(keys)-1 {
			if val == nil {
				if kidx < 0 {
					delete(data, key)
				} else {
					if reflect.TypeOf(data[key]).Kind() != reflect.Slice {
						return fmt.Errorf("value of keyPath %s is not type []interface{}", keyPath)
					}
					var temp []interface{}
					v := reflect.ValueOf(data[key])
					if kidx >= v.Len() {
						return fmt.Errorf("keyPath %s index out of range", keyPath)
					}
					for i := 0; i < v.Len(); i++ {
						if i != kidx {
							temp = append(temp, v.Index(i).Interface())
						}
					}
					data[key] = temp
				}
			} else {
				if kidx < 0 {
					data[key] = val
				} else {
					if reflect.TypeOf(data[key]).Kind() != reflect.Slice {
						return fmt.Errorf("value of keyPath %s is not type []interface{}", keyPath)
					}
					var temp []interface{}
					v := reflect.ValueOf(data[key])
					if kidx >= v.Len() {
						return fmt.Errorf("keyPath %s index out of range", keyPath)
					}
					for i := 0; i < v.Len(); i++ {
						if i != kidx {
							temp = append(temp, v.Index(i).Interface())
						} else {
							temp = append(temp, val)
						}
					}
					data[key] = temp
				}
			}
		} else {
			if kidx < 0 {
				data, ok = data[key].(map[string]interface{})
				if !ok {
					return fmt.Errorf("value of keyPath %s is not type map[string]interface{}", keyPath)
				}
			} else {
				if reflect.TypeOf(data[key]).Kind() != reflect.Slice {
					return fmt.Errorf("value of keyPath %s is not type []interface{}", keyPath)
				}
				v := reflect.ValueOf(data[key])
				if kidx >= v.Len() {
					return fmt.Errorf("keyPath %s index out of range", keyPath)
				}
				data = v.Index(kidx).Interface().(map[string]interface{})
			}
		}
	}

	return nil
}

// convert data to json-like string
func Jsonify(v interface{}) string {
	return string(tojson(nil, v))
}

// get all leaf key paths of the json string
func Jsleaf(jsonData string, separator ...rune) []string {
	out := []string{}
	paths := Jsdig(jsonData, separator...)
	for idx, path := range paths {
		if idx < len(paths)-1 {
			if !stringContainStr(paths[idx+1], paths[idx]+"[") && !stringContainStr(paths[idx+1], paths[idx]+".") {
				out = append(out, path)
			}
		}
	}

	return out
}

// get all key paths of the json string
func Jsdig(jsonData string, separator ...rune) []string {
	sep := '.'
	if len(separator) > 0 {
		sep = separator[0]
	}

	out := []string{}
	mapDig(&out, "", Imapify(jsonData), sep)

	return out
}

// get differences between a and b json strings and return the corresponding key paths
func Jsdiff(a, b string, separator ...rune) []string {
	out := []string{}
	sep := '.'
	if len(separator) > 0 {
		sep = separator[0]
	}
	paths := Jsleaf(a, sep)
	for _, path := range paths {
		if compareInterface(Jsquery(a, path), Jsquery(b, path)) != 0 {
			out = append(out, path)
		}
	}

	return out
}

// set value of the path key from json string
func Jsmodify(jsonData string, keyPath string, val interface{}) string {
	data := Imapify(jsonData)
	_ = Imapset(data, keyPath, val)
	return Jsonify(data)
}

// get value of the path key from json string
func Jsquery(jsonData string, keyPath string) interface{} {
	var val interface{}

	val = Imapify(jsonData)
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

func mapDig(result *[]string, root string, mp map[string]interface{}, sep rune) {
	for k, v := range mp {
		nroot := fmt.Sprintf("%s%c%s", root, sep, stringEscapeSep(k, sep))
		*result = append(*result, nroot)
		switch val := v.(type) {
		case map[string]interface{}:
			mapDig(result, nroot, val, sep)
		case []interface{}:
			for idx, obj := range val {
				sroot := fmt.Sprintf("%s[%d]", nroot, idx)
				*result = append(*result, sroot)
				if nval, ok := obj.(map[string]interface{}); ok {
					mapDig(result, sroot, nval, sep)
				}
			}
		}
	}
}

// func writeFile(logfile string, mode int, data []string) {
// 	if logfile == "" {
// 		fmt.Println("Can not get log file")
// 		return
// 	}
// 	file, err := os.OpenFile(logfile, mode, 0666)
// 	if err != nil {
// 		fmt.Println("Can not open log file: ", err)
// 		return
// 	}
// 	for _, data := range data {
// 		if _, err = file.WriteString(data); err != nil {
// 			fmt.Println("Write log file failed: ", err)
// 			return
// 		}
// 	}
// 	file.Close()
// }

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
		marshaled, err = Marshal(i)
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

func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func compareInterface(a, b interface{}) int {
	if a == nil && b != nil {
		return -1
	} else if a != nil && b == nil {
		return 1
	} else if a == nil && b == nil {
		return 0
	}

	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	switch aVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case aVal.Int() < bVal.Int():
			return -1
		case aVal.Int() == bVal.Int():
			return 0
		case aVal.Int() > bVal.Int():
			return 1
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch {
		case aVal.Uint() < bVal.Uint():
			return -1
		case aVal.Uint() == bVal.Uint():
			return 0
		case aVal.Uint() > bVal.Uint():
			return 1
		}
	case reflect.Float32, reflect.Float64:
		switch {
		case aVal.Float() < bVal.Float():
			return -1
		case aVal.Float() == bVal.Float():
			return 0
		case aVal.Float() > bVal.Float():
			return 1
		}
	case reflect.String:
		switch {
		case aVal.String() < bVal.String():
			return -1
		case aVal.String() == bVal.String():
			return 0
		case aVal.String() > bVal.String():
			return 1
		}
	default:
		ma, _ := json.Marshal(a)
		mb, _ := json.Marshal(b)
		return bytes.Compare(ma, mb)
	}

	return 0
}

func quickSortBool(list []bool, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && (!pvt || list[chigh]) {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && (pvt || !list[clow]) {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortBool(list, low, pivot-1)
		quickSortBool(list, pivot+1, high)
	}
}

func quickSortRune(list []rune, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortRune(list, low, pivot-1)
		quickSortRune(list, pivot+1, high)
	}
}

func quickSortByte(list []byte, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortByte(list, low, pivot-1)
		quickSortByte(list, pivot+1, high)
	}
}

func quickSortBytes(list [][]byte, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && bytes.Compare(pvt, list[chigh]) <= 0 {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && bytes.Compare(pvt, list[clow]) >= 0 {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortBytes(list, low, pivot-1)
		quickSortBytes(list, pivot+1, high)
	}
}

func quickSortInt(list []int, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortInt(list, low, pivot-1)
		quickSortInt(list, pivot+1, high)
	}
}

func quickSortInt8(list []int8, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortInt8(list, low, pivot-1)
		quickSortInt8(list, pivot+1, high)
	}
}

func quickSortInt16(list []int16, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortInt16(list, low, pivot-1)
		quickSortInt16(list, pivot+1, high)
	}
}

func quickSortInt32(list []int32, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortInt32(list, low, pivot-1)
		quickSortInt32(list, pivot+1, high)
	}
}

func quickSortInt64(list []int64, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortInt64(list, low, pivot-1)
		quickSortInt64(list, pivot+1, high)
	}
}

func quickSortUint(list []uint, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortUint(list, low, pivot-1)
		quickSortUint(list, pivot+1, high)
	}
}

func quickSortUint8(list []uint8, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortUint8(list, low, pivot-1)
		quickSortUint8(list, pivot+1, high)
	}
}

func quickSortUint16(list []uint16, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortUint16(list, low, pivot-1)
		quickSortUint16(list, pivot+1, high)
	}
}

func quickSortUint32(list []uint32, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortUint32(list, low, pivot-1)
		quickSortUint32(list, pivot+1, high)
	}
}

func quickSortUint64(list []uint64, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortUint64(list, low, pivot-1)
		quickSortUint64(list, pivot+1, high)
	}
}

func quickSortFloat32(list []float32, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortFloat32(list, low, pivot-1)
		quickSortFloat32(list, pivot+1, high)
	}
}

func quickSortFloat64(list []float64, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortFloat64(list, low, pivot-1)
		quickSortFloat64(list, pivot+1, high)
	}
}

func quickSortTime(list []time.Time, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && !pvt.After(list[chigh]) {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && !pvt.Before(list[clow]) {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortTime(list, low, pivot-1)
		quickSortTime(list, pivot+1, high)
	}
}

func quickSortDuration(list []time.Duration, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortDuration(list, low, pivot-1)
		quickSortDuration(list, pivot+1, high)
	}
}

func quickSortIP(list []net.IP, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && bytes.Compare(pvt, list[chigh]) <= 0 {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && bytes.Compare(pvt, list[clow]) >= 0 {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortIP(list, low, pivot-1)
		quickSortIP(list, pivot+1, high)
	}
}

func quickSortMac(list []net.HardwareAddr, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && bytes.Compare(pvt, list[chigh]) <= 0 {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && bytes.Compare(pvt, list[clow]) >= 0 {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortMac(list, low, pivot-1)
		quickSortMac(list, pivot+1, high)
	}
}

func quickSortString(list []string, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && pvt <= list[chigh] {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && pvt >= list[clow] {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortString(list, low, pivot-1)
		quickSortString(list, pivot+1, high)
	}
}

func quickSortInterface(list []interface{}, low, high int) {
	if high > low {
		clow := low
		chigh := high
		pvt := list[clow]
		for clow < chigh {
			for clow < chigh && compareInterface(pvt, list[chigh]) <= 0 {
				chigh--
			}
			list[clow] = list[chigh]
			for clow < chigh && compareInterface(pvt, list[clow]) >= 0 {
				clow++
			}
			list[chigh] = list[clow]
		}
		list[clow] = pvt
		pivot := clow

		quickSortInterface(list, low, pivot-1)
		quickSortInterface(list, pivot+1, high)
	}
}

// sort various types of slices with quick sort algorithm
func QuickSort(list interface{}) {
	switch v := list.(type) {
	case []byte:
		quickSortByte(v, 0, len(v)-1)
	case [][]byte:
		quickSortBytes(v, 0, len(v)-1)
	case []bool:
		quickSortBool(v, 0, len(v)-1)
	case []int:
		quickSortInt(v, 0, len(v)-1)
	case []int8:
		quickSortInt8(v, 0, len(v)-1)
	case []int16:
		quickSortInt16(v, 0, len(v)-1)
	case []int32:
		quickSortInt32(v, 0, len(v)-1)
	case []int64:
		quickSortInt64(v, 0, len(v)-1)
	case []uint:
		quickSortUint(v, 0, len(v)-1)
	case []uint16:
		quickSortUint16(v, 0, len(v)-1)
	case []uint32:
		quickSortUint32(v, 0, len(v)-1)
	case []uint64:
		quickSortUint64(v, 0, len(v)-1)
	case []float32:
		quickSortFloat32(v, 0, len(v)-1)
	case []float64:
		quickSortFloat64(v, 0, len(v)-1)
	case []string:
		quickSortString(v, 0, len(v)-1)
	case []time.Time:
		quickSortTime(v, 0, len(v)-1)
	case []time.Duration:
		quickSortDuration(v, 0, len(v)-1)
	case []net.IP:
		quickSortIP(v, 0, len(v)-1)
	case []net.HardwareAddr:
		quickSortMac(v, 0, len(v)-1)
	case []interface{}:
		quickSortInterface(v, 0, len(v)-1)
	}
}

// remove duplicate elements in various types of slices
func Uniq(list interface{}) []interface{} {
	out := []interface{}{}
	switch v := list.(type) {
	case []byte:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(byte) != val {
				out = append(out, val)
			}
		}
	case [][]byte:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || string(out[len(out)-1].([]byte)) != string(val) {
				out = append(out, val)
			}
		}
	case []bool:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(bool) != val {
				out = append(out, val)
			}
		}
	case []int:
		sort.Ints(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(int) != val {
				out = append(out, val)
			}
		}
	case []int8:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(int8) != val {
				out = append(out, val)
			}
		}
	case []int16:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(int16) != val {
				out = append(out, val)
			}
		}
	case []int32:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(int32) != val {
				out = append(out, val)
			}
		}
	case []int64:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(int64) != val {
				out = append(out, val)
			}
		}
	case []uint:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(uint) != val {
				out = append(out, val)
			}
		}
	case []uint16:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(uint16) != val {
				out = append(out, val)
			}
		}
	case []uint32:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(uint32) != val {
				out = append(out, val)
			}
		}
	case []uint64:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(uint64) != val {
				out = append(out, val)
			}
		}
	case []float32:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(float32) != val {
				out = append(out, val)
			}
		}
	case []float64:
		sort.Float64s(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(float64) != val {
				out = append(out, val)
			}
		}
	case []string:
		sort.Strings(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(string) != val {
				out = append(out, val)
			}
		}
	case []time.Time:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(time.Time) != val {
				out = append(out, val)
			}
		}
	case []time.Duration:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(time.Duration) != val {
				out = append(out, val)
			}
		}
	case []net.IP:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(net.IP).String() != val.String() {
				out = append(out, val)
			}
		}
	case []net.HardwareAddr:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1].(net.HardwareAddr).String() != val.String() {
				out = append(out, val)
			}
		}
	case []interface{}:
		QuickSort(v)
		for _, val := range v {
			if len(out) == 0 || out[len(out)-1] != val {
				out = append(out, val)
			}
		}
	}

	return out
}

// return the largest element in various types of slices
func Max(list interface{}) interface{} {
	var out interface{}
	switch v := list.(type) {
	case []byte:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(byte) {
				out = val
			}
		}
	case [][]byte:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if string(val) > string(out.([]byte)) {
				var tmp []byte
				copy(tmp, val)
				out = tmp
			}
		}
	case []bool:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val && !out.(bool) {
				out = val
			}
		}
	case []int:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(int) {
				out = val
			}
		}
	case []int8:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(int8) {
				out = val
			}
		}
	case []int16:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(int16) {
				out = val
			}
		}
	case []int32:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(int32) {
				out = val
			}
		}
	case []int64:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(int64) {
				out = val
			}
		}
	case []uint:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(uint) {
				out = val
			}
		}
	case []uint16:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(uint16) {
				out = val
			}
		}
	case []uint32:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(uint32) {
				out = val
			}
		}
	case []uint64:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(uint64) {
				out = val
			}
		}
	case []float32:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(float32) {
				out = val
			}
		}
	case []float64:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(float64) {
				out = val
			}
		}
	case []string:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(string) {
				out = val
			}
		}
	case []time.Time:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val.After(out.(time.Time)) {
				out = val
			}
		}
	case []time.Duration:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val > out.(time.Duration) {
				out = val
			}
		}
	case []net.IP:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val.String() > out.(net.IP).String() {
				out = val
			}
		}
	case []net.HardwareAddr:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val.String() > out.(net.HardwareAddr).String() {
				out = val
			}
		}
	case []interface{}:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if compareInterface(val, out) > 0 {
				out = val
			}
		}
	}

	return out
}

// return minimum element in various types of slices
func Min(list interface{}) interface{} {
	var out interface{}
	switch v := list.(type) {
	case []byte:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(byte) {
				out = val
			}
		}
	case [][]byte:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if string(val) < string(out.([]byte)) {
				var tmp []byte
				copy(tmp, val)
				out = tmp
			}
		}
	case []bool:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if !val && out.(bool) {
				out = val
			}
		}
	case []int:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(int) {
				out = val
			}
		}
	case []int8:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(int8) {
				out = val
			}
		}
	case []int16:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(int16) {
				out = val
			}
		}
	case []int32:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(int32) {
				out = val
			}
		}
	case []int64:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(int64) {
				out = val
			}
		}
	case []uint:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(uint) {
				out = val
			}
		}
	case []uint16:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(uint16) {
				out = val
			}
		}
	case []uint32:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(uint32) {
				out = val
			}
		}
	case []uint64:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(uint64) {
				out = val
			}
		}
	case []float32:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(float32) {
				out = val
			}
		}
	case []float64:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(float64) {
				out = val
			}
		}
	case []string:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(string) {
				out = val
			}
		}
	case []time.Time:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val.Before(out.(time.Time)) {
				out = val
			}
		}
	case []time.Duration:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val < out.(time.Duration) {
				out = val
			}
		}
	case []net.IP:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val.String() < out.(net.IP).String() {
				out = val
			}
		}
	case []net.HardwareAddr:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if val.String() < out.(net.HardwareAddr).String() {
				out = val
			}
		}
	case []interface{}:
		for idx, val := range v {
			if idx == 0 {
				out = val
			} else if compareInterface(val, out) < 0 {
				out = val
			}
		}
	}

	return out
}

// traverse from start to end with step by iterator
func Iter(v ...int) <-chan int {
	start := 0
	end := start
	step := 1
	c := make(chan int)
	if len(v) == 1 {
		end = v[0]
	} else if len(v) == 2 {
		start = v[0]
		end = v[1]
	} else if len(v) > 2 {
		start = v[0]
		end = v[1]
		step = v[2]
	}

	go func() {
		for start < end {
			c <- start
			start += step
		}
		close(c)
	}()
	return c
}

// return integer slice range start to end with step
func IterS(v ...int) []int {
	start := 0
	end := start
	step := 1
	s := []int{}
	if len(v) == 1 {
		end = v[0]
	} else if len(v) == 2 {
		start = v[0]
		end = v[1]
	} else if len(v) > 2 {
		start = v[0]
		end = v[1]
		step = v[2]
	}

	for start < end {
		s = append(s, start)
		start += step
	}

	return s
}

// reverse elements in the slice or string
func Reverse(list interface{}) interface{} {
	val := reflect.ValueOf(list)
	length := val.Len()
	out := make([]interface{}, length)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < length; i++ {
			out[length-1-i] = val.Index(i).Interface()
		}
	} else if val.Kind() == reflect.String {
		r := []rune(list.(string))
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	}
	return out
}

// return all indices of special element in the slice or string
func Index(list interface{}, v interface{}) []int {
	out := []int{}

	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			if val.Index(i).Interface() == v {
				out = append(out, i)
			}
		}
	} else if val.Kind() == reflect.String && reflect.ValueOf(v).Kind() == reflect.String {
		src := list.(string)
		sub := v.(string)
		if len(sub) == 0 {
			out = append(out, 0)
		} else {
			bf := 0
			for {
				if len(src) >= len(sub) {
					if i := stringIndex(src, sub); i >= 0 {
						out = append(out, i+bf)
						bf += len(src[:i]) + len(sub)
						src = src[i+len(sub):]
					} else {
						break
					}
				} else {
					break
				}
			}
		}
	}

	return out
}

// split the slice or string
func Split(list interface{}, v interface{}) []interface{} {
	out := []interface{}{}
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		tmps := []interface{}{}
		for i := 0; i < val.Len(); i++ {
			tmpv := val.Index(i).Interface()
			if tmpv == v {
				var dst []interface{}
				dst = append(dst, tmps...)
				if len(dst) > 0 {
					out = append(out, dst)
				}
				tmps = []interface{}{}
			} else {
				tmps = append(tmps, tmpv)
			}
		}
		if len(tmps) > 0 {
			out = append(out, tmps)
		}
	} else if val.Kind() == reflect.String && reflect.ValueOf(v).Kind() == reflect.String {
		src := list.(string)
		sub := v.(string)
		if len(sub) == 0 {
			for _, r := range src {
				out = append(out, string(r))
			}
		} else {
			for {
				if len(src) >= len(sub) {
					if i := stringIndex(src, sub); i >= 0 {
						if len(src[:i]) > 0 {
							out = append(out, src[:i])
						}
						src = src[i+len(sub):]
					} else {
						break
					}
				} else {
					break
				}
			}
			if len(src) > 0 {
				out = append(out, src)
			}
		}
	}

	return out
}

// return if special element is in the slice or string
func Contain(list interface{}, v interface{}) bool {
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			if val.Index(i).Interface() == v {
				return true
			}
		}
	} else if val.Kind() == reflect.String && reflect.ValueOf(v).Kind() == reflect.String {
		return stringContainStr(list.(string), v.(string))
	}

	return false
}

// remove special element from the slice or string
func Remove(list interface{}, v interface{}) interface{} {
	out := []interface{}{}

	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			if val.Index(i).Interface() != v {
				out = append(out, val.Index(i).Interface())
			}
		}
	} else if val.Kind() == reflect.String && reflect.ValueOf(v).Kind() == reflect.String {
		src := list.(string)
		sub := v.(string)
		if len(sub) == 0 {
			return list
		}
		for {
			if len(src) >= len(sub) {
				if i := stringIndex(src, sub); i >= 0 {
					src = src[:i] + src[i+len(sub):]
				} else {
					break
				}
			} else {
				break
			}
		}
		return src
	}
	return out
}

// remove the elements at special index from the slice or string
func RemoveAt(list interface{}, idx ...int) interface{} {
	offset := 1
	out := []interface{}{}
	index := 0

	if len(idx) > 0 {
		index = idx[0]
	}
	if len(idx) > 1 {
		offset = idx[1]
	}

	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			if i < index || i >= index+offset {
				out = append(out, val.Index(i).Interface())
			}
		}
	} else if val.Kind() == reflect.String {
		src := list.(string)
		if len(src) >= index+offset {
			src = src[:index] + src[index+offset:]
		}
		return src
	}
	return out
}

// return the elements meet the filter function in the slice or string
func Filter(list interface{}, fn func(interface{}) bool) interface{} {
	var out []interface{}
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			if fn(v) {
				out = append(out, v)
			}
		}
	} else if val.Kind() == reflect.String {
		for _, r := range list.(string) {
			if fn(r) {
				out = append(out, string(r))
			}
		}
	}

	return out
}

// return rune length of the string
func Strlen(str string) int {
	return len([]rune(str))
}

// return substring by rune
func Substr(str string, pos, length int) string {
	runes := []rune(str)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// convert to Upper
func Upper(str string) string {
	return toUpper(str)
}

// convert to Lower
func Lower(str string) string {
	return toLower(str)
}

// convert to Capitalize
func Capitalize(str string) string {
	return toCapitalize(str)
}

// convert byte size to human readable format
func HumanSize(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f%ciB", float64(b)/float64(div), "KM1PEZY"[exp])
}

// convert human readable string to byte size
func ByteSize(str string) (uint64, error) {
	i := strings.IndexFunc(str, func(r rune) bool {
		return r != '.' && !unicode.IsDigit(r)
	})
	var multiplier float64 = 1
	var sizeSuffixes = "BKM1PEZY"
	if i > 0 {
		suffix := str[i:]
		multiplier = 0
		for j := 0; j < len(sizeSuffixes); j++ {
			base := string(sizeSuffixes[j])
			// M, MB, or MiB are all valid.
			switch suffix {
			case base, base + "B", base + "iB":
				sz := 1 << uint(j*10)
				multiplier = float64(sz)
				break
			}
		}
		if multiplier == 0 {
			return 0, fmt.Errorf("invalid multiplier suffix %q, expected one of %s", suffix, []byte(sizeSuffixes))
		}
		str = str[:i]
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil || val < 0 {
		return 0, fmt.Errorf("expected a non-negative number, got %q", str)
	}
	val *= multiplier
	return uint64(math.Ceil(val)), nil
}

// check if target is exist
func IsExist(target string) bool {
	if _, err := os.Stat(target); err == nil {
		return true

	} else if os.IsNotExist(err) {
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

// check if target is directory
func IsDir(target string) bool {
	info, err := os.Stat(target)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return true
	} else {
		return false
	}
}

// list all files and directories in root folder
func ListAll(root string) []string {
	var files []string

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	return files
}

// read file line by line
func IterFile(fpath string) <-chan string {
	f, err := os.Open(fpath)
	if err != nil {
		panic(fmt.Sprintf("read file %s fail: %s", fpath, err.Error()))
	}
	//defer f.Close()

	c := make(chan string)
	go func(fl *os.File) {
		buf := bufio.NewScanner(fl)
		defer fl.Close()

		for {
			if !buf.Scan() {
				break
			}
			c <- buf.Text()
		}
		close(c)
	}(f)

	return c
}

// read the entire contents of the file
func ReadFile(fpath string) []byte {
	f, err := os.Open(fpath)
	if err != nil {
		panic(fmt.Sprintf("open file %s fail: %s", fpath, err.Error()))
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("read file %s fail: %s", fpath, err.Error()))
	}

	return bytes
}

// write to file
func WriteFile(fpath string, data []byte, append ...bool) error {
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

// count word frequency
func WordFrequency(fpath string, order bool, analysis func(string) []string) [][2]interface{} {
	var wordFrequencyMap = make(map[string]int)

	for line := range IterFile(fpath) {
		var arr []string
		if analysis != nil {
			arr = analysis(line)
		} else {
			arr = strings.FieldsFunc(line, func(c rune) bool {
				if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
					return true
				}
				return false
			})
		}

		for _, v := range arr {
			if _, ok := wordFrequencyMap[v]; ok {
				wordFrequencyMap[v] = wordFrequencyMap[v] + 1
			} else {
				wordFrequencyMap[v] = 1
			}
		}
	}

	var wordFrequency [][2]interface{}
	for k, v := range wordFrequencyMap {
		wordFrequency = append(wordFrequency, [2]interface{}{k, v})
	}

	if order {
		sort.Slice(wordFrequency, func(i, j int) bool {
			if wordFrequency[i][1].(int) > wordFrequency[j][1].(int) {
				return true
			} else if wordFrequency[i][1].(int) == wordFrequency[j][1].(int) {
				if wordFrequency[i][0].(string) < wordFrequency[j][0].(string) {
					return true
				}
			}
			return false
		})
	}

	return wordFrequency
}

// generate random intergers
func RandInt(count int, v ...int64) []int64 {
	var min int64 = 0
	var max int64 = 100

	if len(v) == 1 {
		max = v[0]
	} else if len(v) > 1 {
		min = v[0]
		max = v[1]
	}

	out := []int64{}
	if min > max {
		return out
	}

	allCount := make(map[int64]struct{})
	maxBigInt := big.NewInt(max)
	for {
		i, _ := rd.Int(rd.Reader, maxBigInt)
		number := i.Int64()
		if i.Int64() >= min {
			_, ok := allCount[number]
			if !ok {
				out = append(out, number)
				allCount[number] = struct{}{}
			}
		}
		if len(out) >= count {
			return out
		}
	}
}

// generate random strings
func RandString(count int, src ...byte) string {
	rand.Seed(time.Now().UnixNano())
	if len(src) == 0 {
		src = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	idxBits := 6
	idxMask := 1<<idxBits - 1
	idxMax := 63 / idxBits
	b := make([]byte, count)

	for i, cache, remain := count-1, rand.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
		}
		if idx := int(cache) & idxMask; idx < len(src) {
			b[i] = src[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}

	return string(b)
}

// ping ip or domain
func Ping(ip string) bool {
	type ICMP struct {
		Type        uint8
		Code        uint8
		Checksum    uint16
		Identifier  uint16
		SequenceNum uint16
	}

	icmp := ICMP{
		Type: 8,
	}

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer

	_ = binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = checkSum(buffer.Bytes())
	buffer.Reset()
	_ = binary.Write(&buffer, binary.BigEndian, icmp)

	Time, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("ip4:icmp", ip, Time)
	if err != nil {
		return exec.Command("ping", ip, "-c", "2", "-i", "1", "-W", "3").Run() == nil
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return false
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	num, err := conn.Read(recvBuf)
	if err != nil {
		return false
	}

	_ = conn.SetReadDeadline(time.Time{})

	return string(recvBuf[0:num]) != ""
}

// get local ipv4 address
func IPv4() []string {
	out := []string{"127.0.0.1"}
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				out = append(out, ipnet.IP.String())
			}
		}
	}
	return out
}

// get all ipv4 addresses in the range of the cidr
func Hosts(cidr string) []string {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); {
		ips = append(ips, ip.String())

		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}

	if len(ips) < 2 {
		return []string{}
	}
	return ips[1 : len(ips)-1]
}

// get current time
func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// get current timestamp
func StampNow() int64 {
	return time.Now().Unix()
}

// convert time to timestamp
func Time2Stamp(t string) int64 {
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", t, time.Local)
	return stamp.Unix()
}

// convert timestamp to time
func Stamp2Time(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// format duration to string
func FormatDuration(d time.Duration) string {
	return (time.Duration(d.Milliseconds()) * time.Millisecond).String()
}

// parse duration string
func ParseDuration(str string) time.Duration {
	d, _ := time.ParseDuration(str)
	return d
}

// gzip compresses the given data
func Gzip(data []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		panic(err)
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// gunzip uncompresses the given data
func Gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

// archive target folder
func Zip(source, target string, filter ...string) (err error) {
	if isAbs := filepath.IsAbs(source); !isAbs {
		source, err = filepath.Abs(source)
		if err != nil {
			return err
		}
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}

	defer func() {
		if err := zipfile.Close(); err != nil {
			return
			// Errorf("file close error: %s, file: %s", err.Error(), zipfile.Name())
		}
	}()

	zw := zip.NewWriter(zipfile)

	defer func() {
		if err := zw.Close(); err != nil {
			return
			// Errorf("zipwriter close error: %s", err.Error())
		}
	}()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if len(filter) > 0 {
			ism, err := filepath.Match(filter[0], info.Name())

			if err != nil {
				return err
			}

			if ism {
				return nil
			}
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			if err := file.Close(); err != nil {
				return
				// Errorf("file close error: %s, file: %s", err.Error(), file.Name())
			}
		}()
		_, err = io.Copy(writer, file)

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

// unzip target archived file
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		unzippath := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(unzippath, file.Mode())
			if err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(unzippath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// retry
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if err.Error() == "retry-stop" {
			return err
		}

		if attempts--; attempts > 0 {
			// Warnf("retry func error: %s. attemps #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

// base64 encode
func Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

// base64 decode
func Decode(src string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// encrypt src data with aes algorithm
func AesEncrypt(src []byte, keyStr string) ([]byte, error) {
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padnum := block.BlockSize() - len(src)%block.BlockSize()
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	src = append(src, pad...)
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src, nil
}

// decrypt src data with aes algorithm
func AesDecrypt(src []byte, keyStr string) []byte {
	key := []byte(keyStr)
	block, _ := aes.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

// generate asymmetric key pair
func GenKeyPair() (privateKey string, publicKey string, e error) {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rd.Reader)
	if err != nil {
		return "", "", err
	}
	ecPrivateKey, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return "", "", err
	}
	privateKey = base64.StdEncoding.EncodeToString(ecPrivateKey)

	X := priKey.X
	Y := priKey.Y
	xStr, err := X.MarshalText()
	if err != nil {
		return "", "", err
	}
	yStr, err := Y.MarshalText()
	if err != nil {
		return "", "", err
	}
	public := string(xStr) + "+" + string(yStr)
	publicKey = base64.StdEncoding.EncodeToString([]byte(public))
	return
}

// build asymmetric private key from private key string
func BuildPrivateKey(privateKeyStr string) (priKey *ecdsa.PrivateKey, e error) {
	bytes, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}
	priKey, err = x509.ParseECPrivateKey(bytes)
	if err != nil {
		return nil, err
	}
	return
}

// build asymmetric public key from public key string
func BuildPublicKey(publicKeyStr string) (pubKey *ecdsa.PublicKey, e error) {
	bytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return nil, err
	}
	split := strings.Split(string(bytes), "+")
	xStr := split[0]
	yStr := split[1]
	x := new(big.Int)
	y := new(big.Int)
	err = x.UnmarshalText([]byte(xStr))
	if err != nil {
		return nil, err
	}
	err = y.UnmarshalText([]byte(yStr))
	if err != nil {
		return nil, err
	}
	pub := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	pubKey = &pub
	return
}

// sign content by private key string
func Sign(content []byte, privateKeyStr string) (signature string, e error) {
	priKey, err := BuildPrivateKey(privateKeyStr)
	if err != nil {
		return "", err
	}
	r, s, err := ecdsa.Sign(rd.Reader, priKey, []byte(Hash(content)))
	if err != nil {
		return "", err
	}
	rt, _ := r.MarshalText()
	st, _ := s.MarshalText()
	signStr := string(rt) + "+" + string(st)
	signature = hex.EncodeToString([]byte(signStr))
	return
}

// verify sign by public key string
func VerifySign(content []byte, signature string, publicKeyStr string) bool {
	decodeSign, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	split := strings.Split(string(decodeSign), "+")
	rStr := split[0]
	sStr := split[1]
	rr := new(big.Int)
	ss := new(big.Int)
	_ = rr.UnmarshalText([]byte(rStr))
	_ = ss.UnmarshalText([]byte(sStr))
	pubKey, err := BuildPublicKey(publicKeyStr)
	if err != nil {
		return false
	}
	return ecdsa.Verify(pubKey, []byte(Hash(content)), rr, ss)
}

// generate sha256 code for data
func Hash(data []byte) string {
	sum := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(sum[:])
}

// generate md5 code for data
func MD5(data []byte) string {
	sum := md5.Sum(data)
	return fmt.Sprintf("%x", sum)
}

// execute command with realtime output
func Exec(command string, args ...string) <-chan string {
	out := make(chan string, 1000)

	if len(args) == 0 && stringIndex(command, " ") >= 0 {
		args = []string{"-c", command}
		command = "bash"
	}

	go func(c string, a ...string) {
		defer close(out)

		cmd := exec.Command(c, a...)

		stdout, _ := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		_ = cmd.Start()

		scanner := bufio.NewScanner(stdout)
		// scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			m := scanner.Text()
			out <- m
			// fmt.Println(m)
		}
		_ = cmd.Wait()
	}(command, args...)

	return out
}

// parse url params
func UrlParams(rawUrl string) (map[string][]string, error) {
	stUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	m := stUrl.Query()
	return m, nil
}

// send http get request
func HttpGet(rawUrl string, args ...string) (int, []byte) {
	if stringIndex(rawUrl, "http") != 0 {
		rawUrl = "http://" + rawUrl
	}

	// rawUrl = url.QueryEscape(rawUrl)
	client := http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, _ := http.NewRequest("GET", rawUrl, nil)

	// req.Header.Del("Cookie")
	// req.Header.Del("Authorization")
	if len(args) > 0 {
		hds := stringSplit(args[0], '|')
		for _, item := range hds {
			kv := stringSplit(item, ':')
			if len(kv) > 1 {
				switch toLower(kv[0]) {
				case "cookie":
					req.Header.Set("Cookie", kv[1])
				case "auth", "basic", "token":
					if v := stringSplit(kv[1], '/'); len(v) > 0 {
						req.SetBasicAuth(v[0], v[1])
					} else {
						req.Header.Add("Authorization", "Bearer "+v[0])
					}
				case "agent":
					req.Header.Add("User-Agent", kv[1])
				}
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 600, []byte(fmt.Sprintf("request %s failed: %s", rawUrl, err.Error()))
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 600, []byte(fmt.Sprintf("read response %s failed: %s", rawUrl, err.Error()))
	}

	return resp.StatusCode, bodyBytes
}

// send http post request
func HttpPost(rawUrl string, jsonData []byte, args ...string) (int, []byte) {
	if stringIndex(rawUrl, "http") != 0 {
		rawUrl = "http://" + rawUrl
	}

	// rawUrl = url.QueryEscape(rawUrl)
	client := http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, _ := http.NewRequest("POST", rawUrl, bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	// req.Header.Del("Cookie")
	// req.Header.Del("Authorization")
	if len(args) > 0 {
		hds := stringSplit(args[0], '|')
		for _, item := range hds {
			kv := stringSplit(item, ':')
			if len(kv) > 1 {
				switch toLower(kv[0]) {
				case "cookie":
					req.Header.Set("Cookie", kv[1])
				case "auth", "basic", "token":
					if v := stringSplit(kv[1], '/'); len(v) > 0 {
						req.SetBasicAuth(v[0], v[1])
					} else {
						req.Header.Add("Authorization", "Bearer "+v[0])
					}
				case "agent":
					req.Header.Add("User-Agent", kv[1])
				}
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 600, []byte(fmt.Sprintf("request %s failed: %s", rawUrl, err.Error()))
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 600, []byte(fmt.Sprintf("read response %s failed: %s", rawUrl, err.Error()))
	}

	return resp.StatusCode, bodyBytes
}
