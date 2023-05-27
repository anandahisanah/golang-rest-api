package models

import "time"

type Item struct {
	ItemId      uint   `gorm:"primaryKey"`
	ItemCode    string `gorm:"not null;type:varchar(191)"`
	Description string `gorm:"type:varchar(191)"`
	Quantity    int    `gorm:"not null;type:int"`
	OrderId     uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
