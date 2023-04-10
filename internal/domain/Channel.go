package domain

type Channel struct {
	Id          uint `gorm:"primarykey;autoIncrement;not null"`
	OwnerId     uint `gorm:"not null; foreignKey:id onUpdate:CASCADE; onDelete:CASCADE"`
	Owner       UserModel	`json:"-"`
	Description string      `gorm:"type:varchar(255);"`
	SocialLink  string      `gorm:"type:varchar(255);"`
	Banner      string      `gorm:"type:varchar(255);"`
	Icon        string      `gorm:"type:varchar(255);"`
	Subscribers []UserModel `gorm:"many2many:channel_subscription;"`
}

type UserChannelPermission struct {
	Permission int64 `gorm:"column:permission"`
	Weight 	   int 	  `gorm:"column:weight"`
}

type Videos struct {
	Id           uint      `gorm:"primarykey;autoIncrement;not null"`
	Name         string    `gorm:"type:varchar(255);"`
	Description  string    `gorm:"type:varchar(255);"`
	Icon         string    `gorm:"type:varchar(255);"`
	VideoURL     string    `gorm:"type:varchar(255);"`
	Views        int       `gorm:"type:integer default:0"`
	ChannelId    uint      `gorm:"foreignKey:id"`
	CreationDate time.Time `gorm:"type:datetime"`
	IsBlock      bool      `gorm:"type:boolean;default:false"`
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
		Where("channels.id = ?", channel.id).
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
		Joins("JOIN roles r ON r.id = role_id").
		Joins("JOIN users u ON u.id = user_model_id").
		Where("channel_id = ? AND u.id = ?", channel.Id, user.Id).
		Select("r.permission, r.weight").
		Find(perms).
		Order("r.weight DESC").
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
