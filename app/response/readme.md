# API
Common return codes.

|        message         |         结果          |
|:----------------------:|:--------------------:|
|        SUCCESS         |         正常          |
|     INTERNAL_ERROR     |     服务器内部错误      |
|    VALIDATION_ERROR    | 验证错误, 详见error字段 |


# Authentication

|       message        |       结果       |
|:--------------------:|:---------------:|
| AUTH_SESSION_EXPIRED |   session超时    |
|   AUTH_NEED_TOKEN    |   未提供token    |

# Permission

|          message           |         结果          |
|:--------------------------:|:--------------------:|
|     PERMISSION_DENIED      |        没有权限        |

## Auth

### Login

|       message      |         结果          |
|:------------------:|:--------------------:|
|   WRONG_USERNAME   |    错误的用户名/邮箱    |
|   WRONG_PASSWORD   |       密码错误        |

### Register

|          message           |         结果          |
|:--------------------------:|:--------------------:|
|      CONFLICT_EMAIL       |       邮箱重复        |
|     CONFLICT_USERNAME     |       用户名重复       |

### EmailRegistered

|          message           |         结果          |
|:--------------------------:|:--------------------:|
|      EMAIL_REGISTERED      |       邮箱已注册       |

## Admin

### User

#### AdminCreateUser

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |
|   PERMISSION_DENIED    |        没有权限        |

#### AdminUpdateUser

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |
|    PERMISSION_DENIED   |        没有权限        |

#### AdminDeleteUser

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|    PERMISSION_DENIED    |       没有权限        |

#### AdminGetUser

|         message        |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminGetUsers

|         message         |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |
|    PERMISSION_DENIED    |        没有权限       |

### Problem

#### AdminCreateProblem
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|    PERMISSION_DENIED    |        没有权限       |

#### AdminGetProblem
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |   无法找到指定problem  |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminGetProblems
|         message         |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |
|    PERMISSION_DENIED    |        没有权限       |

#### AdminUpdateProblem
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |   无法找到指定problem  |
|    PERMISSION_DENIED    |       没有权限         |

#### AdminDeleteProblem
|         message         |         结果          |
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

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |

### GetUser

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|        NOT_FOUND        |     无法找到指定user   |

### GetUsers

|         message         |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |

### ChangePassword

|       message        |         结果          |
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
