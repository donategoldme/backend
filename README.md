# Backend for service donategoldme

## Installation

Create empty directory. Clone this repository

```
git clone https://github.com/donategoldme/backend.git
```

add all enviroments to **env** file.
- [TWITCH_KEY & TWITCH_SECRET](https://www.twitch.tv/kraken/oauth2/clients/new) add new app
- [GPLUS_KEY & GPLUS_SECRET](https://console.developers.google.com/apis/credentials) app for 0auth
- [PEKA2TV_KEY](https://github.com/peka2tv/api/blob/master/oauth.md)
- [GOODGAME_KEY & GOODGAME_SECRET](https://api2.goodgame.ru/oauth/register) need registration on main [portal](https://goodgame.ru)
- CENTRIFUGO_SECRET - what you want
- [YANDEX_SECRET](https://money.yandex.ru/myservices/online.xml?)
- [YOUTUBE_KEY](https://console.developers.google.com/apis/credentials) standard app
- [SPEECHKIT_APIKEY](https://developer.tech.yandex.ru/keys)  get key for speechkit
- [PEKA2TV_TOKEN](http://peka2.tv) get from localStorage
- [GOODGAME_LOGIN & GOODGAME_PASS](https://goodgame.ru)
- [TWITCH_NAME & TWITCH_OAUTH](https://twitch.tv)

Run scripts

1. **run_first.sh**
2. **run_chats_service.sh**
3. **run_backend.sh**


## Depends
- [iris](https://github.com/kataras/iris)
- [gorm](https://github.com/jinzhu/gorm)
- [tarantool](github.com/tarantool/go-tarantool)
- [centrifugo](github.com/centrifugal/gocent)