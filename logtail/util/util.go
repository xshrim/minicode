// Copyright (c) 2015 HPE Software Inc. All rights reserved.
// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

type Logger struct {
	*log.Logger
}

var LOGGER = &Logger{log.New(os.Stderr, "", log.LstdFlags)}

// fatal is like panic except it displays only the current goroutine's stack.
func Fatal(format string, v ...interface{}) {
	// https://github.com/hpcloud/log/blob/master/log.go#L45
	LOGGER.Output(2, fmt.Sprintf("FATAL -- "+format, v...)+"\n"+string(debug.Stack()))
	os.Exit(1)
}

// partitionString partitions the string into chunks of given size,
// with the last chunk of variable size.
func PartitionString(s string, chunkSize int) []string {
	if chunkSize <= 0 {
		panic("invalid chunkSize")
	}
	length := len(s)
	chunks := 1 + length/chunkSize
	start := 0
	end := chunkSize
	parts := make([]string, 0, chunks)
	for {
		if end > length {
			end = length
		}
		parts = append(parts, s[start:end])
		if end == length {
			break
		}
		start, end = end, end+chunkSize
	}
	return parts
}

func getFileType(name string) string {
	filename := name
	// Check if the path requested is a symbolic link
	fi, err := os.Lstat(name)
	if err != nil {
		return ""
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		filename, err = os.Readlink(name)
		if err != nil {
			return ""
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()
	typ, err := getFileContentType(f)
	if err != nil {
		return ""
	}
	return typ
}

func getFileContentType(seeker io.ReadSeeker) (string, error) {
	// At most the first 512 bytes of data are used:
	// https://golang.org/src/net/http/sniff.go?s=646:688#L11
	buff := make([]byte, 512)

	_, err := seeker.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	bytesRead, err := seeker.Read(buff)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Slice to remove fill-up zero values which cause a wrong content type detection in the next step
	buff = buff[:bytesRead]

	return http.DetectContentType(buff), nil
}

func IsFitFile(filename string) bool {
	if filename != "" && filename[0] != '.' && strings.Contains(getFileType(filename), "text/plain") {
		return true
	}
	return false
}
