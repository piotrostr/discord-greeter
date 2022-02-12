package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

func (b *Bot) ReadMessage() error {
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "message.json"))
	if err != nil {
		color.Red("error while opening message.json")
		return err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)

	errr := json.Unmarshal(bytes, &b.Message)
	if errr != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
		color.Red("error while opening config.json")
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

	inviteCode, inviteCodeExists := os.LookupEnv("INVITE_CODE")
	if !inviteCodeExists {
		return fmt.Errorf("invite missing")
	}

	captchaKey, captchaKeyExists := os.LookupEnv("CAPTCHA_KEY")
	if !captchaKeyExists {
		return fmt.Errorf("captchaKey missing")
	}

	b.Proxy = proxy
	b.Token = token
	b.InviteCode = inviteCode
	b.CaptchaKey = captchaKey

	return nil
}
