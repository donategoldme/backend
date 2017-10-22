package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"strings"

	"fmt"

	"donategold.me/components/auth/providers"
	"donategold.me/components/users"
)

func getOrCreateUserProvider(token providers.Tokener, userID uint) (users.Provider, error) {
	username, err := token.GetUsernameUniq()
	if err != nil {
		return users.Provider{}, err
	}
	p := users.Provider{Uid: username, UserID: userID,
		TypeProvider: token.GetProviderName(), AccessToken: token.GetToken(),
		RefreshToken: token.GetRefreshToken(), Expires: token.GetExpires()}
	p.Create()
	return p, nil
}

func GetOrCreateUser(provider, username, token, refreshToken string, expires *time.Time) (users.User, error) {
	var user users.User
	p, exist := GetProviderByTypeAndUid(provider, username)
	if !exist {
		user = users.User{Username: username}
		err := user.CreateUniq()
		if err != nil {
			return user, err
		}
		p = users.Provider{Uid: username, UserID: user.ID, TypeProvider: provider,
			AccessToken: token, RefreshToken: refreshToken, Expires: expires}
		p.Create()
	} else {
		log.Println(token, expires, refreshToken)
		p.AccessToken = token
		p.Expires = expires
		p.RefreshToken = refreshToken
		p.Save()
		user, _ = users.GetUserById(p.UserID)
	}
	if user.ID == 0 {
		// p.Delete()
		return user, errors.New("Cant create user")
	}

	if p.ID == 0 {
		return user, errors.New("Cant create provider")
	}
	return user, nil
}

func SetCookieToken(value string) *http.Cookie {
	c := &http.Cookie{}
	//	c := &fasthttp.Cookie{}

	c.Name = "token"
	c.Value = value
	c.HttpOnly = true
	c.Path = "/"
	c.Expires = time.Now().Add(time.Duration(24*10) * time.Hour)
	return c
}

func checkUserChannel(id uint, ch string) bool {
	channelUser := strings.Split(ch, "/")[0]
	return channelUser == fmt.Sprintf(`$%d`, id)
}

func RefreshToken(provider *users.Provider) (*users.Provider, error) {
	if provider.RefreshTokenNeed() {
		prov, ok := Providers.Get(provider.TypeProvider)
		if !ok {
			return provider, errors.New("provider not supported")
		}
		token, err := prov.RefreshToken(provider.AccessToken)
		if err != nil {
			return provider, err
		}
		provider.AccessToken = token.GetToken()
		provider.Expires = token.GetExpires()
		provider.RefreshToken = token.GetRefreshToken()
		provider.Uid, err = token.GetUsernameUniq()
		if err != nil {
			return provider, err
		}
		provider.Save()
	}
	return provider, nil
}
