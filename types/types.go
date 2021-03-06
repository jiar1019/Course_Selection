package types

import "gorm.io/gorm"

// 说明：
// 1. 所提到的「位数」均以字节长度为准
// 2. 所有的 ID 均为 int64（以 string 方式表现）

// 通用结构

type ErrNo int

const (
	OK                 ErrNo = 0
	ParamInvalid       ErrNo = 1  // 参数不合法
	UserHasExisted     ErrNo = 2  // 该 Username 已存在
	UserHasDeleted     ErrNo = 3  // 用户已删除
	UserNotExisted     ErrNo = 4  // 用户不存在
	WrongPassword      ErrNo = 5  // 密码错误
	LoginRequired      ErrNo = 6  // 用户未登录
	CourseNotAvailable ErrNo = 7  // 课程已满
	CourseHasBound     ErrNo = 8  // 课程已绑定过
	CourseNotBind      ErrNo = 9  // 课程未绑定过
	PermDenied         ErrNo = 10 // 没有操作权限
	StudentNotExisted  ErrNo = 11 // 学生不存在
	CourseNotExisted   ErrNo = 12 // 课程不存在
	StudentHasNoCourse ErrNo = 13 // 学生没有课程
	StudentHasCourse   ErrNo = 14 // 学生有课程

	UnknownError ErrNo = 255 // 未知错误
)

type ResponseMeta struct {
	Code ErrNo
}

type Member struct {
	Deleted  gorm.DeletedAt
	UserID   string   `gorm:"primaryKey;type:bigint UNSIGNED not null AUTO_INCREMENT"`
	Nickname string   `gorm:"type:varchar(32) not null"`
	Username string   `gorm:"type:varchar(32) not null"`
	UserType UserType `gorm:"type:int not null"`
	Password string   `gorm:"type:varchar(32) not null"`
}

type Course struct {
	CourseID  string `gorm:"primaryKey;type:bigint UNSIGNED not null AUTO_INCREMENT"`
	Name      string `gorm:"type:varchar(32) not null"`
	Cap       int    `gorm:"type:int not null"`
	TeacherID string `gorm:"type:varchar(32)"`
}

type SCourse struct {
	CourseID string `gorm:"type:bigint UNSIGNED not null"`
	UserID   string `gorm:"type:bigint UNSIGNED not null"`
}

type TMember struct {
	UserID   string
	Nickname string
	Username string
	UserType UserType
}

type TCourse struct {
	CourseID  string
	Name      string
	TeacherID string
}

// -----------------------------------

// 成员管理

type UserType int

const (
	Admin   UserType = 1
	Student UserType = 2
	Teacher UserType = 3
)

// 系统内置管理员账号
// 账号名：JudgeAdmin 密码：JudgePassword2022

// 创建成员
// 参数不合法返回 ParamInvalid

// 只有管理员才能添加

type CreateMemberRequest struct {
	Nickname string   `form:"Nickname" json:"Nickname" xml:"Nickname"  binding:"required,min=4,max=20"`          // required，不小于 4 位 不超过 20 位
	Username string   `form:"Username" json:"Username" xml:"Username"  binding:"required,alpha,min=8,max=20"`    // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string   `form:"Password" json:"Password" xml:"Password"  binding:"required,alphanum,min=8,max=20"` // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType UserType `form:"UserType" json:"UserType" xml:"UserType"  binding:"required,min=1,max=3"`           // required, 枚举值
}

type CreateMemberResponse struct {
	Code ErrNo
	Data struct {
		UserID string // int64 范围
	}
}

// 获取成员信息

type GetMemberRequest struct {
	UserID string `form:"UserID" json:"UserID" xml:"UserID"  binding:"required"`
}

// 如果用户已删除请返回已删除状态码，不存在请返回不存在状态码

type GetMemberResponse struct {
	Code ErrNo
	Data TMember
}

// 批量获取成员信息

type GetMemberListRequest struct {
	Offset int `form:"Offset" json:"Offset" xml:"Offset"  binding:"required"`
	Limit  int `form:"Limit" json:"Limit" xml:"Limit"  binding:"required"`
}

type GetMemberListResponse struct {
	Code ErrNo
	Data struct {
		MemberList []TMember
	}
}

// 更新成员信息

