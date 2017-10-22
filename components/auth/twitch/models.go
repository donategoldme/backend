package twitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"net/url"

	"donategold.me/components/auth/providers"
)

const twitch = "twitch"
const queryAuthCode = "code"

var cl *http.Client = &http.Client{Timeout: 10 * time.Second}

type TwitchUser struct {
	Identified bool `json:"identified"`
	Token      struct {
		Valid    bool   `json:"valid"`
		Username string `json:"user_name"`
	} `json:"token"`
}

func GetUserFromTwitch(token string) (TwitchUser, error) {
	var twuser TwitchUser
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/kraken", nil)
	req.Header.Add("Accept", `application/vnd.twitchtv.v3+json`)
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))
	res, err := cl.Do(req)
	if err != nil {
		return twuser, err
	}
	defer res.Body.Close()
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return twuser, nil
	}
	err = json.Unmarshal(s, &twuser)
	if err != nil {
		return twuser, err
	}
	return twuser, nil
}

type TwitchProvider struct {
	Key      string
	Secret   string
	Callback string
	Scopes   []string
}

func (t TwitchProvider) Name() string {
	return twitch
}

func (t TwitchProvider) QueryAuthCode() string {
	return queryAuthCode
}
func (t TwitchProvider) getFormatUrl() string {
	return `https://api.twitch.tv/kraken/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s
     &state=`
}
func (t TwitchProvider) getScopesString() string {
	return strings.Join(t.Scopes, " ")
}
func (t TwitchProvider) GetCallbackUrl(host string) (string, error) {
	return fmt.Sprintf(t.getFormatUrl(), t.Key, host+t.Callback, t.getScopesString()), nil
}

func (t TwitchProvider) getFullCallbackUrl(host string) string {
	return host + t.Callback
}

func (t TwitchProvider) getGrantFlowUrl() string {
	return `https://api.twitch.tv/kraken/oauth2/token`
}

func (t TwitchProvider) getGrantFlowFormat(host, code string) string {
	v := url.Values{}
	v.Add("client_id", t.Key)
	v.Add("client_secret", t.Secret)
	v.Add("grant_type", "authorization_code")
	v.Add("redirect_uri", t.getFullCallbackUrl(host))
	v.Add("code", code)
	v.Add("state", "")
	return v.Encode()
}

func (t TwitchProvider) GetToken(host, code string) (providers.Tokener, error) {
	body := t.getGrantFlowFormat(host, code)
	req, _ := http.NewRequest("POST", "https://api.twitch.tv/kraken/oauth2/token", bytes.NewBufferString(body))
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

func (g TwitchProvider) RefreshToken(rt string) (providers.Tokener, error) {
	return Token{}, nil
}

func NewProvider(key, secret, callback string, scopes []string) TwitchProvider {
	t := TwitchProvider{key, secret, callback, scopes}
	return t
}

type Token struct {
	Token string
}

func (t Token) GetToken() string {
	return t.Token
}

func (t Token) GetUsernameUniq() (string, error) {
	twUser, err := GetUserFromTwitch(t.Token)
	if err != nil {
		return "", err
	}
	return twUser.Token.Username, nil
}

func (t Token) GetProviderName() string {
	return twitch
}

func (t Token) GetRefreshToken() string {
	return ""
}

func (t Token) GetExpires() *time.Time {
	return nil
}
