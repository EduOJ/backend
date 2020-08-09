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

|         code         |         结果          |
|:--------------------:|:--------------------:|
| LOGIN_WRONG_USERNAME |    错误的用户名/邮箱    |
| LOGIN_WRONG_PASSWORD |       密码错误        |
### Register

|            code             |         结果          |
|:---------------------------:|:--------------------:|
|  REGISTER_DUPLICATE_EMAIL   |       邮箱重复        |
| REGISTER_DUPLICATE_USERNAME |       用户名重复       |

## Admin

### User

|         code                 |         结果          |
|:----------------------------:|:--------------------:|
|     QUERY_USER_WRONG_ID      |      错误的ID/用户名   |

#### PostUser

|            code              |         结果          |
|:----------------------------:|:--------------------:|
|  POST_USER_DUPLICATE_EMAIL   |       邮箱重复        |
| POST_USER_DUPLICATE_USERNAME |       用户名重复       |

#### PutUser

|         code                 |         结果          |
|:----------------------------:|:--------------------:|
|  PUT_USER_DUPLICATE_EMAIL    |       邮箱重复         |
|  PUT_USER_DUPLICATE_USERNAME |       用户名重复       |

#### GetUsers

|            code                 |         结果          |
|:-------------------------------:|:--------------------:|
| GET_USERS_OFFSET_OUT_OF_BOUNDS  |        偏移量越界      |

## User

### ChangePassword

|         code         |         结果          |
|:--------------------:|:--------------------:|
|    WRONG_PASSWORD    |       密码错误        |