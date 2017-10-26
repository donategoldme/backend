package polls

import (
	"errors"
	"sync"
	"time"
)

func Default() *PoolOfPolls {
	return &PoolOfPolls{Polls: make(map[uint]*Poll)}
}

type Answer struct {
	Text  string `json:"text"`
	Count uint   `json:"count"`
}

func (a *Answer) add() uint {
	a.Count++
	return a.Count
}

type Poll struct {
	Question string          `json:"question"`
	Answers  []Answer        `json:"answers"`
	Count    uint            `json:"count"`
	Voters   map[string]bool `json:"voters"`
	Time     int64           `json:"time"`
	sync.Mutex
}

func (p *Poll) Vote(answer int, username string) (*Poll, error) {
	p.Lock()
	defer p.Unlock()
	if len(p.Answers) < answer || answer < 1 {
		return p, errors.New("need index of answer")
	}
	if _, ok := p.Voters[username]; ok {
		return p, errors.New(username + ", you votted yet")
	}
	p.Voters[username] = true
	index := answer - 1
	p.Answers[index].add()
	p.Count++
	p.Time = time.Now().UnixNano()
	return p, nil
}

type PoolOfPolls struct {
	Polls map[uint]*Poll
}

func (p *PoolOfPolls) Get(userID uint) (Poll, bool) {
	if poll, ok := p.Polls[userID]; ok {
		return *poll, ok
	}
	return Poll{}, false
}

func (p *PoolOfPolls) Vote(userId uint, answerIndex int, username string) (*Poll, error) {
	if _, ok := p.Polls[userId]; !ok {
		return nil, errors.New("No polls for user")
	}
	return p.Polls[userId].Vote(answerIndex, username)
}

func (p *PoolOfPolls) Clear(userID uint) {
	delete(p.Polls, userID)
}

func (p *PoolOfPolls) NewPoll(pfj PollFromJson) *Poll {
	answers := make([]Answer, len(pfj.Answers))
	for i, answer := range pfj.Answers {
		answers[i] = Answer{answer, 0}
	}
	poll := &Poll{Question: pfj.Question, Answers: answers, Count: 0, Voters: map[string]bool{}}
	p.Polls[pfj.UserID] = poll
	return poll
}

type PollFromJson struct {
	UserID   uint     `json:"user_id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}
