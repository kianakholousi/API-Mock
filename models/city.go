package models

import (
	"time"
)

type City struct {
	//gorm.Model
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name      string    `gorm:"column:name;not null;size:255" json:"name"`
	CreatedAt time.Time //`gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time //`gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
