package chatters

var chatter = New()

func Add(channel string, user uint) {
	chatter.Add(channel, user)
}

func Remove(channel string, user uint) int {
	return chatter.Remove(channel, user)
}

func SubsUser(userID uint) []string {
	return chatter.SubsUser(userID)
}

func GetSubsOfChan(channel string) []uint {
	return chatter.GetSubsOfChan(channel)
}
