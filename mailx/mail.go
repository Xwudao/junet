package mailx

import (
	"context"
	"crypto/tls"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/Xwudao/junet/logx"
	"github.com/Xwudao/junet/shutdown"
)

var config = Config{
	host:           "localhost",
	port:           587,
	username:       "email",
	password:       "",
	keepAlive:      true,
	connectTimeout: time.Second * 10,
	sendTimeout:    time.Second * 10,
	msgChan:        make(chan *Message, 20),

	authentication: mail.AuthPlain,
	encryption:     mail.EncryptionNone,
	tlsConfig:      &tls.Config{InsecureSkipVerify: true},
}

type Message struct {
	From    string
	To      []string
	CC      []string
	Subject string
	//text body, html body can only set one
	TextBody string
	//text body, html body can only set one
	HtmlBody string
}

type Config struct {
	//system
	initial bool
	//client
	client  *mail.SMTPClient
	ctx     context.Context
	cancel  context.CancelFunc
	msgChan chan *Message
	//mail server
	host       string
	port       int
	username   string
	password   string
	encryption mail.Encryption

	//config
	authentication mail.AuthType
	keepAlive      bool
	connectTimeout time.Duration
	sendTimeout    time.Duration
	tlsConfig      *tls.Config
}

func (c *Config) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

type Opt func(*Config)

func SetTlsConfig(config *tls.Config) Opt {
	return func(c *Config) {
		c.tlsConfig = config
	}
}
func SetSendTimeout(tm time.Duration) Opt {
	return func(c *Config) {
		c.sendTimeout = tm
	}
}
func SetConnectTimeout(tm time.Duration) Opt {
	return func(c *Config) {
		c.connectTimeout = tm
	}
}
func SetKeepalive(keep bool) Opt {
	return func(c *Config) {
		c.keepAlive = keep
	}
}
func SetAuthentication(auth mail.AuthType) Opt {
	return func(c *Config) {
		c.authentication = auth
	}
}
func SetEncryption(enc mail.Encryption) Opt {
	return func(c *Config) {
		c.encryption = enc
	}
}
func SetPort(port int) Opt {
	return func(c *Config) {
		c.port = port
	}
}
func SetPassword(pass string) Opt {
	return func(c *Config) {
		c.password = pass
	}
}
func SetHost(host string) Opt {
	return func(c *Config) {
		c.host = host
	}
}

func Init(opts ...Opt) {
	for _, opt := range opts {
		opt(&config)
	}
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = config.host
	server.Port = config.port
	server.Username = config.username
	server.Password = config.password
	server.Encryption = config.encryption

	// Since v2.3.0 you can specified authentication type:
	// - PLAIN (default)
	// - LOGIN
	// - CRAM-MD5
	// server.Authentication = mail.AuthPlain

	// Variable to keep alive connection
	server.KeepAlive = config.keepAlive

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = config.connectTimeout

	// Timeout for send the data and wait respond
	server.SendTimeout = config.sendTimeout

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = config.tlsConfig

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		panic(err)
	}

	config.client = smtpClient
	config.ctx, config.cancel = context.WithCancel(context.Background())

	go task()
	config.initial = true
	shutdown.Add(&config)
}

func task() {
	for {
		select {
		case <-config.ctx.Done():
			return
		case m := <-config.msgChan:
			email := mail.NewMSG()
			//From Example <nube@example.com>
			email.SetFrom(m.From).
				AddTo(m.To...).
				AddCc(m.CC...).
				SetSubject(m.Subject)
			if m.HtmlBody != "" {
				email.SetBody(mail.TextHTML, m.HtmlBody)
			} else {
				email.SetBody(mail.TextPlain, m.TextBody)
			}

			// also you can add body from []byte with SetBodyData, example:
			// email.SetBodyData(mail.TextHTML, []byte(htmlBody))

			// add inline
			//email.Attach(&mail.File{FilePath: "/path/to/image.png", Name: "Gopher.png", Inline: true})

			// always check error after send
			if email.Error != nil {
				logx.Error(email.Error)
				return
			}
			err := email.Send(config.client)
			if err != nil {
				logx.Errorf("send email err: %s", err.Error())
				return
			}
		}
	}
}

func SendEmail(m *Message) {
	if !config.initial {
		panic("mail service not initial")
		return
	}
	config.msgChan <- m
}
