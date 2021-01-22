# Helm Chart 仓库维护说明

本仓库一个多版本的Helm Chart仓库, 将作为云昊平台默认的应用商店源.

## Chart源

本仓库的Chart源主要来自三种途径:

- Rancher官方维护的Chart仓库:  https://github.com/rancher/charts/tree/master/charts
- HelmHub收集的各应用的官方或第三方Helm Charts: https://hub.helm.sh/
- 公司内部应用产出的Helm Charts.

## 维护规范

本仓库需严格遵循目录规范, 目录结构如下:

```bash
-- stable
  |-- README.md
  |-- charts
     |-- nginx
     |  |-- logo.png
     |  |-- v6.2.0
     |  |-- v6.2.1
     |-- redis
        |-- logo.png
        |-- v9.0.2
        |-- v10.8.1
```

每个Helm Chart应用一个目录, 目录中包含一个图标文件和应用的各版本目录. 仓库当前维护各应用的最新版Charts, 需遵循以下规范:

- 仓库原则上只需维护各应用的最新版Charts, 除非特定旧版应用是必要的
- 应用图标文件名通常为**logo.png**, 尺寸不宜小于**500*500**
- 应用版本目录均以**vxx.xx.xx**的形式命名
- 应用的每个版本目录下的**Chart.yaml**文件中均需通过`icon: file://../logo.png`字段指定有效的图标文件
- 应用的每个版本目录下均需添加**questions.yml**文件, 该文件中指定了应用分类, 标签和自定义字段渲染规则

## question文件说明

questions.yml文件是对Helm Charts标准结构的补充. 文件中定义的字段将自动在云昊平台应用商店中可视化呈现. 典型的文件样式如下:

```yaml
rancher_max_version: 2.4.5
rancher_min_version: 2.4.5
categories:
- Security
labels:
  io.rancher.certified: partner   #partner
questions:
- variable: defaultImage
  default: true
  description: "Use default Docker image"
  label: Use Default Image
  type: boolean
  show_subquestion_if: false
  group: "Container Images"
  subquestions:
  - variable: image.repository
    default: "quay.io/jetstack/cert-manager-controller"
    description: "Cert-Manager Docker image name"
    type: string
    label: Cert-Manager Docker Image Name
```

该文件一个包括五个字段:

- rancher_min_version: 表示该Helm Charts要求平台至少需要满足的版本
- rancher_max_version: 表示该Helm Charts可用的最高平台版本
- categories: 表示应用分类, 每个应用均可被划分到多个分类中
- labels: 表示应用的额外标签, `io.rancher.certified`标签可将应用标记为`partner`, `experimental`或`operator`, 表示官方认证或实验性的
- questions: 表示应用的自定义字段渲染规则, 其中的字段与values.yaml文件中是一一对应的, 应用商店依此将应用的自定义字段渲染为可视化页面

目前对于所有应用, 只有**categories**字段是必须的, 后续将要求所有应用的**questions**字段都进行补充完善和汉化, 以提供更好的平台使用体验.

rancher官方维护的Chart仓库均提供了完善的`questions.yml`文件, 其他Chart仓库均不提供. 仓库将以rancher仓库为上游进行持续更新和增强.

  ## question变量参考

| 变量                | 类型          | 必需  | 描述                                                         |
| ------------------- | ------------- | ----- | ------------------------------------------------------------ |
| variable            | string        | true  | 定义`values.yaml`文件中指定的变量名, 使用`foo.bar`的形式处理嵌套对象 |
| label               | string        | true  | 定义UI界面上显示的变量标签                                   |
| description         | string        | false | 指定变量描述信息                                             |
| type                | string        | false | 变量类型, 当前支持的类型包括string, multiline, boolean, int, enum, password, storageclass. hostname, pvc和secret, 默认类型为`string` |
| required            | bool          | false | 定义变量是否是必须的(true \| false)                          |
| default             | string        | false | 指定变量的默认值                                             |
| group               | string        | false | 变量分组, 变量在UI上将按照分组进行展示                       |
| min_length          | int           | false | 变量的最小字符串长度                                         |
| max_length          | int           | false | 变量的最大字符串长度                                         |
| min                 | int           | false | 变量的最小值                                                 |
| max                 | int           | false | 变量的最大值                                                 |
| options             | []string      | false | 对于`enum`类型的变量指定其可选项, 如:  options: - "ClusterIP" - "NodePort" - "LoadBalancer" |
| valid_chars         | string        | false | 通过正则表达式验证输入值是否合法                             |
| invalid_chars       | string        | false | 通过正则表达式验证输入值是否不合法                           |
| subquestions        | []subquestion | false | 为变量增加一系列子变量(嵌套)                                 |
| show_if             | string        | false | 当条件变量满足时才显示当前变量, 如: `show_if: "serviceType=Nodeport"` |
| show_subquestion_if | string        | false | 当条件变量满足时才显示当前变量的子变量, 如: `show_subquestion_if: "true"` |

## 应用分类

所有应用均需在**questions.yml**文件中指定其所属分类. 所有可选的分类包括:

- WebServer: web服务器相关应用
- Loadbalancer: 负载均衡相关应用
- DevOps: devops相关应用
- MicroService: 微服务相关应用
- Serverless: serverless相关应用
- BigData: 大数据相关应用
- AI: 人工智能相关应用
- MQ: 消息队列中间件相关应用
- Database: 数据库相关应用, 包括关系数据库, 文档数据库, 时序数据库, KV数据库, 缓存数据库, 数据仓库, 数据库管理软件等
- Storage: 存储相关应用, 包括文件存储, 块存储, 对象存储, 存储驱动, 存储编排等
- Repository: 仓库类应用, 包括制品仓库, 镜像仓库, 代码仓库, 依赖仓库, Chart仓库等
- Observability: 可观察性相关应用, 包括指标度量, 日志聚合, 链路追踪, 混沌工程等
- Security: 安全相关应用, 包括安全检测, 漏洞扫描, 证书管理, 签名验签等
- Other: 其他应用

**不得擅自为应用指定非上述分类的类别.** 如需为应用指定更为清晰明确的关键字, 可以在应用的`Chart.yaml`文件的**keywords**字段中列出, 该字段将有助于在平台页面中检索应用(后续支持). 