package smtp_test

import (
	"os"
	"testing"

	"github.com/adoublef-go/smtp"
	"github.com/stretchr/testify/require"
)

var smtpUrl = os.Getenv("SMTP_URL")

func TestNewClient(t *testing.T) {

	c, err := smtp.NewClient(smtpUrl)
	require.NoError(t, err)
	
	to := "kristopherab@gmail.com"
	
	subject := "Dynamic HTML Email"
	msg := "<h1>Hello World, this is a new message!</h1>"
	
	// **** Send ****
	err = c.Send(subject, msg, to)
	require.NoError(t, err)
	
	// **** Send via TLS ****
	err = c.SendTLS(subject, msg, to)
	require.NoError(t, err)
}
