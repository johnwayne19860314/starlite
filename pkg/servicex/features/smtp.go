package features

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/servicex/utils"
)

type SMTPMessage struct {
	ToUsers []string
	CcUsers []string

	Subject          string
	Content          string
	ContentSignature string

	Options []func(m *gomail.Message)
}

func DoSendSMTPMessage(appCtx appx.AppContext, msg SMTPMessage) error {
	msg.Content = msg.Content + msg.ContentSignature
	if !appCtx.IsPrd() {
		logx.Info("DoSendMail", "subject", msg.Subject, "content", msg.Content, "toUsers", msg.ToUsers, "CcUsers", msg.CcUsers)
	}

	m := gomail.NewMessage()
	for _, opt := range msg.Options {
		opt(m)
	}
	m.SetAddressHeader("From", GetMailFrom(appCtx), "GFSH-BJM")
	m.SetHeader("To", msg.ToUsers...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html;charset=utf-8", msg.Content)
	d := gomail.NewDialer(GetSmtpFrom(appCtx), utils.MailSmtpPort, "", "")
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	if err := d.DialAndSend(m); err != nil {
		return errorx.WithStack(err)
	}
	return nil
}
