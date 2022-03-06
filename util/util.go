package util

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Client = &http.Client{Timeout: 10 * time.Second}

func BasicSender(send string, email string, password string, basic string) int {
	//data, _ := json.Marshal(string{"?email="+email+"&password="+password+"?privileges="})

	//dataR := bytes.NewBuffer(data)
	params := url.Values{}
	params.Add("email", email)
	params.Add("password", password)
	params.Add("privileges", "")

	paramsR := params.Encode()
	paramsE := strings.NewReader(paramsR)
	resp, err := http.NewRequest("POST", send, paramsE)

	resp.Header.Add("Authorization", "Basic "+basic)
	resp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	respR, err := Client.Do(resp)

	defer respR.Body.Close()

	_, err = ioutil.ReadAll(respR.Body)
	if err != nil {
		panic(err)
	}

	return respR.StatusCode
}
