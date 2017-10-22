package polls

import (
    "gopkg.in/kataras/iris.v6"
    "donategold.me/components/users"
)

func Concat(api iris.MuxAPI) {
    api.Get("/", GetPolls)
    api.Post("/", NewPoll)
    api.Delete("/", DeletePoll)
}

func GetPolls(c *iris.Context) {
    userId := c.Get("token").(users.AccessToken).User.ID
    polls := GetPollsByUser(userId)
    c.JSON(200, polls)
}

func NewPoll(c *iris.Context) {
    userId := c.Get("token").(users.AccessToken).User.ID
    var q Question
    c.ReadJSON(&q)
    q.UserID = userId
    q.Create()
    c.JSON(200, q)
}

func DeletePoll(c *iris.Context) {
    userId := c.Get("token").(users.AccessToken).User.ID
    poll_id, err := c.ParamInt("poll_id")
    if err != nil {
        c.EmitError(404)
        return
    }
    DeletePollByUser(userId, poll_id)
    c.WriteString("ok")
}