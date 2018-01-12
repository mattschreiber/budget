package email

import (
  "os"
  "fmt"
  "crypto/tls"
  "log"
  "net/smtp"
  "strings"
)

type Mail struct {
	SenderId string
	ToIds    []string
	Subject  string
	Body     string
}

type SmtpServer struct {
	host string
	port string
}


func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.SenderId)
	if len(mail.ToIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.ToIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	message += "\r\n" + mail.Body

	return message
}

func (mail *Mail) SendMail() {

  // mail := Mail{}
	// mail.senderId = "matt.schreiber01@gmail.com"
	// mail.toIds = []string{"matt.schreiber01@gmail.com"}
	// mail.subject = "New Ledger Entries"
	// mail.body = body

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

	//build an auth
	auth := smtp.PlainAuth("", mail.SenderId, os.Getenv("GMAIL_PASSWORD"), smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Println(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Println(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Println(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.SenderId); err != nil {
		log.Println(err)
	}
	for _, k := range mail.ToIds {
		if err = client.Rcpt(k); err != nil {
			log.Println(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Println(err)
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
	}

	client.Quit()
}
