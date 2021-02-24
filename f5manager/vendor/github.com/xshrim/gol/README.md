# gol日志库说明

## 主要特性

gol是基于Golang标准库log源码实现的增强日志库, 在支持log库的所有功能, 保持简洁轻量的同时, 尽可能多的提供常用的日志特性, gol的新特性主要包括:

- 日志高亮(两种高亮方式)
- 日志调用定位(包名.函数名+文件名:行号)
- 可选输出模式(支持format与json模式输出)
- 链式调用(日志参数可通过链式调用方式配置)
- 多日志级别(ERROR, WARN, NOTIC, INFO, DEBUG, TRACE, FATAL, PANIC)
- 日志钩子(日志拦截与自定义处理)
- 日志上下文(便捷遥测支持)
- 标准http库日志集成(为标准库http服务提供日志输出)
- 日志热加载(动态更新日志级别与显示模式)
- 日志异步持久化(单文件+多文件轮换)
- 多维可定制输出(不同模式多路输出)

## 使用方法

### 基本用法

最简单的使用方式是使用全局默认日志实例:

```go
package main
import "github.com/xshrim/gol"
func main() {
  gol.Info("Hello Golang")
}
```

![1](./images/1.png)

全局默认日志实例的各项参数支持自定义, 参数设置支持链式调用:

```go
package main
import "github.com/xshrim/gol"
func main() {
  gol.Level(gol.TRACE)
  gol.Flag(gol.Ldate | gol.Ltime | gol.Lfile | gol.Lcolor)
  // 以上两行等价于 gol.Level(gol.TRACE).Flag(gol.Ldate | gol.Ltime | gol.Lfile | gol.Lcolor)
  gol.Debug("Hello Golang")
}
```

![2](./images/2.png)

gol支持的日志Flag有:

- Ldate: 打印日期字段
- Ltime: 打印时间字段
- Lmsec: 打印毫秒时间
- Lstack: 打印调用者字段
- Lnolvl: 不打印级别字段
- Lfile: 打印文件和行号字段
- Llfile: 打印文件完整路径
- Ljson: 以json格式打印日志
- Lcolor: 日志级别字段高亮
- Lfcolor: 日志所有字段高亮
- Lutc: 日期和时间使用UTC时间
- Ldefault: 默认打印设置(日期和时间)

此外用户也可以自行创建日志实例, 或基于该实例作进一步封装, 自定义日志实例与全局默认日志实例使用方法相同:

```go
package main
import "os"
import "github.com/xshrim/gol"
func main() {
  logger := gol.New(os.Stdout, "[GO]", gol.INFO, gol.Ldate | gol.Ltime | gol.Lfcolor)
  logger.Error("Hello Golang")
}
```

![3](./images/3.png)

调用gol打印日志的函数分别为`Error`, `Warn`, `Info`, `Debug`, `Trace`, `Fatal`和`Panic`. 所有函数均提供如`Errorf`, `Infof`等格式化输出函数, 日志打印函数均支持不同类型的可变参数, 如:

```go
package main
import "github.com/xshrim/gol"
func main() {
  gol.Level(gol.INFO).Flag(gol.Ldate | gol.Ltime | gol.Lfile | gol.Lstack | gol.Lfcolor)
  gol.Prtln("Pritnln").Warnf("Warn: %s", "invalid format")
  gol.Panic("Error: division by ", 0, gol.Err("error message"))
}
```

gol的日志打印相关函数均是换行打印, 如需不换行打印, 可使用`Log`和`Logf`函数, 用法为: `gol.Log("debug", "print %s without newline, "text")`.

此外, gol还内置了`Prt`, `Prtf`, `Prtln`, `Sprt`, `Sprtf`, `Err`, `Errf`这些常用函数, 分别用于不换行打印, 换行打印, 返回字符串和返回错误等. 类似fmt库中的`Print`, `Errorf`. `Prt`和`Prtln`支持自动识别输出流和格式化输出.

