package bot

import (
	"net/http"
)

type Config struct {
	Timeout               int  `json:"timeout,omitempty"`
	DisableKL             bool `json:"disable_kl,omitempty"`
	CaptchaSolvingEnabled bool `json:"captcha_solving_enabled,omitempty"`
	MaxRejoinAttempts     int  `json:"max_rejoin_attempts,omitempty"`
}

type Bot struct {
	Token      string
	Proxy      string
	InviteCode string
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

type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JoinResponse struct {
	VerificationForm bool  `json:"show_verification_form"`
	GuildObj         Guild `json:"guild"`
}

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
