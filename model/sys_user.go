package model

import (
	"go-admin/global"

	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.GVA_MODEL
	UUID        uuid.UUID      `json:"uuid"      gorm:"index;comment:用户UUID"`                                                // 用户UUID
	Username    string         `json:"userName"  gorm:"index;comment:用户登录名"`                                                 // 用户登录名
	Password    string         `json:"-"         gorm:"comment:用户登录密码"`                                                      // 用户登录密码
	NickName    string         `json:"nickName"  gorm:"default:系统用户;comment:用户昵称"`                                           // 用户昵称
	SideMode    string         `json:"sideMode"  gorm:"default:dark;comment:用户侧边主题"`                                         // 用户侧边主题
	HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	ActiveColor string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`                                      // 活跃颜色
	AuthorityId uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority   SysAuthority   `json:"authority"   gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone       string         `json:"phone"       gorm:"comment:用户手机号"`                     // 用户手机号
	Email       string         `json:"email"       gorm:"comment:用户邮箱"`                      // 用户邮箱
	Enable      int            `json:"enable"      gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

func (SysUser) TableName() string {
	return "sys_users"
}

type ChangeUserInfo struct {
	ID           uint           `gorm:"primarykey"`                                                                              // 主键ID
	NickName     string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                               // 用户昵称
	Phone        string         `json:"phone"    gorm:"comment:用户手机号"`                                                           // 用户手机号
	AuthorityIds []uint         `json:"authorityIds" gorm:"-"`                                                                   // 角色ID
	Email        string         `json:"email"        gorm:"comment:用户邮箱"`                                                        // 用户邮箱
	HeaderImg    string         `json:"headerImg"    gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	SideMode     string         `json:"sideMode"     gorm:"comment:用户侧边主题"`                                                      // 用户侧边主题
	Enable       int            `json:"enable"       gorm:"comment:冻结用户"`                                                        //冻结用户
	Authorities  []SysAuthority `json:"-"            gorm:"many2many:sys_user_authority;"`
}

// Register User register structure
type Register struct {
	Username     string `json:"userName"  example:"用户名"`
	Password     string `json:"passWord"  example:"密码"`
	NickName     string `json:"nickName"  example:"昵称"`
	HeaderImg    string `json:"headerImg" example:"头像链接"`
	AuthorityId  uint   `json:"authorityId"  swaggertype:"string"  example:"int 角色id"`
	Enable       int    `json:"enable"       swaggertype:"string"  example:"int 是否启用"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string"  example:"[]uint 角色id"`
	Phone        string `json:"phone" example:"电话号码"`
	Email        string `json:"email" example:"电子邮箱"`
}

// User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

// Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

// Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

// Modify  user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}
