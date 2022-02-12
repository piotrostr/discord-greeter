package bot

// initialize once per bot
// be careful not to initialize same bot with another ip
// its best to create token - proxy pairs, make sure ip is sticky
// proxy shall be in the form of
// username:password@host:port
//
// sample token
// OTQwNzg3MDY2MTU1OTYyNDA4.YgMetA.ILdepf9Gi1ehwCtzwEysLObpqbo

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/piotrostr/discord-greeter/headers"
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

type friendRequest struct {
	Username string `json:"username"`
	Discrim  int    `json:"discriminator"`
}

type invitePayload struct {
	CaptchaKey string `json:"captcha_key,omitempty"`
}

func (b *Bot) Initialize() error {
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

func (b *Bot) JoinServer(inviteCode string) error {
	var solvedKey string
	var payload invitePayload
	var err error
	for i := 0; i < b.Config.MaxRejoinAttempts; i++ {
		if solvedKey == "" || !b.Config.SolveCaptcha {
			payload = invitePayload{}
		} else {
			payload = invitePayload{
				CaptchaKey: solvedKey,
			}
		}
		payload, err := json.Marshal(payload)
		if err != nil {
			color.Red("error while marshalling payload %v", err)
			err = fmt.Errorf("error while marshalling payload %v", err)
			continue
		}
		url := "https://discord.com/api/v9/invites/" + inviteCode
		req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
		if err != nil {
			color.Red("Error while making http request %v \n", err)
			continue
		}

		cookie, err := b.GetCookieString()
		if err != nil {
			color.Red("[%v] Error while Getting cookies: %v", err)
			continue
		}
		req.Header.Set("Cookie", cookie)
		req = headers.Invite(req)
		resp, err := b.Client.Do(req)
		if err != nil {
			color.Red("Error while sending HTTP request %v \n", err)
			continue
		}

		body, err := ReadBody(*resp)
		if err != nil {
			color.Red("Error while reading body %v \n", err)
			continue
		}
		// TODO finish refactor
		/*

			if strings.Contains(string(body), "captcha_sitekey") {
				if in.Config.CaptchaAPI == "" {
					err = fmt.Errorf("[%v] Captcha detected but no API key provided", time.Now().Format("15:04:05"))
					break
				} else {
					color.Green("[%v] Captcha detected, solving...", time.Now().Format("15:04:05"))
				}
				var resp map[string]interface{}
				err = json.Unmarshal(body, &resp)
				if err != nil {
					color.Red("[%v] Error while Unmarshalling body: %v", time.Now().Format("15:04:05"), err)
					continue
				}
				solvedKey, err = in.SolveCaptcha(resp["captcha_sitekey"].(string))
				if err != nil {
					color.Red("[%v] Error while Solving Captcha: %v", time.Now().Format("15:04:05"), err)
					continue
				}
				if i == in.Config.MaxInvite-1 {
					i--
				}
			}

			var Join joinResponse
			err = json.Unmarshal(body, &Join)
			if err != nil {
				color.Red("Error while unmarshalling body %v %v\n", err, string(body))
				return err
			}
			if resp.StatusCode == 200 {
				color.Green("[%v] %v joint guild", time.Now().Format("15:04:05"), in.Token)
				if Join.VerificationForm {
					if len(Join.GuildObj.ID) != 0 {
						Bypass(in.Client, Join.GuildObj.ID, in.Token, Code)
					}
				}
			}
			if resp.StatusCode != 200 {
				color.Red("[%v] %v Failed to join guild %v", time.Now().Format("15:04:05"), resp.StatusCode, string(body))
			}
		*/
		return nil

	}
	return err
}

func (b *Bot) CheckServer(guildId string) (int, error) {
	url := "https://discord.com/api/v9/guilds/" + guildId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	req = headers.Common(req)
	req.Header.Set("Authorization", b.Token)

	res, err := b.Client.Do(req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}

func (b *Bot) CheckToken() int {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}
	req.Header.Set("authorization", b.Token)

	res, err := b.Client.Do(headers.Common(req))
	if err != nil {
		return -1
	}
	return res.StatusCode
}

func (b *Bot) Friend(username string, discrim int) (*http.Response, error) {
	url := "https://discord.com/api/v9/users/@me/relationships"
	fr := friendRequest{username, discrim}
	jsonx, err := json.Marshal(&fr)
	if err != nil {
		return &http.Response{}, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonx)))
	if err != nil {
		return &http.Response{}, err
	}
	cookie, err := b.GetCookieString()
	if err != nil {
		return &http.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}
	fingerprint, err := b.GetFingerprintString()
	if err != nil {
		return &http.Response{}, fmt.Errorf("error while getting fingerprint %v", err)
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-Fingerprint", fingerprint)
	req.Header.Set("Authorization", b.Token)

	res, err := b.Client.Do(headers.Common(req))
	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}

func (b *Bot) GetCookieString() (string, error) {
	url := "https://discord.com"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("[%v] Error while making request to get cookies %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while making request to get cookie %v", err)
	}

	res, err := b.Client.Do(req)
	if err != nil {
		color.Red("[%v] Error while getting resonse from cookies request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while getting resonse from cookie request %v", err)
	}
	defer res.Body.Close()

	if res.Cookies() == nil {
		color.Red("[%v] Error while getting cookies from resonse %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("there are no cookies in resonse")
	}
	var cookies string
	for _, cookie := range res.Cookies() {
		cookies = cookies + cookie.Name + "=" + cookie.Value + "; "
	}

	return cookies + "locale=en-US", nil
}

// Getting Fingerprint to use in our requests for more legitimate seeming requests.
func (b *Bot) GetFingerprintString() (string, error) {
	url := "https://discord.com/api/v9/experiments"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("[%v] Error while making request to get fingerprint %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while making request to get fingerprint %v", err)
	}
	res, err := b.Client.Do(headers.Register(req))
	if err != nil {
		color.Red("[%v] Error while getting resonse from fingerprint request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while getting resonse from fingerprint request %v", err)
	}

	p, err := ReadBody(*res)
	if err != nil {
		color.Red("[%v] Error while reading body from fingerprint request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while reading body %v", err)
	}

	var Response struct {
		Fingerprint string `json:"fingerprint"`
	}

	err = json.Unmarshal(p, &Response)

	if err != nil {
		color.Red("[%v] Error while unmarshalling body from fingerprint request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while unmarshalling resonse from fingerprint request %v", err)
	}

	return Response.Fingerprint, nil
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

func ReadBody(resp http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipreader, err := zlib.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		gzipbody, err := ioutil.ReadAll(gzipreader)
		if err != nil {
			return nil, err
		}
		return gzipbody, nil
	}

	if resp.Header.Get("Content-Encoding") == "br" {
		brreader := brotli.NewReader(bytes.NewReader(body))
		brbody, err := ioutil.ReadAll(brreader)
		if err != nil {
			fmt.Println(string(brbody))
			return nil, err
		}

		return brbody, nil
	}
	return body, nil
}
