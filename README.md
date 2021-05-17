# GoMiraiBot

## 概要

* 以`mirai`为框架，配合插件`mirai-api-http`和`golang`制作的简易QQ机器人<del>萝卜子</del> 

## Useage

* `golang` 环境
* 在运行目录下添加 `config.yml` 文件，内容如下

```yml
# mirai-api-http 的服务配置
server:
  host: 127.0.0.1 # mirai-api-http 的服务位置
  port: 8848 # mirai-api-http 的端口

# bot 基本信息配置
bot: 
  QQ: 114145 # 目标机器人的qq号
  authKey: AUTHENTICATION_KEY # authKey 由mirai-http-api 提供

# SQL 关系型数据库配置
database:
  # 启动数据库
  enable: true
  # 数据库启动类型 【None|Create|Update】
  mode: Update
  # 数据库名称,目前只支持mysql
  db: mysql
  # 数据库连接相关信息
  dbName: bot_database
  dbUser: bot
  dbPassword: 114145
  dbHost: localhost
  dbPort: 3306
  # orm配置相关
  orm:

# redis NoSQL 数据库配置
redis:
  # 启用redis
  enable: false

# 机器人的功能的额外配置
extraConfig:
  #随机涩图配置 https://api.lolicon.app
  setu:
    url: https://api.lolicon.app/setu/
    apiKey: 1141451919
  #每日60s看世界
  daily:
    url: http://dwz.2xb.cn/zaob
  #明日方舟蹲饼配置
  arknights-dun:
    data: data

```  

* 使用 `go run .`直接运行 或者 `go build .` 编译后运行 `./main`

## 功能支持

**目前仍然在开发中，不稳定**  

### 命令交互

1. #help 显示帮助信息
2. #about 显示关于
3. #setu 随机涩图
4. #s-img 图片以图搜图服务
5. #记账 记账

### 消息类型交互
* XML
* APP

### 事件交互
