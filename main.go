package main

import (
	"dmLaravel/tool"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Watch struct {
	watcher *fsnotify.Watcher
}

//监控目录
func (w *Watch) watchDir(dir string, to string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watcher.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控 : ", path)
		}
		return nil
	})

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				fmt.Println("***********event")
				if !ok {
					logs.Error(event, "错误--event---》")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					pids, _ := process.Pids()
					for _, pida := range pids {
						//logs.Debug(pida, "调试pid")
						pn, _ := process.NewProcess(pida)
						pName, _ := pn.Name()

						// 过滤进程名为main.exe的进程信息
						if 0 == strings.Compare(pName, "php") {
							logs.Debug(pName, "调试kill")
							kill := "kill -9 " + strconv.Itoa(int(pn.Pid))
							res, _ := tool.ExecCommand(kill)
							logs.Debug(res, "调试")
						}
					}
					go func() {
						//启动
						res, _ := tool.ExecCommand(to)
						logs.Error(string(res))
					}()

				}
			case err := <-w.watcher.Errors:
				log.Println("error:", err)
			case <-time.After(2 * time.Second):
				continue
			}

		}
	}()
}

func main() {
	var dir string
	var cm string
	// &user 就是接收命令行中输入 -u 后面的参数值，其他同理
	flag.StringVar(&cm, "c", "", "命令")
	flag.StringVar(&dir, "p", "", "监听文件夹")
	// 解析命令行参数写入注册的flag里
	flag.Parse()

	logs.Debug(dir, "调试dir")
	logs.Debug(cm, "调试命令")
	go func() {
		res, _ := tool.ExecCommand(cm)
		logs.Error(string(res))
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	w := Watch{
		watcher: watcher,
	}

	done2 := make(chan bool)

	w.watchDir(dir, cm)
	if err != nil {
		log.Fatal(err)
	}
	<-done2

}
