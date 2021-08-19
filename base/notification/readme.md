# 通知系统开发文档

## 通知渠道开发
* 初始化工作：写在init()函数中
* 发送消息时的工作：写在符合签名`func(receiver *modules.User, title, message string) error`的发送消息函数中

## 通知渠道启用
* EduOJ启动时调用`notification.Register(渠道名称, 发送消息函数)`

## 发送通知
* 调用`notification.SendMessage`函数
