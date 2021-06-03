package email

import (
	"testing"
)

func TestNewSMTP(t *testing.T) {
	// odftxhxilvfzcbci

	mail := NewMail("xx@qq.com", "xx", "smtp.qq.com",
		"ocean.zhang", 465,false)
	_ = mail.Send("测试用gomail发送邮件", "Good Good Study, Day Day Up!!!!!!", []string{`429472406@qq.com`, "donghai.zhang@cootek.cn"}, []string{})
}