![7](./images/7.png)

### 特性用法

#### 日志高亮

gol支持**无高亮**, **仅Level高亮**和**全高亮**的方式显示日志, 两种高亮方式的启用方式为在日志实例的`flag`字段加入`Lcolor`或`Lfcolor`项, 如: `gol.Flag(gol.Ldate | gol.Ltime | gol.Lfcolor)`, 效果分别如下:

![4](./images/4.png)

#### 调用定位

gol在log库支持显示调用文件及所在行的基础上, 还支持显示调用函数及其所在package, 此特性需要在日志实例的`flag`字段加入`Lstack`项, 如: `gol.Flag(gol.Ldate | gol.Ltime | gol.Lstack | gol.Lfile | gol.Lfcolor)`, 效果如下:

![5](./images/5.png)

#### 输出模式

gol日志可以以普通格式化方式输出, 也可以输出为json数据, 方便在开发或排障过程中直接使用. 默认进行普通格式化输出, 通过在日志实例的`flag`字段加入`Ljson`项将以json格式输出. 以json格式输出时, `Lcolor`和`Lfcolor`配置将失效.

gol支持通过`DateKey`, `StackKey`, `LevelKey`, `CtxKey`, `MsgKey`等函数自定义json格式输出时的各个日志字段的键值.

![6](./images/6.png)

#### 链式调用

对gol日志实例和上下文的配置支持链式调用, 简洁直观.

```go
package main
import "os"
import "github.com/xshrim/gol"
func main() {
  gol.Flag(gol.Lcolor, gol.APPEND).Debug("Debug message").Errorf("%s message", "Error")
  gol.Level(gol.INFO).Out(os.Stdout).Flag(gol.Ldate | gol.Ltime | gol.Lfile | gol.Lstack | gol.Ljson).Warn("Warn: %s", "invalid format")
  gol.New(os.Stdout, "[GO]", gol.TRACE, gol.LstdFlags).Flag(gol.Lstack, gol.APPEND).Saver(gol.NewLogSaverWithRotation("./", 1024*1024, 3)).HotReload().Debug("Debug message").Flush()
  gol.With(nil).Str("foo", "bar").Ints("data", []int{3, 5}).Info("Info Message")
}
```

![8](./images/8.png)

#### 日志级别

gol支持**OFF**, **PANIC**, **FATAL**, **ERROR**, **WARN**, **NOTIC**, **INFO**, **DEBUG**, **TRACE**, **ALL**九种日志级别, 分别对应数字0, 1, 2, 3, 4, 5, 6, 7, 8, 9. 优先级由高到低, 通过`Level`函数设置可显示级别后, 低于该级别的日志将不打印, `OFF`级别将关闭所有日志打印. 如:

```go
package main
import "github.com/xshrim/gol"
func main() {
  gol.Level(gol.INFO).Flag(gol.Ldate | gol.Ltime | gol.Lfile | gol.Lstack | gol.Lcolor)
  gol.Warn("Debug Message")   // 打印
  gol.Trace("Trace Message")  // 不打印
  gol.Fatal("Fatal Message")  // 打印
}
```

![9](./images/9.png)

#### 日志钩子

gol允许在日志实例上定义Hook函数, Hook函数会对日志数据进行拦截, 可以读取或修改日志内容, Hook函数处理完成后, 可以选择是否放行, 放行后将继续执行后续流程(日志输出), 如不放行, 日志将不会输出.

Hook函数形式为: `func(int, *[]byte) bool`, 示例如下:

```go
package main
import "github.com/xshrim/gol"
func modifyHook(lv int, buf *[]byte) bool {
  *buf = []byte("HaHa\n")  // 修改日志内容
  return true  // 放行
}
func main() {
  gol.Level(gol.INFO).Flag(gol.Ldate | gol.Ltime | gol.Lfile | gol.Lstack | gol.Ljson).Hook(modifyHook)
  gol.Warn("Debug Message")   // 将打印HaHa
}
```

