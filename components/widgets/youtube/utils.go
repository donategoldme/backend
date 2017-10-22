package youtube

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

var r = regexp.MustCompile(`PT((?P<hours>\d+)H)?((?P<minutes>\d+)M)?((?P<seconds>\d+)S)?`)
var cl *http.Client = &http.Client{}
var youtubeKey string = os.Getenv("YOUTUBE_KEY")

const urlVideoDetailsPartOne string = "https://content.googleapis.com/youtube/v3/videos?id="
const urlVideoDetailsPartTwo string = "&part=snippet%2C%20statistics%2C%20contentDetails&key="

type VideoDetails struct {
	Items []VideoDetail `json:"items"`
}

type VideoDetail struct {
	ID      string `json:"id"`
	Snippet struct {
		Title                string `json:"title"`
		LiveBroadcastContent string `json:"liveBroadcastContent"`
	} `json:"snippet"`
	ContentDetails struct {
		Duration   string `json:"duration"`
		Definition string `json:"definition"`
	} `json:"contentDetails"`
	Statistics struct {
		ViewCount    string `json:"viewCount"`
		LikeCount    string `json:"likeCount"`
		CommentCount string `json:"commentCount"`
	} `json:"statistics"`
}

func getVideoDetails(videoId string) (VideoDetail, error) {
	u := urlVideoDetailsPartOne + videoId + urlVideoDetailsPartTwo + youtubeKey
	req, _ := http.NewRequest("GET", u, nil)
	res, err := cl.Do(req)
	if err != nil {
		return VideoDetail{}, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return VideoDetail{}, err
	}
	var vd VideoDetails
	err = json.Unmarshal(b, &vd)
	if err != nil {
		return VideoDetail{}, err
	}
	if len(vd.Items) < 1 {
		return VideoDetail{}, errors.New("No video")
	}
	return vd.Items[0], nil
}

func getDuratonToSeconds(duration string) (int, error) {
	k := r.FindStringSubmatch(duration)
	if len(k) < 1 {
		return 0, errors.New("not match")
	}
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			if k[i] == "" {
				result[name] = "0"
			} else {
				result[name] = k[i]
			}
		}
	}

	seconds, err := strconv.Atoi(result["seconds"])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(result["minutes"])
	if err != nil {
		return 0, err
	}
	hours, err := strconv.Atoi(result["hours"])
	if err != nil {
		return 0, err
	}
	dur := seconds + minutes*60 + hours*60*60
	return dur, nil
}

func getVideoIdFromUrl(urlVideo string) string {
	u, err := url.Parse(urlVideo)
	if err != nil {
		return ""
	}
	if u.Host == "www.youtube.com" {
		m, _ := url.ParseQuery(u.RawQuery)
		if len(m["v"]) < 1 {
			return ""
		}
		return m["v"][0]
	} else if u.Host == "youtu.be" {
		if len(u.Path) < 1 {
			return ""
		}
		return u.Path[1:]
	}
	return ""
}
