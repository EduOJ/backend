# EduOJ通知渠道开发手册

- ## 数据库

  1. user表设计

     通知系统为通知渠道开发者预留了`user`表中两列内容，分别是：

     ```go
     	PreferredNoticeMethod string `gorm:"preferred_notice_method"`
     	NoticeAccount string `gorm:"notice_account"`
     ```

     这两列分别记录了用户的通知渠道偏好和通知渠道地址。

     

  2. 修改用户偏好

     可以向后端系统发送路由为

     ```go
     	user.PUT("/user/me", controller.UpdateMe).Name = "user.updateMe"
     ```

     调用函数`updateMe`实现用户偏好更改

     函数内关于通知部分的代码如下

     ```go
     user.PreferredNoticeMethod = req.PreferredNoticeMethod
     	type Noticejson struct {
     		PreferredNoticeMethod string
     		NoticeAccount string
     	}
     	noticejson := Noticejson{
     		PreferredNoticeMethod: req.PreferredNoticeMethod,
     		NoticeAccount: req.NoticeAccount,
     	}
     	noticejson_byte, err := json.Marshal(noticejson)
     	if err != nil {
     		println("could not creat json")
     	}
     	user.NoticeAccount = string(noticejson_byte)
     ```

     该部分可见，前端传来的`PreferredNoticeMethod`会被存入对应的数据库中`user`表中`PreferredNoticeMethod`一列，而前端传来的`PreferredNoticeMethod`和`NoticeAccount`会被序列化为`json`字符串存入`user`表中`NoticeAccount`列

     开发者在得到`user`的信息后，若想知道具体的渠道偏好与通知地址，可以调用`json.UnMarshal`将其反序列化解析

- ## 注册功能模块

  EduOJ为通知渠道开发者预留了功能函数`func Register(name string)`

  函数接收一个字符串作为新的功能模块名

  使用者应该在启用EduOJ后借由其加载需要用的通知渠道

- ## 注册路由

  EduOJ提供了一个`Event`系统，为通知模块开发者提供加载其功能的服务。

  基于`Event`模块，EduOJ的通知系统在`router`中为通知模块开发者预留了注册新路由的服务。

  ```go
  _, _ = event.FireEvent("register_router", e)
  ```

  而通知渠道开发者只需要在启动自己通知渠道时，注册一个关于该`register_router`的监听器，追加私有的新路由。

  使用范例：

  ​	功能：balabal

  ​	启动时

  ```go
  func init() {
  	event.RegisterListener("register_rouer", func(e *echo.Echo) {
  		e.POST("/balabal", ...)
  	})
  }
  ```

- ## 发送消息

  EduOJ通知系统包含一个广义兼容的发送通知的请求，该方法会调用通知渠道开发人员开发的具体通知方法对用户进行通知。

  该方法参数包括但不限于：通知对象的具体属性，通知内容，通知标题，

  该方法会触发接收人通知偏好对应的一个通知渠道，通过`fireevent`调用通知接口，实现具体通知功能。

  

  请注意，调用通知接口时，会使用事件名称作为`fireevent`对象，所以通知渠道开发者开发的发送通知的函数监听器名称，应当遵守命名规范。

  命名规范应当遵守：`%s_send_message`, 其中`%s`代表具体通知渠道名称，应当与数据库字段`PreferredNoticeMethod`相一致。

  具体范例如下：

  ```go
  email_send_message
  ```

  

  该函数返回值解析：

  ​	发送失败的原因有以下两点：

   1. 触发具体接口失败，触发系统panic
      可能的原因包括但不限于：

       	1. 开发者启用自己模块功能时忘记注册对应名称监听器
       	2. 开发者注册监听器的路由错误
       	3. 开发者监听器名称不合法或与上述范例命名不匹配
       	4. ......

   2. 通知失败

       1. 该处失败代表成功的调用了具体通知渠道，但是信息发送不成功

       2. 可能的原因包括但不限于：

           1. 管理员通知账号被封

              例如：

              ​	管理员邮箱账户未开启smtp服务

              ​	管理员短信服务欠费

             	2. 收件人通知地址错误
           	
             	3. ......

- ## 通知系统其他功能

  以上功能是开发者实现一个具体通知最基本需要掌握的方法

  EduOJ通知系统也为通知渠道开发者提供其他通知函数方便开发使用，请阅读EduOJ通知系统设计方案学习使用。