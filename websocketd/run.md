# README

Turn any program that uses STDIN/STDOUT into a WebSocket server

## 运行

- 运行后台

```bash
websocketd --port=8181 bash ./count.sh   # 支持指定地址(默认bindall), 端口, 日志级别等, 后接任何命令
```

- 打开网页

```bash
open ./count.html   # 自动调用默认浏览器(linux下为xdg-open), 也可手动点击打开
```

## 项目

https://github.com/joewalnes/websocketd

## 说明

websocketd --help 查看使用说明