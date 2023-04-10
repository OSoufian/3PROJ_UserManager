package domain

type Role struct {
	Id          uint `gorm:"primarykey;autoIncrement;not null"`
	ChannelId   int
	Channel     Channel     `gorm:"foreignKey:ChannelId; onUpdate:CASCADE; onDelete:CASCADE"`
	User        []UserModel `gorm:"many2many:user_roles; onUpdate:CASCADE; onDelete:CASCADE"`
	Permission  uint64      `gorm:"type:bigint"`
	Name        string      `gorm:"type:varchar(255);"`
	Description string      `gorm:"type:varchar(255);"`
}

func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) Create() *Role {
	tx := Db.Create(r)

	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}

func (c *Channel) GetRoles() []Role {
	r := []Role{}
	tx := Db.Where("channel_id = ?", c.Id).Find(r)

	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}

func (u *UserModel) GetRoles() []Role {
	r := []Role{}
	tx := Db.Where("userid = ?", u.Id).Find(r)

	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}

func (r *Role) Get() *Role {
	tx := Db.Where("id = ?", r.Id).First(r)

	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}

func (r *Role) Update() *Role {
	tx := Db.Save(&r)
	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}

func (r *Role) Delete() *Role {
	tx := Db.Delete(r)
	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}
