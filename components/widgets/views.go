package widgets

import (
	"donategold.me/components/widgets/chats"
	"donategold.me/components/widgets/goals"
	"donategold.me/components/widgets/standard"
	"donategold.me/components/widgets/youtube"
	"gopkg.in/kataras/iris.v6"
)

func Concat(api *iris.Router) {
	goals.Concat(api.Party("/goals"))
	youtube.Concat(api.Party("/youtube"))
	standard.Concat(api.Party("/standard"))
	chats.Concat(api.Party("/chats"))
	// subscribers.Concat(api.Party("/subs"))
}
