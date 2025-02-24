package services

import (
	"bytes"
	"context"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"html/template"
	"net/smtp"
)

type EmailService struct {
	from     string
	password string
	host     string
	port     string
}

func NewEmailService(from, password, host, port string) api.EmailApi {
	return &EmailService{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (service *EmailService) Send(ctx context.Context, receiver, subject, text string) error {
	auth := smtp.PlainAuth("", service.from, service.password, service.host)

	message := "From: " + service.from + "\n" +
		"To: " + receiver + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		text

	return smtp.SendMail(service.host+":"+service.port, auth, service.from, []string{receiver}, []byte(message))
}

func (service *EmailService) SendLoginEmail(ctx context.Context, receiver, device, datetime, code string) error {
	auth := smtp.PlainAuth("", service.from, service.password, service.host)

	registerTemplate, err := template.ParseFiles("templates/login.html")
	if err != nil {
		panic(err)
	}

	data := struct {
		Device   string
		DateTime string
		Code     string
	}{device, datetime, code}

	buf := new(bytes.Buffer)
	if err := registerTemplate.Execute(buf, data); err != nil {
		return err
	}

	message := "From: " + service.from + "\n" +
		"To: " + receiver + "\n" +
		"Subject: Подтверждение входа — " + code + "\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		string(buf.Bytes())

	return smtp.SendMail(service.host+":"+service.port, auth, service.from, []string{receiver}, []byte(message))
}

func (service *EmailService) SendRegisterEmail(ctx context.Context, receiver, device, datetime, code string) error {
	auth := smtp.PlainAuth("", service.from, service.password, service.host)

	registerTemplate, err := template.ParseFiles("templates/register.html")
	if err != nil {
		panic(err)
	}

	data := struct {
		Device   string
		DateTime string
		Code     string
	}{device, datetime, code}

	buf := new(bytes.Buffer)
	if err := registerTemplate.Execute(buf, data); err != nil {
		return err
	}

	message := "From: " + service.from + "\n" +
		"To: " + receiver + "\n" +
		"Subject: Создание аккаунта — " + code + "\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		string(buf.Bytes())

	return smtp.SendMail(service.host+":"+service.port, auth, service.from, []string{receiver}, []byte(message))
}
