package bot

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/fatih/color"
	"github.com/piotrostr/discord-greeter/headers"
)

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
