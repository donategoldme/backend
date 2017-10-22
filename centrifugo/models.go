package centrifugo

type ChannelsAuth struct {
    Channels []string `json:"channels"`
    Client string
}

type ChannelSign struct {
    Sign string `json:"sign"`
    Info string `json:"info"`
}

type ChannelsSigns map[string]ChannelSign