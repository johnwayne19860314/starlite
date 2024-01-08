package features

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/servicex/utils"
)

type EmailOptFunc func(m *gomail.Message)

func WithCC(ccUsers []string) EmailOptFunc {
	return func(m *gomail.Message) {
		m.SetHeader("Cc", ccUsers...)
	}
}

func DoSendMail(appCtx appx.AppContext, toUsers []string, subject, content, endpoint string) error {
	return DoSendEmailWithOptions(appCtx, toUsers, subject, content, endpoint)
}

func DoSendEmailWithOptions(appCtx appx.AppContext, toUsers []string, subject, content, endpoint string, opts ...EmailOptFunc) error {
	content = EmailContentAppendSignatureDefault(content)
	if !appCtx.IsPrd() {
		logx.Info("DoSendMail", "subject", subject, "content", content, "toUsers", toUsers)
	}

	m := gomail.NewMessage()
	for _, opt := range opts {
		opt(m)
	}
	m.SetAddressHeader("From", GetMailFrom(appCtx), "GFSH-BJM")
	m.SetHeader("To", toUsers...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html;charset=utf-8", content)
	d := gomail.NewDialer(endpoint, utils.MailSmtpPort, "", "")
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	if err := d.DialAndSend(m); err != nil {
		return errorx.WithStack(err)
	}
	return nil
}

func GetSmtpFrom(appCtx appx.AppContext) string {
	if appCtx.IsLocal() {
		return utils.MailSmtpIntHost
	}
	return utils.MailSmtpExtHost
}

func GetMailFrom(appCtx appx.AppContext) string {
	if !appCtx.IsPrd() {
		return utils.MailFromTest
	}
	return utils.MailFrom
}

const EmailContentSignature = `
<br>
<br>
<br>
声明：此邮件中包含的信息及附件为特斯拉保密信息。本邮件只能用于收件人，不得向收件人以外的人披露或者被其使用。
如果收到本邮件的不是收件人，我们特此通知您，严禁传播、分发或复制此邮件或附件。如果您错误地收到此邮件，请联系发件人，并立即删除此邮件。<br>
Notice : The information and attachments contained in this email are xxx’s confidential information.This email is 
intended solely for the use of the addressee.The contents may not be disclosed to or used by anyone other than the 
addressee. If you are not the intended recipient, we hereby inform you that you are strictly prohibited to disseminate, 
distribute or copy this message or its attachments.If you receive this email by mistake, please contact the sender and 
delete this email immediately.
`

func EmailContentAppendSignatureDefault(content string) string {
	return content + EmailContentSignature
}
