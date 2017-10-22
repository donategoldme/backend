package goodgame

import "errors"

var (
	// Ошибка связи с сервером гудгейма
	ErrorBadResponseGoodgame = errors.New("Error goodgame response")
	ErrorReadResponse        = errors.New("Error read response goodgame")
	ErrorJsonUnmarshal       = errors.New("Error parse json to struct response goodgame")
	ErrorUserNotFound        = errors.New("Error: user not found in response goodgame")
)