#### 日志上下文

日志和链路追踪是可观察性遥测技术的重要组成部分, 为了便于在日志数据和调用链路之间相互追溯, 通常在一个或者关联的处理逻辑中加入相同的追踪ID(traceid)等信息, 该逻辑中所有的输出日志都附带这些信息. 这样一组相关日志携带关联数据的情况就是日志上下文.

gol提供了简单易用的上下文功能支持. 通过`With`或者`NewContext`函数即可创建一个非线程安全的上下文, 通过`WithSafe`或者`NewSafeContext`函数即可创建一个线程安全的上下文, 上下文的所有日志都将在保持原日志实例配置的前提下携带上下文特有信息, 该日志实例的其他日志不受影响. 上下文同样支持链式调用. 通过`Field`函数可以设置或删除上下文字段内容.

例如对于一个Web API请求, 请求头或请求体中可能会携带诸如请求ID或追踪ID以及用户ID等信息, 一个请求可能会产生多条日志, 这些日志就可以通过上下文关联起来.

```go
package main
import (
  "fmt"
  "github.com/xshrim/gol"
  "net/http"
)
func hello(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "ParseForm() err: %v", err)
    return
  }
  traceid := r.FormValue("traceid")
  userid := r.FormValue("userid")
  username := r.FormValue("username")
  ctx := gol.Level(gol.DEBUG).With(nil).Field(map[string]interface{}{"traceid": traceid, "userid": userid})
  ctx.Infof("%s is requested by user %s", "/", username)
  ctx.Debugf("User %s requests %s with method %s", username, "/", r.Method)
  // 等同于:
  // gol.Level(gol.DEBUG).With(map[string]interface{}{"traceid": traceid, "userid": userid}).Infof("%s is requested by user %s", "/", username).Debugf("User %s requests %s with method %s", username, "/", r.Method)
  fmt.Fprintf(w, "Welcome %s!\n", username)
}

func main() {
  http.HandleFunc("/", hello)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    gol.Fatal(err)
  }
}
```

在终端输入以下命令模拟请求:

```bash
curl -XPOST -d 'traceid=12345&userid=001&username=admin' http://127.0.0.1:8080
```

设置上下文字段有以下几种方式:

```go
// tk.M与map[string]interface{}等价
gol.With(tk.M{"url": "www.demo.com", "duration": time.Second}) // 直接作为With函数参数
gol.With(nil).Field(map[string]interface{}{"url": "www.demo.com", "duration": time.Second}) // 通过Field函数设置
gol.With(nil).Str("url", "www.demo.com").Dur("duration", time.Secomd) // 链式字段设置
```

几种方式可结合使用, 上下文字段会不断追加, 上下文内容会以json格式输出.

![10](./images/10.png)

#### 标准http库日志集成

gol为标准http库提供日志集成, 通过`HttpHandler`和`HttpHandlerFunc`函数对用户的后台请求响应逻辑(`Handler`或`HandlerFunc`)进行无侵入式封装, 从而自动为基于http标准库的Web服务提供日志输出. gol将依据日志实例的`Ljson`和`Lcolor`标记自动选择http日志的输出格式(json, 高亮, 无高亮). 向`HttpHandler`和`HttpHandlerFunc`函数追加不定长参数还可以输出更多http请求头部信息(参数必须是http header中的字段大写).

```go
package main
import (
  "fmt"
  "net/http"
  "github.com/xshrim/gol"
)
func Index(w http.ResponseWriter, r *http.Request) {
  if r.RequestURI != "/" {
    w.WriteHeader(404)
    return
  }
  fmt.Fprintf(w, "Hello Golang")
}
func main() {
  gol.Flag(gol.Lcolor, gol.APPEND)
  http.HandleFunc("/", gol.HttpHandlerFunc(Index, "Content-Type", "User-Agent"))
  http.ListenAndServe(":8000", nil)
}
```

