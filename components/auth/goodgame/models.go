package goodgame

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"net/url"

	"donategold.me/components/auth/providers"
)

const ggURL = "https://api2.goodgame.ru/oauth/authorize"
const ggTokenReqURL = "https://api2.goodgame.ru/oauth"
const ggUserInfo = "https://api2.goodgame.ru/info"

const ggName = "goodgame"
const queryAuthCode = "code"

type client struct{ http.Client }

func (c client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.Do(req)
}

var cl = &client{http.Client{Timeout: 10 * time.Second}}

type ATReq struct {
	RedirectUri  string `json:"redirect_uri,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	GrantType    string `json:"grant_type"`
}

func (a ATReq) JsonString() ([]byte, error) {
	return json.Marshal(a)
}

type GoodgameProvider struct {
	Key      string
	Secret   string
	Callback string
	Scopes   []string
}

func (g GoodgameProvider) Name() string {
	return ggName
}

func (g GoodgameProvider) QueryAuthCode() string {
	return queryAuthCode
}

func (g GoodgameProvider) getScopeString() (scopes string) {
	return strings.Join(g.Scopes, " ")
}

func (g GoodgameProvider) GetCallbackUrl(host string) (string, error) {
	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", g.Key)
	values.Set("redirect_uri", host+g.Callback)
	values.Set("scope", g.getScopeString())
	values.Set("state", "fuStateRequired")
	return ggURL + "?" + values.Encode(), nil
}

func (g GoodgameProvider) requestDo(body []byte) (providers.Tokener, error) {
	req, err := http.NewRequest("POST", ggTokenReqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil
	}
	res, err := cl.do(req)
	token := Token{}
	if err != nil {
		return token, err
	}
	defer res.Body.Close()
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return token, err
	}
	log.Println(string(body))
	log.Println(string(s))
	err = json.Unmarshal(s, &token)
	if err != nil {
		return token, err
	}
	if token.AccessToken == "" {
		return token, errors.New("no token")
	}
	return token, nil
}

func (g GoodgameProvider) GetToken(host, code string) (providers.Tokener, error) {
	body, _ := ATReq{RedirectUri: host + g.Callback, ClientID: g.Key,
		ClientSecret: g.Secret, GrantType: "authorization_code", Code: code}.JsonString()
	return g.requestDo(body)
}

func (g GoodgameProvider) RefreshToken(rt string) (providers.Tokener, error) {
	if rt == "" {
		return nil, errors.New("refresh token required")
	}
	body, _ := ATReq{RefreshToken: rt, ClientID: g.Key, GrantType: "refresh_token", ClientSecret: g.Secret}.JsonString()
	return g.requestDo(body)
}

func NewProvider(key, secret, callback string, scopes []string) GoodgameProvider {
	t := GoodgameProvider{key, secret, callback, scopes}
	return t
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (t Token) GetToken() string {
	return t.AccessToken
}

func (t Token) GetUsernameUniq() (string, error) {
	user, err := GetUserFromGoodgame(t.AccessToken)
	return user.User.Username, err
}

func (t Token) GetProviderName() string {
	return ggName
}

func (t Token) GetExpires() *time.Time {
	var tt *time.Time
	if t.ExpiresIn > 0 {
		tn := time.Now().Add(time.Second * time.Duration(t.ExpiresIn))
		tt = &tn
	}
	return tt
}

func (t Token) GetRefreshToken() string {
	return t.RefreshToken
}

type User struct {
	User struct {
		Username string `json:"username"`
	} `json:"user"`
}

func GetUserFromGoodgame(token string) (User, error) {
	req, err := http.NewRequest("GET", ggUserInfo, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := cl.do(req)
	if err != nil {
		return User{}, ErrorBadResponseGoodgame
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return User{}, ErrorReadResponse
	}
	log.Println(string(b))
	var user User
	err = json.Unmarshal(b, &user)
	if err != nil {
		return user, ErrorJsonUnmarshal
	}
	if user.User.Username == "" {
		return user, ErrorUserNotFound
	}
	return user, nil
}
