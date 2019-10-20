package main

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"./analysis"
	"./notify"
	"./xlog"
)

func main() {
	xlog.Level = xlog.DEBUG

	imode := flag.String("mode", "single", "代码分析模式(single|daemon)")
	icodeDir := flag.String("codedir", ".", "要扫描的代码目录")
	ioutputFile := flag.String("output", "", "解析结果保存到该文件中")
	iignoreDir := flag.String("ignore", "", "无需扫描和解析的目录(代码目录相对路径, 逗号分隔)")

	gopathDir := os.Getenv("GOPATH")

	if gopathDir == "" {
		xlog.Fatalln("GOPATH目录不存在")
	}

	if len(os.Args) == 1 {
		xlog.Fatalln("使用例子\n" +
			os.Args[0] + " --codedir /appdev/gopath/src/github.com/contiv/netplugin --output /tmp/result")
	}

	flag.Parse()

	mode := strings.TrimSpace(*imode)
	codeDir := strings.TrimSpace(*icodeDir)
	output := strings.TrimSpace(*ioutputFile)
	ignoreDir := strings.TrimSpace(*iignoreDir)

	if codeDir == "" {
		xlog.Fatalln("代码目录不能为空")
	}

	codeDir, err := filepath.Abs(codeDir)
	if err != nil {
		xlog.Fatalln("代码目录解析出错: ", err.Error())
	}

	if output != "" {
		output, err = filepath.Abs(output)
		if err != nil {
			xlog.Fatalln("输出文件解析出错: ", err.Error())
		}
	}

	originDir := codeDir

	ctx, cancel := context.WithCancel(context.Background()) // 监控和中断上下文

	// 中断协程在程序中断时通过sigfn回调函数发出上下文(ctx)的cancel消息
	// 目录监控和监控同步协程均传入相同的上下文(ctx), 中断协程发出一次cancel消息, 两个协程均可接收到
	// 目录监控分析函数
	aysfn := func() {
		// TODO 代码分析
		// TODO 分析结果入库(内存库)
		// TODO API查询代码分析结果
		// time.Sleep(5 * time.Second)
	}

	// 中断操作回调函数
	sigfn := func() {
		cancel()
	}

	// 代码目录不在gopath下, 代码分析程序将无法正常分析, 需将代码复制到gopath下的临时目录下并进行两个目录的实时同步
	/*
		if !strings.HasPrefix(codeDir, gopathDir) {
			dir, fpath := filepath.Split(codeDir)

			tmpDir := path.Join(gopathDir, "/src/tmp-"+utils.GetRandomString(6))

			if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
				os.RemoveAll(tmpDir)
			}

			os.MkdirAll(tmpDir, os.ModePerm)

			utils.CopyDir(codeDir, tmpDir)

			codeDir = path.Join(tmpDir, fpath)

			originDir := path.Join(dir, path.Base(codeDir))

			// defer os.RemoveAll(tmpDir)

			sigfn = func() {
				cancel()             // 向上下文(ctx)发出cancel消息, 关闭其他协程
				os.RemoveAll(tmpDir) // 删除临时目录
			}

			if mode != "single" {
				// 目录监控同步函数
				sycfn := func() {
					// 目录同步
					utils.SyncDir(originDir, codeDir)
				}

				go notify.Notify(ctx, originDir, sycfn) // 开启目录监控(同步)
			}
		}
	*/

	var ignoreDirs []string
	if ignoreDir != "" {
		ignoreDirs = strings.Split(ignoreDir, ",")
	}

	for idx, sdir := range ignoreDirs {
		ignoreDirs[idx] = path.Join(codeDir, sdir)
	}

	config := analysis.Config{
		CodeDir:    codeDir,
		OriginDir:  originDir,
		GopathDir:  gopathDir,
		VendorDir:  path.Join(codeDir, "vendor"),
		IgnoreDirs: ignoreDirs,
	}

	aysfn = func() {
		result := analysis.AnalysisCode(config)

		xlog.Infoln("================================================")
		result.Output(output)
	}

	sig := []os.Signal{os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2}
	go notify.Signal(sig, sigfn) // 开启中断监测

	if mode != "single" {
		go notify.Notify(ctx, codeDir, aysfn) // 开启目录监控(分析)

		xlog.Infoln("后台进程正在进行实时代码分析")
		select {}
	} else {
		aysfn()
		sigfn()
		// fmt.Println("Finish")
	}

	xlog.Errorln("abc", "ok", "yes")
	xlog.Errorf("%s->%s", "name", "age")
	xlog.Debugln("debug")

	// fmt.Println(utils.ColorString(4, "color", "good", "(/home/xshrim)", "[10,2]"))
}

// go run main.go --codedir /home/xshrim/code/go/m --ignore vendor
// go run main.go --codedir ./testdata
