# SMTP

This package wraps the `net/smtp` package in a simpler API

## Basic Usage

```go
package main

import (
    "github.com/hyphengolang/smtp"
)

var connString = "google://golang:password@google.com:547"

func main() {
    //...

    client := smtp.NewClient(connString)

    // Can send mail to recipients
    to := "send@to.com"

    subj := "Golang is awesome!!"

    msg := `
    <h1>This is a HTML Document</h1>
    <p>HTML is correctly parsed and formatted for the recipients</p>
    `

    err = client.Send(subj, msg, to)

    // Can choose to send securely instead
    err = client.SendTLS(subj, msg, to)

    //...
}
```

The connection string works similarly to a database connection string. This takes the form `<scheme>://<username>:<password>@<host>:<port>` where:

- **schema** is the provider name. currently only Google is configured with more to be added and tested at a later date.
- **username** is the client's first part of their email address 
- **password** is the client's app password. To configure a Google email address can find details [here](https://support.google.com/mail/answer/185833?hl=en-GB)
- **host** is the email provider domain. Currently `google.*` and google workspace emails are tested to work, with more to come.
- **port** is the port to communicate over.