package gplus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"donategold.me/components/auth/providers"
)

const youtube = "gplus"
const queryCode = "code"

const urlDetails string = "https://www.googleapis.com/youtube/v3/channels?part=snippet&mine=true&access_token="

var cl *http.Client = &http.Client{}

type YoutubeUser struct {
	Items []Channel    `json:"items"`
	Error YoutubeError `json:"error"`
}

type Channel struct {
	Id      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	Title string `json:"title"`
}

type YoutubeError struct {
	Code int `json:"code"`
}

func GetUserFromGoggle(token string) (YoutubeUser, error) {
	var ytuser YoutubeUser
	req, _ := http.NewRequest("GET", urlDetails+token, nil)
	res, err := cl.Do(req)
	if err != nil {
		return ytuser, err
	}
	defer res.Body.Close()
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ytuser, err
	}
	err = json.Unmarshal(s, &ytuser)
	if err != nil {
		return ytuser, err
	}
	return ytuser, nil
}

type YoutubeProvider struct {
	ClientID string
	Secret   string
	Callback string
	Scopes   string
}

func (y YoutubeProvider) Name() string {
	return youtube
}

func (y YoutubeProvider) getFormatUrl() string {
	return `https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&access_type=offline&include_granted_scopes=true&redirect_uri=%s&response_type=code&scope=%s`
}

func (y YoutubeProvider) GetCallbackUrl(host string) (string, error) {
	return fmt.Sprintf(y.getFormatUrl(), y.ClientID, host+y.Callback, y.Scopes), nil
}

func (y YoutubeProvider) getValuesCode(host, code string) string {
	v := url.Values{}
	v.Add("client_id", y.ClientID)
	v.Add("client_secret", y.Secret)
	v.Add("grant_type", "authorization_code")
	v.Add("redirect_uri", host+y.Callback)
	v.Add("code", code)
	return v.Encode()
}

func (y YoutubeProvider) GetToken(host, code string) (providers.Tokener, error) {
	body := y.getValuesCode(host, code)
	req, _ := http.NewRequest("POST", "https://www.googleapis.com/oauth2/v4/token", bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
	m := make(map[string]interface{})
	err = json.Unmarshal(s, &m)
	if err != nil {
		return token, err
	}
	if m["access_token"] == "" {
		return token, errors.New("no token")
	}
	token.Token, _ = m["access_token"].(string)
	return token, nil
}

func (y YoutubeProvider) QueryAuthCode() string {
	return queryCode
}

func (g YoutubeProvider) RefreshToken(rt string) (providers.Tokener, error) {
	return Token{}, nil
}

type Token struct {
	Token string
}

func (t Token) GetToken() string {
	return t.Token
}

func (t Token) GetUsernameUniq() (string, error) {
	user, err := GetUserFromGoggle(t.Token)
	if err != nil {
		return "", err
	}
	return user.Items[0].Id, nil
}

func (t Token) GetProviderName() string {
	return youtube
}

func (t Token) GetExpires() *time.Time {
	return nil
}

func (t Token) GetRefreshToken() string {
	return ""
}

func NewProvider(clientID, secret, callback string, scopes string) YoutubeProvider {
	t := YoutubeProvider{clientID, secret, callback, scopes}
	return t
}
