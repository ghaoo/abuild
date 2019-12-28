# golang 自动编译工具

`autobuild` 是一个 golang 的自动编译工具，解放手动编译。当文件发生变动时自动执行 `go build` 命令并编译项目，再也不用修改一行文件执行一次 `go run xx.go` 或者 `go build` 了！

## 安装
```go
go get -u github.com/ghaoo/autobuild
```

## 使用

**请确保 `github.com/ghaoo/autobuild` 已经编译安装**

命令行参数
```go
Usage:
  	-h    显示当前帮助信息
  	-f	  指定main文件
  	-o    执行编译后的可执行文件名
  	-r    是否搜索子目录，默认为true
```

请在需要执行自动编译的文件夹下使用：

```bash
// 执行 autobuild 命令
 autobuild -f xxx.go -o xxx -r=false
```
> 如果需要监听多个目录请在 `autobuild` 命令下直接填写目录地址，如监听目录下的父目录：`autobuild ../`

自动编译默认只监听以 `.go` 为扩展名的文件，如需监听其他文件需创建文件 `.extensions`，并在文件中添加需要监听的文件扩展名
如需增加以 `.conf` 为后缀的文件，在文件中填写:
```go
.go
.conf
```

当存在 `.extensions` 并且文件中存在数据时只监听 `.extensions` 配置中的扩展名


