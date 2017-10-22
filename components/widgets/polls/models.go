package polls

import "donategold.me/db"

func init() {
    db.DB.AutoMigrate(&Question{}, &Choice{})
}

type Question struct {
    ID uint `json:"id"`
    UserID uint `json:"user_id"`
    Name string `json:"name"`
    Description string `json:"description"`
    Choices []Choice `json:"choices"`
    Ended bool `json:"-"`
}

func (q *Question) Create() {
    db.DB.Create(q)
}

type Choice struct {
    ID uint `json:"id"`
    QuestionId uint `json:"question_id"`
    Text string `json:"text"`
    Votes uint `json:"votes"`
}

func GetPollsByUser(userId uint) []Question {
    var polls []Question
    db.DB.Preload("Choices").Where("user_id = ?", userId).Find(&polls)
    return polls
}

func DeletePollByUser(userId uint, pollId int) {
    db.DB.Where("user_id = ? and id = ?", userId, pollId).Delete(Question{})
}