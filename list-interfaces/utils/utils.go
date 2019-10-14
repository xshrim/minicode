package utils

import (
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
