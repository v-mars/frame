package utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func UrlEncode(param map[string]string) string {
	params := url.Values{}
	for key, val := range param {
		params.Add(key, val)
	}
	return params.Encode()
}

func Request(url, method string, param interface{}, headers []map[string]string, client *http.Client,
	out interface{}) ([]byte, error) {
	var jsonBytes []byte
	var err error
	if param != nil {
		log.Debugf("request url[%s] param convert jsonBytes start", url)
		jsonBytes, err = json.Marshal(param)
		if err != nil {
			log.Errorf("request param convert jsonBytes err, %s", err)
			return nil, err
		}
	}

	//client := &http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Errorf("request new http request err, %s", err)
		return nil, err
	}

	//增加header选项
	//request.Header.Add("Content-Type", "application/json")
	for _, h := range headers {
		h := h
		request.Header.Add(h["key"], h["value"])
	}

	//处理返回结果
	log.Debugf("request url[%s] start", url)
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("request do err, %s", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("request Body.Close do err, %s", err)
		}
	}(resp.Body)

	log.Debugf("read data bytes start")
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read data bytes err, %s", err)
		return nil, err
	}
	//fmt.Println("request content:", string(content))
	err = json.Unmarshal(content, out)
	if err != nil {
		log.Errorf("request content convert interface err, %s", err)
		return content, err
	}
	return content, err
}
