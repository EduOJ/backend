# API
Common return codes.

|        message         |         结果          |
|:----------------------:|:--------------------:|
|        SUCCESS         |         正常          |
|     INTERNAL_ERROR     |     服务器内部错误      |
|    VALIDATION_ERROR    | 验证错误, 详见error字段 |
|        NOT_FOUND       |      找不到指定对象     |

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

### Log

#### AdminGetLogs

|   message     |    结果     |
|:-------------:|:----------:|
| INVALID_LEVEL | 非法的level |
### User

#### AdminCreateUser

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |

#### AdminUpdateUser

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复        |
|    CONFLICT_USERNAME   |       用户名重复       |

#### AdminDeleteUser

#### AdminGetUser

#### AdminGetUsers

|         message         |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |

## User

### GetMe

### UpdateMe

|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     CONFLICT_EMAIL     |        邮箱重复         |
|    CONFLICT_USERNAME   |       用户名重复        |
|    METHOD_NOT_FOUND    |     用户偏好设置不存在    |
|    INVALID_ACCOUNT     |       非法的账户信息     |

### GetUser

### GetUsers

|         message         |         结果         |
|:-----------------------:|:-------------------:|
|      INVALID_ORDER      |     无效的排序设置     |

### ChangePassword

|       message        |         结果          |
|:--------------------:|:--------------------:|
|    WRONG_PASSWORD    |       密码错误        |

## Problem

### CreateProblem

### GetProblem

### GetRandomProblem

### GetProblemAttachmentFile

### GetProblems
|         message         |         结果         |
|:-----------------------:|:-------------------:|
|     INVALID_STATUS      |     无效的状态设置     |

### UpdateProblem

### DeleteProblem

### CreateTestCase
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|      INVALID_FILE       |        缺少文件       |

### UpdateTestCase

### DeleteTestCase

### DeleteTestCases

### GetTestCaseInputFile

### GetTestCaseOutputFile

## Image
### CreateImage
|     code     |  结果   |
|:------------:|:------:|
| ILLEGAL_TYPE | 类型非法 |

## Submission

### CreateSubmission
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     INVALID_LANGUAGE    |       无效的语言       |
|       INVALID_FILE      |        缺少文件       |

### GetSubmission

### GetSubmissions

### GetSubmissionCode

### GetRunOutput
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|  SUBMISSION_NOT_FOUND   |   无法找到submission   |
|   JUDGEMENT_UNFINISHED  |       评测未完成       |

### GetRunCompilerOutput
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|  SUBMISSION_NOT_FOUND   |   无法找到submission   |
|   JUDGEMENT_UNFINISHED  |       评测未完成       |

### GetRunComparerOutput
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|  SUBMISSION_NOT_FOUND   |   无法找到submission   |
|   JUDGEMENT_UNFINISHED  |       评测未完成       |

# Judger
## UpdateRun

|      message      |                  结果                  |
|:-----------------:|:-------------------------------------:|
|   WRONG_RUN_ID    | 发起请求的judger与获取道当前run的judger不同 |
| ALREADY_SUBMITTED |          一个run被提交了两次结果          |

## Class

### CreateClass

### GetClass

### GetClassesIManage

### GetClassesITake

### UpdateClass

### RefreshInviteCode

### AddStudents

### DeleteStudents

### JoinClass
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|    WRONG_INVITE_CODE    |      错误的邀请码       |
|    ALREADY_IN_CLASS     |   用户已是该class学生   |

### DeleteClass

## ProblemSet

### CreateProblemSet

### CloneProblemSet

### GetProblemSet

### UpdateProblemSet

### AddProblemsToSet

### DeleteProblemsFromSet

### DeleteProblemSet

### GetProblemSetProblem

### GetProblemSetProblemInputFile

### GetProblemSetProblemOutputFile

### RefreshGrades

## ProblemSetSubmission

### ProblemSetCreateSubmission
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|     INVALID_LANGUAGE    |       无效的语言       |
|       INVALID_FILE      |        缺少文件       |

### ProblemSetGetSubmission

### ProblemSetGetSubmissions

### ProblemSetGetSubmissionCode

### ProblemSetGetRunOutput
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|   JUDGEMENT_UNFINISHED  |       评测未完成       |

### ProblemSetGetRunInput

### ProblemSetGetRunCompilerOutput
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|   JUDGEMENT_UNFINISHED  |       评测未完成       |

### ProblemSetGetRunComparerOutput
|         message         |         结果          |
|:-----------------------:|:--------------------:|
|   JUDGEMENT_UNFINISHED  |       评测未完成       |
