package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"frp-admin/config"
	"frp-admin/logger"
	template2 "html/template"
	"net"
	"net/smtp"
	"os"
)

var TemplateMap = make(map[string]string)

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

func SendTestMail(toMail string) {
	content := make(map[string]string)
	content["Title"] = "Test mail"
	content["Content"] = "This is a test mail."
	content["BtnLink"] = "https://www.google.com"
	content["BtnText"] = "Click to Google"
	content["Author"] = "ShawnGao"
	content["Note"] = "You received this email because you signed up for our services.If you did not, please ignore this email."
	body, err := FillTemplate("example-template", content)
	if err != nil {
		logger.LogErr("Send failed. %s", err)
		return
	}
	SendMail(toMail, "Test mail", body)
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
	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)
	if err != nil {
		logger.LogErr("%s", err)
	}
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

func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	c, err := Dial(addr)
	if err != nil {
		logger.LogErr("Create SMTP client error: %s", err)
		return err
	}
	defer func(c *smtp.Client) {
		err := c.Close()
		if err != nil {
			logger.LogErr("%s", err)
		}
	}(c)
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				logger.LogErr("Error during AUTH: %s", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
