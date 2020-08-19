# API
Common return codes.

|        message         |         结果          |
|:----------------------:|:--------------------:|
|        SUCCESS         |         正常          |
|     INTERNAL_ERROR     |     服务器内部错误      |
|    VALIDATION_ERROR    | 验证错误, 详见error字段 |


# Authentication

|         code         |       结果       |
|:--------------------:|:---------------:|
| AUTH_TOKEN_NOT_FOUND |   未找到token    |
| AUTH_SESSION_EXPIRED |   session超时    |
|   AUTH_NEED_TOKEN    |   未提供token    |


## Auth

### Login

|         code       |         结果          |
|:------------------:|:--------------------:|
|   NOT_FOUND   |    错误的用户名/邮箱    |
|   WRONG_PASSWORD   |       密码错误        |
### Register

|            code            |         结果          |
|:--------------------------:|:--------------------:|
|      DUPLICATE_EMAIL       |       邮箱重复        |
|     DUPLICATE_USERNAME     |       用户名重复       |

### EmailRegistered

|            code            |         结果          |
|:--------------------------:|:--------------------:|
|      EMAIL_REGISTERED      |       邮箱已注册       |

## Admin

### User

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|     DUPLICATE_EMAIL     |        邮箱重复        |
|    DUPLICATE_USERNAME   |       用户名重复       |
