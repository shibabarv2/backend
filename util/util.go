package util

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"shiba-backend/structs"
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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func GetDomains() []string {
	key := os.Getenv("ADMIN_KEY")

	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("API_URL")+"/admin/mail/users?format=json", nil)

	if err != nil {
		return nil
	}

	req.Header.Add("Authorization", "Basic "+key)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	var r structs.StatsResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil
	}

	var domains []string

	for _, v := range r {
		if Contains(domains, v.Domain) == false {
			domains = append(domains, v.Domain)
		}
	}

	return domains
}
