package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameGoDbUserV2 = "go_db_user_v2"

// GoDbUser mapped from table <go_db_user>
type GoDbUserV2 struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	UUID      string         `gorm:"column:uuid;not null" json:"uuid"`
	UserName  string         `gorm:"column:user_name" json:"user_name"`
	IsActive  bool           `gorm:"column:is_active" json:"is_active"`
}

// TableName GoDbUser's table name
func (*GoDbUserV2) TableName() string {
	return TableNameGoDbUserV2
}
