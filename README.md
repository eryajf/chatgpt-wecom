<p align='center'>
<br>
    🚀 ChatGPT WeCom 🚀
</p>

<p align='center'>🌉 基于GO语言实现的体验最好的企微应用集成ChatGPT项目 🌉</p>

<div align="center">

[![Auth](https://img.shields.io/badge/Auth-eryajf-ff69b4)](https://github.com/eryajf)
[![Go Version](https://img.shields.io/github/go-mod/go-version/eryajf/chatgpt-wecom)](https://github.com/eryajf/chatgpt-wecom)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/eryajf/chatgpt-wecom)](https://github.com/eryajf/chatgpt-wecom/pulls)
[![GitHub Pull Requests](https://img.shields.io/github/stars/eryajf/chatgpt-wecom)](https://github.com/eryajf/chatgpt-wecom/stargazers)
[![HitCount](https://views.whatilearened.today/views/github/eryajf/chatgpt-wecom.svg)](https://github.com/eryajf/chatgpt-wecom)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/eryajf/chatgpt-wecom)](https://hub.docker.com/r/eryajf/chatgpt-wecom)
[![Docker Pulls](https://img.shields.io/docker/pulls/eryajf/chatgpt-wecom)](https://hub.docker.com/r/eryajf/chatgpt-wecom)
[![GitHub license](https://img.shields.io/github/license/eryajf/chatgpt-wecom)](https://github.com/eryajf/chatgpt-wecom/blob/main/LICENSE)

</div>

<img src="https://camo.githubusercontent.com/82291b0fe831bfc6781e07fc5090cbd0a8b912bb8b8d4fec0696c881834f81ac/68747470733a2f2f70726f626f742e6d656469612f394575424971676170492e676966" width="800"  height="3">
</div><br>
<img src='https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_215927.jpg' alt='' />


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 前言

本项目可以助你将GPT机器人集成到企微应用当中。当前默认模型为`gpt-3.5`，支持`gpt-4`，同时支持Azure-OpenAI。

>- `📢 注意`：当下部署以及配置流程都已非常成熟，文档和issue中基本都覆盖到了，因此不再回答任何项目安装部署与配置使用上的问题，如果完全不懂，可考虑通过 **[邮箱](mailto:eryajf@163.com)** 联系我进行付费的技术支持。


🥳 **欢迎关注我的其他开源项目：**
>
> - [Go-Ldap-Admin](https://github.com/eryajf/go-ldap-admin)：🌉 基于Go+Vue实现的openLDAP后台管理项目。
> - [learning-weekly](https://github.com/eryajf/learning-weekly)：📝 周刊内容以运维技术和Go语言周边为主，辅以GitHub上优秀项目或他人优秀经验。
> - [HowToStartOpenSource](https://github.com/eryajf/HowToStartOpenSource)：🌈 GitHub开源项目维护协同指南。
> - [read-list](https://github.com/eryajf/read-list)：📖 优质内容订阅，阅读方为根本
> - [awesome-github-profile-readme-chinese](https://github.com/eryajf/awesome-github-profile-readme-chinese)：🦩 优秀的中文区个人主页搜集

🚜 我还创建了一个项目 **[awesome-chatgpt-answer](https://github.com/eryajf/awesome-chatgpt-answer)** ：记录那些问得好，答得妙的时刻，欢迎提交你与ChatGPT交互过程中遇到的那些精妙对话。

⚗️ openai官方提供了一个 **[状态页](https://status.openai.com/)** 来呈现当前openAI服务的状态，同时如果有问题发布公告也会在这个页面，如果你感觉它有问题了，可以在这个页面看看。

## 功能介绍

- [x] 🚀 帮助菜单：通过发送 `帮助` 将看到帮助列表
- [x] 🙋 单聊模式：每次对话都是一次新的对话，没有历史聊天上下文联系
- [x] 🗣 串聊模式：带上下文理解的对话模式
- [x] 🎭 角色扮演：支持场景模式，通过 `#周报` 的方式触发内置prompt模板
- [x] 🧑‍💻 频率限制：通过配置指定，自定义单个用户单日最大对话次数
- [x] 🔗 自定义api域名：通过配置指定，解决国内服务器无法直接访问openai的问题
- [x] 🪜 添加代理：通过配置指定，通过给应用注入代理解决国内服务器无法访问的问题
- [x] 👐 默认模式：支持自定义默认的聊天模式，通过配置化指定
- [x] 👹 白名单机制：通过配置指定，支持指定群组名称和用户名称作为白名单，从而实现可控范围与机器人对话
- [x] 💂‍♀️ 管理员机制：通过配置指定管理员，部分敏感操作，以及一些应用配置，管理员有权限进行操作
- [x] ㊙️ 敏感词过滤：通过配置指定敏感词，提问时触发，则不允许提问，回答的内容中触发，则以 🚫 代替
- [ ] 🎨 图片生成：通过发送 `#图片`关键字开头的内容进行生成图片
- [ ] 📝 查询对话：通过发送`#查对话 username:xxx`查询xxx的对话历史，可在线预览，可下载到本地

## 使用前提

* 有Openai账号，并且创建好`api_key`，关于这个的申请创建，这里不过多赘述。
* 在企微开发者后台创建应用，需要如下五个配置信息：
  * `corp_id：`企业ID
  * `agent_id:` 应用ID
  * `agent_secret:` 应用秘钥
  * `receive_msg_token:` API接收消息的token
  * `receive_msg_key:` API接收消息的key


接下来跟随如下教程，保你一次配置成功。

## 使用教程

### 第一步，创建应用

1. [点我登陆](https://work.weixin.qq.com/wework_admin/frame)企业微信管理后台。

   此时点击我的企业可以获取到企业ID`corp_id`。
    <details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_210053.jpg">
    </details>

2. 点击应用，创建一个新的应用。

   <details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_210600.jpg">
    </details>

3. 创建成功之后，能够获取到应用的`agent_id`与`agent_secret`。

   <details>
         <summary>🖼 点我查看示例图</summary>
         <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_210807.jpg">
       </details>

4. 点击接收消息的`设置 API 接收`，进入回调地址配置页面，点击两个随机获取，可以获得`receive_msg_token`与`receive_msg_key`。

   <details>
         <summary>🖼 点我查看示例图</summary>
         <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_211253.jpg">
       </details>

- `📢 注意：`此时点击保存会提示 `openapi回调地址请求不通过`，是因为保存配置的时候企微会向服务发送请求进行校验。所以此时这个页面先放这儿，去部署应用。应用部署起来之后，再回来保存配置就能成功了。


### 第二步，部署项目

#### docker部署

推荐你使用docker-compose快速运行本项目。

```yaml
version: '3'
services:
  chatgpt:
    container_name: chatgpt
    image: registry.cn-hangzhou.aliyuncs.com/ali_eryajf/chatgpt-wecom
    restart: always
    environment:
      LOG_LEVEL: "info" # 应用的日志级别 info/debug
      CORP_ID: "" # 企业ID
      AGENT_ID: "" # 应用ID
      AGENT_SECRET: "" # 应用秘钥
      RECEIVE_MSG_TOKEN: "" # API接收消息的token
      RECEIVE_MSG_KEY: "" # API接收消息的key
      APIKEY: xxxxxx  # 你的 api_key
      BASE_URL: ""  # 如果你使用官方的接口地址 https://api.openai.com，则留空即可，如果你想指定请求url的地址，可通过这个参数进行配置，注意需要带上 http 协议
      MODEL: "gpt-3.5-turbo" # 指定模型，默认为 gpt-3.5-turbo , 可选参数有： "gpt-4-32k-0613", "gpt-4-32k-0314", "gpt-4-32k", "gpt-4-0613", "gpt-4-0314", "gpt-4", "gpt-3.5-turbo-16k-0613", "gpt-3.5-turbo-16k", "gpt-3.5-turbo-0613", "gpt-3.5-turbo-0301", "gpt-3.5-turbo"，如果使用gpt-4，请确认自己是否有接口调用白名单，如果你是用的是azure，则该配置项可以留空或者直接忽略
      SESSION_TIMEOUT: 600 # 会话超时时间,默认600秒,在会话时间内所有发送给机器人的信息会作为上下文
      MAX_QUESTION_LEN: 2048 # 最大问题长度，默认4096 token，正常情况默认值即可，如果使用gpt4-8k或gpt4-32k，可根据模型token上限修改。
      MAX_ANSWER_LEN: 2048 # 最大回答长度，默认4096 token，正常情况默认值即可，如果使用gpt4-8k或gpt4-32k，可根据模型token上限修改。
      MAX_TEXT: 4096 # 最大文本 = 问题 + 回答, 接口限制，默认4096 token，正常情况默认值即可，如果使用gpt4-8k或gpt4-32k，可根据模型token上限修改。
      HTTP_PROXY: http://host.docker.internal:15777 # 指定请求时使用的代理，如果为空，则不使用代理，注意需要带上 http 协议 或 socks5 协议
      DEFAULT_MODE: "单聊" # 指定默认的对话模式，可根据实际需求进行自定义，如果不设置，默认为单聊，即无上下文关联的对话模式
      MAX_REQUEST: 0 # 单人单日请求次数上限，默认为0，即不限制
      PORT: 8090 # 指定服务启动端口，默认为 8090，容器化部署时，不需要调整，一般在二进制宿主机部署时，遇到端口冲突时使用，如果run_mode为stream模式，则可以忽略该配置项
      SERVICE_URL: ""  # 指定服务的地址，就是当前服务可供外网访问的地址(或者直接理解为你配置在企微回调那里的地址)，用于生成图片时给企微做渲染
      # 以下 ALLOW_USERS、DENY_USERS、VIP_USERS、ADMIN_USERS 配置中填写的是用户的userid
      # 比如 ["1301691029702722","1301691029702733"]，这个信息需要在企微管理后台的通讯录当中获取：hhttps://work.weixin.qq.com/wework_admin/frame#contacts
      # 哪些用户可以进行对话，如果留空，则表示允许所有用户，如果要限制，则列表中写用户的userid
      ALLOW_USERS: "" # 哪些用户可以进行对话，如果留空，则表示允许所有用户，如果要限制，则填写用户的userid
      DENY_USERS: "" # 哪些用户不可以进行对话，如果留空，则表示允许所有用户（如allow_user有配置，需满足相应条件），如果要限制，则列表中写用户的userid，黑名单优先级高于白名单
      VIP_USERS: "" # 哪些用户可以进行无限对话，如果留空，则表示只允许管理员（如max_request配置为0，则允许所有人），如果要针对指定VIP用户放开限制（如max_request配置不为0），则列表中写用户的userid
      ADMIN_USERS: "" # 指定哪些人为此系统的管理员，如果留空，则表示没有人是管理员，如果要限制，则列表中写用户的userid
      SENSITIVE_WORDS: "" # 敏感词，提问时触发，则不允许提问，回答的内容中触发，则以 🚫 代替
      AZURE_ON: "false" # 是否走Azure OpenAi API, 默认false ,如果为true，则需要配置下边的四个参数
      AZURE_API_VERSION: "" # Azure OpenAi API 版本，比如 "2023-03-15-preview"
      AZURE_RESOURCE_NAME: "" # Azure OpenAi API 资源名称，比如 "openai"
      AZURE_DEPLOYMENT_NAME: "" # Azure OpenAi API 部署名称，比如 "openai"
      AZURE_OPENAI_TOKEN: "" # Azure token
      HELP: "Commands:\n\n=================================\n\n🙋 单聊 👉 单独聊天，缺省\n\n🗣 串聊 👉 带上下文聊天\n\n🔃 重置 👉 重置带上下文聊天\n\n🚀 帮助 👉 显示帮助信息\n\n=================================\n\n💪 Power By [eryajf/chatgpt-wecom](https://github.com/eryajf/chatgpt-wecom)" # 帮助信息，放在配置文件，可供自定义
    volumes:
      - ./data:/app/data
    ports:
      - "8000:8000"
    extra_hosts:
      - host.docker.internal:host-gateway
```

启动服务：

```sh
$ docker compose up -d
```

#### 二进制部署


如果你想通过命令行直接部署，可以直接下载release中的[压缩包](https://github.com/eryajf/chatgpt-wecom/releases) ，请根据自己系统以及架构选择合适的压缩包，下载之后直接解压运行。

下载之后，在本地解压，即可看到可执行程序，与配置文件：

```sh
$ tar xf chatgpt-wecom-v0.0.4-darwin-arm64.tar.gz
$ cd chatgpt-wecom-v0.0.4-darwin-arm64
$ cp config.example.yml  config.yml
$ ./chatgpt-wecom  # 直接运行

# 如果要守护在后台运行
$ nohup ./chatgpt-wecom &> run.log &
$ tail -f run.log
```


### 第三步，完善企微配置

此时服务如正常启动，可再回到企微管理后台的配置页面，尝试保存配置。如果服务正常，应该就能保存成功了，从服务日志当中，也可以看到企微发起的回调：

```sh
[GIN] 2023/07/02 - 21:33:53 | 200 |      78.469µs |  113.108.92.100 | GET      "/ai/callback?msg_signature=fb23b8490965c74600dcb08b2e8b86d2aff664e4&timestamp=1688304833&nonce=1688533588&echostr=I5D3M3C%2Fk7AqBFRkACk8eHZAzt%2Fjx14IKk8wXUpA85xsQ2aU67lxEhgHVudLSrCWEPRFapeQ3EcYbni0Bqj01Q%3D%3D"
```

此时还需要将部署服务的IP加入到企微的白名单当中：

<details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_213710.jpg">
    </details>

到这里配置步骤就完成了，可以尽情享用与智能机器人的交互了。

来到企业微信，点击工作台，然后可以看到我们添加的应用，点击应用，即可开始对话。

<details>
      <summary>🖼 点我查看示例图</summary>
      <img src="https://cdn.staticaly.com/gh/eryajf/tu/main/img/image_20230702_213915.jpg">
    </details>

## 配置文件说明

```yaml
# 应用的日志级别，info or debug
log_level: "info"
# openai api_key,如果你是用的是azure，则该配置项可以留空或者直接忽略
api_key: "xxxxxxxxx"
# 如果你使用官方的接口地址 https://api.openai.com，则留空即可，如果你想指定请求url的地址，可通过这个参数进行配置，注意需要带上 http 协议，如果你是用的是azure，则该配置项可以留空或者直接忽略
base_url: ""
# 指定模型，默认为 gpt-3.5-turbo , 可选参数有： "gpt-4-32k-0613", "gpt-4-32k-0314", "gpt-4-32k", "gpt-4-0613", "gpt-4-0314", "gpt-4", "gpt-3.5-turbo-16k-0613", "gpt-3.5-turbo-16k", "gpt-3.5-turbo-0613", "gpt-3.5-turbo-0301", "gpt-3.5-turbo"，如果使用gpt-4，请确认自己是否有接口调用白名单，如果你是用的是azure，则该配置项可以留空或者直接忽略
model: "gpt-3.5-turbo"
# 会话超时时间,默认600秒,在会话时间内所有发送给机器人的信息会作为上下文
session_timeout: 600
# 最大问题长度
max_question_len: 2048
# 最大回答长度
max_answer_len: 2048
# 最大上下文文本长度，通常该参数可设置为与模型Token限制相同
max_text: 4096
# 指定请求时使用的代理，如果为空，则不使用代理，注意需要带上 http 协议 或 socks5 协议，如果你是用的是azure，则该配置项可以留空或者直接忽略
http_proxy: ""
# 指定默认的对话模式，可根据实际需求进行自定义，如果不设置，默认为单聊，即无上下文关联的对话模式
default_mode: "单聊"
# 单人单日请求次数上限，默认为0，即不限制
max_request: 0
# 指定服务启动端口，默认为 8090，一般在二进制宿主机部署时，遇到端口冲突时使用，如果run_mode为stream模式，则可以忽略该配置项
port: "8090"
# 指定服务的地址，就是当前服务可供外网访问的地址(或者直接理解为你配置在企微回调那里的地址)，用于生成图片时给企微做渲染，最新版本中将图片上传到了企微服务器，理论上你可以忽略该配置项，如果run_mode为stream模式，则可以忽略该配置项
service_url: "http://xxxxxx"
# 以下 allow_users、deny_users、vip_users、admin_users 配置中填写的是用户的userid，outgoing机器人模式下不适用这些配置
# 比如 ["1301691029702722","1301691029702733"]，这个信息需要在企微管理后台的通讯录当中获取：https://oa.dingtalk.com/contacts.htm#/contacts
# 哪些用户可以进行对话，如果留空，则表示允许所有用户，如果要限制，则列表中写用户的userid
allow_users: []
# 哪些用户不可以进行对话，如果留空，则表示允许所有用户（如allow_user有配置，需满足相应条件），如果要限制，则列表中写用户的userid，黑名单优先级高于白名单
deny_users: []
# 哪些用户可以进行无限对话，如果留空，则表示只允许管理员（如max_request配置为0，则允许所有人）
# 如果要针对指定VIP用户放开限制（如max_request配置不为0），则列表中写用户的userid
vip_users: []
# 指定哪些人为此系统的管理员，如果留空，则表示没有人是管理员，如果要限制，则列表中写用户的userid
admin_users: []
# 敏感词，提问时触发，则不允许提问，回答的内容中触发，则以 🚫 代替
sensitive_words: []
# 帮助信息，放在配置文件，可供自定义
help: "Commands:\n\n=================================\n\n🙋 单聊 👉 单独聊天，缺省\n\n🗣 串聊 👉 带上下文聊天\n\n🔃 重置 👉 重置带上下文聊天\n\n🚀 帮助 👉 显示帮助信息\n\n=================================\n\n💪 Power By [eryajf/chatgpt-wecom](https://github.com/eryajf/chatgpt-wecom)"

# Azure OpenAI 配置
azure_on: false # 如果是true，则会走azure的openai接口
azure_resource_name: "eryajf" # 对应你的主个性域名
azure_deployment_name: "gpt-35-turbo" # 对应的是 /deployments/ 后边跟着的这个值
azure_api_version: "2023-03-15-preview" # 对应的是请求中的 api-version 后边的值
azure_openai_token: "xxxxxxx"
```

## 常见问题

- 企微只支持与应用一对一聊天，不支持群聊当中添加机器人的对话形式。

## 感谢

这个项目能够成立，离不开这些开源项目：

- [go-resty/resty](https://github.com/go-resty/resty)
- [patrickmn/go-cache](https://github.com/patrickmn/go-cache)
- [solywsh/chatgpt](https://github.com/solywsh/chatgpt)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [avast/retry-go](https://github.com/avast/retry-go)
- [sashabaranov/go-openapi](https://github.com/sashabaranov/go-openai)
- [charmbracelet/log](https://github.com/charmbracelet/log)
- [xen0n/go-workwx](https://github.com/xen0n/go-workwx)

## 赞赏

如果觉得这个项目对你有帮助，你可以请作者[喝杯咖啡 ☕️](https://wiki.eryajf.net/reward/)

## 贡献者列表

<div align="center">
<!-- readme: collaborators,contributors -start -->
<!-- readme: collaborators,contributors -end -->
</div>

