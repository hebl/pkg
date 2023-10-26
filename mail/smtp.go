package mail

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"strings"
)

// SMTP Client
type SMTPClient struct {
	Server   string
	Port     string
	From     string
	Password string
	SSL      bool // SSL or not
}

// EmailData
type Data struct {
	Subject string
	Body    string
}

func NewSMTPClient(server, port, from, passwd string, ssl bool) *SMTPClient {
	return &SMTPClient{
		Server:   server,
		Port:     port,
		From:     from,
		Password: passwd,
		SSL:      ssl,
	}
}

func (m *SMTPClient) SendMail(to []string, data Data) error {
	if m.SSL {
		return m.sendMailSSL(to, data)
	}
	return m.sendMail(to, data)
}

// Normal send
func (m *SMTPClient) sendMail(to []string, data Data) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Server)

	message := "From: " + m.From + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + data.Subject + "\n\n" +
		data.Body

	return smtp.SendMail(m.Server+":"+m.Port, auth, m.From, to, []byte(message))
}

// SSL 发送邮件
func (m *SMTPClient) sendMailSSL(to []string, data Data) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Server)
	//
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.Server,
	}

	// 建立连接
	conn, err := tls.Dial("tcp", m.Server+":"+m.Port, tlsConfig)
	if err != nil {

		return errors.New("Dial TCP connect error")
	}

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, m.Server)
	if err != nil {
		return errors.New("Create SMTP client error")
	}

	// 使用认证登录
	err = client.Auth(auth)
	if err != nil {
		return errors.New("Login auth error")
	}

	// 设置发件人
	err = client.Mail(m.From)
	if err != nil {
		return errors.New("Set from error")
	}

	// 设置多个收件人
	for _, toEmail := range to {
		err = client.Rcpt(toEmail)
		if err != nil {
			return errors.New("Set mailto error")
		}
	}

	message := "From: " + m.From + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + data.Subject + "\n\n" +
		data.Body

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return errors.New("Send mail content error")
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return errors.New("Write mail content error")
	}

	// 关闭邮件内容写入
	err = writer.Close()

	if err != nil {
		return errors.New("Close mail content error")
	}

	// 退出客户端
	client.Quit()

	// 关闭连接
	conn.Close()

	return nil
}
