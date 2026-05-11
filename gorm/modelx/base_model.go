package modelx

import "time"

type BaseModel struct {
	CreateUserName string    `gorm:"column:create_user_name" json:"create_user_name"`
	CreateUserID   uint64    `gorm:"column:create_user_id" json:"create_user_id"`
	CreateTime     time.Time `gorm:"column:create_time;type:timestamptz(6)" json:"create_time"`
	UpdateUserName string    `gorm:"column:update_user_name" json:"update_user_name"`
	UpdateUserID   uint64    `gorm:"column:update_user_id" json:"update_user_id"`
	UpdateTime     time.Time `gorm:"column:update_time;type:timestamptz(6)" json:"update_time"`
	BaseModelDel
}

type BaseModelNoDel struct {
	CreateUserName string    `gorm:"column:create_user_name" json:"create_user_name"`
	CreateUserID   uint64    `gorm:"column:create_user_id" json:"create_user_id"`
	CreateTime     time.Time `gorm:"column:create_time;type:timestamptz(6)" json:"create_time"`
	UpdateUserName string    `gorm:"column:update_user_name" json:"update_user_name"`
	UpdateUserID   uint64    `gorm:"column:update_user_id" json:"update_user_id"`
	UpdateTime     time.Time `gorm:"column:update_time;type:timestamptz(6)" json:"update_time"`
}

type BaseModelDel struct {
	DeleteUserName string    `gorm:"column:delete_user_name" json:"delete_user_name"`
	DeleteUserID   uint64    `gorm:"column:delete_user_id" json:"delete_user_id"`
	DeleteTime     time.Time `gorm:"column:delete_time;type:timestamptz(6)" json:"delete_time"`
}

type BaseModelTime struct {
	CreateTime time.Time `gorm:"column:create_time;type:timestamptz(6)" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamptz(6)" json:"update_time"`
}

var (
	DBTableOmits = []string{"create_user_name", "create_user_id",
		"delete_user_name", "delete_user_id", "delete_time",
		"update_user_name", "update_user_id", "update_time"}

	DBTableOmitsUpdataDelete = []string{
		"delete_user_name", "delete_user_id", "delete_time",
		"update_user_name", "update_user_id", "update_time"}

	DBTableOmitsUpdate = []string{"update_user_name", "update_user_id", "update_time"}
)
