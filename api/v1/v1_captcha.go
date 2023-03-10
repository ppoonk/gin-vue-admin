package api

import (
	"log"
	"time"

	"go-admin/global"
	"go-admin/model"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
// var store = captcha.NewDefaultRedisStore()
var store = base64Captcha.DefaultMemStore

type BaseApi struct{}

// Captcha
// @Tags      Base
// @Summary   生成验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysCaptchaResponse,msg=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func Captcha(c *gin.Context) {
	log.Println("开始Captcha")
	// 判断验证码是否开启
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.BlackCache.Get(key) ///???
	log.Println("开始v", v)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool

	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}

	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.GVA_CONFIG.Captcha.ImgHeight, global.GVA_CONFIG.Captcha.ImgWidth, global.GVA_CONFIG.Captcha.KeyLong, 0.7, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		//global.GVA_LOG.Error("验证码获取失败!", zap.Error(err))
		model.Result(7, "验证码获取失败", nil, c)
		return
	}
	log.Println("验证码获取成功id:", id)
	model.Result(200, "验证码获取成功", gin.H{
		"captchaId":     id,
		"picPath":       b64s,
		"captchaLength": global.GVA_CONFIG.Captcha.KeyLong,
		"openCaptcha":   oc,
	}, c)
}

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}
