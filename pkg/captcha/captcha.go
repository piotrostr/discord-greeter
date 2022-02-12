package captcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/piotrostr/discord-greeter/pkg/headers"
)

var captchaBaseUrl string = "https://api.capmonster.cloud/"

type Payload struct {
	ClientKey string `json:"clientKey"`
	Task      Task   `json:"task"`
	ErrorID   int    `json:"ErrorId"`
	TaskID    int    `json:"taskId"`
}

type Task struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	UserAgent  string `json:"userAgent"`
}

type Response struct {
	TaskID   int      `json:"taskID"`
	ErrorID  int      `json:"ErrorId"`
	Status   string   `json:"status"`
	Solution Solution `json:"solution"`
}

type Solution struct {
	Answer string `json:"gRecaptchaResponse"`
}

// Function to use a captcha solving service and return a solved captcha key
func Solve(siteKey string, clientKey string) (string, error) {
	jsonx := Payload{
		ClientKey: clientKey,
		Task: Task{
			Type:       "HCaptchaTaskProxyless",
			WebsiteURL: "https://discord.com/channels/@me",
			WebsiteKey: siteKey,
			UserAgent:  headers.UserAgent,
		},
	}

	bytes, err := json.Marshal(jsonx)
	if err != nil {
		return "", fmt.Errorf("error marshalling json [%v]", err)
	}
	res, err := http.Post(
		captchaBaseUrl+"/createTask",
		"application/json",
		strings.NewReader(string(bytes)))
	if err != nil {
		return "", fmt.Errorf("error creating the request for captcha [%v]", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the resonse body [%v]", err)
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling the resonse body [%v]", err)
	}
	switch response.ErrorID {
	case 0:
		// Poling server for the solved captcha
		jsonx = Payload{
			ClientKey: clientKey,
			TaskID:    response.TaskID,
		}
		y, err := json.Marshal(jsonx)
		if err != nil {
			return "", fmt.Errorf("error marshalling json [%v]", err)
		}
		// anti captcha documentation prescribes to use a delay of
		// 5 seconds before requesting the captcha and 3 seconds delays after that.
		time.Sleep(5 * time.Second)
		p := 0
		for {
			if p > 50 {
				// max retries
				break
			}
			res, err := http.Post(
				captchaBaseUrl+"getTaskResult",
				"application/json",
				strings.NewReader(string(y)))
			if err != nil {
				return "", fmt.Errorf("error creating the request for captcha [%v]", err)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return "", fmt.Errorf("error reading the resonse body [%v]", err)
			}
			var resonse Response
			err = json.Unmarshal(body, &resonse)
			if err != nil {
				return "", fmt.Errorf("error unmarshalling the resonse body [%v]", err)
			}
			if resonse.ErrorID == 16 {
				return "", fmt.Errorf("error getting captcha [%v]", resonse.ErrorID)
			}
			if resonse.Status == "ready" {
				return resonse.Solution.Answer, nil
			} else if resonse.Status == "processing" {
				p++ // Incrementing the counter
				time.Sleep(3 * time.Second)
			}

		}
		return "", fmt.Errorf("max captcha retries reached [%v]", err)
	case 2:
		color.Red("no available captcha workers. Sleeping 10 seconds")
		time.Sleep(10 * time.Second)
		return "", fmt.Errorf("no captcha workers were available, retrying")
	case 3:
		return "", fmt.Errorf("captcha you are uploading is less than 100 bytes.")
	case 4:
		return "", fmt.Errorf("captcha you are uploading is greater than 500,000 bytes.")

	case 10:
		return "", fmt.Errorf("you have zero or negative captcha API balance")
	case 11:
		return "", fmt.Errorf("captcha was unsolvable.")
	default:
		return "", fmt.Errorf("unknown error [%v]", response.ErrorID)
	}
}
