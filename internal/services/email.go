package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
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

	fmt.Println(message)

	return smtp.SendMail(service.host+":"+service.port, auth, service.from, []string{receiver}, []byte(message))
}

func (service *EmailService) SendRegisterMail(ctx context.Context, receiver string, tempAccount *domain.TempAccount) error {
	auth := smtp.PlainAuth("", service.from, service.password, service.host)

	registerTemplate, err := template.ParseFiles("templates/register.html")
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	if err := registerTemplate.Execute(buf, tempAccount); err != nil {
		return err
	}

	message := "From: " + service.from + "\n" +
		"To: " + receiver + "\n" +
		"Subject: Complete register\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		string(buf.Bytes())

	return smtp.SendMail(service.host+":"+service.port, auth, service.from, []string{receiver}, []byte(message))
}
