package bot

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/fatih/color"
	"github.com/piotrostr/discord-greeter/headers"
)

type bypassInformation struct {
	Version    string      `json:"version"`
	FormFields []FormField `json:"form_fields"`
}

type FormField struct {
	FieldType   string   `json:"field_type"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Values      []string `json:"values"`
	Response    bool     `json:"response"`
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

func Bypass(client *http.Client, serverid string, token string, invite string) error {
	// First we require to get all the rules to send in the request
	site := "https://discord.com/api/v9/guilds/" + serverid + "/member-verification?with_guild=false&invite_code=" + invite
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		return err
	}
	req = headers.Common(req)
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ReadBody(*resp)
	if err != nil {
		return err
	}

	var bypassInfo bypassInformation
	err = json.Unmarshal(body, &bypassInfo)
	if err != nil {
		return err
	}

	// Now we have all the rules, we can send the request along with our response
	for i := 0; i < len(bypassInfo.FormFields); i++ {
		// We set the response to true because we accept the terms as the good TOS followers we are
		bypassInfo.FormFields[i].Response = true
	}

	jsonData, err := json.Marshal(bypassInfo)
	if err != nil {
		return err
	}
	url := "https://discord.com/api/v9/guilds/" + serverid + "/requests/@me"

	req, err = http.NewRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		color.Red("Error while making http request %v \n", err)
		return err
	}

	req.Header.Set("Authorization", token)
	resp, err = client.Do(headers.Common(req))
	if err != nil {
		color.Red("Error while sending HTTP request bypass %v \n", err)
		return err
	}
	body, err = ReadBody(*resp)
	if err != nil {
		color.Red("[%v] Error while reading body %v \n", time.Now().Format("15:04:05"), err)
		return err
	}

	if resp.StatusCode == 201 || resp.StatusCode == 204 {
		color.Green("[%v] Successfully bypassed token %v", time.Now().Format("15:04:05"), token)
	} else {
		color.Red("[%v] Failed to bypass Token %v %v %v", time.Now().Format("15:04:05"), token, resp.StatusCode, string(body))
	}
	return nil
}

func Snowflake() int64 {
	snowflake := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(snowflake, 2, 64)
	return nonce
}
