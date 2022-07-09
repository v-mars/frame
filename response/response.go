package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/v-mars/frame/pkg/convert"
	"github.com/v-mars/frame/pkg/logger"
	"github.com/v-mars/frame/pkg/utils"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

type Data struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

type CmdData struct {
	Data
	ExitCode int    `json:"exit_code"`
	Stderr   string `json:"stderr"`
	Stdout   string `json:"stdout"`
}

type PageDataList struct {
	List     interface{} `json:"list"`
	PageSize int         `json:"pageSize"`
	Page     int         `json:"page"`
	Total    int64       `json:"total"`
}

func InitPageData(c *gin.Context, list interface{}, all bool) PageDataList {
	pageSize := 10
	page := 1
	exportExcel := ""
	if c != nil {
		pageSize = GetPageLimit(c)
		page = GetPageIndex(c)
		exportExcel = GetPageExport(c)
	}
	if exportExcel == "xlsx" || exportExcel == "csv" {
		all = true
	}
	if all {
		pageSize = 0
	}
	var pageData = PageDataList{Page: page, PageSize: pageSize, List: list}
	return pageData
}

// PageSuccess 响应成功
func PageSuccess(c *gin.Context, v PageDataList) {
	ret := Data{Code: SUCCESS_CODE, Status: "success", Message: "ok", Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

// CreateSuccess 创建响应成功
func CreateSuccess(c *gin.Context, v ...interface{}) {
	ret := Data{Code: SUCCESS_CODE, Status: "success", Message: "创建成功.", Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

// UpdateSuccess 更新响应成功
func UpdateSuccess(c *gin.Context, v ...interface{}) {
	ret := Data{Code: SUCCESS_CODE, Status: "success", Message: "更新成功.", Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

// DeleteSuccess 删除响应成功
func DeleteSuccess(c *gin.Context, v ...interface{}) {
	ret := Data{Code: SUCCESS_CODE, Status: "success", Message: "删除成功.", Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

// Error 响应失败
func Error(c *gin.Context, err error, v ...interface{}) {
	ret := Data{Code: FAIL_CODE, Status: "error", Message: err.Error(), Data: v}
	c.Set("errorMsg", err)
	stack := Stack(2)
	logger.Error(string(stack))
	ResJSON(c, http.StatusOK, &ret)
}

func ErrorNoStack(c *gin.Context, err error, v ...interface{}) {
	ret := Data{Code: FAIL_CODE, Status: "error", Message: err.Error(), Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

func ParamFailed(c *gin.Context, err error, v ...interface{}) {
	errStr := err.Error()
	// Key: 'QueryParam.K8sClusterID' Error:Field validation for 'K8sClusterID' failed on the 'required' tag
	comp := regexp.MustCompile(`Key: '.+' Error:Field validation for '(.+)?' failed on the '(.+)?' .+`)
	subMatches := comp.FindAllStringSubmatch(errStr, -1)
	// 报错格式化"
	if len(subMatches) > 0 {
		if len(subMatches[0]) >= 3 {
			errStr = fmt.Sprintf("请求参数验证错误：%s参数【%s】", subMatches[0][2], subMatches[0][1])
			errStr = strings.Replace(errStr, "required", "缺少", -1)
		} else {
			errStr = fmt.Sprintf("请求参数验证错误：%s", subMatches[0])
		}
	}
	ret := Data{Code: FAIL_CODE, Status: "error", Message: errStr, Data: v}
	c.Set("errorMsg", err)
	stack := Stack(2)
	logger.Error(string(stack))
	ResJSON(c, http.StatusOK, &ret)
}

func SqlFailed(c *gin.Context, err error, v ...interface{}) {
	errStr := err.Error()
	// Error 1062: Duplicate entry '1.1.1.2' for key 'idx_name_code'
	comp := regexp.MustCompile(`Error 1062: Duplicate entry '(.+)?' for key '(.+)?'`)
	subMatches := comp.FindAllStringSubmatch(errStr, -1)
	// "记录重复。具体报错："
	if len(subMatches) > 0 {
		if len(subMatches[0]) >= 2 {
			errStr = fmt.Sprintf("记录重复。具体报错：%s", subMatches[0][1])
		} else {
			errStr = fmt.Sprintf("记录重复。具体报错：%s", subMatches[0])
		}
	}
	ret := Data{Code: FAIL_CODE, Status: "error", Message: errStr, Data: v}
	c.Set("errorMsg", err)
	stack := Stack(2)
	logger.Error(string(stack))
	ResJSON(c, http.StatusOK, &ret)
}

func NoPermission(c *gin.Context, v ...interface{}) {
	ret := Data{Code: ErrNoPerm, Status: "error", Message: "无权限访问", Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

func FailedMsg(c *gin.Context, msg error, v ...interface{}) {
	ret := Data{Code: FAIL_CODE, Status: "error", Message: msg.Error(), Data: v}
	c.Set("errorMsg", msg.Error())
	stack := Stack(2)
	logger.Error(string(stack))
	ResJSON(c, http.StatusOK, &ret)
}

// FailedCode 响应失败 code
func FailedCode(c *gin.Context, code int, msg error, v ...interface{}) {
	ret := Data{Code: code, Status: "error", Message: msg.Error(), Data: v}
	c.Set("errorMsg", msg.Error())
	stack := Stack(2)
	logger.Error(string(stack))
	ResJSON(c, http.StatusOK, &ret)
}

// FailedCodeRecovery 响应失败 code
func FailedCodeRecovery(c *gin.Context, code int, msg error, RecoveryErr error) {
	ret := Data{Code: code, Status: "error", Message: msg.Error(), Data: nil}
	c.Set("errorMsg", msg.Error())
	if RecoveryErr != nil {
		c.Set("stack", fmt.Errorf("%s", RecoveryErr))
	}
	ResJSON(c, http.StatusInternalServerError, &ret)
}

func Success(c *gin.Context, v interface{}) {
	ret := Data{Code: SUCCESS_CODE, Status: "success", Message: "ok", Data: v}
	export, _ := c.GetQuery("export_excel")
	if export == "csv" || export == "xlsx" {
		fileName := fmt.Sprintf("%s_download.csv", convert.GetNowTimeNoFormatStr())
		filePath := fmt.Sprintf("./%s", fileName)
		err := createExport(filePath, v)
		if err != nil {
			Error(c, err)
			return
		}
		//c.Header("Content-Type", "text/csv;charset=utf-8")
		//c.Header("Content-Type", "application/vnd.ms-excel;charset=utf-8")
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment;filename="+fileName)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(filePath)
		//c.Data(http.StatusOK, "text/csv;charset=utf-8", fileContent)
		return
	} else {
		ResJSON(c, http.StatusOK, &ret)
	}
}

func Download(c *gin.Context, fileName, filePath string) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment;filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}

func Res(c *gin.Context, code int, msg string, v interface{}) {
	ret := Data{Code: code, Status: "success", Message: msg, Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

func SuccessMsg(c *gin.Context, msg string, v interface{}) {
	ret := Data{Code: SUCCESS_CODE, Status: "success", Message: msg, Data: v}
	ResJSON(c, http.StatusOK, &ret)
}

func SuccessCmd(c *gin.Context, stdout, stderr string, exitCode int, v interface{}) {
	ret := CmdData{Data: Data{Code: SUCCESS_CODE, Status: "success", Message: "ok", Data: v},
		ExitCode: exitCode, Stdout: stdout, Stderr: stderr}
	ResJSON(c, http.StatusOK, &ret)
}

func ErrorCmd(c *gin.Context, stdout, stderr string, exitCode int, err error, v interface{}) {
	ret := CmdData{Data: Data{Code: FAIL_CODE, Status: "error", Message: err.Error(), Data: v},
		ExitCode: exitCode, Stdout: stdout, Stderr: stderr}
	c.Set("errorMsg", err)
	stack := Stack(2)
	logger.Error(string(stack))
	ResJSON(c, FAIL_CODE, &ret)
}

func Html(c *gin.Context, htmlContent string) {
	var templateName = "text_html"
	var pageInstance render.Render
	pageContent := template.Must(template.New(templateName).Parse(htmlContent))
	htmlRender := render.HTMLProduction{Template: pageContent}
	pageInstance = htmlRender.Instance(templateName, map[string]interface{}{})
	c.Render(http.StatusOK, pageInstance)
}

func createExport(filePath string, data interface{}) error {
	pageData, ok := data.(PageDataList)
	if !ok {
		return fmt.Errorf("需要导出的数据格式不正确")
	}
	//var t = reflect.TypeOf(pageData.List)
	//var V = reflect.ValueOf(pageData.List)
	fmt.Println()
	//_, _ = pageData.List.(t)
	//var t = reflect.TypeOf(s1)
	//var v = reflect.ValueOf(s1)
	var array []interface{}
	//var array []map[string]interface{}
	err := utils.AnyToAny(pageData.List, &array)
	if err != nil {
		return fmt.Errorf("需要导出的数据转换为Map数组报错，%s", err)
	}
	if len(array) == 0 {
		return fmt.Errorf("需要导出的数据不能为空")
	}
	//for _, v := range array {
	//	var t = reflect.TypeOf(v)
	//	var vv = reflect.ValueOf(v)
	//	fmt.Println("t.NumField():", t.NumField())
	//	for i := 0; i < t.NumField(); i++ {
	//		tField := t.Field(i)
	//		vField := vv.Field(i)
	//		fmt.Printf("Name:%v\n", tField.Name)
	//		fmt.Printf("Name:%v\n", tField.Type)
	//		fmt.Printf("Kind:%v\n", vField.Kind())
	//	}
	//
	//}

	//var headline []string

	//for k := range array[0] {
	//	headline = append(headline, k)
	//}
	//
	//file, err := os.Create(filePath)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//writer := csv.NewWriter(file)
	//writer.Comma = ','
	////headline := []string{"URL", "Title", "Keywords", "Description"}
	//err = writer.Write(headline)
	//if err != nil {
	//	return err
	//}
	//var rows = make([][]string, 0)
	//for _, dict := range array {
	//	dict := dict
	//	var temp = make([]string, 0)
	//	for _, value := range dict {
	//		value := value
	//		temp = append(temp, convert.ToString(value))
	//	}
	//	err = writer.Write(temp)
	//	if err != nil {
	//		return err
	//	}
	//	rows = append(rows, temp)
	//}
	//writer.Flush()
	return nil
}
