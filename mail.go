package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Service struct {
	conf Config
}

type Config struct {
	identity string
	username string
	password string
	hostname string
	port     string
}

func NewService(c Config) *Service {
	return &Service{
		conf: c,
	}
}

func (s *Service) Send(to string, msg []byte) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(s.conf.identity, s.conf.username, s.conf.password, s.conf.hostname)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.conf.hostname,
	}

	// call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", s.conf.hostname+":"+s.conf.port, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, s.conf.hostname)
	if err != nil {
		return err
	}

	// Auth
	if err := c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err := c.Mail(s.conf.username); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		defer func() {
			fmt.Println("panic", recover(), "is recovered")
		}()
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	c.Quit()

	return nil
}