![11](./images/11.png)

同样的处理逻辑, 我们也可以将gol作为其他web服务框架的日志中间件. 以**gin** WEB框架为例:

```go
package main
import (
  "time"
  "github.com/gin-gonic/gin"
  "github.com/xshrim/gol"
)
func logger() gin.HandlerFunc {
  return func(c *gin.Context) {
    startTime := time.Now()
    c.Next()
    gol.Infof("| %3d | %10v |",
      c.Writer.Status(),
      time.Now().Sub(startTime),
    )
  }
}
func main() {
  r := gin.New()
  r.Use(gin.Recovery())
  r.Use(logger())
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  _ = r.Run()
}
```

#### 日志热加载

gol允许在不中断应用的情况下动态更新日志级别和输出样式, 实现运行时输出日志的详细度和可视友好度调整, 为应用排障提供便利. 此外此特性还可以实现日志输出的动态关闭与开启. 通过`HotReload`函数开启热加载功能后, 在操作系统终端执行`echo <msg> > /tmp/.gol`即可动态调整日志级别和输出样式, 可以在`HotReload`函数参数中自定义需要监听的热加载文件的路径.

`<msg>`有四类值, 其作用分别如下:

1. error, info, debug等字符串: 表示限制输出的日志级别, 也可以使用1, 2, 3, 4等对应的数字
2. color, fcolor, nocolor字符串: 表示三种日志高亮的输出模式
3. format, json字符串: 表示以格式化或者json格式的方式输出日志
4. stackfile, nostackfile字符串: 表示是否显示调用定位

四类值可自由组合, 示例如下:

```go
package main
import "time"
import "github.com/xshrim/gol"
func main() {
  gol.Level(gol.INFO).Flag(gol.Ldate | gol.Ltime).HotReload()
  for {
    gol.Error("Error Message")
    gol.Warn("Debug Message")
    gol.Info("Info Message")
    gol.Debug("Debug Message")
    gol.Trace("Trace Message")
    time.Sleep(time.Millisecond * 500)
  }
}
```

运行后在终端执行:

```bash
echo debug > /tmp/.gol   # windows下对应文件为 C:\.gol
echo color > /tmp/.gol
echo "warn json" > /tmp/.gol
echo "format fcolor" > /tmp/.gol
echo "format stackfile" > /tmp/.gol
```

程序将按照指定的配置输出日志, 执行`echo off > /tmp/.gol`将关闭日志输出.

**[注意]:**

1. 删除.gol文件或执行`echo 0 > /tmp/.gol`将还原日志原配置.
2. 基于文件的热加载功能的主要目的是临时性调试, 请不要用作gol配置文件, 临时动态调整后建议删除该文件, 以免日志输出与程序内设置不符

![12](./images/12.png)

#### 日志持久化

gol支持将日志持久化输出到日志文件中, 既可以指定单个日志文件, 也可以指定日志目录和相关轮换参数进行日志文件的**自动轮换**. 而且允许将日志同时输出到标准输出设备和日志文件中. gol的日志持久化操作是异步完成的. 使用方法为通过`NewLogSaverWithLogFile`或`NewLogSaverWithRotation`生成LogSaver并通过`Saver`绑定到日志实例.

```go
package main
import "github.com/xshrim/gol"
func main() {
  gol.Level(gol.TRACE).Flag(gol.Ldate | gol.Ltime | gol.Lstack | gol.Lfile | gol.Lfcolor).Saver(gol.NewLogSaverWithRotation("./", 1024*1024, 3)) // 参数依次为: 轮换日志目录, 日志文件大小上限, 轮换日志数量
  gol.Info("日志持久化")
  defer gol.Flush() // 保证日志全部落盘
}
```

**[注意]:**

