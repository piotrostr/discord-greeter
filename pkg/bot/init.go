package bot

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

func (b *Bot) Initialize() error {
	b.ReadEnv()
	b.ReadConfig()
	b.ReadMessage()

	proxyUrl, err := url.Parse("http://" + b.Proxy)
	if err != nil {
		return fmt.Errorf("could not parse proxy: http://%s", b.Proxy)
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
	color.Red("Websocket closed %v %v", err, b.Token)
	/* TODO add error handling here, exit if needed
	in.Ws, err = in.NewConnection(in.wsFatalHandler)
	if err != nil {
		b.fatal <- fmt.Errorf("failed to create websocket connection: %s", err)
		return
	}
	*/
}
