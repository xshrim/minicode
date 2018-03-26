# Go语言

## 基础

### 概览

- [Go语言](#go%E8%AF%AD%E8%A8%80)
    - [基础](#%E5%9F%BA%E7%A1%80)
        - [概览](#%E6%A6%82%E8%A7%88)
        - [编译](#%E7%BC%96%E8%AF%91)
        - [变量和数据类型](#%E5%8F%98%E9%87%8F%E5%92%8C%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B)
        - [流程控制](#%E6%B5%81%E7%A8%8B%E6%8E%A7%E5%88%B6)
        - [函数](#%E5%87%BD%E6%95%B0)
    - [面向对象](#%E9%9D%A2%E5%90%91%E5%AF%B9%E8%B1%A1)
        - [类](#%E7%B1%BB)
        - [接口](#%E6%8E%A5%E5%8F%A3)
        - [错误处理](#%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86)
    - [并发](#%E5%B9%B6%E5%8F%91)
        - [并发模型](#%E5%B9%B6%E5%8F%91%E6%A8%A1%E5%9E%8B)
        - [goroutine](#goroutine)
        - [channel](#channel)
        - [select](#select)
        - [多核并发](#%E5%A4%9A%E6%A0%B8%E5%B9%B6%E5%8F%91)
        - [锁和全局唯一性操作](#%E9%94%81%E5%92%8C%E5%85%A8%E5%B1%80%E5%94%AF%E4%B8%80%E6%80%A7%E6%93%8D%E4%BD%9C)

### 编译

- 编译生成可执行文件（不同平台生成的可执行文件不同）或者不生成可执行文件直接运行
    ```bash
    # 生成可执行文件
    go build hello.go
    # 直接运行
    go run hello.go
    ```
- Go是一个命令行工具，提供众多源代码管理命令，真正的Go编译器和连接器被隐藏在Go命令行工具之后,6g和8g分别是64位和32位编译器，6l和8l同理
    ```bash
    6g hello.go
    6l hello.6
    ./6.out
    ```
- Go 1.5版之后，编译器6g/8g命令统一合并为go tool compile, 汇编器和连接器也分别合并为go tool asm和go tool link，目标文件统一以.o后缀
    ```bash
    go tool compile hello.go
    go tool link hello.o -o hello
    ./hello
    ```
- Go编译生成的二进制程序可以直接使用gdb调试

    ```bash
    gdb hello
    ```

### 变量和数据类型

- 变量声明方式
    ```go
    var v1 int              //普通类型变量
    var s1, s2 string       //多个普通类型变量
    var a1 [3]int           //数组变量
    var l1 []int            //切片变量
    var p1 *int             //指针变量
    var c1 chan int         //channel变量
    var t1 struct {         //结构体变量
        f int
    }
    var m1 map[string]int   //map变量
    var f1 func(x int) int  //函数变量
    var (                   //变量批量声明
        a int
        b string
    )
    ```
- 变量初始化
    ```go
    //三种形式，后两种形式go编译器会自动推导变量类型
    var v1 int = 10
    var v2 = 10
    v3 := 10
    //三种形式均支持批量初始化
    var v1, v2, v3 int = 1, 2, 3
    var v1, v2, v3 = 1, 2, 3
    v1, v2, v3 := 1, 2, "hello"        //使用:=作批量初始化时要求左侧变量至少有一个是未声明过的新变量(_占位符不属于未声明的新变量)
    //数组初始化可以不指定大小
    a1 := [...]int{1, 2, 3}
    ```
- 匿名变量
    ```go
    //使用匿名变量_可以丢弃不需要的代码返回值
    var a, b = 1, 2
    _, v1 := a, b
    ```
- 变量赋值
    ```go
    //go提供多重赋值功能
    i, j = j, i
    ```
- 常量
    ```go
    //常量赋值支持表达式，但由于是编译期行为，不能将任何需要运行期才能得出结果的表达式赋给常量
    const Pi flast64 = 3.1415926
    const (
        size int64 = 1024
        eof = -1
    )
    const a, b, c = 3, 4, "hello"
    const mask = 1 << 3               //正确
    const Home = os.GetEnv("HOME")    //错误
    //go语言有三个预定义常量：true、false和iota
    //iota是一个可以被编译器修改的常量，每一个const关键字出现之前都被重置为0，在下一个const出现前，每次出现iota，其代表的数字都会自动增1
    const (
        c0 = iota             // c0 == 0
        c1 = iota             // c1 == 1
        c2 = iota             // c2 == 2
    )
    const (
        d0 = iota             // d0 == 0
        d1 = iota             // d1 == 1
        d2 = iota             // d2 == 2
    )
    ```
- 数据类型
  - 基础类型：
    - 布尔：bool
    - 整型：int8、byte、int16、int32、uint8、int、uint、uintptr等
    - 浮点：float32、float64
    - 复数：complex64、complex128
    - 字符串：string
    - 字符：rune
    - 错误：error
  - 复合类型：
    - 指针(pointer)
    - 数组(array)
    - 切片(slice)
    - 字典(map)
    - 通道(chan)
    - 结构体(struct)
    - 接口(interface)
- 类型说明
  - Go语言中的变量在声明后都有默认值，可以不必赋值直接使用，整型和浮点型默认值是0，布尔型默认值是false， 复数型默认值是0+0i,字符串默认值是""，错误、指针、切片、字典、channel、接口默认值是nil，数组和struct默认值依其内部数据类型而定
    ```go
    var (
        v1 int
        v2 bool
        v3 error
        v4 [3]int
        v5 [3]*int
    )
    fmt.Println(v1, v2, v3, v4, v5)
    //输出     0 false <nil> [0 0 0] [<nil> <nil> <nil>]
    ```
  - Go的常量变量函数名以大写字母开头的在包外可见，小写字母开头的在包外不可见
    ```go
    package ext
    var Aa int = 1
    var bb int = 2
    func Add(x, y int) int {
        return x + y
    }
    func minus(x, y int) int {
        return x - y
    }
    ```
    ```go
    package main
    import (
        "fmt"
        "ext"
    )
    func main() {
        x, y := 5, 3
        i := ext.Aa             //正确
        j := ext.bb             //错误
        r1 := ext.Add(x, y)     //正确
        r2 := ext.minus(x, y)   //错误
    }
    ```
  - 布尔类型不支持自动或强制类型转换
    ```go
    var v1 bool
    v1 = true                   //正确
    v1 = (1 == 2)               //正确
    v1 = 1                      //错误
    v1 = bool(1)                //错误
    ```
  - int、uint和uintptr是平台相关的封装类型，封装类型和int8等基础类型不能自动转换，通常使用这些封装类型即可
    ```go
    var v1 int32
    v2 := 12                    //v2被自动推导为int类型
    v1 = v2                     //错误
    v1 = int32(v2)              //正确
    ```
  - float32和float64分别等价于c语言的float和double类型。自动推导的浮点数初始化赋值会被go推导为float64，浮点数的比较不建议直接使用==，而应使用math.Fdim(f1, f2) < p
    ```go
    var f1 float32
    f2 := 10.0                  //f2被自动推导为float64
    f1 = f2                     //错误
    f1 = float32(f2)            //正确
    ```
    ```go
    import "math"
    //p为自定义比较精度，如0.00001
    func IsEqual(f1, f2, p float64) bool {
        return math.Fdim(f1, f2) < p
    }
    ```
  - 字符串中子串可以像数组一样通过下标切片获取，用+拼接，用len函数获取长度，字符串不允许修改
    ```go
    var str1 = "hello"
    var str2 = str1[1:3]             //正确
    var str3 = str1[:2] + str1[3:4]  //正确
    c := len(str1)                   //正确
    str1[1] = 'x'                    //错误
    ```
  - Go语言支持两种字符类型，其中中bype代表UTF-8字符串单个字节的值，实际上是uint8的别名，rune代表单个Unicode字符
  - 数组是值类型，因此数组作为参数传递的时候无法在函数内修改原数组；切片是一个增强的数组，但是是引用类型，其内部结构可以抽象为三个变量：指向数组的指针、切片中的元素个数、切片已分配的存储空间。切片可以动态改变大小，可以基于数组创建，也可以直接使用make或者[]初始化创建，对数组的操作均支持切片，可以使用copy对切片作内容复制。
  - go数据类型分为值类型和引用类型，除slice、map、channel和interface是引用类型外，其他均是值类型（包括函数和struct）。相应的值传递和引用传递存在是否会改变原变量的值的问题。通过指针传递值类型变量的地址可以实现改变原变量的值
    ```go
    //值传递
    var arr1, arr2 [3]int
    var pt1 *[3]int
    arr1 = [3]int{1, 2, 3}
    arr2 = arr1
    arr2[1] = 0
    fmt.Println(arr1, arr2)     //输出  [1 2 3] [1 0 3]
    pt1 = &arr1
    (*pt1)[1] = 4
    fmt.Println(arr1, arr2)     //输出  [1 4 3] [1 0 3]
    //引用传递
    var sli1, sli2 []int
    sli1 = []int{1, 2, 3}
    sli2 = sli1
    sli2[1] = 0
    fmt.Println(sli1, sli2)     //输出  [1 0 3] [1 0 3]
    ```
    ```go
    sli1 := []int{1, 2, 3, 4, 5}
    sli2 := []int{5, 4, 3}
    copy(sli2, sli1)            //只会复制sli1的前3个元素到sli2中
    copy(sli1, sli2)            //只会复制sli2的3个元素到sli1的前3个位置，后2个位置的值不变
    ```
  - map类型可以直接通过value， ok := myMap["a"]的方式判断键a是否存在，这种语法同样支持判断channel变量是否已关闭和判断接口的对象实例是否实现了某个接口
    ```go
    if value, ok := myMap["a"]; ok {
        fmt.Println(value)
    }
    if file2, ok := file1.(IStream); ok {
        ...                   //file2成为实现了IStream接口的对象
    }
    ```

### 流程控制

- 条件语句
  - if后的条件无需括号
  - 无论语句体内有几条语句，花括号都是必须的
  - 左花括号必须与if或者else在同一行
  - 在if之后，条件语句之前，可以添加变量初始化语句，使用;间隔
  - 在有返回值的函数中，不允许将最终的return语句包含在if...else...结构中
    ```go
    func exam (x int) bool {}
        if a, b := 1, 2; a < i {
            return true
        } else {
            return false
        }
        return false               // 没有最后的return语句将报错
    }
    ```
- 选择语句
  - 左花括号必须与switch在同一行
  - 条件表达式不限制为常量或者整数
  - 单个case中可以有多个结果选项
  - 不需要用break明确退出一个case
  - 只有在case中明确添加fallthrough关键字，才会继续执行下一个case
  - switch后可以不需要表达式
    ```go
    switch i {
        case 0:
            fmt.Println("0")
        case 1:
            fallthrough
        case 2, 3, 4:
            fmt.Println("2, 3, 4")
        default:
            fmt.Println("Default")
    }
    switch {
        case i == 0:
            fmt.Println("0")
        case i == 1:
            fallthrough
        case i == 2 || i == 3 || i == 4:
            fmt.Println("2, 3, 4")
        default:
            fmt.Println("Default")
    }
    ```
  - switch语句可以用于进行变量类型查询（包括接口变量）
    ```go
    var v1 interface{} = ...
    switch v := v1.(type) {            //v变成具体的类型
        case int:
        case string:
        ...
    }
    ```
- 循环语句
  - go语言只有for循环，没有do while循环，无条件的for循环可以实现do循环
  - 左花括号必须与for在同一行
  - for循环条件定义中可以初始化变量，但多个赋值语句只能通过平行赋值实现
  - for循环支持break和continue，break支持选择要中断的循环
  - for和range配合可以实现复合数据类型的遍历
  - go语言支持goto跳转语句
    ```go
    var str = [4]string{"tom", "jack", "luna", "lucy"}
    for _, value := range str {          //丢弃数组元素的索引值
        for i, j := 0, 10; i < j; i ++ {
            if i % 2 == 0 {
                continue
            } else {
                if i == j - 1 {
                    break JLoop
                } else {
                    fmt.Println(i, value)
                }
            }
        }
    }
    JLoop:
    //...
    ```
    ```go
    func f() {
        i := 0
        HERE:
        fmt.Println(i)
        i++
        if i < 10 {
            goto HERE
        }
    }
    ```

### 函数

- 函数支持多返回值，调用时不需要的返回值可以使用_进行占位丢弃，多返回值可以用于返回错误
- 函数返回值可仅指定类型也可同时指定返回值变量名，返回值变量名可以直接在函数内部使用，指定返回值变量名后return时可以不指定返回内容
    ```go
    func getName() (firestname, lastName string, err error) {
        firstName = "May"
        lastName = "Chen"
        err = nil
        return
    }
    _, lastName, _ := getName()
    ```

- 函数是值类型，可以作为参数传递，支持匿名函数和闭包
    ```go
    f := func() func(int, int) int {
        return func(x, y int) int {
            return x + y
        }
    }
    df := func(x, y int) int {
        return 3 * f()(x, y)
    }
    res := func(x, y int, tf func(int, int) int) int {
        return tf(x, y)
    }(1, 2, df)
    fmt.Println(res)
    ```
- 函数支持不定参数和任意类型的不定参数,fmt.Println函数就利用了这一特性
    ```go
    var v1 = []int{1, 2, 3, 4}
    var v2 = []interface{}{1, 2.5, "a", nil}  //因为可以认为任意类型都实现了空接口，所以空接口类型可以支持任意类型的变量
    func f1(args ...int) {
        for _, arg := range args {
            fmt.Println(arg)
        }
    }
    func f2(args ...interface{}) {
        for _, arg := range args{
            switch arg.(type) {
                case int:
                    fmt.Println(arg, "is an int value")
                case string:
                    fmt.Println(arg, "is a string value")
                default:
                    fmt.Println(arg, "is an unknown type")
            }
        }
    }
    f1(v1...)                 //原样传递
    f2(v2[:4]...)             //传递片段
    ```

## 面向对象

### 类

Go语言中没有专门封装的面向对象类型，而是依赖内置和自定义类型，是一种更加开放式的面向对象，因此可以给除指针外的其他类型添加方法，定义新的类型，即类。Go的类没有this指针，而是在成员方法中直接将this指针暴露出来。Go语言中类和类的成员变量以及成员方法的在包外的可见性也和变量和函数的包外可见性规则一样。

    ```go
    type Int int                       //int类型可以自动转换为Int类型
    func (a Int) Less(b Int) bool {    //成员方法无法修改对象本身
        return a < b
    }
    func (a *Int) Add(b Int) bool {    //成员方法可以修改对象本身
        *a += b
    }
    func main() {
        var a Int = 1
        a.Add(2)                       //对象的值被修改
        if a.Less(5) {
            fmt.Println("YES")
        }
    }
    ```
    ```go
    //List :实现python中的List数据类型部分特性（删除指定位置元素）
    type List []interface{}           //自定义空接口切片类型
    func (a *List) RemoveAt(idx int) {
        *a = append((*a)[:idx], (*a)[idx+1:]...)
    }
    var s List = []interface{}{1, 2, 3, 4.5, "a"}
    s.RemoveAt(2)
    ```
- 类对象实例初始化
    ```go
    type Rect struct {
        x, y float64
        width, height float64
    }
    var rect1 *Rect =  new(Rect)
    rect2 := new(Rect)
    rect3 := &Rect{}
    rect4 := &Rect{0, 10, 100, 200}  //顺序为成员变量初始化值，剩余成员变量使用默认值
    rect5 := &Rect{width: 100, height: 200}    //为指定成员变量赋值，未指定的成员变量使用默认值
    //也可以创建单独的对象构建函数来创建对象
    func NewRect(x, y, width, height float64) *Rect {
        return &Rect(x, y, width, height)
    }
    ```
- Go语言通过匿名组合的方式实现类的继承，子类可以直接使用父类的成员变量和方法，也可以重写父类方法
    ```go
    type Person struct {
        Name string
        Age int
    }
    func (p *Person) GetName() string {
        return p.Name
    }
    func (p *Person) Show() {
        fmt.Println(p.Name, p.Age)
    }

    type Student struct {
        Person
        College string
    }
    func (s *Student) GetCollege() string {
        return s.College
    }
    func (s *Student) Show() {
        fmt.Println(s.Name, s.Age, s.College)
    }
    var stu *Student = new(Student)
    stu.Name = "TOM"
    stu.Age = 21
    stu.College = "UU"
    stu.Show()                     //调用的是子类重写后的Show函数
    fmt.Println(stu.GetName())     //直接使用父类方法，等同于stu.Person.GetName()
    ```

### 接口

- Go语言的接口是非侵入式的，类不需要指定实现某个接口再逐一实现该接口的方法，而是只要类实现了某个接口的所有方法，就表示实现了该接口，即类和接口不需要显式关联
    ```go
    type IFile interface {
        Read(buf []byte) (n int, err error)
        Write(buf []byte) (n int, err error)
        Seek(off int64, whence int) (pos int64, err error)
        Close() error
    }
    type IReader interface {
        Read(buf []byte) (n int, err error)
    }
    type IWriter interface {
        Write(buf []byte) (n int, err error)
    }
    type File struct {
        // ...
    }
    func (f *File) Read(buf []byte) (n int, err error)
    func (f *File) Write(buf []byte) (n int, err error)
    func (f *File) Seek(off int64, whence int) (pos int64, err error)
    func (f *File) Close() error
    //File类实现了IFile、IReader和IWriter三个接口的所有方法，就表示File类实现了这三个接口，就可以将File类的对象赋值给接口
    var file1 IFile = new(File)
    var file2 IReader = new(File)
    var file3 IWriter = new(File)
    ```
- 接口赋值有两种情况：将对象实例赋值给接口和将一个接口赋值给另一个接口。将对象实例赋值给接口时应当传递对象的地址，即赋指针，如果直接传递对象本身可能导致编译失败（当对象所在的类包含指针方法时）；具有相同方法的接口是等价的，接口只可以赋值给等价或子集接口
    ```go
    type Int int
    func (a Int) Less(b Int) bool {
        return a < b
    }
    func (a *Int) Add(b Int) {
        *a += b
    }
    type Lesser interface {
        Less(b Int) bool
    }
    type Adder interface {
        Add(b Int)
    }
    type LessAdder interface {
        Less(b Int) bool
        Add(b Int)
    }
    var num Int = 3
    var l1 Lesser = num              //正确
    var l2 Lesser = &num             //正确
    var a1 Adder = num               //错误
    var a2 Adder = &num              //正确
    var la1 LessAdder = num          //错误
    var la2 LessAdder = &num         //正确
    a1 = la1                         //正确
    la1 = a1                         //错误
    ```
- 接口和类一样支持匿名组合，接口可以支持和map、channel一样的查询文法，空接口可以接受任意类型的对象实例，可以作类型查询
    ```go
    //接口组合
    type IReader interface {
        Read(buf []byte) (n int, err error)
    }
    type IWriter interface {
        Write(buf []byte) (n int, err error)
    }
    type IReadWriter interface {
        IReader
        IWriter
    }
    //接口查询
    var rw = IReadWriter = ...
    if i, ok := rw.(type); ok {
        ...
    }
    //空接口和类型查询
    var v1 = []interface{}{1, 10.1, "aa", nil}
    switch v := v1.(type) {
        case int:
        case string:
        ...
    }
    ```

### 错误处理

Go语言通过利用error接口实现自定义错误类型来处理各种错误，使用panic函数抛出异常，使用defer延迟关键字和recover函数捕获异常。

- defer
    关键字用于标记函数体最后执行的语句，在函数关闭前调用，多个defer定义遵循FILO原则，类似于栈操作，最先定义的最后执行。
    ```go
    // 代码一
    func funcA() int {
        x := 5
        defer func() {
            x += 1
        }()
        return x
    }
    // 代码二
    func funcB() (x int) {
        defer func() {
            x += 1
        }()
        return 5
    }
    // 代码三
    func funcC() (y int) {
        x := 5
        defer func() {
            x += 1
        }()
        return x
    }
    // 代码四
    func funcD() (x int) {
        defer func(x int) {
            x += 1
        }(x)
        return 5
    }
    ```
    明确以下几点：

  - return xxx语句并非原子指令，实际执行时是先执行返回变量=xxx，然后执行return
  - defer在函数关闭前执行，即在return之前，而不是在return xxx之前
  - 普通变量赋值和传递参数都是值传递

    上述四个函数可以作以下解析：
    ```go
    // 解析代码一：返回temp的值，在将x赋值给temp后，temp未发生改变，最终返回值为5
    func funcA() int {
        x := 5
        temp=x      //temp变量表示未显示声明的return变量
        func() {
            x += 1
        }()
        return
    }
    // 解析代码二：返回x的值，先对其复制5，接着函数中改变为6，最终返回值为6
    func funcB() (x int) {
        x = 5
        func() {
            x += 1
        }()
        return
    }
    // 解析代码三：返回y的值，在将x赋值给y后，y未发生改变，最终返回值为5
    func funcC() (y int) {
        x := 5
        y = x      //这里是值拷贝
        func() {
            x += 1
        }()
        return
    }
    // 解析代码四：返回x的值，传递x到匿名函数中执行时，传递的是x的拷贝，不影响外部x的值，最终返回值为5
    func funcD() (x int) {
        x := 5
        func(x int) { //这里是值拷贝
            x += 1
        }(x)
        return
    }
    ```
- Error
    Go函数支持利用多返回值特性将错误和正常返回值一起返回，通过实现error接口自定义错误类型处理错误信息
    ```go
    //Go内置error接口定义
    type error interface {
        Error() string
    }
    ```
    ```go
    package main
    import (
        "fmt"
    )
    // 定义一个 DivideError 结构
    type DivideError struct {
        dividee int
        divider int
    }
    // 实现  `error` 接口
    func (de *DivideError) Error() string {
        strFormat := `
            Cannot proceed, the divider is zero.
            dividee: %d
            divider: 0`
        return fmt.Sprintf(strFormat, de.dividee)
    }
    // 定义 `int` 类型除法运算的函数
    func Divide(varDividee int, varDivider int) (result int, errorMsg string) {
        if varDivider == 0 {
            dData := DivideError{
                dividee: varDividee,
                divider: varDivider,
            }
            errorMsg = dData.Error()
            return
        } else {
            return varDividee / varDivider, ""
        }
    }
    func main() {
        // 正常情况
        if result, errorMsg := Divide(100, 10); errorMsg == "" {
            fmt.Println("100/10 = ", result)
        }
        // 当被除数为零的时候会返回错误信息
        if _, errorMsg := Divide(100, 0); errorMsg != "" {
            fmt.Println("errorMsg is: ", errorMsg)
        }
    }
    ```
    如果不需要对错误作特殊处理，可以定义一个专门的错误处理函数，避免频繁的书写if err != nil代码块
    ```go
    func checkError(err error) {
        if err != nil {
            fmt.Println("Error: ", err)
            os.Exit(-1)
        }
    }
    ```
- panic和recover
    panic和recover是两个内置函数，用于处理运行时异常，程序发生运行时异常时会在函数体内自动抛出panic，也可以使用panic函数手动抛出panic，发生panic后函数会退出运行，在退出前会执行defer代码块，所以可以在defer代码块内利用recover函数捕获panic异常。
    ```go
    //panic和recover函数的内置定义
    func panic(interface{})
    func recover() interface{}
    ```
    ```go
    package main
    import (
        "fmt"
    )
    // 最简单的例子
    func SimplePanicRecover() {
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("Panic info is: ", err)
            }
        }()
        panic("SimplePanicRecover function panic-ed!")
    }
    // 当 defer 中也调用了 panic 函数时，最后被调用的 panic 函数的参数会被后面的 recover 函数获取到
    // 一个函数中可以定义多个 defer 函数，按照 FILO 的规则执行
    func MultiPanicRecover() {
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("Panic info is: ", err)
            }
        }()
        defer func() {
            panic("MultiPanicRecover defer inner panic")
        }()
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("Panic info is: ", err)
            }
        }()
        panic("MultiPanicRecover function panic-ed!")
    }
    // recover 函数只有在 defer 函数中被直接调用的时候才可以获取 panic 的参数
    func RecoverPlaceTest() {
        // 下面一行代码中 recover 函数会返回 nil，但也不影响程序运行
        defer recover()
        // recover 函数返回 nil
        defer fmt.Println("recover() is: ", recover())
        defer func() {
            func() {
                // 由于不是在 defer 调用函数中直接调用 recover 函数，recover 函数会返回 nil
                if err := recover(); err != nil {
                    fmt.Println("Panic info is: ", err)
                }
            }()
        }()
        defer func() {
            if err := recover(); err != nil {
                fmt.Println("Panic info is: ", err)
            }
        }()
        panic("RecoverPlaceTest function panic-ed!")
    }
    // 如果函数没有 panic，调用 recover 函数不会获取到任何信息，也不会影响当前进程。
    func NoPanicButHasRecover() {
        if err := recover(); err != nil {
            fmt.Println("NoPanicButHasRecover Panic info is: ", err)
        } else {
            fmt.Println("NoPanicButHasRecover Panic info is: ", err)
        }
    }
    func main() {
        SimplePanicRecover()
        MultiPanicRecover()
        RecoverPlaceTest()
        NoPanicButHasRecover()
    }
    ```
    如果不需要对panic作特殊处理，可以定义一个专门处理panic的函数供defer调用，避免在函数体中加入大量的recover代码块
    ```go
    func CallRecover() {
        if err := recover(); err != nil {
            fmt.Println("Panic info is: ", err)
        }
    }
    func Divide(x, y int) int {
        defer CallRecover()
        // panic("RecoverInOutterFunc function panic-ed!")
        return x / y
    }
    fmt.Println(Divide(3, 0))
    ```

## 并发

### 并发模型

相比传统的串行编程，并发能更客观的表现问题模型，能够充分利用CPU核心的优势和CPU与IO固有的异步性，提高程序的执行效率。主流的并发模型有四种：

- 多进程。实现简单，进程间相互独立，相互竞争，所有进程由内核维护，上下文切换开销大，进程间通信实现麻烦
- 多线程。同一进程的线程间共享全局变量和内存，线程间通信比较简单，系统开销比多进程小，但高并发模式下，开销依然不小，影响执行效率
- 异步回调。通过事件驱动的方式使用异步IO，系统开销极低，但编程比多线程复杂，且不利于反映问题模型
- 协程。协程是一种用户态的轻量级线程，有自己的寄存器上下文和栈，无需线程切换，由程序本身控制上下文的切换，开销小，高并发下可支持大量协程，且编程简单，结构清晰，协程是不同于回调的另一种异步方式，代码直观感受是以看似同步的代码实现异步

### goroutine

协程本质上是线程内部通过中断进行上下文切换，在很多语言中并非原生支持，php、python等语言中通过迭代器/生成器的可中断特性来模拟协程，具体实现需要编码完成，而go语言原生支持协程，直接提供了实现好的协程工具来满足并发需求，即是goroutine，可以直接使用go关键字将函数以协程运行。类似的，在python中有一个第三方协程库greenlet，但由于python并非原生支持协程，所以并不是所有代码都可以使用greenlet直接以协程的方式运行，约束较多。go语言中任何函数（包括匿名函数）都可以协程运行，且无需对函数编写专门的代码，使用非常方便。

    ```go
    package main
    import "fmt"
    func f(from string) {
        for i := 0; i < 3; i++ {
            fmt.Println(from, ":", i)
        }
    }
    func main() {

        f("direct")                 //直接运行f函数

        go f("goroutine")           //协程运行f函数

        go func(msg string) {       //协程运行匿名函数
            fmt.Println(msg)
        }("going")
        //匿名函数的运行结果可能先于f函数输出，这表明goruntine协程是异步运行的
        var input string
        fmt.Scanln(&input)          //等待键盘输入
        fmt.Println("done")
    }
    ```

### channel

上例中如果删除fmt.Scanln(&input)语句，goroutine协程输出将不会显示，因为goruntine是异步执行的，主程序发起goruntine后不会等待其执行完成，主程序继续运行直到退出，goruntine也因主程序退出而终断。并发通信的两种主流模型是消息机制和共享内存，go语言为了解决并发编程中不同goruntine之间的通信问题，也提供了这两种方式，但是go更倾向于消息机制，在语言级别为消息通信机制提供了一种专门的数据类型即channel(通道)用于goruntine间的消息传递。

- channel是类型相关的引用类型，即一个channel只能传递一种类型的消息数据，作为参数传递后其值可以被修改
- channel默认读写都是阻塞的，只有当读写双方都就绪才能完成消息传递，程序才能继续执行
- channel类似unix系统中的管道(pipe)，作为goruntine之间的消息传递通道，也可以支持缓冲，当channel通道缓冲已满时，写操作会阻塞，当通道中没有数据时，读操作会阻塞
- channel作为参数传递的时候可以指定为单向channel来保证程序的类型安全
- channel可以使用close函数关闭
    ```go
    //channel的声明和定义
    var chan int ch                           //channel声明， channel是引用类型，默认值是nil
    var chan string msg = make(chan string)   //channel定义
    msg := make(chan string)                  //channel也支持多种定义方式
    msg := make(chan string, 1)               //支持缓冲的channel
    ```
    ```go
    //单向channel
    func ping(pings chan<- string, msg string) {             //只写
        pings <- msg
    }
    func pong(pings <-chan string, pongs chan<- string) {    //只读
        msg := <-pings
        pongs <- msg
    }
    pings := make(chan string, 1)             //注意这里如果定义为无缓冲的channel，程序会报错，所有goruntine均是asleep状态，go就认为发生了死锁
    pongs := make(chan string, 1)
    ping(pings, "passed message")
    pong(pings, pongs)
    fmt.Println(<-pongs)
    ```
    ```go
    //goruntine和channel实现的生产者-消费者模型
    package main
    import "fmt"
    func producer(c chan int) {
        for i := 0; i < 1000; i++ {
            fmt.Printf("Put product, ID is : %d \n", i)
            c <- i
            //time.Sleep(0.5 * 1e9)
        }
        defer close(c)                    //生产完成后关闭channel
    }
    func consumer(c, r chan int) {
        var p int
        hm := true
        for hm {
            if p, hm = <-c; hm {          //只有当channel关闭且channel缓冲内已经没有数据，hm值才为false，channel关闭后缓冲中的数据仍然可以被读取
                //time.Sleep(0.5 * 1e9)
                fmt.Printf("Get product, ID is : %d \n", p)
            }
        }
        r <- 1
    }

    func main() {
        c := make(chan int, 10)                   //如果c是无缓冲channel，输出会有什么不同？
        r := make(chan int)                       //r用于确认消费者是否将所有产品消费完
        go producer(c)
        go consumer(c, r)
        <-r
    }
    ```

### select

除了goruntine和channel之外，go语言还提供select(选择器)用于支持同时监控/等待多个channel操作,goruntine、channel、select的结合是Go的一个强大特性。

- select用法类似switch，与case结合使用实现多channel同时监控
- 与select结合的每个case语句的条件都必须是一个channel操作
- select默认只处理第一个已经就绪的channel操作，利用这一特性可以实现超时处理
- 如果有default子句且没有任何case的channel就绪，则直接执行default子句的内容，利用这一特性可以实现非阻塞的channel操作
- 如果没有default字句，select将阻塞，直到某个channel可以运行
    ```go
    package main
    import "time"
    import "fmt"
    func main() {
        c1 := make(chan string)
        c2 := make(chan string)
        go func() {
            time.Sleep(time.Second * 1)
            c1 <- "one"
        }()
        go func() {
            time.Sleep(time.Second * 2)
            c2 <- "two"
        }()
        for i := 0; i < 2; i++ {        //因为select只处理一次channel操作，如果去掉for循环，必然有一个goruntine的channel操作无限阻塞，产生死锁
            select {
                case msg1 := <-c1:
                    fmt.Println("received", msg1)
                case msg2 := <-c2:
                    fmt.Println("received", msg2)
            }
        }
    }
    ```
    ```go
    //使用select+channel实现超时
    package main
    import "time"
    import "fmt"
    func main() {
        c1 := make(chan string, 1)
        go func() {
            time.Sleep(time.Second * 2)
            c1 <- "result 1"
        }()
        select {
            case res := <-c1:
                fmt.Println(res)
            case <-time.After(time.Second * 1):                  //time.After函数会在指定时间后返回一个chan Time对象
                fmt.Println("timeout 1")
        }
        c2 := make(chan string, 1)
        go func() {
            time.Sleep(time.Second * 2)
            c2 <- "result 2"
        }()
        select {
            case res := <-c2:
                fmt.Println(res)
            case <-time.After(time.Second * 3):
                fmt.Println("timeout 2")
        }
    }
    ```
    ```go
    //select+channel实现非阻塞channel
    package main
    import "fmt"
    func main() {
        messages := make(chan string)
        signals := make(chan bool)
        select {
            case msg := <-messages:             //读就绪而写没有就绪，直接进入default子句，messages通道是无缓冲的，写操作没有执行成功
                fmt.Println("received message", msg)
            default:
                fmt.Println("no message received")
        }
        msg := "hi"
        select {
            case messages <- msg:               //上一个select进入了default子句，其messages通道的读就绪也就不存在了，所以这里写就绪而读没有就绪
                fmt.Println("sent message", msg)
            default:
                fmt.Println("no message sent")
        }
        select {
            case msg := <-messages:
                fmt.Println("received message", msg)
            case sig := <-signals:
                fmt.Println("received signal", sig)
            default:
                fmt.Println("no activity")
        }
    }
    ```

### 多核并发

- 现代CPU都是以线程作为最小调度单位的
- Go语言中的goruntine是线程内部的代码中断和上下文切换
- go的runtime会产生多个内核线程，每一个线程内都有一个goruntine调度器和运行器
- 每一个运行器上同时刻只能有一个goruntine运行，调度器负责线程内部多个goruntine之间的调度切换(未运行的goruntine位于线程内的任务队列中)
- 如有需要，goruntine还可以被分配到其他线程中
- 进程和线程由操作系统内核负责调度，goruntine由go运行时和线程内调度器调度
- 程序员编写并发和多核并发代码时无需关心线程（go语言不允许在代码中创建线程），只需要编写业务逻辑然后生成goruntine运行即可

### 锁和全局唯一性操作

- 对于goruntine并发，Go语言也提供了同步锁：sync.Mutex和sync.RWMutex。Mutex是最简单的排它锁，当一个goruntine获得Mutex，其他goruntine就只能等待其释放Mutex，RWMutex则支持读写锁分离，多个goruntine可同时获得读锁RLock，但写锁WLock会阻止任何其他goruntine获得锁
- Go语言提供了一个Once类型用于保证全局唯一性操作，即多个goruntine都运行一段代码，Once类型可以保证指定的代码段只会执行一次
    ```go
    package main
    import (
        "fmt"
        "sync"
    )
    var a int = 1
    var once sync.Once
    func setup() {
        a++
    }
    func doprint() {
        once.Do(setup)
        fmt.Println(a)
    }
    func twoprint() {
        go doprint()
        go doprint()
    }
    func main() {
        var input string
        twoprint()
        fmt.Scanln(&input)
    }
    ```
