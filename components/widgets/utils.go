package widgets

import (
	"donategold.me/components/widgets/standard"
	"donategold.me/components/widgets/youtube"
)

func PushWidget(widget string, widgetID uint, userID uint, money int, transID uint) error {
	var err error
	switch widget {
	case "YT":
		err = youtube.CreateFromDonate(widgetID, userID, transID, money)
	case "SD":
		err = standard.CreateFromDonate(widgetID, userID, money, transID)
	default:
		err = standard.CreateFromDonate(widgetID, userID, money, transID)
	}
	err = callBacks(widget, widgetID, userID, money, transID)
	return err
}

func callBacks(widget string, widgetID uint, userID uint, money int, transID uint) error {
	return nil

}
