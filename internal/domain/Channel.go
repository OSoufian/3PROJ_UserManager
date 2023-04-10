package domain

type Channel struct {
	Id          uint `gorm:"primarykey;autoIncrement;not null"`
	OwnerId     uint `gorm:"not null; foreignKey:id onUpdate:CASCADE; onDelete:CASCADE"`
	Owner       UserModel
	Description string      `gorm:"type:varchar(255);"`
	SocialLink  string      `gorm:"type:varchar(255);"`
	Banner      string      `gorm:"type:varchar(255);"`
	Icon        string      `gorm:"type:varchar(255);"`
	Subscribers []UserModel `gorm:"many2many:channel_subscription;"`
}

type UserChannelPermission struct {
	Permission uint64 `gorm:"column:permission"`
}

func (channel *Channel) TableName() string {
	return "channels"
}

func (channel *Channel) Get() *Channel {
	tx := Db.Where("id = ?", channel.Id).First(channel)

	if tx.RowsAffected == 0 {
		return nil
	}

	return channel
}

func (channel *Channel) GetByOwer() *Channel {
	tx := Db.Where("OwerId = ?", channel.OwnerId).First(channel)

	if tx.RowsAffected == 0 {
		return nil
	}

	return channel
}

func (channel *Channel) GetByVideoId(videoId uint) *Channel {
	err := Db.Joins("JOIN video_info ON channels.id = video_info.channel_id").
		Where("video_info.id = ?", videoId).
		First(channel).Error
	if err != nil {
		return nil
	}

	return channel
}

func (channel *Channel) GetUserRole(user UserModel) (UserChannelPermission, error) {
	var perms UserChannelPermission
	err := Db.
		Joins("JOIN roles r ON r.id = role_id").
		Joins("JOIN users u ON u.id = user_model_id").
		Where("channel_id = ? AND u.id = ?", channel.Id, user.Id).
		Select("r.permission").
		Find(perms).
		Error

	if err != nil {
		return perms, err
	}

	return perms, nil

}

func (channel *Channel) Create() *Channel {
	tx := Db.Create(channel)

	if tx.RowsAffected == 0 {
		return nil
	}

	return channel
}

func (channel *Channel) Update() *Channel {
	tx := Db.Save(&channel)
	if tx.RowsAffected == 0 {
		return nil
	}

	return channel
}

func (channel *Channel) Delete() *Channel {
	tx := Db.Delete(channel)
	if tx.RowsAffected == 0 {
		return nil
	}

	return channel
}
