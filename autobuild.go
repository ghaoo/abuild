package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const usage = `
 Usage:
  	-h    显示当前帮助信息
  	-f	  指定main文件
  	-o    执行编译后的可执行文件名
  	-r    是否搜索子目录，默认为true
`

var (
	eventTime    = make(map[string]int64)
	scheduleTime time.Time
)

// 监听器
type watch struct {
	path      string    // 监听目录，默认当前目录
	appName   string    // 输出的程序文件
	appCmd    *exec.Cmd // appName的命令行包装引用，方便结束其进程。
	goCmdArgs []string  // 传递给go build的参数
}

func main() {
	// 初始化flag
	var showHelp, recursive bool
	var watchPath, outputName, mainFiles string

	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")
	flag.BoolVar(&recursive, "r", true, "是否查找子目录")
	flag.StringVar(&watchPath, "p", "./", "指定监听目录")
	flag.StringVar(&outputName, "o", "", "指定输出名称")
	flag.StringVar(&mainFiles, "f", "", "指定需要编译的文件")
	flag.Usage = func() {
		fmt.Printf(usage)
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		ColorLog("[ERRO] 获取当前工作目录时，发生错误: [%s] \n", err)
		return
	}

	// 初始化goCmd的参数
	args := []string{"build", "-o", outputName}
	if len(mainFiles) > 0 {
		args = append(args, mainFiles)
	}

	w := &watch{
		appName:   getAppName(outputName, wd),
		goCmdArgs: args,
	}

	w.watcher(recursivePath(recursive, append(flag.Args(), wd)))

	go w.build()

	done := make(chan bool)
	<-done
}

// 读取.autowatcher文件，用于设置需要监听的文件后缀
func watchExtensions() []string {
	var exts []string
	f, err := os.OpenFile(".autowatcher", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		exts = append(exts, line)
	}

	return exts
}

func (w *watch) watcher(paths []string) {

	ColorLog("[TRAC] 初始化文件监视器... \n")
	//初始化监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ColorLog("[ERRO] 初始化监视器失败: [%s] \n", err)
		os.Exit(2)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				build := true
				if !checkIfWatch(event.Name) {
					fmt.Println("忽略 ", event.Name)
					continue
				}

				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					ColorLog("[SKIP] 跳过 [%s] \n", event)
					continue
				}

				mt := w.getFileModTime(event.Name)
				if t := eventTime[event.Name]; mt == t {
					ColorLog("[SKIP] 跳过 [%s] \n", event.String())
					build = false
				}

				eventTime[event.Name] = mt

				if build {
					go func() {
						scheduleTime = time.Now().Add(time.Second)
						for {
							time.Sleep(scheduleTime.Sub(time.Now()))
							if time.Now().After(scheduleTime) {
								break
							}
							return
						}
						ColorLog("[INFO] 触发编译事件: <%s> \n", event)
						w.build()
					}()
				}

			case err := <-watcher.Errors:
				ColorLog("[ERRO] 监控失败 [%s] \n", err)
			}
		}
	}()

	for _, path := range paths {
		ColorLog("[TRAC] 监视文件夹: (%s) \n", path)
		err = watcher.Add(path)
		if err != nil {
			ColorLog("[ERRO] 监视文件夹失败: [%s] \n", err)
			os.Exit(2)
		}
	}
}

// 开始编译代码
func (w *watch) build() {
	ColorLog("[TRAC] 编译代码... \n")

	goCmd := exec.Command("go", w.goCmdArgs...)
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout

	if err := goCmd.Run(); err != nil {
		ColorLog("[ERRO] 编译失败: [%s] \n", err)
		return
	}

	ColorLog("[SUCC] 编译成功... \n")

	w.restart()
}

func (w *watch) restart() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Kill.recover -> ", err)
		}
	}()

	if w.appCmd != nil && w.appCmd.Process != nil {
		ColorLog("[TRAC] 终止旧进程... \n")
		if err := w.appCmd.Process.Kill(); err != nil {
			ColorLog("[ERRO] 终止进程失败 [%s] ...\n", err)
		}
		ColorLog("[TRAC] 旧进程被终止! \n")
	}

	ColorLog("[TRAC] 重新启动 [%s] \n", w.appName)
	if strings.Index(w.appName, "./") == -1 {
		w.appName = "./" + w.appName
	}

	ColorLog("[SUCC] 启动新进程... \n")
	w.appCmd = exec.Command(w.appName)
	w.appCmd.Stderr = os.Stderr
	w.appCmd.Stdout = os.Stdout
	if err := w.appCmd.Start(); err != nil {
		ColorLog("[ERRO] 启动进程时出错: [%s] \n", err)
	}
}

func checkIfWatch(name string) bool {
	exts := watchExtensions()
	fmt.Println(exts)
	if len(exts) > 0 {
		for _, s := range exts {
			if strings.HasSuffix(name, s) {
				return true
			}
		}
	}
	return false
}

func (w *watch) getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {

		ColorLog("[ERRO] 文件打开失败 [%s]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		ColorLog("[ERRO] 获取不到文件信息 [%s]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func getAppName(outputName, wd string) string {
	if len(outputName) == 0 {
		outputName = filepath.Base(wd)
	}
	if runtime.GOOS == "windows" && !strings.HasSuffix(outputName, ".exe") {
		outputName += ".exe"
	}
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {

	}

	return outputName
}

// 根据recursive值确定是否递归查找paths每个目录下的子目录。
func recursivePath(recursive bool, paths []string) []string {
	if !recursive {
		return paths
	}

	var ret []string

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			ColorLog("[ERRO] 遍历监视目录错误: [%s] \n", err)
		}

		if fi.IsDir() && strings.Index(path, "/.") < 0 {
			ret = append(ret, path)
		}
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			ColorLog("[ERRO] 遍历监视目录错误: [%s] \n", err)
		}
	}

	return ret
}
