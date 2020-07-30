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
| AUTH_TOKEN_NOT_FOUND | Token not found |
| AUTH_SESSION_EXPIRED |   session超时    |
|   AUTH_NEED_TOKEN    |   未提供token    |


## Auth

### Login

|         code         |         结果          |
|:--------------------:|:--------------------:|
| LOGIN_WRONG_USERNAME |    错误的用户名/邮箱    |
| LOGIN_WRONG_PASSWORD |       密码错误        |
### Register

|            code             |         结果          |
|:---------------------------:|:--------------------:|
|  REGISTER_DUPLICATE_EMAIL   |       邮箱重复        |
| REGISTER_DUPLICATE_USERNAME |       用户名重复       |
