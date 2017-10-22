package peka2tv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"donategold.me/components/auth/providers"
)

const pekaCodeURL = "http://peka2.tv/api/oauth/request"
const pekaOAuthURL = "http://funstream.tv/oauth/"
const pekaTokenURL = "http://peka2.tv/api/oauth/exchange"

const peka2tv = "peka2tv"
const queryAuthCode = "code"

var cl *http.Client = &http.Client{Timeout: 10 * time.Second}

func GetUserFromPeka(token string) {

}

type CodeResp struct {
	Code string `json:"code"`
}

type Peka2tvProvider struct {
	Key      string
	Callback string
	Scopes   string
}

func (p Peka2tvProvider) Name() string {
	return peka2tv
}

func (p Peka2tvProvider) QueryAuthCode() string {
	return queryAuthCode
}

func (p Peka2tvProvider) getCode() (code string, err error) {
	body := fmt.Sprintf(`{"key": "%s"}`, p.Key)
	req, err := http.NewRequest("POST", pekaCodeURL, bytes.NewBufferString(body))
	if err != nil {
		return
	}
	res, err := cl.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	var codeResp CodeResp
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println(string(b))
	err = json.Unmarshal(b, &codeResp)
	if err != nil {
		return
	}
	if codeResp.Code == "" {
		err = errors.New("No code")
		return
	}
	return codeResp.Code, nil
}

func (p Peka2tvProvider) GetCallbackUrl(host string) (string, error) {
	code, err := p.getCode()
	if err != nil {
		return "", err
	}
	return pekaOAuthURL + code + "?code=" + code, nil
}

func (p Peka2tvProvider) GetToken(host, code string) (providers.Tokener, error) {
	body := fmt.Sprintf(`{"code": "%s"}`, code)
	req, err := http.NewRequest("POST", pekaTokenURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, nil
	}
	res, err := cl.Do(req)
	token := Token{}
	if err != nil {
		return token, err
	}
	defer res.Body.Close()
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return token, err
	}
	log.Println(string(s))
	err = json.Unmarshal(s, &token)
	if err != nil {
		return token, err
	}
	if token.Token == "" {
		return token, errors.New("no token")
	}
	if token.User.Guest {
		return token, errors.New("User not authenticated on peka2.tv")
	}
	return token, nil
}

func (g Peka2tvProvider) RefreshToken(rt string) (providers.Tokener, error) {
	return Token{}, nil
}

func NewProvider(key, callback string, scopes string) Peka2tvProvider {
	t := Peka2tvProvider{key, callback, scopes}
	return t
}

type Token struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (t Token) GetToken() string {
	return t.Token
}

func (t Token) GetUsernameUniq() (string, error) {
	return t.User.Name, nil
}

func (t Token) GetProviderName() string {
	return peka2tv
}

func (t Token) GetExpires() *time.Time {
	return nil
}

func (t Token) GetRefreshToken() string {
	return ""
}

type User struct {
	ID    uint
	Guest bool   `json:"guest"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}
