package api

import (
	"errors"
	"fmt"
	"go-admin/global"
	"go-admin/model"
	"go-admin/service"
	"go-admin/utils"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser = cases.Title(language.English)

// PreviewTemp
// @Tags      AutoCode
// @Summary   预览创建后的代码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.struct                                      true  "预览创建代码"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "预览创建后的代码"
// @Router    /autoCode/preview [post]
func PreviewTemp(c *gin.Context) {
	var a model.AutoCodeStruct
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoCodeVerify); err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	a.Pretreatment() // 处理go关键字
	a.PackageT = caser.String(a.Package)
	autoCode, err := service.PreviewTemp(a)
	if err != nil {

		model.Result(7, "预览失败", nil, c)
	} else {
		model.Result(0, "预览成功", gin.H{"autoCode": autoCode}, c)
	}
}

// CreateTemp
// @Tags      AutoCode
// @Summary   自动代码模板
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.struct  true  "创建自动代码"
// @Success   200   {string}  string                 "{"success":true,"data":{},"msg":"创建成功"}"
// @Router    /autoCode/createTemp [post]
func CreateTemp(c *gin.Context) {
	var a model.AutoCodeStruct
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoCodeVerify); err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	a.Pretreatment()
	var apiIds []uint
	if a.AutoCreateApiToSql {
		if ids, err := service.AutoCreateApi(&a); err != nil {
			c.Writer.Header().Add("success", "false")
			c.Writer.Header().Add("msg", url.QueryEscape("自动化创建失败!请自行清空垃圾数据!"))
			return
		} else {
			apiIds = ids
		}
	}
	a.PackageT = caser.String(a.Package)
	err := service.CreateTemp(a, apiIds...)
	if err != nil {
		if errors.Is(err, model.ErrAutoMove) {
			c.Writer.Header().Add("success", "true")
			c.Writer.Header().Add("msg", url.QueryEscape(err.Error()))
		} else {
			c.Writer.Header().Add("success", "false")
			c.Writer.Header().Add("msg", url.QueryEscape(err.Error()))
			_ = os.Remove("./ginvueadmin.zip")
		}
	} else {
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ginvueadmin.zip")) // fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Writer.Header().Add("success", "true")
		c.File("./ginvueadmin.zip")
		_ = os.Remove("./ginvueadmin.zip")
	}
}

// GetDB
// @Tags      AutoCode
// @Summary   获取当前所有数据库
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取当前所有数据库"
// @Router    /autoCode/getDatabase [get]
func GetDB(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbs, err := service.Database(businessDB).GetDB(businessDB)
	var dbList []map[string]interface{}
	for _, db := range global.GVA_CONFIG.DBList {
		var item = make(map[string]interface{})
		item["aliasName"] = db.AliasName
		item["dbName"] = db.Dbname
		item["disable"] = db.Disable
		item["dbtype"] = db.Type
		dbList = append(dbList, item)
	}
	if err != nil {
		model.Result(7, "获取失败", nil, c)
	} else {
		model.Result(0, "获取成功", gin.H{"dbs": dbs, "dbList": dbList}, c)
	}
}

// GetTables
// @Tags      AutoCode
// @Summary   获取当前数据库所有表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取当前数据库所有表"
// @Router    /autoCode/getTables [get]
func GetTables(c *gin.Context) {
	dbName := c.DefaultQuery("dbName", global.GVA_CONFIG.Mysql.Dbname)
	businessDB := c.Query("businessDB")
	tables, err := service.Database(businessDB).GetTables(businessDB, dbName)
	if err != nil {
		model.Result(7, "查询table失败", nil, c)
	} else {
		model.Result(0, "获取成功", gin.H{"tables": tables}, c)
	}
}

// GetColumn
// @Tags      AutoCode
// @Summary   获取当前表所有字段
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取当前表所有字段"
// @Router    /autoCode/getColumn [get]
func GetColumn(c *gin.Context) {
	businessDB := c.Query("businessDB")
	dbName := c.DefaultQuery("dbName", global.GVA_CONFIG.Mysql.Dbname)
	tableName := c.Query("tableName")
	columns, err := service.Database(businessDB).GetColumn(businessDB, tableName, dbName)
	if err != nil {
		model.Result(7, "获取失败", nil, c)
	} else {
		model.Result(0, "获取成功", gin.H{"columns": columns}, c)
	}
}

// CreatePackage
// @Tags      AutoCode
// @Summary   创建package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAutoCode                                         true  "创建package"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "创建package成功"
// @Router    /autoCode/createPackage [post]
func CreatePackage(c *gin.Context) {
	var a model.SysAutoCode
	_ = c.ShouldBindJSON(&a)
	if err := utils.Verify(a, utils.AutoPackageVerify); err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err := service.CreateAutoCode(&a)
	if err != nil {
		model.Result(7, "创建失败", nil, c)
	} else {
		model.Result(0, "创建成功", nil, c)
	}
}

// GetPackage
// @Tags      AutoCode
// @Summary   获取package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "创建package成功"
// @Router    /autoCode/getPackage [post]
func GetPackage(c *gin.Context) {
	pkgs, err := service.GetPackage()
	if err != nil {
		model.Result(7, "获取失败", nil, c)
	} else {
		model.Result(0, "获取成功", gin.H{"pkgs": pkgs}, c)
	}
}

// DelPackage
// @Tags      AutoCode
// @Summary   删除package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAutoCode                                         true  "创建package"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "删除package成功"
// @Router    /autoCode/delPackage [post]
func DelPackage(c *gin.Context) {
	var a model.SysAutoCode
	_ = c.ShouldBindJSON(&a)
	err := service.DelPackage(a)
	if err != nil {
		model.Result(7, "删除失败", nil, c)
	} else {
		model.Result(0, "删除成功", nil, c)
	}
}

// AutoPlug
// @Tags      AutoCode
// @Summary   创建插件模板
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAutoCode                                         true  "创建插件模板"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "创建插件模板成功"
// @Router    /autoCode/createPlug [post]
func AutoPlug(c *gin.Context) {
	var a model.AutoPlugReq
	err := c.ShouldBindJSON(&a)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	a.Snake = strings.ToLower(a.PlugName)
	a.NeedModel = a.HasRequest || a.HasResponse
	err = service.CreatePlug(a)
	if err != nil {
		model.Result(7, "预览失败", nil, c)
		return
	}
	model.Result(0, "操作成功", nil, c)

}

// InstallPlugin
// @Tags      AutoCode
// @Summary   安装插件
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     plug  formData  file                                              true  "this is a test file"
// @Success   200   {object}  response.Response{data=[]interface{},msg=string}  "安装插件成功"
// @Router    /autoCode/createPlug [post]
func InstallPlugin(c *gin.Context) {
	header, err := c.FormFile("plug")
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	web, server, err := service.InstallPlugin(header)
	webStr := "web插件安装成功"
	serverStr := "server插件安装成功"
	if web == -1 {
		webStr = "web端插件未成功安装，请按照文档自行解压安装，如果为纯后端插件请忽略此条提示"
	}
	if server == -1 {
		serverStr = "server端插件未成功安装，请按照文档自行解压安装，如果为纯前端插件请忽略此条提示"
	}
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	model.Result(0, "查询成功", []interface{}{
		gin.H{
			"code": web,
			"msg":  webStr,
		},
		gin.H{
			"code": server,
			"msg":  serverStr,
		},
	}, c)
}
