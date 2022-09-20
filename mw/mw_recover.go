package mw

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/notify/email"
	"github.com/v-mars/frame/response"
)

var (
	AppName         = ""
	Username        = ""
	Password        = ""
	SmtpServer      = ""
	SmtpPort        = 25
	Tls             = false
	OpsMailReceiver []string
)

// RecoveryMiddleware 崩溃恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//var data string
		//if c.Request.Method != http.MethodGet { // 如果不是GET请求，则读取body
		//	body, err := c.GetRawData() 		// body 只能读一次，读出来之后需要重置下 Body
		//	if err != nil {logger.Error(err)}
		//	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重置body
		//	data = string(body)
		//}
		//start := time.Now()
		defer func() {
			if err := recover(); err != any(nil) {
				//cost := time.Since(start)
				//userName, ok := c.Get("nickname")
				//if !ok {userName = "nil"}
				stack := response.Stack(3)
				//errMsg :=fmt.Sprintf(
				//	"用户: %s, 方法: %s, URL: %s, CODE: %d, 耗时: %dms, Body数据: %s, \nERROR: %s, \n堆栈信息: \n%s",
				//	userName,c.Request.Method,c.Request.URL.Path,c.Writer.Status(),cost.Milliseconds(),data, err,stack)
				// 这里会打印出错栈信息 time.Now().Format("2006-01-02 15:04:05")
				//logger.Error(errMsg)

				if len(Username) > 0 && len(Password) > 0 && len(SmtpServer) > 0 && len(OpsMailReceiver) > 0 {
					var title = fmt.Sprintf("%s [Recovery]", AppName)
					var mailText = fmt.Sprintf(`
<pre style="margin: 1px;white-space: pre-wrap;word-wrap: break-word;">
用户：%s 
%s请求URL: %s
堆栈信息:
%s</pre>`, c.GetString("nickname"), c.Request.Method, c.FullPath(), string(stack))
					notify := email.NewMail(Username, Password,
						SmtpServer, Username, SmtpPort, Tls)
					if e := notify.Send(title, mailText, OpsMailReceiver, []string{}); e != nil {
					}
				}
				//errMsg :=fmt.Sprintf(
				//	"用户: %s, 方法: %s, URL: %s, CODE: %d, 耗时: %dms, Body数据: %s, \nERROR: %s, \n堆栈信息: \n%s",
				//	userName,c.Request.Method,c.Request.URL.Path,c.Writer.Status(),cost.Milliseconds(),data, err,stack)
				// 这里会打印出错栈信息 time.Now().Format("2006-01-02 15:04:05")
				//logger.Error(errMsg)
				response.FailedCodeRecovery(c, 5009, fmt.Errorf("[Recovery]: %s", err), fmt.Errorf("%s", stack))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
