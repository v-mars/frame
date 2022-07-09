package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/convert"
)

const (
	SUCCESS_CODE          = 200        //成功的状态码
	WARNING_CODE          = 300        //成功的状态码
	FAIL_CODE             = 400        //失败的状态码
	MD5_PREFIX            = "jkfldfsf" //MD5加密前缀字符串
	TOKEN_KEY             = "X-Token"  //页面token键名
	USER_ID_Key           = "X-USERID" //页面用户ID键名
	USER_UUID_Key         = "X-UUID"   //页面UUID键名
	SUPER_ADMIN_ID uint64 = 956986     // 超级管理员账号ID
	//SUCCESS_CODE          = 2000      //成功的状态码
	//FAIL_CODE             = 4000      //失败的状态码
	ErrNoPerm          = 401 //无访问权限
	ErrNotFound        = 404 //资源不存在
	ErrMethodNotAllow  = 405 //方法不被允许
	ErrTooManyRequests = 429 //请求过于频繁
	ErrInternalServer  = 500 //服务器发生错误
	ErrInvalidToken    = 999 //令牌失效
)

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, v interface{}) {
	c.JSON(status, v)
	c.Abort()
}

// GetPageIndex 获取页码
func GetPageIndex(c *gin.Context) int {
	return GetQueryToInt(c, "page")
}

// GetPageLimit 获取每页记录数
func GetPageLimit(c *gin.Context) int {
	limit := GetQueryToInt(c, "pageSize")
	//if limit > 500 {
	//	limit = 500
	//}
	return limit
}

// GetPageExport 获取导出参数
func GetPageExport(c *gin.Context) string {
	return GetQueryToStr(c, "export_excel")
}

// GetPageSort 获取排序信息
func GetPageSort(c *gin.Context) string {
	return GetQueryToStr(c, "sort")
}

// GetPageKey 获取搜索关键词信息
func GetPageKey(c *gin.Context) string {
	return GetQueryToStr(c, "key")
}

func GetQueryToStrE(c *gin.Context, key string) (string, error) {
	str, ok := c.GetQuery(key)
	if !ok {
		return "", fmt.Errorf("没有这个值[%s]传入", key)
	}
	return str, nil
}

func GetQueryToStr(c *gin.Context, key string, defaultValues ...string) string {
	var defaultValue string
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}
	str, err := GetQueryToStrE(c, key)
	if str == "" || err != nil {
		return defaultValue
	}
	return str
}

// GetQueryToUintE QueryToUintE
func GetQueryToUintE(c *gin.Context, key string) (uint, error) {
	str, err := GetQueryToStrE(c, key)
	if err != nil {
		return 0, err
	}
	return convert.ToUintE(str)
}

// GetQueryToUint QueryToUint
func GetQueryToUint(c *gin.Context, key string, defaultValues ...uint) uint {
	var defaultValue uint
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}
	val, err := GetQueryToUintE(c, key)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetQueryToUint64E QueryToUintE
func GetQueryToUint64E(c *gin.Context, key string) (uint64, error) {
	str, err := GetQueryToStrE(c, key)
	//fmt.Println(str)
	if err != nil {
		return 0, err
	}
	return convert.ToUint64E(str)
}

// GetQueryToIntE QueryToUintE
func GetQueryToIntE(c *gin.Context, key string) (int, error) {
	str, err := GetQueryToStrE(c, key)
	//fmt.Println(str)
	if err != nil {
		return 0, err
	}
	return convert.ToInt(str), nil
}

// GetQueryToUint64 QueryToUint
func GetQueryToUint64(c *gin.Context, key string, defaultValues ...uint64) uint64 {
	var defaultValue uint64
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}
	val, err := GetQueryToUint64E(c, key)
	if err != nil {
		return defaultValue
	}
	return val
}

// GetQueryToInt QueryToInt
func GetQueryToInt(c *gin.Context, key string, defaultValues ...int) int {
	var defaultValue int
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}
	val, err := GetQueryToIntE(c, key)
	if err != nil {
		return defaultValue
	}
	return val
}
