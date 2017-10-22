package components

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
	"os"
	"fmt"
	"time"
)

var centrifugoSecret []byte = []byte(os.Getenv("CENTRIFUGO_SECRET"))

func GenerateChannelSign(client, channel, channelData string) string {
	sign := hmac.New(sha256.New, centrifugoSecret)
	sign.Write([]byte(client))
	sign.Write([]byte(channel))
	sign.Write([]byte(channelData))
	return hex.EncodeToString(sign.Sum(nil))
}

func GenerateClientToken(user, timestamp, info string) string {
	token := hmac.New(sha256.New, centrifugoSecret)
	token.Write([]byte(user))
	token.Write([]byte(timestamp))
	token.Write([]byte(info))
	return hex.EncodeToString(token.Sum(nil))
}

func GetTokenForCentrifugo(user uint) (string, string) {
	userId := fmt.Sprintf("%d", user)
	t := fmt.Sprintf("%d", time.Now().Unix())
	token := GenerateClientToken(userId, t, "")
	return t, token
}