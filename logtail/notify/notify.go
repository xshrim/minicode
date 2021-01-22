package notify

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hpcloud/tail"
	"github.com/hpcloud/tail/util"
)

func args2config() (tail.Config, int64) {
	config := tail.Config{Follow: true}
	n := int64(0)
	maxlinesize := int(0)
	logger := false
	flag.Int64Var(&n, "n", 0, "tail from the last Nth location")
	flag.IntVar(&maxlinesize, "max", 0, "max line size")
	flag.BoolVar(&config.Follow, "f", false, "wait for additional data to be appended to the file")
	flag.BoolVar(&config.ReOpen, "F", false, "follow, and track file rename/rotation")
	flag.BoolVar(&config.Poll, "p", false, "use polling, instead of inotify")
	flag.BoolVar(&logger, "l", false, "enable logger")
	flag.StringVar(&config.Target, "t", "/var/log/container", "watch target")
	flag.Parse()
	if config.ReOpen {
		config.Follow = true
	}
	if !logger {
		config.Logger = tail.DiscardingLogger
	}
	config.MaxLineSize = maxlinesize
	return config, n
}

type Notifier struct {
	watch    *fsnotify.Watcher
	tailList map[string]*tail.Tail
	config   tail.Config
	wg       sync.WaitGroup
}

func NewNotifier() *Notifier {
	config, n := args2config()
	if flag.NFlag() < 1 {
		fmt.Println("need one or more files as arguments")
		os.Exit(1)
	}

	if n != 0 {
		config.Location = &tail.SeekInfo{-n, os.SEEK_END}
	}

	w := new(Notifier)
	w.watch, _ = fsnotify.NewWatcher()
	w.tailList = make(map[string]*tail.Tail)
	w.config = config
	return w
}

//监控目录
func (n *Notifier) WatchDir() {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(n.config.Target, func(path string, info os.FileInfo, err error) error {
		//判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = n.watch.Add(path)
			if err != nil {
				return err
			}
			//fmt.Println("监控 : ", path)
		} else {
			if util.IsFitFile(path) {
				if _, ok := n.tailList[path]; !ok {
					t, err := tail.TailFile(path, n.config)
					if err != nil {
						// fmt.Println("Tail error: ", err)
					} else {
						n.tailList[path] = t
						n.wg.Add(1)
						go t.DoTail(&n.wg)
					}
				}
			}
		}
		return nil
	})

	go n.WatchEvent() //协程
}

func (n *Notifier) WatchEvent() {
	for {
		select {
		case ev := <-n.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					//	fmt.Println("创建文件 : ", ev.Name)
					//获取新创建文件的信息，如果是目录，则加入监控中
					file, err := os.Stat(ev.Name)
					if err == nil {
						if file.IsDir() {
							n.watch.Add(ev.Name)
							//	fmt.Println("添加监控 : ", ev.Name)
						} else {
							if fpath, err := filepath.Abs(ev.Name); err == nil {
								if util.IsFitFile(fpath) {
									if _, ok := n.tailList[fpath]; !ok {
										t, err := tail.TailFile(fpath, n.config)
										if err != nil {
											// fmt.Println("Tail error: ", err)
										} else {
											n.tailList[fpath] = t
											n.wg.Add(1)
											go t.DoTail(&n.wg)
										}
									}
								}
							}
						}
					}

				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					//fmt.Println("写入文件 : ", ev.Name)
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					//fmt.Println("删除文件 : ", ev.Name)
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					// if !fi.IsDir() {
					// 	print(ev.Name)
					// 	// if _, ok := flist[ev.Name]; ok {
					// 	// 	delete(flist, ev.Name)
					// 	// }
					// }
					if fpath, err := filepath.Abs(ev.Name); err == nil {
						n.tailList[fpath].Cleanup()
						delete(n.tailList, fpath)
					}

					if err == nil {
						if fi.IsDir() {
							n.watch.Remove(ev.Name)
							//fmt.Println("删除监控 : ", ev.Name)
						}
					}
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					//fmt.Println("重命名文件 : ", ev.Name)
					n.watch.Remove(ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					//	fmt.Println("修改权限 : ", ev.Name)
				}
			}
		case <-n.watch.Errors:
			{
				//fmt.Println("error : ", err)
				return
			}
		}
	}
}
