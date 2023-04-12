package domain

type Channel struct {
	Id          uint        `gorm:"primarykey;autoIncrement;not null"`
	OwnerId     uint        `gorm:"not null; foreignKey:id onUpdate:CASCADE; onDelete:CASCADE"`
	Owner       UserModel   `json:"-"`
	Description string      `gorm:"type:varchar(255);"`
	SocialLink  string      `gorm:"type:varchar(255);"`
	Banner      string      `gorm:"type:varchar(255);"`
	Icon        string      `gorm:"type:varchar(255);"`
	Subscribers []UserModel `gorm:"many2many:channel_subscription;"`
}

type UserChannelPermission struct {
	Permission int64 `gorm:"column:permission"`
	Weight     int   `gorm:"column:weight"`
}

type Videos struct {
	Id          uint   `gorm:"primarykey;autoIncrement;not null"`
	Name        string `gorm:"type:varchar(255);"`
	Description string `gorm:"type:varchar(255);"`
	Icon        string `gorm:"type:varchar(255);"`
	VideoURL    string `gorm:"type:varchar(255);"`
	Views       int    `gorm:"type:integer default:0"`
	channId     uint   `gorm:"foreignKey:id"`
	CreatedAt   string `gorm:"type:time without time zone"`
	IsBlock     bool   `gorm:"type:boolean;default:false"`
}

func (videos Videos) TableName() string {
	return "video_info"
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

func (channel *Channel) GetAllVideos() []Videos {
	var videos []Videos
	Db.Joins("video_info vo ON vo.channel_id = channels.id").
		Where("channels.id = ?", channel.Id).
		Find(videos)

	return videos
}

func (channel *Channel) GetByOwner() *Channel {
	tx := Db.Where("owner_id = ?", channel.OwnerId).First(channel)

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

func (channel *Channel) GetUserRole(user UserModel) ([]UserChannelPermission, error) {
	var perms []UserChannelPermission
	err := Db.
		Table("channels").
		Select("r.permission, r.weight").
		Joins("JOIN roles r ON r.channel_id = channels.id").
		Joins("JOIN user_roles ur ON ur.role_id = r.id").
		Joins("JOIN users u ON u.id = ur.user_model_id").
		Where("channels.id = ? AND u.id = ?", channel.Id, user.Id).
		Order("r.weight DESC").
		Scan(&perms).
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
