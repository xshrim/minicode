package exception

import "fmt"

//  实现 java中的try...catch...

// 调用recover函数将会捕获到当前的panic（如果有的话），被捕获到的panic就不会向上传递了

// 不过要注意的是，recover之后，逻辑并不会恢复到panic那个点去，函数还是会在defer之后返回。

func Try(tryBlock func(), catchBlock func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err :", err)
			catchBlock(err)
		}
	}()
	tryBlock()
}

/*
func Excep() {
	Try(func() {
		panic("foo")
	}, func(e interface{}) {
		fmt.Println("doing something.")
		print(e)
	})
}*/
