package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"./codeanalysis"
)

func CopyDir(src, dst string) error {
	cmd := exec.Command("sh", "-c", "cp -a -r "+src+" "+dst)
	// log.Printf("Running cp -a")
	return cmd.Run()
}

func main() {

	codeDir := flag.String("codedir", ".", "要扫描的代码目录")
	outputFile := flag.String("output", "list-interface-output.txt", "解析结果保存到该文件中")
	ignoreDir := flag.String("ignore", "", "需要排除的目录,不需要扫描和解析")

	gopathDir := os.Getenv("GOPATH")

	if gopathDir == "" {
		fmt.Println("GOPATH目录不存在")
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		fmt.Println("使用例子\n" +
			os.Args[0] + " --codedir /appdev/gopath/src/github.com/contiv/netplugin --output /tmp/result")
		os.Exit(1)
	}

	flag.Parse()

	if *codeDir == "" {
		fmt.Println("代码目录不能为空")
		os.Exit(1)
	}

	if gopathDir == "" {
		fmt.Println("GOPATH目录不能为空")
		os.Exit(1)
	}

	tmpDir := ""
	if !strings.HasPrefix(*codeDir, gopathDir) {
		// panic(fmt.Sprintf("代码目录%s,必须是GOPATH目录%s的子目录", *codeDir, *gopathDir))
		tmpDir = path.Join(gopathDir, "/src/tmp/")

		if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
			os.RemoveAll(tmpDir)
		}

		os.Mkdir(tmpDir, os.ModePerm)

		CopyDir(*codeDir, tmpDir)

		*codeDir = tmpDir
		fmt.Println(*codeDir)
	}

	var ignoreDirs []string
	if *ignoreDir != "" {
		ignoreDirs = strings.Split(*ignoreDir, ",")
	}

	for _, dir := range ignoreDirs {
		if !strings.HasPrefix(dir, *codeDir) {
			// panic(fmt.Sprintf("需要排除的目录%s,必须是代码目录%s的子目录", dir, *codeDir))
			// os.Exit(1)
			print(dir)
		}
	}

	config := codeanalysis.Config{
		CodeDir:    *codeDir,
		GopathDir:  gopathDir,
		VendorDir:  "", // path.Join(*codeDir, "vendor"),
		IgnoreDirs: ignoreDirs,
	}

	result := codeanalysis.AnalysisCode(config)

	result.OutputToFile(*outputFile)

	if tmpDir != "" {
		os.RemoveAll(tmpDir)
	}

}