type UpdateMemberRequest struct {
	UserID   string `form:"UserID" json:"UserID" xml:"UserID"  binding:"required"`
	Nickname string `form:"Nickname" json:"Nickname" xml:"Nickname"  binding:"required"`
}

type UpdateMemberResponse struct {
	Code ErrNo
}

// 删除成员信息
// 成员删除后，该成员不能够被登录且不应该不可见，ID 不可复用

type DeleteMemberRequest struct {
	UserID string `form:"UserID" json:"UserID" xml:"UserID"  binding:"required"`
}

type DeleteMemberResponse struct {
	Code ErrNo
}

// ----------------------------------------
// 登录

type LoginRequest struct {
	Username string `form:"Username" json:"Username" xml:"Username"  binding:"required"`
	Password string `form:"Password" json:"Password" xml:"Password"  binding:"required"`
}

// 登录成功后需要 Set-Cookie("camp-session", ${value})
// 密码错误范围密码错误状态码

type LoginResponse struct {
	Code ErrNo
	Data struct {
		UserID string
	}
}

// 登出

type LogoutRequest struct{}

// 登出成功需要删除 Cookie

type LogoutResponse struct {
	Code ErrNo
}

// WhoAmI 接口，用来测试是否登录成功，只有此接口需要带上 Cookie

type WhoAmIRequest struct {
}

// 用户未登录请返回用户未登录状态码

type WhoAmIResponse struct {
	Code ErrNo
	Data TMember
}

// -------------------------------------
// 排课

// 创建课程
// Method: Post
type CreateCourseRequest struct {
	Name string `form:"Name" json:"Name" xml:"Name"  binding:"required"`
	Cap  int    `form:"Cap" json:"Cap" xml:"Cap"  binding:"required"`
}

type CreateCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseID string
	}
}

// 获取课程
// Method: Get
type GetCourseRequest struct {
	CourseID string `form:"CourseID" json:"CourseID" xml:"CourseID"  binding:"required"`
}

type GetCourseResponse struct {
	Code ErrNo
	Data TCourse
}

// 老师绑定课程
// Method： Post
// 注：这里的 teacherID 不需要做已落库校验
// 一个老师可以绑定多个课程 , 不过，一个课程只能绑定在一个老师下面
type BindCourseRequest struct {
	CourseID  string `form:"CourseID" json:"CourseID" xml:"CourseID"  binding:"required"`
	TeacherID string `form:"TeacherID" json:"TeacherID" xml:"TeacherID"  binding:"required"`
}

type BindCourseResponse struct {
	Code ErrNo
}

// 老师解绑课程
// Method： Post
type UnbindCourseRequest struct {
	CourseID  string `form:"CourseID" json:"CourseID" xml:"CourseID"  binding:"required"`
	TeacherID string `form:"TeacherID" json:"TeacherID" xml:"TeacherID"  binding:"required"`
}

type UnbindCourseResponse struct {
	Code ErrNo
}

// 获取老师下所有课程
// Method：Get
type GetTeacherCourseRequest struct {
	TeacherID string `form:"TeacherID" json:"TeacherID" xml:"TeacherID"  binding:"required"`
}

type GetTeacherCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseList []*TCourse
	}
}

// 排课求解器，使老师绑定课程的最优解， 老师有且只能绑定一个课程
// Method： Post
type ScheduleCourseRequest struct {
	TeacherCourseRelationShip map[string][]string `form:"TeacherCourseRelationShip" json:"TeacherCourseRelationShip" xml:"TeacherCourseRelationShip"  binding:"required"` // key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
}

type ScheduleCourseResponse struct {
	Code ErrNo
	Data map[string]string // key 为 teacherID , val 为老师最终绑定的课程 courseID
}

type BookCourseRequest struct {
	StudentID string `form:"StudentID" json:"StudentID" xml:"StudentID"  binding:"required"`
	CourseID  string `form:"CourseID" json:"CourseID" xml:"CourseID"  binding:"required"`
}

// 课程已满返回 CourseNotAvailable

type BookCourseResponse struct {
	Code ErrNo
}

type GetStudentCourseRequest struct {
	StudentID string `form:"StudentID" json:"StudentID" xml:"StudentID"  binding:"required"`
}

type GetStudentCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseList []TCourse
	}
}
