package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/piotrostr/discord-greeter/pkg/captcha"
	"github.com/piotrostr/discord-greeter/pkg/headers"
)

func (b *Bot) JoinServer() error {
	var solvedKey string
	var payload invitePayload
	var err error

	for i := 0; i < b.Config.MaxRejoinAttempts; i++ {
		if solvedKey == "" || !b.Config.CaptchaSolvingEnabled {
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
		url := "https://discord.com/api/v9/invites/" + b.InviteCode
		fmt.Print(payload)
		fmt.Print(url)
		req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
		fmt.Printf("%v+", req)
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
		if strings.Contains(string(body), "captcha_sitekey") {
			if b.CaptchaKey == "" {
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
			solvedKey, err = captcha.Solve(resp["captcha_sitekey"].(string), b.CaptchaKey)
			if err != nil {
				color.Red("[%v] Error while Solving Captcha: %v", time.Now().Format("15:04:05"), err)
				continue
			}
		}
		var Join JoinResponse
		err = json.Unmarshal(body, &Join)
		if err != nil {
			color.Red("Error while unmarshalling body %v %v\n", err, string(body))
			return err
		}
		if resp.StatusCode == 200 {
			color.Green("[%v] %v joint guild", time.Now().Format("15:04:05"), b.Token)
			if Join.VerificationForm {
				if len(Join.GuildObj.ID) != 0 {
					Bypass(b.Client, Join.GuildObj.ID, b.Token, b.InviteCode)
				}
			}
		}
		if resp.StatusCode != 200 {
			color.Red(
				"[%v] %v Failed to join guild %v", time.Now().Format("15:04:05"), resp.StatusCode, string(body))
		}
		return nil

	}
	return err
}

func (b *Bot) SendFriendRequest(username string, discrim int) (*http.Response, error) {
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
