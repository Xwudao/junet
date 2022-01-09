package mailx

import (
	"testing"
	"time"
)

func TestSendEmail(t *testing.T) {
	Init(SetHost("smtp.mxhichina.com"),
		SetPassword(""),
		SetUsername(""),
		SetPort(80),
	)

	SendEmail(&Message{
		From:     "",
		To:       []string{""},
		CC:       nil,
		Subject:  "这是一封测试邮件",
		TextBody: "hello kuan",
		HtmlBody: "",
	})

	time.Sleep(time.Second * 20)
}
