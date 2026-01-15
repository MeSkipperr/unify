package mailer

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"strconv"
	"strings"
)

func loadSenderFromEnv() (Sender, string, int, error) {
	host :=  os.Getenv("SMTP_HOST")
	port, err :=  strconv.Atoi(	os.Getenv("SMTP_PORT"))

	email := os.Getenv("SMTP_EMAIL")
	pass :=  os.Getenv("SMTP_PASS")
	name :=  os.Getenv("SMTP_NAME")

	if host == "" || email == "" || pass == "" || name == "" || err != nil {
		return Sender{}, "", 0, fmt.Errorf("SMTP env not configured")
	}

	return Sender{
		Email: email,
		Pass:  pass,
		Name:  name,
	}, host, port, nil
}

func replaceTemplate(body string, firstName, lastName string) string {
	body = strings.ReplaceAll(body, "{{firstName}}", firstName)
	body = strings.ReplaceAll(body, "{{lastName}}", lastName)
	body = strings.ReplaceAll(body, "{{PROPERTY}}", os.Getenv("PROPERTY"))
	return body
}

func SendEmailSMTP(data EmailStructure) error {
	sender, smtpHost, smtpPort, err := loadSenderFromEnv()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth(
		"",
		sender.Email,
		sender.Pass,
		smtpHost,
	)

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	for _, r := range data.Recipients {

		// Safety
		firstName := ""
		lastName := ""
		if r.FirstName != nil {
			firstName = *r.FirstName
		}
		if r.LastName != nil {
			lastName = *r.LastName
		}

		// Replace template per user
		body := replaceTemplate(
			data.EmailData.BodyTemplate,
			firstName,
			lastName,
		)

		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		headers := map[string]string{
			"From":         fmt.Sprintf("%s <%s>", sender.Name, sender.Email),
			"To":           r.Email,
			"Subject":      data.EmailData.Subject,
			"MIME-Version": "1.0",
			"Content-Type": "multipart/mixed; boundary=" + writer.Boundary(),
		}

		for k, v := range headers {
			buf.WriteString(k + ": " + v + "\r\n")
		}
		buf.WriteString("\r\n")

		// Body
		bodyHeader := textproto.MIMEHeader{}
		bodyHeader.Set("Content-Type", "text/plain; charset=utf-8")
		bodyHeader.Set("Content-Transfer-Encoding", "quoted-printable")


		bodyPart, err := writer.CreatePart(bodyHeader)
		if err != nil {
			return err
		}

		bodyPart.Write([]byte(body))

		// Attachment (optional)
		for _, file := range data.EmailData.FileAttachment {
			if err := AddAttachment(writer, file); err != nil {
				return err
			}
		}

		writer.Close()

		// Send per recipient
		if err := smtp.SendMail(
			addr,
			auth,
			sender.Email,
			[]string{r.Email},
			buf.Bytes(),
		); err != nil {
			return err
		}
	}

	return nil
}
