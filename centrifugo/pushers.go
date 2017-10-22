package centrifugo

import (
	"fmt"
)

func Publish(userId uint, data []byte) error {
	_, err := Cgo.Publish(fmt.Sprintf("$%d/", userId), data)
	return err
}

func PublishChan(userID uint, ch string, data []byte) error {
	_, err := Cgo.Publish(fmt.Sprintf("$%d/%s", userID, ch), data)
	return err
}

func PublishChanAndMain(userID uint, ch string, data []byte) error {
	err := Publish(userID, data)
	if err != nil {
		return err
	}
	err = PublishChan(userID, ch, data)
	if err != nil {
		return err
	}
	return nil
}

func AddMoney(userId uint, money int) error {
	data := []byte(fmt.Sprintf(`{"type": "add_to_balance", "money": %d}`, money))
	err := Publish(userId, data)
	return err
}

func AddPoll(userId uint, pollId int, answer int) error {
	data := []byte(fmt.Sprintf(`{"type": "add_polls", "poll_id": %d, "answer": %d}`, pollId, answer))
	return Publish(userId, data)
}


func AddGoal(userId uint, goalId uint, money int) error {
	data := []byte(fmt.Sprintf(`{"type": "add_goal", "goal_id": %d, "money": %d}`, goalId, money))
	return Publish(userId, data)
}
