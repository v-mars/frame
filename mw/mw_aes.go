package mw

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/frame/response"
	"io/ioutil"
)

func AesEncDec() gin.HandlerFunc {
	return func(c *gin.Context) {
		encryptTag := c.Request.Header.Get("EncryptData")
		if encryptTag == "yes" {
			rawData, err := c.GetRawData() // body 只能读一次，读出来之后需要重置下 Body
			if err != nil {
				response.FailedCode(c, 6001, fmt.Errorf("%s", err))
				return
			}
			data := utils.DeTxtByAes(string(rawData), response.DefaultAesKey)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(data))) // 重置body
		}
		// 处理请求
		c.Next()
	}
}
