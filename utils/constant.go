package utils

var AllowedNotificationReferenceType = []interface{}{"thread", "mention", "comment"}
var AllowedNotificationType = []interface{}{"thread", "mention", "comment"}
var AllowedNotificationPriority = []interface{}{"low", "medium", "high"}

var NotificationType = map[string]string{
	"thread":  "THREAD",
	"comment": "COMMENT",
	"mention": "MENTION",
}

var NotificationReferenceType = map[string]string{
	"thread":  "THREAD",
	"comment": "COMMENT",
	"mention": "MENTION",
}

var NotificationPriority = map[string]int32{
	"low":    1,
	"medium": 2,
	"high":   3,
}
