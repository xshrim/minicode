package gol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// check and deal with multiple output arguments
func formatData(v ...interface{}) string {
	if v == nil || len(v) < 1 {
		return ""
	}
	val := ""

	if len(v) > 1 {
		if format, ok := v[0].(string); ok && isFormatString(format) {
			str := fmt.Sprintf(format, v[1:]...)
			if !bytes.Contains([]byte(str), []byte("%!(")) {
				val = str
			}
		}
	}

	if val == "" {
		val = fmt.Sprint(v...)
	}
	return val
}

// output colorful message to stdout using default logger
func Cprt(color string, v ...interface{}) (int, error) {
	v = append([]interface{}{color}, v...)
	v = append(v, ColorOff)
	return Fprt(os.Stdout, v...)
}

// output colorful format message to stdout using default logger
func Cprtf(color string, format string, v ...interface{}) (int, error) {
	format = "%s" + format + "%s"
	v = append([]interface{}{color}, v...)
	v = append(v, ColorOff)
	return Fprtf(os.Stdout, format, v...)
}

// output colorful message to stdout with newline using default logger
func Cprtln(color string, v ...interface{}) (int, error) {
	line := fmt.Sprintln(v...)
	line = color + line[:len(line)-1] + ColorOff
	return Fprtln(os.Stdout, line)
}

// output message to stdout using default logger
func Prt(v ...interface{}) (int, error) {
	return Fprt(os.Stdout, v...)
}

// output format message to stdout using default logger
func Prtf(format string, v ...interface{}) (int, error) {
	return Fprtf(os.Stdout, format, v...)
}

// output message to stdout with newline using default logger
func Prtln(v ...interface{}) (int, error) {
	return Fprtln(os.Stdout, v...)
}

// return message using default logger
func Sprt(v ...interface{}) string {
	return fmt.Sprint(v...)
}

// return format message using default logger
func Sprtf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

// return format message with newline using default logger
func Sprtln(v ...interface{}) string {
	return fmt.Sprintln(v...)
}

// output message to the writer using default logger
func Fprt(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprint(w, v...)
}

// output format message to the writer using default logger
func Fprtf(w io.Writer, format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, v...)
}

// output format message to the writer with newline using default logger
func Fprtln(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprintln(w, v...)
}

// return error message using default logger
func Err(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}

// return format error message using default logger
func Errf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}

// output colorful message to stdout
func (l *Logger) Cprt(color string, v ...interface{}) (int, error) {
	v = append([]interface{}{color}, v...)
	v = append(v, ColorOff)
	return Fprt(os.Stdout, v...)
}

// output colorful format message to stdout
func (l *Logger) Cprtf(color string, format string, v ...interface{}) (int, error) {
	format = "%s" + format + "%s"
	v = append([]interface{}{color}, v...)
	v = append(v, ColorOff)
	return Fprtf(os.Stdout, format, v...)
}

// output colorful message to stdout with newline
func (l *Logger) Cprtln(color string, v ...interface{}) (int, error) {
	line := fmt.Sprintln(v...)
	line = color + line[:len(line)-1] + ColorOff
	return Fprtln(os.Stdout, line)
}

// output message to stdout
func (l *Logger) Prt(v ...interface{}) (int, error) {
	return Fprt(os.Stdout, v...)
}

// output format message to stdout
func (l *Logger) Prtf(format string, v ...interface{}) (int, error) {
	return Fprtf(os.Stdout, format, v...)
}

// output message to stdout with newline
func (l *Logger) Prtln(v ...interface{}) (int, error) {
	return Fprtln(os.Stdout, v...)
}

// return message
func (l *Logger) Sprt(v ...interface{}) string {
	return fmt.Sprint(v...)
}

// return format message
func (l *Logger) Sprtf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

// return format message with newline
func (l *Logger) Sprtln(v ...interface{}) string {
	return fmt.Sprintln(v...)
}

// output message to the writer
func (l *Logger) Fprt(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprint(w, v...)
}

// output format message to the writer
func (l *Logger) Fprtf(w io.Writer, format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, v...)
}

// output format message to the writer with newline
func (l *Logger) Fprtln(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprintln(w, v...)
}

// return error message
func (l *Logger) Err(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}

// return format error message
func (l *Logger) Errf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}
