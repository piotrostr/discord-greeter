package bot

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type Config struct {
	Timeout               int  `json:"timeout,omitempty"`
	DisableKL             bool `json:"disable_kl,omitempty"`
	CaptchaSolvingEnabled bool `json:"captcha_solving_enabled,omitempty"`
	MaxRejoinAttempts     int  `json:"max_rejoin_attempts,omitempty"`
}

type Bot struct {
	GuildId    string
	Token      string
	Proxy      string
	Message    Message
	Client     *http.Client
	Config     Config
	CaptchaKey string
	fatal      chan error
}

type Message struct {
	Content string `json:"content,omitempty"`
	Author  User   `json:"author,omitempty"`
	GuildId string `json:"guild_id,omitempty"`
}

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type friendRequest struct {
	Username string `json:"username"`
	Discrim  int    `json:"discriminator"`
}

type invitePayload struct {
	CaptchaKey string `json:"captcha_key,omitempty"`
}

func (b *Bot) ReadConfig() error {
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "config.json"))
	if err != nil {
		color.Red("Error while Opening config.json")
		return err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &b.Config)
	if errr != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (b *Bot) ReadEnv() error {
	proxy, proxyExists := os.LookupEnv("PROXY")
	if !proxyExists {
		return fmt.Errorf("proxy missing")
	}

	token, tokenExists := os.LookupEnv("TOKEN")
	if !tokenExists {
		return fmt.Errorf("token missing")
	}

	guildId, guildIdExists := os.LookupEnv("GUILD")
	if !guildIdExists {
		return fmt.Errorf("guildId missing")
	}

	captchaKey, captchaKeyExists := os.LookupEnv("CAPTCHA_KEY")
	if !captchaKeyExists {
		return fmt.Errorf("captchaKey missing")
	}

	b.Proxy = proxy
	b.Token = token
	b.GuildId = guildId
	b.CaptchaKey = captchaKey

	return nil
}

func (b *Bot) Initialize() error {
	b.ReadEnv()

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
