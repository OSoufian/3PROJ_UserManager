package domain

type Channel struct {
	Id          uint `gorm:"primarykey;autoIncrement;not null"`
	OwerId      int
	Owner       UserModel   `gorm:"foreignKey:OwerId"`
	Description string      `gorm:"type:varchar(255);"`
	SocialLink  string      `gorm:"type:varchar(255);"`
	Banner      string      `gorm:"type:varchar(255);"`
	Icon        string      `gorm:"type:varchar(255);"`
	Subscribers []UserModel `gorm:"many2many:channel_subscription;"`
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
	tx := Db.Where("OwerId = ?", channel.OwerId).First(channel)

	if tx.RowsAffected == 0 {
		return nil
	}

	return channel
}

func (channel *Channel) GetUserRole(user UserModel) ([]Role, error) {
	var users []UserModel
	var roles []Role
	err := Db.
		Model(&UserModel{}).
		Preload("Channel").
		Find(&users).
		Where("username == ?", user.Username).
		Model(&UserModel{}).
		Preload("roles").
		Find(&roles)
	if err != nil {
		return roles, err.Error
	}

	return roles, nil

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
