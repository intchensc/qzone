# QZone SDK

[![Go Version](https://img.shields.io/github/go-mod/go-version/intchensc/qzone)](https://github.com/intchensc/qzone)
[![GitHub license](https://img.shields.io/github/license/intchensc/qzone)](https://github.com/intchensc/qzone/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/intchensc/qzone)](https://goreportcard.com/report/github.com/intchensc/qzone)

> 一个强大的 QQ 空间 Go 语言开发工具包，提供简单易用的接口来操作 QQ 空间功能。

## 📁 项目结构

```
qzone/
├── api/                    # API 实现目录
│   ├── common/            # 公共功能
│   ├── friend/            # 好友相关 API
│   ├── group/             # 群组相关 API
│   ├── history/           # 历史记录相关 API
│   ├── shuoshuo/          # 说说相关 API
│   └── api.go             # API 聚合器
├── auth/                   # 认证相关
│   ├── base.go            # 基础认证接口
│   ├── cookie.go          # Cookie 认证实现
│   └── qrcode.go          # 二维码认证实现
├── test/                   # 测试用例
└── qzone.go               # 主入口文件
```

## ✨ 特性

- 🔐 支持扫码登录，安全便捷
- 🚀 模块化的 API 设计，接口清晰
- 📝 完整的说说操作支持
- 👥 好友与群组管理功能
- 🔄 异步操作支持
- 📦 零第三方存储依赖
- 🛡️ 稳定可靠的错误处理机制

## 🚀 快速开始

### 安装

```bash
go get -u github.com/intchensc/qzone
```

### 基础使用示例

```go
package main

import (
    "fmt"
    "github.com/intchensc/qzone"
)

func main() {
    // 创建 QZone 实例
    q := qzone.New(&auth.QrAuth{})
    
    // 登录
    if err := q.Login(); err != nil {
        panic(err)
    }
    
    // 使用说说 API
    shuoshuo := q.API.ShuoShuo()
    // 获取说说列表
    list, err := shuoshuo.List()
    if err != nil {
        panic(err)
    }
    
    // 使用好友 API
    friend := q.API.Friend()
    // 获取好友列表
    friends, err := friend.List()
    if err != nil {
        panic(err)
    }
}
```

## 📚 API 文档

### 认证相关

```go
// 创建实例并使用扫码登录
q := qzone.New(nil)
err := q.Login()

// 使用 Cookie 登录
q := qzone.New(&auth.CookieAuth{Cookie: "your-cookie"})
```

### 说说操作 (ShuoShuoAPI)

```go
api := q.API.ShuoShuo()

// 获取说说列表
list, err := api.List()

// 发布说说
err := api.Publish(content, images)

// 获取说说评论
comments, err := api.Comments(tid)
```

### 好友操作 (FriendAPI)

```go
api := q.API.Friend()

// 获取好友列表
list, err := api.List()

// 获取好友详情
detail, err := api.Detail(uin)
```

### 群组操作 (GroupAPI)

```go
api := q.API.Group()

// 获取群组列表
list, err := api.List()

// 获取群成员
members, err := api.Members(groupId)
```

### 历史记录 (HistoryAPI)

```go
api := q.API.History()

// 获取历史记录
history, err := api.Get()
```

## 🔧 高级配置

### 自定义认证实现

你可以通过实现 `auth.BaseAuth` 接口来创建自己的认证方式：

```go
type BaseAuth interface {
    Login() error
    GetCookie() string
}
```

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建您的特性分支 (git checkout -b feature/AmazingFeature)
3. 提交您的更改 (git commit -m 'Add some AmazingFeature')
4. 推送到分支 (git push origin feature/AmazingFeature)
5. 打开一个 Pull Request

## 📝 开发计划

- [x] 基础接口封装
- [x] 扫码登录
- [ ] 规范接口返回字段
- [ ] 接口的统一分页设计
- [ ] 便捷功能封装
- [ ] "与我相关"推送接口
- [ ] 多种登录方式支持



## 🌟 Star History

如果这个项目对您有帮助，请给我们一个 star！您的支持是我们持续改进的动力。

---

> 📢 注意：本项目仍在积极开发中，API 可能会有重大变更。建议在生产环境使用前关注版本更新。