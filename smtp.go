package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/url"
	"strings"
)

const (
	GoogleHost = "smtp.gmail.com"
)

// Sender is an interface for sending emails.
type Sender interface {
	// Send sends an email to the given recipients.
	Send(subject string, msg string, to ...string) error
	// SendTLS sends an email to the given recipients using TLS.
	SendTLS(subject string, msg string, to ...string) error
}

// Client is a client for sending emails.
// All fields are exported in-case, you want to set them manually.
type Client struct {
	// Username is client's email address Username
	Username string
	// Password is client's app Password for their email
	// currently supporting gmail.
	Password string
	// Hostname is the Hostname of the smtp server
	// currently supporting `smtp.gmail.com`.
	Hostname string
	// Port is the Port that the smtp server is listening on
	Port string
}

// NewClient returns a new Client based on the connection string
// which should be in the format:
//
//	`<scheme>://<username>:<password>@<host>:<port>`
func NewClient(connString string) (Sender, error) {
	u, err := url.Parse(connString)
	if err != nil {
		return nil, err
	}

	username := u.User.Username() + "@" + u.Hostname()

	password, _ := u.User.Password()
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	port := u.Port()

	switch u.Scheme {
	case "google":
		return &Client{username, password, GoogleHost, port}, nil
	// TODO: add support for other email providers
	// case "yahoo":
	// 	return NewYahooClient(url)
	// case "outlook":
	// 	return NewOutlookClient(url)
	default:
		return nil, fmt.Errorf("unknown scheme: %s", u.Scheme)
	}
}

// Deprecated: use SendTLS instead
func (c *Client) Send(subject, msg string, to ...string) error {
	var buf strings.Builder
	{
		// **** From: <email address> \r\n ****
		buf.WriteString(fmt.Sprintf("From: %s", c.Username))
		buf.WriteString("\r\n")
		// **** To: <email address> \r\n ****
		buf.WriteString(fmt.Sprintf("To: %s", strings.Join(to, ",")))
		buf.WriteString("\r\n")
		// **** Subject: <subject> \r\n ****
		buf.WriteString(fmt.Sprintf("Subject: %s", subject))
		buf.WriteString("\r\n")
		// **** header fields ****
		hdr := []string{"MIME-version: 1.0;", "Content-Type: text/html; charset=\"UTF-8\";"}
		buf.WriteString(strings.Join(hdr, "\r\n"))
		// **** \r\n\r\n ****
		buf.WriteString("\r\n")
		buf.WriteString("\r\n")
		// **** <message> ****
		buf.WriteString(msg)
	}
	msg = buf.String()

	addr := fmt.Sprintf("%s:%s", c.Hostname, c.Port)

	a := smtp.PlainAuth("", c.Username, c.Password, c.Hostname)

	return smtp.SendMail(addr, a, c.Username, to, []byte(msg))
}

func (c *Client) SendTLS(subject, msg string, to ...string) error {
	s, err := smtp.Dial(c.Hostname + ":" + c.Port)
	if err != nil {
		return err
	}
	defer s.Close()

	tls := &tls.Config{InsecureSkipVerify: true, ServerName: c.Hostname}
	if err = s.StartTLS(tls); err != nil {
		return err
	}

	a := smtp.PlainAuth("", c.Username, c.Password, c.Hostname)
	if err := s.Auth(a); err != nil {
		return err
	}

	if err := s.Mail(c.Username); err != nil {
		return err
	}

	for _, r := range to {
		if err := s.Rcpt(r); err != nil {
			return err
		}
	}

	w, err := s.Data()
	if err != nil {
		return err
	}

	var buf strings.Builder
	{
		// **** From: <email address> \r\n ****
		buf.WriteString(fmt.Sprintf("From: %s", c.Username))
		buf.WriteString("\r\n")
		// **** To: <email address> \r\n ****
		buf.WriteString(fmt.Sprintf("To: %s", strings.Join(to, ",")))
		buf.WriteString("\r\n")
		// **** Subject: <subject> \r\n ****
		buf.WriteString(fmt.Sprintf("Subject: %s", subject))
		buf.WriteString("\r\n")
		// **** header fields ****
		hdr := []string{"MIME-version: 1.0;", "Content-Type: text/html; charset=\"UTF-8\";"}
		buf.WriteString(strings.Join(hdr, "\r\n"))
		// **** \r\n\r\n ****
		buf.WriteString("\r\n")
		buf.WriteString("\r\n")
		// **** <message> ****
		buf.WriteString(msg)
	}
	msg = buf.String()

	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return s.Quit()
}
