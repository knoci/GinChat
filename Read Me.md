# Read Me

## 项目简介

​	基于Gin框架的IM即时通讯小demo，实现了用户注册，添加好友，创建和加入群聊；好友聊天，群聊天；可以自定义个人信息，群信息；在聊天中可以发送文字、图片、表情、语音。

## 技术使用

​	前端：Vue

​	后端：Gin、GORM、Mysql、Redis、Swagger、Viper、Websocket、govalidator、MD5

## 运行

​	在config/app.yml中配置好mysql和redis参数，以及运行的端口和通信的服务端口。

​	在默认情况下，访问127.0.0.1/:8848即可进入