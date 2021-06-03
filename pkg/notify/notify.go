package notify

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Sender it send notify to user
type Sender interface {
	Send(to []string, title string, content string) error
}

// alarm notify
// mail
// chat
// dingding
// slack
// telegram
// server jiang
// lark

//JSONPost Post req json data to url
func JSONPost(method, url string, data interface{}, client *http.Client) ([]byte, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		//log.Error("client.Do",err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
