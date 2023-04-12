package domain

type Role struct {
	Id          uint `gorm:"primarykey;autoIncrement;not null"`
	ChannelId   int
	Channel     Channel     `gorm:"foreignKey:channel_id; onUpdate:CASCADE; onDelete:CASCADE"`
	Users       []UserModel `gorm:"many2many:user_roles; onUpdate:CASCADE; onDelete:CASCADE"`
	Weight      int         `gorm:"integer"`
	Permission  int64       `gorm:"type:bigint"`
	Name        string      `gorm:"type:varchar(255);"`
	Description string      `gorm:"type:varchar(255);"`
}

const (
	DefaultRoleName        = "everyone"
	DefaultRoleDescription = "default permissions"
	DefaultRolePermissions = 4607
)

func CreateDefaultRole(channId int) *Role {
	role := &Role{
		ChannelId:   channId,
		Permission:  DefaultRolePermissions,
		Name:        DefaultRoleName,
		Description: DefaultRoleDescription,
	}

	return role.Create()
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

func (r *Role) Get() (*Role, error) {
	err := Db.Where("id = ?", r.Id).First(r).Error

	if err != nil {
		return nil, err
	}

	return r, nil
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
