package auth

import (
	"os"

	"donategold.me/centrifugo"
	"donategold.me/components"
	"donategold.me/components/auth/goodgame"
	"donategold.me/components/auth/gplus"
	"donategold.me/components/auth/peka2tv"
	"donategold.me/components/auth/providers"
	"donategold.me/components/auth/twitch"
	"donategold.me/components/users"

	"log"

	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	api.Post("/centrifugo", CentrifugeAuth)
	api.Get("/:provider/callback", AuthCallback)
	api.Get("/:provider", Auth)
}

var Providers = providers.Providers{}

func init() {
	Providers.Add(twitch.NewProvider(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "/auth/twitch/callback", []string{"user_subscriptions", "user_read", "chat_login", "channel_editor"}))
	Providers.Add(gplus.NewProvider(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), "/auth/gplus/callback", "https://www.googleapis.com/auth/youtube+https://www.googleapis.com/auth/youtube.force-ssl"))
	Providers.Add(peka2tv.NewProvider(os.Getenv("PEKA2TV_KEY"), "/auth/peka2tv/callback", ""))
	Providers.Add(goodgame.NewProvider(os.Getenv("GOODGAME_KEY"), os.Getenv("GOODGAME_SECRET"), "/auth/goodgame/callback", []string{"channel.subscribers", "channel.premiums", "channel.donations", "chat.token"}))
}

func AuthCallback(c *iris.Context) {
	// user, err := gothic.CompleteUserAuth(c)
	provider, ok := Providers.Get(c.Param("provider"))
	if !ok {
		log.Println("No provider in provider map")
		c.EmitError(404)
		return
	}
	token, err := provider.GetToken("http://"+c.Host(), c.URLParam(provider.QueryAuthCode()))
	if err != nil {
		log.Println(err.Error())
		c.JSON(501, "token error")
		return
	}

	if tokenUser, ok := c.Get("token").(users.AccessToken); ok && tokenUser.User.ID != 0 {
		userProvider, err := getOrCreateUserProvider(token, tokenUser.User.ID)
		if err != nil {
			log.Println(err)
			c.JSON(501, err.Error())
			return
		}
		err = publishAddProvider(userProvider)
		if err != nil {
			log.Println(err.Error())
			c.JSON(501, err.Error())
			return
		}
		c.JSON(201, "close window")
		return
	}
	username, err := token.GetUsernameUniq()
	if err != nil {
		log.Println(err.Error())
		c.EmitError(400)
		return
	}
	user, err := GetOrCreateUser(token.GetProviderName(), username,
		token.GetToken(), token.GetRefreshToken(), token.GetExpires())
	if err != nil {
		log.Println(err.Error())
		c.EmitError(404)
		return
	}
	userToken := users.NewTokenForUser(user)
	c.SetCookie(SetCookieToken(userToken.Token))

	timestamp, tokenCentrgo := components.GetTokenForCentrifugo(userToken.User.ID)
	userProviders := users.GetProvidersByUser(userToken.User.ID)
	centrifugo := map[string]string{"timestamp": timestamp, "token": tokenCentrgo}
	c.JSON(200, map[string]interface{}{"user": userToken, "centrifugo": centrifugo, "providers": userProviders})
}

func Auth(c *iris.Context) {
	//err := gothic.BeginAuthHandler(c)
	provider, ok := Providers.Get(c.Param("provider"))
	if !ok {
		c.JSON(404, "Provider not supported!")
		return
	}
	url, err := provider.GetCallbackUrl("http://" + c.Host())
	if err != nil {
		c.JSON(501, "error")
		return
	}
	// url, err := gothic.GetAuthURL(c)
	// if err != nil {
	// 	c.EmitError(400)
	// }
	log.Println(url)
	c.JSON(200, url)
}

func CentrifugeAuth(c *iris.Context) {
	u := c.Get("token").(users.AccessToken)
	if !u.IsAuthenticated() {
		c.EmitError(401)
		return
	}
	var channels centrifugo.ChannelsAuth
	c.ReadJSON(&channels)
	if len(channels.Channels) < 1 {
		c.EmitError(404)
		return
	}
	signs := centrifugo.ChannelsSigns{}
	for _, ch := range channels.Channels {
		if !centrifugo.CheckAllowChannel(ch) || !checkUserChannel(u.User.ID, ch) {
			c.EmitError(401)
			return
		}
		sign := centrifugo.ChannelSign{Sign: components.GenerateChannelSign(channels.Client, ch, "")}
		signs[ch] = sign
	}
	c.JSON(200, signs)
}
