package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"frp-admin/config"
	"frp-admin/logger"
	"github.com/goccy/go-json"
	template2 "html/template"
	"net"
	"net/smtp"
	"os"
)

type MailContent struct {
	Title   string
	Content string
	BtnLink string
	BtnText string
	Author  string
	Note    string
}

var TemplateMap = make(map[string]string)

var DefaultFooterNote = "You received this email because you signed up for our services.If you did not, please ignore this email."

func InitEmailTemplate() {
	templateList := config.Conf.Mail.Template
	for _, template := range templateList {
		if _, err := os.Stat(template.Path); os.IsNotExist(err) {
			logger.LogErr("Template file [%s - %s] is not exist!", template.Name, template.Path)
			continue
		}
		TemplateMap[template.Name] = template.Path
		logger.LogSuccess("Mail Template [%s - %s] - Loaded.", template.Name, template.Path)
	}
}

func SendDefaultMail(toMail string, subject string, content *MailContent) {
	var contentMap = make(map[string]string)
	Struct2Map(content, &contentMap)
	body, err := FillTemplate("example-template", contentMap)
	if err != nil {
		logger.LogErr("Send failed. %s", err)
		return
	}
	SendMail(toMail, subject, body)
}

func FillTemplate(templateName string, data map[string]string) (string, error) {
	path := TemplateMap[templateName]
	template, err := template2.ParseFiles(path)
	if err != nil {
		logger.LogErr("Template file [%s - %s] file is not exist!", templateName, path)
		return "", err
	}
	buf := new(bytes.Buffer)
	err = template.Execute(buf, data)
	if err != nil {
		logger.LogErr("Fill error: %s", err)
		return "", err
	}
	return buf.String(), nil
}

func SendMail(toMail string, subject string, htmlBody string) {
	host := config.Conf.Mail.Host
	port := config.Conf.Mail.Port
	email := config.Conf.Mail.Mail
	password := config.Conf.Mail.Password
	toEmail := toMail
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s<%s>", config.Conf.Mail.NickName, config.Conf.Mail.FromMail)
	header["To"] = toEmail
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)
	SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)
}

func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		logger.LogErr("Dialing Error: %s", err)
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) {
	c, err := Dial(addr)
	if err != nil {
		logger.LogErr("Create SMTP client error: %s", err)
		return
	}
	isClosed := false
	defer func() {
		if !isClosed {
			err := c.Quit() // 优雅地关闭SMTP会话
			if err != nil {
				// 如果出现错误，仅记录它，不影响err的返回值，因为邮件可能已经发送
				logger.LogErr("Close SMTP client error: %s", err)
			}
		}
	}()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				logger.LogErr("Error during AUTH: %s", err)
				return
			}
		}
	}
	if err = c.Mail(from); err != nil {
		logger.LogErr("Error mail: %s", err)
		return
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			logger.LogErr("Error Rcpt: %s", err)
			return
		}
	}
	w, err := c.Data()
	if err != nil {
		logger.LogErr("Error data: %s", err)
		return
	}
	_, err = w.Write(msg)
	if err != nil {
		logger.LogErr("Error write msg: %s", err)
		return
	}
	err = w.Close()
	if err != nil {
		logger.LogErr("Error close: %s", err)
		return
	}
	err = c.Quit()
	if err != nil {
		logger.LogErr("Quit SMTP client error: %s", err)
		err = c.Close()
		if err != nil {
			logger.LogErr("Close SMTP client error: %s", err)
		}
	}
	// 设置标志表明连接已关闭
	isClosed = true
}

func Struct2Map(obj any, resultMap *map[string]string) {
	b, _ := json.Marshal(&obj)
	_ = json.Unmarshal(b, resultMap)
}
