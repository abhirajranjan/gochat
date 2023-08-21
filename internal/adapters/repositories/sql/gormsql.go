package sql

import (
	"gorm.io/gorm"
)

// give user's channel with their recent message
func getUserChannelswithMessage(tx *gorm.DB, userid string) ([]channelBanner, error) {
	var (
		arrChannelBanner []channelBanner
	)

	joinedtable := tx.Model(&Messages{}).
		Select("messages.*").
		Joins("JOIN user_channels On messages.channel_id = user_channels.channel_id").
		Where("user_channels.user_id = ?", userid)

	windowtable := tx.Table("(?) AS messages", tx.Raw("Select temp.*, row_number() over(partition by channel_id order by created_at DESC) as rn from (?) AS temp", joinedtable)).
		Select("messages.*").
		Where("messages.rn = 1")

	channelMessage := tx.Table("(?) as messages", windowtable).
		Joins("INNER JOIN channels ON messages.channel_id = channels.id").
		Select("messages.id as message_id, messages.user_id as message_user_id, messages.content as message_content, messages.type as message_type, messages.created_at as message_created_at, " +
			"channels.id as channel_id, channels.name as channel_name, channels.picture as channel_picture")

	res := tx.Table("(?) AS cm", channelMessage).
		Joins("INNER JOIN users ON cm.message_user_id = users.id").
		Order("cm.message_created_at DESC").
		Select("users.id as user_id, users.name_tag as user_name_tag, users.given_name as user_given_name, users.family_name as user_family_name, users.picture as user_picture, " +
			"message_id, message_content, message_type, message_created_at, " +
			"channel_id, channel_name, channel_picture").
		Find(&arrChannelBanner)

	return arrChannelBanner, res.Error
}
