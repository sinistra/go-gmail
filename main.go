package smtp_email

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"strings"
)

type SmtpServer struct {
	Host     string
	Port     string
	User     string
	Password string
}

func main() {

	godotenv.Load()
	port, _ := os.LookupEnv("SMTP_PORT")
	host, _ := os.LookupEnv("SMTP_HOST")
	user, _ := os.LookupEnv("SMTP_USER")
	password, _ := os.LookupEnv("SMTP_PASS")

	smtpServer := SmtpServer{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	// The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{"abc@gmail.com", "xyz@gmail.com", "larrypage@googlemail.com"}
	Subject := "Testing HTLML Email from golang"
	message := `
	<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
	<html>
	<head><meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1"></head>
	<body>This is the body<br>
	<div class="moz-signature">
	<i><br><br>
	Regards<br>
	Alex<br>
	<i></div>
	</body>
	</html>
	`
	bodyMessage := smtpServer.WriteHTMLEmail(Receiver, Subject, message)
	smtpServer.SendMail(Receiver, Subject, bodyMessage)
}

func (s SmtpServer) SendMail(Dest []string, Subject, bodyMessage string) {

	msg := "From: " + s.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(s.Host+":"+s.Port,
		smtp.PlainAuth("", s.User, s.Password, s.Host),
		s.User, Dest, []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return
	}

	fmt.Println("Mail sent successfully!")
}

func (s SmtpServer) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {

	header := make(map[string]string)
	header["From"] = s.User

	recipient := ""

	for _, user := range dest {
		recipient = recipient + user
	}

	header["To"] = recipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

func (s *SmtpServer) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {
	return s.WriteEmail(dest, "text/html", subject, bodyMessage)
}

func (s *SmtpServer) WritePlainEmail(dest []string, subject, bodyMessage string) string {
	return s.WriteEmail(dest, "text/plain", subject, bodyMessage)
}
