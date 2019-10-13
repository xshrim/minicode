使用ast库分析go代码，获取代码中有哪些interface，有哪些struct实现了该interface。

# 使用例子
```
list-interfaces --codedir /appdev/gopath/src/github.com/contiv/netplugin \ 
--gopath /appdev/gopath \
--outputfile  /tmp/result

参数说明
--codedir 要分析的代码目录
--gopath GOPATH环境变量目录
--outputfile 分析结果保存到该文件
```

# 输出样例
```
interface item 在文件/appdev/gopath/src/github.com/contiv/netplugin/vendor/google.golang.org/grpc/transport/transport.go中
有2个struct实现了接口
struct windowUpdate 在文件/appdev/gopath/src/github.com/contiv/netplugin/vendor/google.golang.org/grpc/transport/control.go中
struct settings 在文件/appdev/gopath/src/github.com/contiv/netplugin/vendor/google.golang.org/grpc/transport/control.go中
```