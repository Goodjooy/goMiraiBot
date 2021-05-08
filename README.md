# GoMiraiBot

## 概要
* 以`mirai`为框架，配合插件`mirai-api-http`和`golang`制作的简易QQ机器人<del>萝卜子</del> 

## Useage
* `golang` 环境
* 在运行目录下添加 `config.yml` 文件，内容如下

```yml
server:
  host: 127.0.0.1 # mirai-api-http 的服务位置
  port: 8848 # mirai-api-http 的端口

bot: 
  QQ: 114145 # 目标机器人的qq号
  authKey: AUTHENTICATION_KEY # authKey

```
* 使用 `go run .`直接运行 或者 `go build .` 编译后运行 `./main`

## 功能支持

**目前仍然在开发中，不稳定**

### 单次交互

1. #help 显示帮助信息
2. #about 显示关于

### 上下文交互

1. #s-img 图片以图搜图服务