1. 进行日志轮换时, 日志文件格式为`<程序名>.<序号>.log`
2. 由于日志持久化是异步完成的, 因此为了保证所有日志都写入文件中, 建议在程序退出前调用`Flush`函数
3. gol默认将日志输出到**Stderr**(可自定义), 如果通过`Saver`函数指定了日志持久化, 则日志默认会同时输出到默认输出和日志文件中. 可以通过`Writer()`或`UnWriter()`关闭默认输出

#### 多维可定制输出

gol支持三个维度的输出流以满足不同的场景需求.

- MultiWriter: 一个Logger实例下利用io.MultiWriter实现多路输出, 所有输出的格式都是相同的
- MultiLogger: 利用日志上下文实现一个Context实例下多个Logger, 不同Logger的输出格式可以不同
- LogSaver: 一个Logger实例下除多路输出流外, 还可额外定义一个日志持久化输出, 二者互不影响

LogSaver是异步方式输出日志, 另外两种是同步输出. 三者可同时使用.

MultiWriter示例:

```go
package main
import "os"
import "github.com/xshrim/gol"
func main() {
  logger := gol.New(nil, "", gol.INFO, gol.Ldate | gol.Ltime | gol.Lfile)
  file, _ := os.OpenFile("demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  logger.Writer(os.Stdout, file)  // 定义多路输出
  logger.Info("Info Message") // 将同时在终端和文件中输出
}
```

MultiLogger示例:

```go
package main
import "os"
import "github.com/xshrim/gol"
func main() {
  stdout := gol.New(os.Stdout, "", gol.INFO, gol.Ldate|gol.Ltime|gol.Lfile|gol.Lfcolor)
  file, _ := os.OpenFile("demo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  fileout := gol.New(file, "", gol.DEBUG, gol.Ldate|gol.Ltime|gol.Lstack)
  ctx := gol.NewContext(nil, stdout, fileout)
  ctx.Info("Info Message") // 将在终端和文件中输出不同内容
  ctx.Debug("Debug Message") // 将只在文件中输出
}
```

LogSaver已在日志持久化部分介绍.

### 工具集

gol日志库内置了**Prt**, **Prtf**, **Prtln**, **Sprtf**, **Errf**等常用输出函数(包括提供彩色输出的**Cprt**, **Cprtf**, **Cprtln**等函数). 此外其中的**tk**库还提供了大量的常用工具函数:

- **Jsonify**: 将任意数据结构转换为json字符串(内置的**M**类型为map[string]interface{}别名, 同样支持转换)
- **Jsquery**: 根据字符串路径提取json字符串中的特定key的value, 无需构建结构体, 字符串路径支持高级格式, 如"user*#0[0].[1-3]"表示json字符串中第一个带user前缀的key的value的第一个元素的第二和第三个子元素.
- **Imapify**: 将json字符串转换为map[string]interface{}(内置的**M**类型为map[string]interface{}别名)
- **Imapset**: 修改map[string]interface{}指定key的value, 字符串路径支持下标, 如user[0]表示修改key为user的第一个子元素.
- **etc...**: 其他诸如Iter, Remove, QuickSort, Uniq, Exec, HttpGet等常用工具函数.

使用举例:

```go
package main
import "github.com/xshrim/gol"
import "github.com/xshrim/gol/tk"
type Person struct {
  Name string
  Age  int
}
func main() {
  jsdata := `{
    "name": "tom",
    "age": 25,
    "hobbies": ["football", "movie", "read", "music"],
    "location": [{
      "prov.city": "foo",
      "street": "bar"
    }]
  }`
  gol.Prtln(tk.Jsquery(jsdata, "loc*.[0].prov\\.city"))
  gol.Prtln(tk.Jsquery(jsdata, "hobbies[odd].[0]"))
  user := tk.M{
    "name":     "tom",
    "age":      25,
    "children": []Person{Person{"Lucy", 5}, Person{"Lily", 4}},
    "hobbies":    []string{"football", "movie", "read", "music"},
    "location": []tk.M{tk.M{
      "city":   "foo",
      "street": "bar",
    }},
  }
  gol.Prtln(tk.Jsonify(user))
  gol.Prtln(tk.Jsonify([]map[string]interface{}{map[string]interface{}{"a": 1, "b": "demo"}, map[string]interface{}{"c": 2, "d": []string{"foo", "bar"}}}))
  tk.Imapset(user, "hobbies[0]", []string{"game"})
  gol.AddFlag(gol.Ljson).With(user).Info(user)
}
```

