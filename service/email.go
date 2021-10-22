package service

import (
	"gopkg.in/gomail.v2"
	"log"
)


// MailboxConf 邮箱配置
type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

func SendEmail() {
	var mailConf MailboxConf
	mailConf.Title = "测试用gomail发送邮件"
	mailConf.Body = "Good Good Study, Day Day Up!!!!!!"
	mailConf.RecipientList = []string{`xxx@qq.com`}
	mailConf.Sender = `xxx@qq.com`
	mailConf.SPassword = "xxx"
	mailConf.SMTPAddr = `smtp.qq.com`
	mailConf.SMTPPort = 25

	m := gomail.NewMessage()
	m.SetHeader(`From`, mailConf.Sender)
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, mailConf.Body)
	m.Attach("./Dockerfile")   //添加附件
	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		log.Fatalf("Send Email Fail, %s", err.Error())
		return
	}
	log.Printf("Send Email Success")
}

