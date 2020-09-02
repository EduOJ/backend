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
| AUTH_SESSION_EXPIRED |   session超时    |
|   AUTH_NEED_TOKEN    |   未提供token    |

# Permission

|            code            |         结果          |
|:--------------------------:|:--------------------:|
|     PERMISSION_DENIED      |        没有权限        |

## Auth

### Login

|         code       |         结果          |
|:------------------:|:--------------------:|
|   WRONG_USERNAME   |    错误的用户名/邮箱    |
|   WRONG_PASSWORD   |       密码错误        |

### Register

|            code            |         结果          |
|:--------------------------:|:--------------------:|
|      CONFLICT_EMAIL       |       邮箱重复        |
|     CONFLICT_USERNAME     |       用户名重复       |

### EmailRegistered

|            code            |         结果          |
|:--------------------------:|:--------------------:|
|      EMAIL_REGISTERED      |       邮箱已注册       |

## Admin

### User

#### AdminCreateUser

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |
|   PERMISSION_DENIED    |        没有权限        |

#### AdminUpdateUser

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |
|    PERMISSION_DENIED   |        没有权限        |

#### AdminDeleteUser

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|    PERMISSION_DENIED    |       没有权限        |

#### AdminGetUser

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminGetUsers

|           code          |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |
|    PERMISSION_DENIED    |        没有权限       |

### Problem

#### AdminCreateProblem
|           code          |         结果          |
|:-----------------------:|:--------------------:|
|      CONFLICT_NAME      |       名称重复         |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminGetProblem
|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |   无法找到指定problem  |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminGetProblems
|           code          |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminUpdateProblem
|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |   无法找到指定problem  |
|      CONFLICT_NAME      |       名称重复         |
|    PERMISSION_DENIED    |       没有权限         |

#### AdminDeleteProblem
|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |   无法找到指定problem  |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminCreateTestCase
#### AdminGetTestCase
#### AdminGetTestCases
#### AdminUpdateTestCase
#### AdminDeleteTestCase

## User

### GetMe

### UpdateMe

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |

### GetUser

|           code          |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |

### GetUsers

|           code          |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |

### ChangePassword

|         code         |         结果          |
|:--------------------:|:--------------------:|
|    WRONG_PASSWORD    |       密码错误        |

## Problem

### GetProblem
### GetProblems

### GetTestCase
### GetTestCases
## Image
### CreateImage
|     code     |  结果   |
|:------------:|:------:|
| ILLEGAL_TYPE | 类型非法 |
