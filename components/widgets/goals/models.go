package goals

import (
    "time"
    "donategold.me/db"
)

func init() {
    db.DB.AutoMigrate(&Goal{})
}

type Goal struct {
    ID uint `json:"id"`
    UserId uint
    Name string `json:"name"`
    Description string `json:"description"`
    DateFrom time.Time `json:"date_from"`
    DateTo *time.Time `json:"date_to"`
}

func (g *Goal) Create() {
    db.DB.Create(g)
}

func GetGoalsByUserId(user_id uint) []Goal {
    var goals []Goal
    db.DB.Where("user_id = ?", user_id).Find(&goals)
    return goals
}

func GetGoalById(goal_id uint) (Goal, bool) {
    var goal Goal
    var exist bool
    db.DB.Where("id = ?", goal_id).First(&goal)
    if goal.ID != 0 {
        exist = true
    }
    return goal, exist
}