package domain

import "errors"

type Role struct {
	Id          uint `gorm:"primarykey;autoIncrement;not null"`
	ChannelId   int
	Channel     Channel     `gorm:"foreignKey:channel_id; onUpdate:CASCADE; onDelete:CASCADE"`
	Users       []UserModel `gorm:"many2many:user_roles; onUpdate:CASCADE; onDelete:CASCADE"`
	Weight      int         `gorm:"integer;"`
	Permission  int64       `gorm:"type:bigint"`
	Name        string      `gorm:"type:varchar(255);"`
	Description string      `gorm:"type:varchar(255);"`
}

const (
	DefaultRoleName        = "everyone"
	DefaultRoleDescription = "default permissions"
	DefaultRolePermissions = 1380863
)

func CreateDefaultRole(channId int) *Role {
	role := &Role{
		ChannelId:   channId,
		Permission:  DefaultRolePermissions,
		Name:        DefaultRoleName,
		Description: DefaultRoleDescription,
	}

	role.Create()

	return role
}

func (r *Role) TableName() string {
	return "roles"
}

func GetHighestRoleWeight() (int, error) {
	var maxWeight int
	err := Db.Model(&Role{}).Select("MAX(weight)").Row().Scan(&maxWeight)
	if err != nil {
		return 0, err
	}
	return maxWeight, nil
}

func (r *Role) Create() (*Role, error) {

	highestWeight, err := GetHighestRoleWeight()
	if err != nil {
		return nil, err
	}
	r.Weight = highestWeight + 1
	tx := Db.Create(r)
	if tx.RowsAffected == 0 {
		return nil, errors.New("failed to create role")
	}
	return r, nil
}

func (c *Channel) GetRoles() ([]Role, error) {
	var r []Role
	err := Db.Where("channel_id = ?", c.Id).Find(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
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

func (r *Role) Update() (*Role, error) {
	var role *Role
	err := Db.Where("id = ?", r.Id).First(role).Error

	if err != nil {
		return nil, err
	}

	if r.Weight == role.Weight {
		err := Db.Where("id = ?", r.Id).Updates(r).Error

		if err != nil {
			return nil, err
		}

		return r, nil
	}

	// Begin a transaction to ensure consistency
	tx := Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get all roles with a weight greater than or equal to the new weight
	var roles []Role
	if err := tx.Where("weight >= ?", r.Weight).Order("weight asc").Find(&roles).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(r).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Adjust the weights of the other roles
	for _, role := range roles {
		if role.Id != r.Id {

			role.Weight++
			if err := tx.Save(&role).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return r, nil
}

func (r *Role) Delete() *Role {
	tx := Db.Delete(r)
	if tx.RowsAffected == 0 {
		return nil
	}

	return r
}
