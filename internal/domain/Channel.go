package domain

// import "time"

type Channel struct {
	Id          uint        `gorm:"primarykey;autoIncrement;not null"`
	OwnerId     uint        `gorm:"not null; foreignKey:id onUpdate:CASCADE; onDelete:CASCADE"`
	Owner       UserModel   `json:"-"`
	Name        string      `gorm:"type:varchar(255);"`
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
	Id            uint   `gorm:"primarykey;autoIncrement;not null"`
	Name          string `gorm:"type:varchar(255);"`
	Description   string `gorm:"type:varchar(1500);"`
	Icon          string `gorm:"type:varchar(255);"`
	VideoURL      string `gorm:"type:varchar(255);"`
	Views         int    `gorm:"type:integer"`
	Size          int64  `gorm:"type:integer"`
	ChannelId     uint   `gorm:"foreignKey:id"`
	Channel       Channel
	CreatedAt     string `gorm:"column:created_at"`
	CreationDate  string `gorm:"column:creation_date"`
	IsBlock       bool   `gorm:"type:boolean;default:false"`
	IsHide        bool   `gorm:"type:boolean;default:false"`
}

func (video Videos) TableName() string {
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

// func (channel *Channel) DeleteAllVideos() error {
// 	// Retrieve all videos of the channel
// 	var videos []Videos
// 	if err := Db.Where("SELECT video_info ON video_info.channel_id = channels.id").
// 		Where("channels.id = ?", channel.Id).
// 		Find(&videos).Error; err != nil {
// 		return err
// 	}

// 	// Delete each video
// 	for _, video := range videos {
// 		if err := Db.Delete(&video).Error; err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (channel *Channel) GetByOwner() (*Channel, error) {
	err := Db.Where("owner_id = ?", channel.OwnerId).First(channel).Error

	if err != nil {
		return nil, err
	}

	return channel, nil
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
	tx := Db.Where("id = ?", channel.Id).Updates(&channel)
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
