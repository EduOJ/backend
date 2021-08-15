# EduOJ通知系统设计方案

- ## 数据库设计

  对user表新增加两列

  原有数据库设计：

  ```go
  type User struct {
  	ID       uint   `gorm:"primaryKey" json:"id"`
  	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
  	Nickname string `gorm:"index:nickname" json:"nickname"`
  	Email    string `gorm:"unique_index" json:"email"`
  	Password string `json:"-"`
  
  	Roles      []UserHasRole `json:"roles"`
  	RoleLoaded bool          `gorm:"-" json:"-"`
  
  	ClassesManaging []*Class `json:"class_managing" gorm:"many2many:user_manage_classes"`
  	ClassesTaking   []*Class `json:"class_taking" gorm:"many2many:user_in_classes"`
  
  	Grades []*Grade `json:"grades"`
  
  	CreatedAt time.Time      `json:"created_at"`
  	UpdatedAt time.Time      `json:"-"`
  	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
  	// TODO: bio
  
  	Credentials []WebauthnCredential
  }
  ```

  修改后的新表：

  ```go
  type User struct {
  	ID       uint   `gorm:"primaryKey" json:"id"`
  	Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5,username"`
  	Nickname string `gorm:"index:nickname" json:"nickname"`
  	Email    string `gorm:"unique_index" json:"email"`
  	Password string `json:"-"`
  	PreferedNoticeMethod string `gorm:"prefered_notice_method"`
  	NoticeAddress string `gorm:"notice_address"`
  
  	Roles      []UserHasRole `json:"roles"`
  	RoleLoaded bool          `gorm:"-" json:"-"`
  
  	ClassesManaging []*Class `json:"class_managing" gorm:"many2many:user_manage_classes"`
  	ClassesTaking   []*Class `json:"class_taking" gorm:"many2many:user_in_classes"`
  
  	Grades []*Grade `json:"grades"`
  
  	CreatedAt time.Time      `json:"created_at"`
  	UpdatedAt time.Time      `json:"-"`
  	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
  	// TODO: bio
  
  	Credentials []WebauthnCredential
  }
  ```

  通知系统为通知渠道开发者预留了`user`表中两列内容，分别是：

  ```go
  	PreferedNoticeMethod string `gorm:"prefered_notice_method"`
  	NoticeAddress string `gorm:"notice_address"`
  ```

  这两列分别记录了用户的通知渠道偏好和通知渠道地址。

- ## 更新通知偏好

  修改了`UpdateMe`函数实现更新用户的通知属性功能。

  ```go
  func UpdateMe(c echo.Context) error {
  	user, ok := c.Get("user").(models.User)
  	if !ok {
  		panic("could not convert my user into type models.User")
  	}
  	if !user.RoleLoaded {
  		user.LoadRoles()
  	}
  	req := request.UpdateMeRequest{}
  	err, ok := utils.BindAndValidate(&req, c)
  	if !ok {
  		return err
  	}
  	count := int64(0)
  	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
  	if count > 1 || (count == 1 && user.Email != req.Email) {
  		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_EMAIL", nil))
  	}
  	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
  	if count > 1 || (count == 1 && user.Username != req.Username) {
  		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_USERNAME", nil))
  	}
  	user.Username = req.Username
  	user.Nickname = req.Nickname
  	user.Email = req.Email
  	user.PreferedNoticeMethod = req.PreferedNoticeMethod
  	type Noticejson struct {
  		PreferedNoticeMethod string
  		NoticeAddress string
  	}
  	noticejson := Noticejson{
  		PreferedNoticeMethod: req.PreferedNoticeMethod,
  		NoticeAddress: req.NoticeAddress,
  	}
  	noticejson_byte, err := json.Marshal(noticejson)
  	if err != nil {
  		println("could not creat json")
  	}
  	user.NoticeAddress = string(noticejson_byte)
  	utils.PanicIfDBError(base.DB.Omit(clause.Associations).Save(&user), "could not update user")
  	return c.JSON(http.StatusOK, response.UpdateMeResponse{
  		Message: "SUCCESS",
  		Error:   nil,
  		Data: struct {
  			*resource.UserForAdmin `json:"user"`
  		}{
  			resource.GetUserForAdmin(&user),
  		},
  	})
  }
  ```

  其中

  ```go
  user.PreferedNoticeMethod = req.PreferedNoticeMethod
  	type Noticejson struct {
  		PreferedNoticeMethod string
  		NoticeAddress string
  	}
  	noticejson := Noticejson{
  		PreferedNoticeMethod: req.PreferedNoticeMethod,
  		NoticeAddress: req.NoticeAddress,
  	}
  	noticejson_byte, err := json.Marshal(noticejson)
  	if err != nil {
  		println("could not creat json")
  	}
  	user.NoticeAddress = string(noticejson_byte)
  ```

  该部分可见，前端传来的`PreferedNoticeMethod`会被存入对应的数据库中`user`表中`PreferedNoticeMethod`一列，而前端传来的`PreferedNoticeMethod`和`NoticeAddress`会被序列化为`json`字符串存入`user`表中`NoticeAddress`列

  开发者在得到`user`的信息后，若想知道具体的渠道偏好与通知地址，可以调用`json.UnMarshal`将其反序列化解析

  

- ## 注册新的通知功能

1. 设计思路

   - 在`notification`包下的全局变量`RegistedPreferedNoticeMethod`用于记录已经注册启用的通知模块

   - 在注册时应该校验传入数据，防止类似同名方式出现

     ```go
     func Register(name string) {
     	RegistedPreferedNoticedMethod = append(RegistedPreferedNoticedMethod, name)
     	//..
     }
     ```

- ## 展示各个通知方式的使用情况

1. 设计思路

   - 遍历数据库表`user`，记录每个用户的`PreferedNoticeMethod`，将该json字段使用`json.Unmarshal`解析存入一个结构

   - 输出通知方式使用用户数量

     ```
     func ShowUsedMethod(db base.DB,) slice,error{
     
     }
     ```

     

- ## 删除禁用已经注册过的通知方式

1. 设计思路

   - 删除`RegistedPreferedNoticeMethod`中部分字段

   - 对整个`user`表进行遍历，修改使用被删通知渠道的用户的通知渠道为默认渠道

   - 

     ```go
     func DeleteRegistedMethod() {
     
     }
     ```

     全局变量，无需传参

- ## 得到现有已经注册过的通知方法个数

1. 设计思路

   ```GO
   len(ShowUsedMethod())
   ```


- ## 发送消息

  1. 设计思路

     - 参数设计：接收者，发送标题，发送内容，...

     - 查询接收者的`PreferedNoticeMethod`

     - 解析接收者的`NoticeDetail`得到收件地址

     - 触发新的事件：`XXX_send_message`，XXX代表通知渠道，事件应该由渠道开发者书写名称，如果名称不匹配，会造成FireEvent的panic，错误信息返回在err当中，向渠道开发者返回触发事件失败log，帮助其调整修复代码。

     - 若事件触发成功，但由于一些原因，包含但不限于：邮箱被封、邮箱空号、网路问题。返回的报错信息会在函数返回值中呈现，应该对其解析理解，向系统使用者报出具体的失败原因。

       ```go
       func SendMessage(receiver *models.User, title string, message string) {
       	method := receiver.PreferedNoticeMethod
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

       

