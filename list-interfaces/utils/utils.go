package utils

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

// 目录复制
func CopyDir(src, dst string) error {
	cmd := exec.Command("sh", "-c", "cp -a -r "+src+" "+dst)
	// log.Printf("Running cp -a")
	return cmd.Run()
}

// 目录同步
func SyncDir(src, dst string) error {
	cmd := exec.Command("sh", "-c", "rsync -avz "+src+" "+dst)
	// log.Printf("Running rsync -avz")
	return cmd.Run()
}

// 生成随机数
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func ColorString(depth int, kind, name, content, location string) string {
	if kind == "color" {
		switch depth {
		case 1:
			name = fmt.Sprintf("\033[1;31;40m%s\033[0m", name)
		case 2:
			name = fmt.Sprintf("\033[1;34;40m%s\033[0m", name)
		case 3:
			name = fmt.Sprintf("\033[1;36;40m%s\033[0m", name)
		case 4:
			name = fmt.Sprintf("\033[1;32;40m%s\033[0m", name)
		}
	}

	prefix := ""
	for depth > 1 {
		prefix = prefix + "|   "
		depth--
	}
	prefix = prefix + "|── "

	if kind == "color" {
		prefix = fmt.Sprintf("\033[1;35;40m%s\033[0m", prefix)

		content = fmt.Sprintf("\033[1;33;40m%s\033[0m", content)

		location = fmt.Sprintf("\033[1;37;40m%s\033[0m", location)
	}
	// TODO 检测平台
	return prefix + name + content + location
}
