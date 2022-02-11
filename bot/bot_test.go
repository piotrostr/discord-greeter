package bot_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/piotrostr/discord-greeter/bot"
)

var _ = Describe("Bot", func() {
	var bot *Bot

	BeforeEach(func() {
		bot = &Bot{}
		err := bot.Initialize()
		Expect(err).To(BeNil())
	})

	It("stands behind the proxy", func() {
		server := makeMockServer()
		defer server.Close()
		req, err := http.NewRequest("GET", server.URL, nil)
		Expect(err).To(BeNil())
		res, err := bot.Client.Do(req)
		Expect(err).To(BeNil())
		defer res.Body.Close()
		Expect(res).NotTo(BeNil())
	})
})

func makeMockServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				resp, _ := json.Marshal(map[string]string{
					"ip": GetIP(r),
				})
				w.Write(resp)
			}))
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
