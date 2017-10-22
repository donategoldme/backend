package public

type Youtube struct {
	Username string `json:"username"`
	VideoURL string `json:"video_url"`
	Nickname string `json:"nickname"`
	Money    int    `json:"money"`
}

type Standard struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
	Money    int    `json:"money"`
}
