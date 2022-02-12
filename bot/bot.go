package bot

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Timeout   int  `json:"timeout,omitempty"`
	DisableKL bool `json:"disable_kl,omitempty"`
}

type Bot struct {
	Token   string
	Proxy   string
	Message Message
	Client  *http.Client
	Config  Config
	fatal   chan error
}

type Message struct {
	Content string `json:"content,omitempty"`
	Author  User   `json:"author,omitempty"`
	GuildID string `json:"guild_id,omitempty"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (b *Bot) Initialize() error {
	// initialize once per bot
	// be careful not to initialize same bot with another ip
	// its best to create token - proxy pairs, make sure ip is sticky
	// proxy shall be in the form of
	// username:password@host:port
	//
	// sample token
	// OTQwNzg3MDY2MTU1OTYyNDA4.YgMetA.ILdepf9Gi1ehwCtzwEysLObpqbo
	proxy, proxyExists := os.LookupEnv("PROXY")
	token, tokenExists := os.LookupEnv("TOKEN")
	if !(proxyExists || tokenExists) {
		return fmt.Errorf("proxy or token missing")
	}
	b.Proxy = proxy
	b.Token = token

	proxyUrl, err := url.Parse("http://" + proxy)
	if err != nil {
		return fmt.Errorf("could not parse proxy: http://%s", proxy)
	}

	b.Client = &http.Client{
		Timeout: time.Second * time.Duration(b.Config.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					0x1301, 0x1303, 0x1302,
					0xc02b, 0xc02f, 0xcca9,
					0xcca8, 0xc02c, 0xc030,
					0xc00a, 0xc009, 0xc013,
					0xc014, 0x009c, 0x009d,
					0x002f, 0x0035,
				},
				InsecureSkipVerify: true,
				CurvePreferences: []tls.CurveID{
					tls.CurveID(0x001d), tls.CurveID(0x0017),
					tls.CurveID(0x0018), tls.CurveID(0x0019),
					tls.CurveID(0x0100), tls.CurveID(0x0101),
				},
			},
			DisableKeepAlives: b.Config.DisableKL,
			ForceAttemptHTTP2: true,
			Proxy:             http.ProxyURL(proxyUrl),
		},
	}

	return nil
}

func (b *Bot) FatalHandler(err error) {
	if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == 4004 {
		b.fatal <- fmt.Errorf("Authentication failed, try using a new token")
		return
	}
	fmt.Printf("Websocket closed %v %v", err, b.Token)
	/* TODO add error handling here, exit if needed
	in.Ws, err = in.NewConnection(in.wsFatalHandler)
	if err != nil {
		b.fatal <- fmt.Errorf("failed to create websocket connection: %s", err)
		return
	}
	*/
}
