# EduOJ通知系统设计

## 需求清单

1. 用户修改通知属性

   1. 设计思路

      - 修改model/user，增加数据库中新的字段

      - `NoticeDetail` 存入一个`json`字符串记录用户的通知信息

        ```go
        PreferredNoticeMethod string `gorm:"preferred_notice_method"`
        NoticeDetail string
        ```

      - 修改request/user中结构体 `updateme`,添加两个string字段，分别为`PreferredNoticeMethod` 和 `NoticeAccount`

        ```go
        type UpdateMeRequest struct {
        	Username string `json:"username" form:"username" query:"username" validate:"required,max=30,min=5,username"`
        	Nickname string `json:"nickname" form:"nickname" query:"nickname" validate:"required,max=30,min=1"`
        	Email    string `json:"email" form:"email" query:"email" validate:"required,email,max=320,min=5"`
        	PreferredNoticeMethod string `json:"preferrednoticemethod" form:"preferrednoticemethod" query:"preferrednoticemethod"`
        	NoticeAccount string `json:"noticeaccount" form:"noticeaccount" query:"noticeaccount"`
        }
        ```

      - `PreferredNoticeMethod`和`NoticeAccount`将会以`json`字符串的格式存入数据库`user`表中`PreferredNoticeMethod`列，需要用到 `json.Marshal`将结构体序列化，取用时`json.UnMarshal`将其解析。

      - 以下代码中`NoticeDetail`仅包含一部分信息，应该修改代码存入更长的json

        ```go
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
        	user.PreferredNoticeMethod = req.PreferredNoticeMethod
        user.NoticeDetail = string(notice_byte)
        ```

      - 当然这里的设计是不足的，应该在修改属性存入数据库前对于`Method`进行校验，系统进行检测是否是已经注册过的启用了的通知功能，否则返回报错信息

2. 注册新的通知功能

   1. 设计思路

      - 在`notification`包下的全局变量`RegistedPreferredNoticeMethod`用于记录已经注册启用的通知模块

      - 在注册时应该校验传入数据，防止类似同名方式出现

        ```go
        func register(name string) {
        	RegistedPreferredNoticedMethod = append(RegistedPreferredNoticedMethod, name)
        	//..
        }
        ```

        

3. 展示各个通知方式的使用情况

   1. 设计思路
      - 遍历数据库表`user`，记录每个用户的`PreferredNoticeMethod`，将该json字段使用`json.Unmarshal`解析并统计，实现统计各个方式的具体使用情况
      - 输出通知方式使用用户数量

4. 初始化路由

   1. 设计思路

      - event模块通过go `reflect`来实现事务功能

      - 添加路由：调用函数`event.FireEvent`在`router`中触发“regist_router"事件

        2. 使用范例：使用时在init时应该调用`RegistListener`注册"register_router"

        ```go
        func init() {
        	event.RegisterListener("register_rouer", func(e *echo.Echo) {
        		e.POST("/bind_sms", ...)
        	})
        }
        ```

        

9. 发送消息

   1. 设计思路

      - 参数设计：接收者，发送标题，发送内容，...

      - 查询接收者的`PreferredNoticeMethod`

      - 解析接收者的`NoticeDetail`得到收件地址

      - 触发新的事件：`XXX_send_message`，XXX代表通知渠道，事件应该由渠道开发者书写名称并且**注册**（调用函数`registlistener`)，如果名称不匹配，会造成FireEvent的panic，错误信息返回在err当中，向渠道开发者返回触发事件失败log，帮助其调整修复代码。

      - 若事件触发成功，但由于一些原因，包含但不限于：邮箱被封、邮箱空号、网路问题。返回的报错信息会在函数返回值中呈现，应该对其解析理解，向系统使用者报出具体的失败原因。

        ```go
        func SendMessage(receiver *models.User, title string, message string) {
        	method := receiver.PreferredNoticeMethod
        	result, err := event.FireEvent(fmt.Sprintf("%s_send_message", method), receiver, title, message)
        	if err != nil {
        		//panic
        		//事件不存在？
        	}
        	if mErr := result[0][0].(error); mErr != nil {
        		//手机号不存在？
        	}
        }
        ```

        

