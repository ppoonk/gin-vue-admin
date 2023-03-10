package api

import (
	"go-admin/global"
	"go-admin/model"
	"go-admin/service"
	"go-admin/utils"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// 用户注册
func Register(c *gin.Context) {
	var r model.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		model.Result(7, "参数错误", err.Error(), c)

		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		model.Result(7, "参数校验错误", err.Error(), c)

		return
	}
	var authorities []model.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, model.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &model.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userReturn, err := service.Register(*user)
	if err != nil {
		model.Result(7, "注册失败", err.Error(), c)
		return
	}
	model.Result(7, "注册成功", userReturn, c)

}

// 用户登录
func Login(c *gin.Context) {
	var l model.Login
	err := c.ShouldBindJSON(&l) //请求头使用的是form-data而后端接受使用的json
	log.Println("接受的json", l)

	key := c.ClientIP()

	if err != nil {
		model.Result(7, "参数错误", err.Error(), c)
		return
	}
	//校验参数
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		model.Result(7, "参数格式错误", err.Error(), c)
		return
	}

	//判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)

	if !oc || store.Verify(l.CaptchaId, l.Captcha, true) { //验证码
		u := &model.SysUser{Username: l.Username, Password: l.Password} //
		user, err := service.Login(u)
		if err != nil {

			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			model.Result(7, "用户名不存在或者密码错误", err.Error(), c)

			return
		}
		if user.Enable != 1 {

			// 验证码次数+1
			global.BlackCache.Increment(key, 1)
			model.Result(0, "用户被禁止登录", err.Error(), c)
			return
		}
		TokenNext(c, *user)
		//登录以后签发jwt
		return
	}
	// 验证码次数+1
	global.BlackCache.Increment(key, 1)
	model.Result(7, "验证码错误", nil, c)

}

// TokenNext 登录以后签发jwt
func TokenNext(c *gin.Context, user model.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(model.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	log.Println("token", token)
	if err != nil {
		model.Result(7, "获取token失败", nil, c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint { //多点登陆拦截
		model.Result(0, "登录成功", gin.H{
			"user":      user,
			"token":     token,
			"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
		}, c)
		return
	}
	//
	if jwtStr, err := service.GetRedisJWT(user.Username); err == redis.Nil { //
		if err := service.SetRedisJWT(token, user.Username); err != nil {
			model.Result(7, "设置登录状态失败", nil, c)
			return
		}
		model.Result(0, "登录成功", gin.H{
			"user":      user,
			"token":     token,
			"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
		}, c)

	} else if err != nil {
		model.Result(7, "设置登录状态失败", nil, c)
	} else {
		var blackJWT model.JwtBlacklist //
		blackJWT.Jwt = jwtStr
		if err := service.JsonInBlacklist(blackJWT); err != nil {
			model.Result(7, "jwt作废失败", nil, c)
			return
		}
		if err := service.SetRedisJWT(token, user.Username); err != nil {
			model.Result(7, "设置登录状态失败", nil, c)
			return
		}
		model.Result(0, "登录成功", gin.H{
			"user":      user,
			"token":     token,
			"expiresAt": claims.StandardClaims.ExpiresAt * 1000,
		}, c)
	}
}

// 用户修改密码
func ChangePassword(c *gin.Context) {
	var req model.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify) //校验参数
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	uid := utils.GetUserID(c) //从Gin的Context中获取从jwt解析出来的用户ID
	u := &model.SysUser{GVA_MODEL: global.GVA_MODEL{ID: uid}, Password: req.Password}
	_, err = service.ChangePassword(u, req.NewPassword)
	if err != nil {
		model.Result(7, "修改失败，原密码与当前账户不符", nil, c)
		return
	}
	model.Result(0, "修改成功", nil, c)
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := service.GetUserInfo(uuid)
	if err != nil {
		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{"userInfo": ReqUser}, c)
}

// @Summary   设置用户信息
func SetUserInfo(c *gin.Context) {
	var user model.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(user, utils.IdVerify) //校验参数
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}

	if len(user.AuthorityIds) != 0 {
		err = service.SetUserAuthorities(user.ID, user.AuthorityIds) //设置一个用户的权限
		if err != nil {
			model.Result(7, "设置失败", nil, c)
			return
		}
	}
	err = service.SetUserInfo(model.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		model.Result(7, "设置失败", nil, c)
		return
	}
	model.Result(0, "设置成功", nil, c)
}

// 设置用户权限组
func SetUserAuthorities(c *gin.Context) {
	var sua model.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		model.Result(7, "修改失败", nil, c)
		return
	}
	model.Result(0, "修改成功", nil, c)
}

// "设置用户权限"

func SetUserAuthority(c *gin.Context) {
	var sua model.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil { //参数校验
		model.Result(7, UserVerifyErr.Error(), nil, c)
		return
	}
	userID := utils.GetUserID(c)
	err = service.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {

		model.Result(7, err.Error(), nil, c)
		return
	}
	claims := utils.GetUserInfo(c) //从Gin的Context中获取从jwt解析出来的用户角色

	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // 唯一签名
	claims.AuthorityId = sua.AuthorityId                                  //更新角色id

	if token, err := j.CreateToken(*claims); err != nil { //设置新token
		model.Result(7, err.Error(), nil, c)
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt, 10))
		model.Result(0, "修改成功", nil, c)
	}
}

// 设置用户自身信息
func SetSelfInfo(c *gin.Context) {
	var user model.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	user.ID = utils.GetUserID(c)
	err = service.SetUserInfo(model.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		model.Result(7, "设置失败", nil, c)
		return
	}
	model.Result(0, "设置成功", nil, c)
}

// "重置用户密码"
func ResetPassword(c *gin.Context) {
	var user model.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {

		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.ResetPassword(user.ID)
	if err != nil {
		model.Result(7, "重置失败"+err.Error(), nil, c)

		return
	}
	model.Result(0, "重置成功", nil, c)
}

// "删除用户"
func DeleteUser(c *gin.Context) {
	var reqId model.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify) //参数校验
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		model.Result(7, "删除失败, 自杀失败", nil, c)
		return
	}
	err = service.DeleteUser(reqId.ID)
	if err != nil {
		model.Result(7, "删除失败", nil, c)
		return
	}
	model.Result(7, "删除成功", nil, c)
}

//分页获取用户列表
// "分页获取用户列表,返回包括列表,总数,页码,每页数量"

func GetUserList(c *gin.Context) {
	var pageInfo model.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	list, total, err := service.GetUserInfoList(pageInfo)
	if err != nil {
		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"list":     list,
		"total":    total,
		"page":     pageInfo.Page,
		"pageSize": pageInfo.PageSize,
	}, c)

}
