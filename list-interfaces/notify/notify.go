package notify

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

func Notify(ctx context.Context, fpath string, fn func()) {
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	//添加要监控的对象，文件或文件夹
	err = watch.Add(fpath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Start Watch: ", fpath)
	//我们另启一个goroutine来处理监控对象的事件
	//go func() {
	for {
		select {
		case ev := <-watch.Events:
			{
				//判断事件发生的类型，如下5种
				// Create 创建
				// Write 写入
				// Remove 删除
				// Rename 重命名
				// Chmod 修改权限
				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("重命名文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("修改权限 : ", ev.Name)
				}
			}
		case err := <-watch.Errors:
			{
				log.Println("error : ", err)
				return
			}
		case <-ctx.Done():
			log.Println("Stop Watch: ", fpath)
			return
		}
	}
	//}()

	//循环
	//select {}
}