## 性能测试

### 测试说明

本次选取**标准库log**, **logrus**, **zap**, **zapsugar**和**zerolog**这五种常用的日志库作为对照, 对以下几种使用场景进行性能测试:

- Normal:  默认日志配置下输出单个字符串
- Format: 默认日志配置下输出格式化字符串
- DiscardWriter: 默认日志配置下丢弃输出
- WithoutFlags: 最简日志配置下仅输出单个字符串
- WithDebugLevel: 默认日志配置下输出debug日志
- WithFields: 默认日志配置下携带3个字段的上下文输出单个字符串
- WithFields:默认日志配置下携带3个字段的上下文输出格式化字符串

> 已调整各日志库的默认日志输出仅包括**日期**, **时间**, **日志级别**和**日志字符串**. 最简日志输出仅包括**日志字符串**.

测试机配置和磁盘空间所限, 每项测试仅持续5s. 测试结果仅供参考. 测试命令如下:

```bash
go test -bench=. -benchtime=10s -timeout 10m -benchmem -run=none
```

### 测试结果

#### 每次操作执行时间(ns/op)


|          | Normal | Format | DiscardWriter | WithoutFlags | WithDebugLevel | WithFields | WithFieldsFormat |
| -------- | ------ | ------ | ------------- | ------------ | -------------- | ---------- | ---------------- |
| log      | 2688   | 2895   | 741           | 2321         | -              | -          | -                |
| logrus   | 7254   | 8849   | 4786          | 4429         | 8453           | 12751      | 13890            |
| gol      | 2750   | 2926   | 768           | 2450         | 2764           | 4585       | 4657             |
| zap      | 3903   | -      | 978           | 3019         | 3991           | 5602       | -                |
| zapsugar | 4489   | 4627   | 1283          | 2895         | 3823           | 10791      | 11214            |
| zerolog  | 3203   | 3751   | 999           | 2164         | 3094           | 3579       | 4194             |

#### 每次操作分配内存(B/op)


|          | Normal | Format | DiscardWriter | WithoutFlags | WithDebugLevel | WithFields | WithFieldsFormat |
| -------- | ------ | ------ | ------------- | ------------ | -------------- | ---------- | ---------------- |
| log      | 16     | 36     | 16            | 16           | -              | -          | -                |
| logrus   | 405    | 505    | 405           | 277          | 469            | 955        | 1048             |
| gol      | 16     | 36     | 16            | 16           | 16             | 132        | 149              |
| zap      | 0      | -      | 0             | 0            | 0              | 192        | -                |
| zapsugar | 16     | 33     | 16            | 0            | 0              | 336        | 360              |
| zerolog  | 0      | 35     | 0             | 0            | 0              | 0          | 34               |

#### 每次操作内存分配次数(allocs/op)


|          | Normal | Format | DiscardWriter | WithoutFlags | WithDebugLevel | WithFields | WithFieldsFormat |
| -------- | ------ | ------ | ------------- | ------------ | -------------- | ---------- | ---------------- |
| log      | 1      | 2      | 1             | 1            | -              | -          | -                |
| logrus   | 13     | 18     | 13            | 8            | 16             | 19         | 24               |
| gol      | 1      | 2      | 1             | 1            | 1              | 5          | 6                |
| zap      | 0      | -      | 0             | 0            | 0              | 1          | -                |
| zapsugar | 1      | 2      | 1             | 0            | 0              | 4          | 5                |
| zerolog  | 0      | 2      | 0             | 0            | 0              | 0          | 2                |
