package main

import (
	"label.admission/xlog"
)

func test() {
	xlog.Warn("hhhhaddf")
}

func main() {
	// fpath := "." + string(os.PathSeparator) + "aa.log"
	// fmt.Println(fpath)
	// f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(f)
	// 	fmt.Println(os.Stat(fpath))
	// }

	// log.Println(f)
	// fname := os.Args[0]
	// var by []byte
	// for i := len(fname) - 1; i >= 0; i-- {
	// 	if fname[i] != os.PathSeparator {
	// 		by = append([]byte{fname[i]}, by...)
	// 	} else {
	// 		break
	// 	}
	// }

	// fmt.Println(fmt.Sprintf("%s", by))

	xlog.SetLevel(xlog.TRACE)
	xlog.SetFlag(xlog.Ldate | xlog.Ltime | xlog.Lstack | xlog.Lshortfile | xlog.Lfullcolor)
	//xlog.SetSaver(xlog.NewLogSaverWithRotation("./", 1024*1024, 1))

	xlog.Log("error", "abc%s", "AMD")
	test()

	for i := 0; i < 10; i++ {
		xlog.Debug("This is log %d", i)
	}
	xlog.Flush()

	//time.Sleep(time.Second * 1)
}
